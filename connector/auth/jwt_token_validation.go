package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/lestrrat-go/jwx/jwk"
)

var metadataURL = "https://login.botframework.com/v1/.well-known/openidconfiguration"

type TokenValidator interface {
	ValidateToken(authHeader string, credentials CredentialProvider, channelService string, channelID string) ClaimsIdentity
}

type JwtTokenValidation struct {
	Activity   schema.Activity
	AuthHeader string
}

func (jv JwtTokenValidation) AuthenticateRequest(ctx context.Context, activity schema.Activity, authHeader string, credentials CredentialProvider, channelService string) (ClaimsIdentity, error) {
	if authHeader == "" {
		if credentials.IsAuthenticationDisabled() {
			return nil, nil
		}
		return nil, errors.New("Unauthorized Access. Request is not authorized")
	}

	// if IsTokenFromEmulator(authHeader) {
	// 	return nil
	// }

	identity, err := jv.getIdentity(authHeader)
	if err != nil || !identity.IsAuthenticated() {
		return nil, err
	}

	// validate serviceURL
	// This is done outside validateIdentity method to have provision for channel based authentication in future.
	if identity.GetClaimValue("serviceurl") != activity.ServiceUrl {
		return nil, errors.New("Unauthorized, service_url claim is invalid")
	}

	return identity, nil
}

func (jv JwtTokenValidation) getIdentity(authHeader string) (ClaimsIdentity, error) {

	jwksURL, err := jv.getJwkURL(metadataURL)
	if err != nil {
		return nil, err
	}

	getKey := func(token *jwt.Token) (interface{}, error) {

		set, err := jwk.FetchHTTP(jwksURL)
		if err != nil {
			return nil, err
		}

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("Expecting JWT header to have string kid")
		}

		if key := set.LookupKeyID(keyID); len(key) == 1 {
			return key[0].Materialize()
		}

		return nil, errors.New("Could not find public key")
	}

	// TODO: Add options verify_aud and verify_exp
	token, err := jwt.Parse(strings.Split(authHeader, " ")[1], getKey)
	if err != nil {
		return nil, err
	}

	// Check allowed signing algorithms
	alg := token.Header["alg"]
	isAllowed := func() bool {
		for _, allowed := range ALLOWED_SIGNING_ALGORITHMS {
			if allowed == alg {
				return true
			}
		}
		return false
	}()

	if !isAllowed {
		return nil, errors.New("Unauthorized. Invalid signing algorithm")
	}

	claims := token.Claims.(jwt.MapClaims)
	return NewClaimIdentity(claims, true), nil
}

func (jv JwtTokenValidation) validateIdentity(identity ClaimsIdentity, credentials CredentialProvider) error {
	// check issuer
	if identity.GetClaimValue(ISSUER_CLAIM) != TO_BOT_FROM_CHANNEL_TOKEN_ISSUER {
		return errors.New("Unautorized, invlid token issuer")
	}

	// check App ID
	if !credentials.IsValidAppID(identity.GetClaimValue(AUDIENCE_CLAIM)) {
		return errors.New("Unauthorized. Invalid AppId passed on token")
	}

	return nil
}

type metadata struct {
	JwksURI string `json:"jwks_uri"`
}

func (jv JwtTokenValidation) getJwkURL(metadataURL string) (string, error) {
	response, err := http.Get(metadataURL)
	if err != nil {
		return "", errors.New("Error getting metadata document")
	}

	data := metadata{}
	err = json.NewDecoder(response.Body).Decode(&data)
	return data.JwksURI, err
}

func (jv JwtTokenValidation) ValidateAuthHeader(ctx context.Context, authHeader string, channelService, channelID, serviceURL string) (ClaimsIdentity, error) {
	return nil, errors.New("NotImplemented")
}

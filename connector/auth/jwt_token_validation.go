package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/lestrrat-go/jwx/jwk"
)

var metadataURL = "https://login.botframework.com/v1/.well-known/openidconfiguration"

type jwkCache struct {
	Keys   jwk.Set
	Expiry time.Time
}

func (cache *jwkCache) IsExpired() bool {
	if diff := time.Now().Sub(cache.Expiry).Hours(); diff > 0 {
		return true
	}
	return false
}

var cache *jwkCache

// TokenValidator  provides functionanlity to authenticate a request from the connector service.
type TokenValidator interface {
	AuthenticateRequest(ctx context.Context, activity schema.Activity, authHeader string, credentials CredentialProvider, channelService string) (ClaimsIdentity, error)
}

// JwtTokenValidator is the default implementation of TokenValidator.
type JwtTokenValidator struct {
	Activity   schema.Activity
	AuthHeader string
}

// NewJwtTokenValidator return a new TokenValidator value.
func NewJwtTokenValidator() TokenValidator {
	return &JwtTokenValidator{}
}

// AuthenticateRequest autheticates received request from connector service.
//
// The Bearer token is validated for the correct issuer, audience, serviceURL expiry and the signature is verified using the public JWK fetched from BotFramework API.
func (jv *JwtTokenValidator) AuthenticateRequest(ctx context.Context, activity schema.Activity, authHeader string, credentials CredentialProvider, channelService string) (ClaimsIdentity, error) {
	if authHeader == "" {
		if credentials.IsAuthenticationDisabled() {
			return nil, nil
		}
		return nil, errors.New("Unauthorized Access. Request is not authorized")
	}

	identity, err := jv.getIdentity(authHeader)
	if err != nil || !identity.IsAuthenticated() {
		return nil, err
	}

	// Validate serviceURL
	// This is done outside validateIdentity method to have provision for channel based authentication in future.
	if identity.GetClaimValue("serviceurl") != activity.ServiceURL {
		return nil, errors.New("Unauthorized, service_url claim is invalid")
	}

	err = jv.validateIdentity(identity, credentials)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (jv *JwtTokenValidator) getIdentity(authHeader string) (ClaimsIdentity, error) {

	getKey := func(token *jwt.Token) (interface{}, error) {

		if cache == nil || cache.IsExpired() {

			jwksURL, err := jv.getJwkURL(metadataURL)
			if err != nil {
				return nil, err
			}

			set, err := jwk.FetchHTTP(jwksURL)
			if err != nil {
				return nil, err
			}

			cache = &jwkCache{
				Keys:   *set,
				Expiry: time.Now().Add(time.Hour * 24 * 5),
			}
		}

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("Expecting JWT header to have string kid")
		}

		if key := cache.Keys.LookupKeyID(keyID); len(key) == 1 {
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
		for _, allowed := range AllowedSigningAlgorithms {
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

func (jv *JwtTokenValidator) validateIdentity(identity ClaimsIdentity, credentials CredentialProvider) error {
	// check issuer
	if identity.GetClaimValue(IssuerClaim) != ToBotFromChannelTokenIssuer {
		return errors.New("Unautorized, invlid token issuer")
	}

	// check App ID
	if !credentials.IsValidAppID(identity.GetClaimValue(AudienceClaim)) {
		return errors.New("Unauthorized. Invalid AppId passed on token")
	}

	return nil
}

type metadata struct {
	JwksURI string `json:"jwks_uri"`
}

func (jv *JwtTokenValidator) getJwkURL(metadataURL string) (string, error) {
	response, err := http.Get(metadataURL)
	if err != nil {
		return "", errors.New("Error getting metadata document")
	}

	data := metadata{}
	err = json.NewDecoder(response.Body).Decode(&data)
	return data.JwksURI, err
}

// Copyright (c) 2020 InfraCloud Technologies
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/infracloudio/msbotbuilder-go/connector/cache"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/lestrrat-go/jwx/jwk"
)

var metadataURL = "https://login.botframework.com/v1/.well-known/openidconfiguration"

// Timeout for calls fetching the metadata and JWK URLs
const fetchTimeout = 20

// Regular expression to validate the auth header
// The "Bearer " prefix is made optional here
var authHeaderMatch = regexp.MustCompile("^(?:Bearer )?([A-Za-z0-9-_=]+\\.[A-Za-z0-9-_=]+\\.[A-Za-z0-9-_.+/=]*)$")

var httpClient *http.Client

// Init function for the package
func init() {
	// Create a HTTP client with a timeout
	httpClient = &http.Client{
		Timeout: fetchTimeout * time.Second,
	}
}

// TokenValidator provides functionality to authenticate a request from the connector service.
type TokenValidator interface {
	AuthenticateRequest(ctx context.Context, activity schema.Activity, authHeader string, credentials CredentialProvider, channelService string) (ClaimsIdentity, error)
}

// JwtTokenValidator is the default implementation of TokenValidator.
type JwtTokenValidator struct {
	cache.AuthCache
}

// NewJwtTokenValidator returns a new TokenValidator value with an empty cache
func NewJwtTokenValidator() TokenValidator {
	return &JwtTokenValidator{cache.AuthCache{}}
}

// AuthenticateRequest authenticates the received request from connector service.
//
// The Bearer token is validated for the correct issuer, audience, serviceURL expiry and the signature is verified using the public JWK fetched from BotFramework API.
func (jv *JwtTokenValidator) AuthenticateRequest(ctx context.Context, activity schema.Activity, authHeader string, credentials CredentialProvider, channelService string) (ClaimsIdentity, error) {
	// Check the format of the auth header
	match := authHeaderMatch.FindStringSubmatch(strings.TrimSpace(authHeader))
	if len(match) < 2 {
		if credentials.IsAuthenticationDisabled() {
			return nil, nil
		}
		return nil, errors.New("Unauthorized Access. Request is not authorized")
	}

	identity, err := jv.getIdentity(match[1])
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

func (jv *JwtTokenValidator) getIdentity(jwtString string) (ClaimsIdentity, error) {

	getKey := func(token *jwt.Token) (interface{}, error) {

		jwksURL, err := jv.getJwkURL(metadataURL)
		if err != nil {
			return nil, err
		}

		// Get new JWKs if the cache is expired
		if jv.AuthCache.IsExpired() {
			ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout*time.Second)
			defer cancel()
			set, err := jwk.Fetch(ctx, jwksURL)
			if err != nil {
				return nil, err
			}
			// Update the cache
			// The expiry time is set to be of 5 days
			jv.AuthCache = cache.AuthCache{
				Keys:   set,
				Expiry: time.Now().Add(time.Hour * 24 * 5),
			}
		}
		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("Expecting JWT header to have string kid")
		}
		// Return cached JWKs
		key, ok := jv.AuthCache.Keys.(jwk.Set).LookupKeyID(keyID)
		if ok {
			var rawKey interface{}
			err := key.Raw(&rawKey)
			if err != nil {
				return nil, err
			}
			return rawKey, nil
		}

		return nil, errors.New("Could not find public key")
	}

	// TODO: Add options verify_aud and verify_exp
	token, err := jwt.Parse(jwtString, getKey)
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
		return errors.New("Unauthorized: invalid token issuer")
	}

	// check App ID
	if !credentials.IsValidAppID(identity.GetClaimValue(AudienceClaim)) {
		return errors.New("Unauthorized: invalid AppId passed on token")
	}

	return nil
}

type metadata struct {
	JwksURI string `json:"jwks_uri"`
}

func (jv JwtTokenValidator) getJwkURL(metadataURL string) (string, error) {
	response, err := httpClient.Get(metadataURL)
	if err != nil {
		return "", errors.New("Error getting metadata document")
	}

	data := metadata{}
	err = json.NewDecoder(response.Body).Decode(&data)
	return data.JwksURI, err
}

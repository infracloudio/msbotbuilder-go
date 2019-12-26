package auth

import (
	"context"

	"github.com/pkg/errors"

	"github.com/infracloudio/msbotbuilder-go/schema"
)

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

	claimsID := jv.ValidateAuthHeader(ctx, authHeader, credentials, channelService, activity.ChannelId, activity.ServiceUrl)
	// TODO: perform error validation
	return claimsID, nil
}

func (jv JwtTokenValidation) ValidateAuthHeader(ctx context.Context, authHeader string, credentials CredentialProvider, channelService, channelID, serviceURL string) ClaimsIdentity {
	if IsTokenFromEmulator(authHeader) {
		return nil
	}
	return nil
}

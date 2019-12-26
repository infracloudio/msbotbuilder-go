package adapter

import (
	"context"

	"github.com/infracloudio/msbotbuilder-go/connector/auth"
	"github.com/infracloudio/msbotbuilder-go/schema"
)

const (
	OAUTH_ENDPOINT = "https://api.botframework.com"
)

type Setting struct {
	AppID              string
	AppPassword        string
	ChannelAuthTenant  string
	OauthEndpoint      string
	OpenIDMetadata     string
	ChannelService     string
	CredentialProvider auth.CredentialProvider
}

type BotFrameworkAdapter struct {
	Setting
}

func New(settings Setting) *BotFrameworkAdapter {
	// TODO: Support other credential providers - OpenID, MicrosoftApp, Government
	settings.CredentialProvider = auth.SimpleCredentialProvider{
		AppID:    settings.AppID,
		Password: settings.AppPassword,
	}

	if settings.ChannelService == "" {
		settings.ChannelService = auth.CHANNEL_SERVICE
	}
	return &BotFrameworkAdapter{settings}
}

func (bf *BotFrameworkAdapter) ProcessActivity(ctx context.Context, req schema.Activity, headers string) error {
	return bf.AuthenticateRequest(ctx, req, headers)
}

func (bf *BotFrameworkAdapter) AuthenticateRequest(ctx context.Context, req schema.Activity, headers string) error {
	return nil
}

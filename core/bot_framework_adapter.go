package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/infracloudio/msbotbuilder-go/connector/auth"
	"github.com/infracloudio/msbotbuilder-go/connector/client"
	"github.com/infracloudio/msbotbuilder-go/core/activity"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/pkg/errors"
)

type Adapter interface {
	ParseRequest(ctx context.Context, req *http.Request) (schema.Activity, error)
	ProcessActivity(ctx context.Context, req schema.Activity, handler activity.Handler) error
}

type AdapterSetting struct {
	AppID              string
	AppPassword        string
	ChannelAuthTenant  string
	OauthEndpoint      string
	OpenIDMetadata     string
	ChannelService     string
	CredentialProvider auth.CredentialProvider
}

type BotFrameworkAdapter struct {
	AdapterSetting
}

func NewBotAdapter(settings AdapterSetting) Adapter {
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

func (bf *BotFrameworkAdapter) ProcessActivity(ctx context.Context, req schema.Activity, handler activity.Handler) error {
	
	turnContext := &activity.TurnContext{
		Activity: req,
	}

	replyActivity, err := activity.PrepareActivityContext(handler, turnContext)
	if err != nil {
		return err
	}

	connectorClient, err := bf.prepareConnectorClient()
	if err != nil {
		return err
	}

	response, err := activity.NewActivityResponse(connectorClient)
	if err != nil {
		return err
	}

	return response.SendActivity(replyActivity)
}

func (bf *BotFrameworkAdapter) ParseRequest(ctx context.Context, req *http.Request) (schema.Activity, error) {
	activity := schema.Activity{}
	// Find auth headers
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return activity, errors.New("Authentication headers are missing in the request")
	}

	// Parse request body
	err := json.NewDecoder(req.Body).Decode(&activity)
	if err != nil {
		return activity, errors.Wrap(err, "Error while parsing Bot request")
	}
	return activity, bf.authenticateRequest(ctx, activity, authHeader)
}

func (bf *BotFrameworkAdapter) authenticateRequest(ctx context.Context, req schema.Activity, headers string) error {
	jwtValidation := auth.NewJwtTokenValidator()

	_, err := jwtValidation.AuthenticateRequest(ctx, req, headers, bf.CredentialProvider, bf.ChannelService)

	return err
}

func (bf *BotFrameworkAdapter) prepareConnectorClient() (client.Client,error){
	
	clientConfig, err := client.NewClientConfig(bf.AdapterSetting.CredentialProvider, auth.TO_CHANNEL_FROM_BOT_LOGIN_URL[0])
	if err != nil {
		return nil, err
	}

	connectorClient, err := client.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return connectorClient, nil
}

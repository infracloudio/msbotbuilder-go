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

// Adapter is the primary interface for the user program to perform operations with
// the connector service.
type Adapter interface {
	ParseRequest(ctx context.Context, req *http.Request) (schema.Activity, error)
	ProcessActivity(ctx context.Context, req schema.Activity, handler activity.Handler) error
}

// AdapterSetting is the configuration for the Adapter.
type AdapterSetting struct {
	AppID              string
	AppPassword        string
	ChannelAuthTenant  string
	OauthEndpoint      string
	OpenIDMetadata     string
	ChannelService     string
	CredentialProvider auth.CredentialProvider
}

// BotFrameworkAdapter implements Adapter and is currently the only implementation returned to the user program.
type BotFrameworkAdapter struct {
	AdapterSetting
}

// NewBotAdapter creates and reuturns a new BotFrameworkAdapter with the specified AdapterSettings.
func NewBotAdapter(settings AdapterSetting) Adapter {
	// TODO: Support other credential providers - OpenID, MicrosoftApp, Government
	settings.CredentialProvider = auth.SimpleCredentialProvider{
		AppID:    settings.AppID,
		Password: settings.AppPassword,
	}

	if settings.ChannelService == "" {
		settings.ChannelService = auth.ChannelService
	}
	return &BotFrameworkAdapter{settings}
}

// ProcessActivity receives an activity, processes it as specified in by the 'handler' and
// sends it to the connector service.
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

// ParseRequest parses the received activity in a HTTP reuqest to:
//
// 1. Validate the structure.
//
// 2. Authenticate the request (using authenticateRequest())
//
// Returns an Activity value on successfull parsing.
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

func (bf *BotFrameworkAdapter) prepareConnectorClient() (client.Client, error) {

	clientConfig, err := client.NewClientConfig(bf.AdapterSetting.CredentialProvider, auth.ToChannelFromBotLoginURL[0])
	if err != nil {
		return nil, err
	}

	connectorClient, err := client.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return connectorClient, nil
}

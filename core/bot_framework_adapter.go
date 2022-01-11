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
	ProactiveMessage(ctx context.Context, ref schema.ConversationReference, handler activity.Handler) error
	DeleteActivity(ctx context.Context, activityID string, ref schema.ConversationReference) error
	UpdateActivity(ctx context.Context, activity schema.Activity) error
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
	AuthClient         *http.Client
	ReplyClient        *http.Client
}

// BotFrameworkAdapter implements Adapter and is currently the only implementation returned to the user program.
type BotFrameworkAdapter struct {
	AdapterSetting
	auth.TokenValidator
	client.Client
}

// NewBotAdapter creates and reuturns a new BotFrameworkAdapter with the specified AdapterSettings.
func NewBotAdapter(settings AdapterSetting) (Adapter, error) {
	// TODO: Support other credential providers - OpenID, MicrosoftApp, Government
	settings.CredentialProvider = auth.SimpleCredentialProvider{
		AppID:    settings.AppID,
		Password: settings.AppPassword,
	}

	if settings.ChannelService == "" {
		settings.ChannelService = auth.ChannelService
	}

	// Prepare new config and Client
	clientConfig, err := client.NewClientConfig(settings.CredentialProvider, auth.ToChannelFromBotLoginURL[0])
	if err != nil {
		return nil, err
	}

	if settings.AuthClient != nil {
		clientConfig.AuthClient = settings.AuthClient
	}

	if settings.ReplyClient != nil {
		clientConfig.ReplyClient = settings.ReplyClient
	}

	connectorClient, err := client.NewClient(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Connector Client.")
	}

	return &BotFrameworkAdapter{settings, auth.NewJwtTokenValidator(), connectorClient}, nil
}

// ProcessActivity receives an activity, processes it as specified in by the 'handler' and
// sends it to the connector service.
func (bf *BotFrameworkAdapter) ProcessActivity(ctx context.Context, req schema.Activity, handler activity.Handler) error {
	turnContext := &activity.TurnContext{
		Activity: req,
	}

	replyActivity, err := activity.PrepareActivityContext(handler, turnContext)
	if err != nil {
		return errors.Wrap(err, "Failed to create Activity context.")
	}

	response, err := activity.NewActivityResponse(bf.Client)
	if err != nil {
		return errors.Wrap(err, "Failed to create response object.")
	}

	return response.SendActivity(ctx, replyActivity)
}

// ProactiveMessage sends activity to a conversation.
// This methods is used for Bot initiated conversation.
func (bf *BotFrameworkAdapter) ProactiveMessage(ctx context.Context, ref schema.ConversationReference, handler activity.Handler) error {
	// Prepare activity with conversation reference
	activity := activity.ApplyConversationReference(schema.Activity{Type: schema.Message}, ref, true)
	return bf.ProcessActivity(ctx, activity, handler)
}

// DeleteActivity Deletes an existing activity by Activity ID
func (bf *BotFrameworkAdapter) DeleteActivity(ctx context.Context, activityID string, ref schema.ConversationReference) error {
	// Prepare activity with conversation reference
	req := activity.ApplyConversationReference(schema.Activity{Type: schema.Message}, ref, true)
	req.ID = activityID

	response, err := activity.NewActivityResponse(bf.Client)
	if err != nil {
		return errors.Wrap(err, "Failed to create response object.")
	}

	return response.DeleteActivity(ctx, req)
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

	_, err := bf.TokenValidator.AuthenticateRequest(ctx, req, headers, bf.CredentialProvider, bf.ChannelService)

	return errors.Wrap(err, "Authentication failed.")
}

// UpdateActivity Updates an existing activity
func (bf *BotFrameworkAdapter) UpdateActivity(ctx context.Context, req schema.Activity) error {
	response, err := activity.NewActivityResponse(bf.Client)

	if err != nil {
		return errors.Wrap(err, "Failed to create response object.")
	}
	return response.UpdateActivity(ctx, req)
}

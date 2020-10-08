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

package core_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infracloudio/msbotbuilder-go/connector/auth"
	"github.com/infracloudio/msbotbuilder-go/connector/client"
	"github.com/infracloudio/msbotbuilder-go/core"
	"github.com/infracloudio/msbotbuilder-go/core/activity"
	"github.com/infracloudio/msbotbuilder-go/schema"

	"github.com/stretchr/testify/assert"
)

func serverMock(t *testing.T) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/v3/conversations/abcd1234/activities", msTeamsMockMock)
	h1 := &msTeamsActivityUpdateMock{t: t}
	handler.Handle("/v3/conversations/TestActivityUpdate/activities", h1)
	srv := httptest.NewServer(handler)

	return srv
}

func msTeamsMockMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("{\"id\":\"1\"}"))
}

type msTeamsActivityUpdateMock struct {
	t *testing.T
}

func (th *msTeamsActivityUpdateMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	assert.Equal(th.t, "PUT", r.Method, "Expect PUT method")
	activity := schema.Activity{}
	err := json.NewDecoder(r.Body).Decode(&activity)
	assert.Equal(th.t, "TestLabel", activity.Label, "Expect PUT method")
	assert.Nil(th.t, err, fmt.Sprintf("Failed with error %s", err))
	_, _ = w.Write([]byte("{\"id\":\"1\"}"))
}

// Create a handler that defines operations to be performed on respective events.
// Following defines the operation to be performed on the 'message' event.
var customHandler = activity.HandlerFuncs{
	OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
		return turn.SendActivity(activity.MsgOptionText("Echo: " + turn.Activity.Text))
	},
}

func TestExample(t *testing.T) {
	srv := serverMock(t)
	// activity depicts a request as received from a client
	activity := schema.Activity{
		Type: schema.Message,
		From: schema.ChannelAccount{
			ID:   "12345678",
			Name: "Pepper's News Feed",
		},
		Conversation: schema.ConversationAccount{
			ID:   "abcd1234",
			Name: "Convo1",
		},
		Recipient: schema.ChannelAccount{
			ID:   "1234abcd",
			Name: "SteveW",
		},
		Text:       "Message from Teams Client",
		ReplyToID:  "5d5cdc723",
		ServiceURL: srv.URL,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		setting := core.AdapterSetting{
			AppID:       "asdasd",
			AppPassword: "cfg.MicrosoftTeams.AppPassword",
		}
		setting.CredentialProvider = auth.SimpleCredentialProvider{
			AppID:    setting.AppID,
			Password: setting.AppPassword,
		}
		clientConfig, err := client.NewClientConfig(setting.CredentialProvider, auth.ToChannelFromBotLoginURL[0])
		assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
		connectorClient, err := client.NewClient(clientConfig)
		assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
		adapter := core.BotFrameworkAdapter{setting, &core.MockTokenValidator{}, connectorClient}
		act, err := adapter.ParseRequest(ctx, req)
		assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
		err = adapter.ProcessActivity(ctx, act, customHandler)
		assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
	})
	rr := httptest.NewRecorder()
	bodyJSON, err := json.Marshal(activity)
	assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
	bodyBytes := bytes.NewReader(bodyJSON)
	req, err := http.NewRequest(http.MethodPost, "/api/messages", bodyBytes)
	assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
	req.Header.Set("Authorization", "Bearer abc123")
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 200, "Expect 200 response status")
}

func TestActivityUpdate(t *testing.T) {
	srv := serverMock(t)

	activity := schema.Activity{
		Type: schema.Message,
		From: schema.ChannelAccount{
			ID:   "12345678",
			Name: "Pepper's News Feed",
		},
		Conversation: schema.ConversationAccount{
			ID:   "TestActivityUpdate",
			Name: "Convo1",
		},
		Recipient: schema.ChannelAccount{
			ID:   "1234abcd",
			Name: "SteveW",
		},
		Text:       "Message from Teams Client",
		ReplyToID:  "5d5cdc723",
		ServiceURL: srv.URL,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		setting := core.AdapterSetting{
			AppID:       "asdasd",
			AppPassword: "cfg.MicrosoftTeams.AppPassword",
		}
		setting.CredentialProvider = auth.SimpleCredentialProvider{
			AppID:    setting.AppID,
			Password: setting.AppPassword,
		}
		clientConfig, err := client.NewClientConfig(setting.CredentialProvider, auth.ToChannelFromBotLoginURL[0])
		assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
		connectorClient, err := client.NewClient(clientConfig)
		assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
		adapter := core.BotFrameworkAdapter{setting, &core.MockTokenValidator{}, connectorClient}
		act, err := adapter.ParseRequest(ctx, req)
		act.Label = "TestLabel"
		err = adapter.UpdateActivity(ctx, act)
		assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
	})
	rr := httptest.NewRecorder()
	bodyJSON, err := json.Marshal(activity)
	assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
	bodyBytes := bytes.NewReader(bodyJSON)
	req, err := http.NewRequest(http.MethodPost, "/api/messages", bodyBytes)
	assert.Nil(t, err, fmt.Sprintf("Failed with error %s", err))
	req.Header.Set("Authorization", "Bearer abc123")
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 200, "Expect 200 response status")
}

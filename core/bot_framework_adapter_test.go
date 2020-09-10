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
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/infracloudio/msbotbuilder-go/core"
	"github.com/infracloudio/msbotbuilder-go/core/activity"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/stretchr/testify/assert"
)

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/v3/conversations/abcd1234/activities", msTeamsMockMock)

	srv := httptest.NewServer(handler)

	return srv
}

func msTeamsMockMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("{\"id\":\"1\"}"))
}

// Create a handler that defines operations to be performed on respective events.
// Following defines the operation to be performed on the 'message' event.
var customHandler = activity.HandlerFuncs{
	OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
		return turn.SendActivity(activity.MsgOptionText("Echo: " + turn.Activity.Text))
	},
}

func processMessage(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	setting := core.AdapterSetting{
		AppID:       "asdasd",
		AppPassword: "cfg.MicrosoftTeams.AppPassword",
	}
	adapter, err := core.NewBotAdapter(setting)
	act, err := adapter.ParseRequest(ctx, req)
	err = adapter.ProcessActivity(ctx, act, customHandler)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func TestExample(t *testing.T) {
	srv := serverMock()
	// Load settings from environment variables to AdapterSetting.
	setting := core.AdapterSetting{
		AppID:       os.Getenv("APP_ID"),
		AppPassword: os.Getenv("APP_PASSWORD"),
	}

	// Make an adapter to perform operations with the Bot Framework using this library.
	adapter, err := core.NewBotAdapter(setting)
	if err != nil {
		log.Fatal(err)
	}

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

	// Pass the activity and handler to the adapter for proecssing
	ctx := context.Background()
	err = adapter.ProcessActivity(ctx, activity, customHandler)
	if err != nil {
		fmt.Println("Failed to process request", err)
	}
	handler := http.HandlerFunc(processMessage)
	rr := httptest.NewRecorder()
	bodyJson, _ := json.Marshal(activity)
	bodyBytes := bytes.NewReader(bodyJson)
	req, _ := http.NewRequest(http.MethodPost, "/api/messages", bodyBytes)
	req.Header.Set("Authorization", "Bearer abc123")
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 200, "Expect 200 response status")
}

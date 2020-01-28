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
	"context"
	"fmt"
	"log"
	"os"

	"github.com/infracloudio/msbotbuilder-go/core"
	"github.com/infracloudio/msbotbuilder-go/core/activity"
	"github.com/infracloudio/msbotbuilder-go/schema"
)

func Example() {

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

	// Create a handler that defines operations to be performed on respective events.
	// Following defines the operation to be performed on the 'message' event.
	var customHandler = activity.HandlerFuncs{
		OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
			return turn.SendActivity(activity.MsgOptionText("Echo: " + turn.Activity.Text))
		},
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
		Text:      "Message from Teams Client",
		ReplyToID: "5d5cdc723",
	}

	// Pass the activity and handler to the adapter for proecssing
	ctx := context.Background()
	err = adapter.ProcessActivity(ctx, activity, customHandler)
	if err != nil {
		fmt.Println("Failed to process request", err)
		return
	}
}

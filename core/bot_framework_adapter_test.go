package core_test

import (
	"context"
	"fmt"
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
	adapter := core.NewBotAdapter(setting)

	// Create a handler that defines operations to be performed on respective events.
	// Following defines the operation to be performed on the 'message' event.
	var customHandler = activity.HandlerFuncs{
		OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
			activity := turn.Activity
			activity.Text = "Echo: " + activity.Text
			return turn.TextMessage(activity), nil
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
	err := adapter.ProcessActivity(ctx, activity, customHandler)
	if err != nil {
		fmt.Println("Failed to process request", err)
		return
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/infracloudio/msbotbuilder-go/core"
	"github.com/infracloudio/msbotbuilder-go/core/activity"
	"github.com/infracloudio/msbotbuilder-go/schema"
)

// Card content
// Visit: https://adaptivecards.io/explorer to build your own card format
var cardJSON = []byte(`{
  "$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
  "type": "AdaptiveCard",
  "version": "1.0",
  "body": [
    {
      "type": "TextBlock",
      "text": "This is some text",
      "size": "large"
    },
    {
      "type": "TextBlock",
      "text": "It doesn't wrap by default",
      "weight": "bolder"
    },
    {
      "type": "TextBlock",
      "text": "So set **wrap** to true if you plan on showing a paragraph of text",
      "wrap": true
    },
    {
      "type": "TextBlock",
      "text": "You can also use **maxLines** to prevent it from getting out of hand",
      "wrap": true,
      "maxLines": 2
    },
    {
      "type": "TextBlock",
      "text": "You can even draw attention to certain text with color",
      "wrap": true,
      "color": "attention"
    }
  ]
}`)

var customHandler = activity.HandlerFuncs{
	OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
		var obj map[string]interface{}
		err := json.Unmarshal(cardJSON, &obj)
		if err != nil {
			return schema.Activity{}, err
		}
		attachments := []schema.Attachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content:     obj,
			},
		}
		return turn.SendActivity(activity.MsgOptionText("Echo: "+turn.Activity.Text), activity.MsgOptionAttachments(attachments))
	},
}

// HTTPHandler handles the HTTP requests from then connector service
type HTTPHandler struct {
	core.Adapter
}

func (ht *HTTPHandler) processMessage(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	activity, err := ht.Adapter.ParseRequest(ctx, req)
	if err != nil {
		fmt.Println("Failed to parse request.", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ht.Adapter.ProcessActivity(ctx, activity, customHandler)
	if err != nil {
		fmt.Println("Failed to process request.", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Request processed successfully.")
}

func main() {
	setting := core.AdapterSetting{
		AppID:       os.Getenv("APP_ID"),
		AppPassword: os.Getenv("APP_PASSWORD"),
	}

	adapter, err := core.NewBotAdapter(setting)
	if err != nil {
		log.Fatal("Error creating adapter: ", err)
	}

	httpHandler := &HTTPHandler{adapter}

	http.HandleFunc("/api/messages", httpHandler.processMessage)
	fmt.Println("Starting server on port:3978...")
	http.ListenAndServe(":3978", nil)
}

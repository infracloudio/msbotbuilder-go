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

var customHandler = activity.HandlerFuncs{
	OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
		if turn.Activity.Text == "getCard" {
			sJSON := `{
				"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
				"type": "AdaptiveCard",
				"version": "1.0",
				"body": [
				  {
					"type": "TextBlock",
					"text": "Sample"
				  },
				  {  
					"type": "Input.Text",  
					"id": "firstName",  
					"placeholder": "What is your first name?"  
				  }
				],
				"actions": [
				  {
					"type": "Action.Submit",
					"title": "Submit"
				  }
				]
			  }`
			var obj map[string]interface{}
			err := json.Unmarshal(([]byte(sJSON)), &obj)
			if err != nil {
				return schema.Activity{}, nil
			}
			attachments := []schema.Attachment{
				{
					ContentType: "application/vnd.microsoft.card.adaptive",
					Content:     obj,
				},
			}
			return turn.SendActivity(activity.MsgOptionAttachments(attachments))
		}
		if turn.Activity.Value != nil {
			fmt.Println("Activity=", turn.Activity.Value)
			activityID = turn.Activity.ReplyToID
		}

		return turn.SendActivity(activity.MsgOptionText("Echo: " + turn.Activity.Text))
	},
}

var activityID string

// HTTPHandler handles the HTTP requests from then connector service
type HTTPHandler struct {
	core.Adapter
}

func (ht *HTTPHandler) processMessage(w http.ResponseWriter, req *http.Request) {

	ctx := context.Background()
	activityInstance, err := ht.Adapter.ParseRequest(ctx, req)
	if err != nil {
		fmt.Println("Failed to parse request.", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ht.Adapter.ProcessActivity(ctx, activityInstance, customHandler)
	if err != nil {
		fmt.Println("Failed to process request", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	conversationRef := activity.GetCoversationReference(activityInstance)
	act := activity.ApplyConversationReference(schema.Activity{Type: schema.Message}, conversationRef, true)
	if activityID != "" {
		act.Text = "Changed Activity"
		act.ID = activityID
		err = ht.Adapter.UpdateActivity(ctx, act)
		if err != nil {
			fmt.Println("Failed to process request", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err.Error())
			return
		}
		activityID = conversationRef.ActivityID
	} else {
		activityID = conversationRef.ActivityID
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
	fmt.Println("Starting server on port:8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error creating adapter: ", err)
	}
}

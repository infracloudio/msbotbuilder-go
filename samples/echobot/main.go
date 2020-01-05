package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/infracloudio/msbotbuilder-go/core/adapter"
	"github.com/infracloudio/msbotbuilder-go/schema"
)

var Adapter *adapter.BotFrameworkAdapter

func init() {
	setting := adapter.Setting{
		AppID:       os.Getenv("APP_ID"),
		AppPassword: os.Getenv("APP_PASSWORD"),
	}
	Adapter = adapter.New(setting)
}

func main() {
	http.HandleFunc("/api/messages", processMessage)
	http.ListenAndServe(":3978", nil)
}

func processMessage(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	activity := schema.Activity{}
	err := json.NewDecoder(req.Body).Decode(&activity)
	if err != nil {
		fmt.Println("Failed to read body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("ACTIVITY:: %#v\n", activity)

	authHeader := req.Header.Get("Authorization")
	err = Adapter.AuthenticateRequest(ctx, activity, authHeader)
	if err != nil {
		fmt.Println("Failed to authenticate request", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

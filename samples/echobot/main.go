package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/infracloudio/msbotbuilder-go/core"
)

var adapter core.Adapter

func init() {
	setting := core.AdapterSetting{
		AppID:       os.Getenv("APP_ID"),
		AppPassword: os.Getenv("APP_PASSWORD"),
	}
	adapter = core.NewBotAdapter(setting)
}

func main() {
	http.HandleFunc("/api/messages", processMessage)
	fmt.Println("Starting server on port:3978...")
	http.ListenAndServe(":3978", nil)
}

func processMessage(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	activity, err := adapter.ParseRequest(ctx, req)
	if err != nil {
		fmt.Println("Failed to parse request.", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("ACTIVITY:: %#v\n", activity)
}

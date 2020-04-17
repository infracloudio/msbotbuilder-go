package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/infracloudio/msbotbuilder-go/core"
	"github.com/infracloudio/msbotbuilder-go/core/activity"
	"github.com/infracloudio/msbotbuilder-go/schema"
)

func putRequest(u string, data []byte) error {
	client := &http.Client{}
	dec, err := url.QueryUnescape(u)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, dec, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	size := fmt.Sprintf("%d", len(data))
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", size)
	req.Header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", len(data)-1, len(data)))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return fmt.Errorf("failed to upload file with status %d", resp.StatusCode)
	}
	return nil
}

func (ht *HTTPHandler) cleanupConsents(ID string, ref schema.ConversationReference) {
	fmt.Printf("Deleting activity %s\n", ID)
	if err := ht.DeleteActivity(context.TODO(), ID, ref); err != nil {
		log.Printf("Failed to delete activity. %s", err.Error())
	}
}

// HTTPHandler handles the HTTP requests from then connector service
type HTTPHandler struct {
	core.Adapter
}

func (ht *HTTPHandler) processMessage(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	act, err := ht.Adapter.ParseRequest(ctx, req)
	if err != nil {
		fmt.Println("Failed to parse request.", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	aj, _ := json.MarshalIndent(act, "", "  ")
	fmt.Printf("Incoming Activity:: \n%s\n", aj)

	customHandler := activity.HandlerFuncs{
		// handle message events
		OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
			fi, err := os.Stat("data.txt")
			if err != nil {
				return schema.Activity{}, fmt.Errorf("failed to read file: %s", err.Error())
			}

			// send file upload request
			attachments := []schema.Attachment{
				{
					ContentType: "application/vnd.microsoft.teams.card.file.consent",
					Name:        "data.txt",
					Content: map[string]interface{}{
						"description": "Sample data",
						"sizeInBytes": fi.Size(),
					},
				},
			}
			return turn.SendActivity(activity.MsgOptionText("Echo: "+turn.Activity.Text), activity.MsgOptionAttachments(attachments))
		},
		// handle invoke events
		// https://developer.microsoft.com/en-us/microsoft-teams/blogs/working-with-files-in-your-microsoft-teams-bot/
		OnInvokeFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
			ht.cleanupConsents(turn.Activity.ReplyToID, activity.GetCoversationReference(turn.Activity))
			data, err := ioutil.ReadFile("data.txt")
			if err != nil {
				return schema.Activity{}, fmt.Errorf("failed to read file: %s", err.Error())
			}
			if turn.Activity.Value["type"] != "fileUpload" {
				return schema.Activity{}, nil
			}
			if turn.Activity.Value["action"] != "accept" {
				return schema.Activity{}, nil
			}

			// parse upload info from invoke accept response
			uploadInfo := schema.UploadInfo{}
			infoJSON, err := json.Marshal(turn.Activity.Value["uploadInfo"])
			if err != nil {
				return schema.Activity{}, err
			}
			err = json.Unmarshal(infoJSON, &uploadInfo)
			if err != nil {
				return schema.Activity{}, err
			}

			// upload file
			err = putRequest(uploadInfo.UploadURL, data)
			if err != nil {
				return schema.Activity{}, fmt.Errorf("failed to upload file: %s", err.Error())
			}

			// notify user about uploaded file
			fileAttach := []schema.Attachment{
				{
					ContentType: "application/vnd.microsoft.teams.card.file.info",
					ContentURL:  uploadInfo.ContentURL,
					Name:        uploadInfo.Name,
					Content: map[string]interface{}{
						"uniqueId": uploadInfo.UniqueID,
						"fileType": uploadInfo.FileType,
					},
				},
			}

			return turn.SendActivity(activity.MsgOptionAttachments(fileAttach))
		},
	}

	err = ht.Adapter.ProcessActivity(ctx, act, customHandler)
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

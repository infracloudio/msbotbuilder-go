package adapter

import (
	"github.com/infracloudio/msbotbuilder-go/connector/client"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"net/url"
	"strings"
	"path"
	"errors"
	"fmt"
)

const replyURL = "v3/conversations/{conversationId}/activities/{activityId}"

type ConversationOperation struct {
	Client client.ConnectorClient
}

// SendActivity Send an activity to the BOT connector service
func (op *ConversationOperation) SendActivity(conversation Conversation) error {
	if(conversation.Type == schema.MESSAGE) {
		return op.sendTextMessage(conversation)
	}

	return errors.New("No operation for specified Activity type")
}

func (op *ConversationOperation) sendTextMessage(conversation Conversation) error {
	url, err := op.prepareURL(conversation.Activity)
	if err != nil {
		return err
	}
	activity := op.prepareActivity(conversation.Activity)
	err = op.Client.Post(url, activity)
	if err != nil {
		return fmt.Errorf("Error sending reply: %v", err)
	}
	return nil
}


func (op *ConversationOperation) prepareActivity(activity schema.Activity) schema.Activity {
	return schema.Activity{
		Text : activity.Text,
		Type : activity.Type,
		From : activity.Recipient,
		Conversation: activity.Conversation,
		Recipient : activity.From,
		Id : activity.Id,
	}
}

func (op *ConversationOperation) prepareURL(activity schema.Activity) (url.URL, error) {

	u, err := url.Parse(activity.ServiceUrl)

	if err != nil {
		return *u,err
	}
	
	r := strings.NewReplacer("{conversationId}", activity.Conversation.Id,
			"{activityId}", activity.Id)
			
	u.Path = path.Join(u.Path, r.Replace(replyURL))

	return *u, nil
}
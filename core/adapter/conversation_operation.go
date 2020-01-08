package adapter

import (
	"github.com/infracloudio/msbotbuilder-go/connector/client"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"net/url"
	"strings"
	"path"
)

const replyURL = "v3/conversations/{conversationId}/activities/{activityId}"

type ConversationOperation struct {
	Client client.ConnectorClient
}

func (op *ConversationOperation) SendActivity(conversation Conversation) {
	if(conversation.Type == schema.MESSAGE) {
		op.sendTextMessage(conversation)
	}
}

func (op *ConversationOperation) sendTextMessage(conversation Conversation) {
	op.Client.Post(op.prepareURL(conversation.Activity), op.prepareActivity(conversation.Activity))
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

func (op *ConversationOperation) prepareURL(activity schema.Activity) url.URL {

	u, _ := url.Parse(activity.ServiceUrl)
	
	r := strings.NewReplacer("{conversationId}", activity.Conversation.Id,
			"{activityId}", activity.Id)
			
	u.Path = path.Join(u.Path, r.Replace(replyURL))

	return *u
}
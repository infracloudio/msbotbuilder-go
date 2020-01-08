package adapter

import (
	"github.com/infracloudio/msbotbuilder-go/schema"
)

type TurnContext struct {
	Activity schema.Activity
}

func (t *TurnContext) TextMessage(message string) Conversation {
	activity := t.Activity
	activity.Text = message

	return Conversation{
		Type:     schema.MESSAGE,
		Activity: activity,
	}
}

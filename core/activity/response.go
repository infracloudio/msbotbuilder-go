package activity

import (
	"errors"
	"net/url"
	"path"
	"strings"

	"github.com/infracloudio/msbotbuilder-go/connector/client"
	"github.com/infracloudio/msbotbuilder-go/schema"
)

// Response provides functionalities to send activity to the connector service.
type Response interface {
	SendActivity(activity schema.Activity) error
}

const replyToAcitivityURL = "v3/conversations/{conversationId}/activities/{activityId}"

// DefaultResponse is the default implementation of Response.
type DefaultResponse struct {
	Client client.Client
}

// SendActivity sends an activity to the BOT connector service.
func (response *DefaultResponse) SendActivity(activity schema.Activity) error {
	if activity.Type == schema.Message {
		return response.sendTextMessage(activity)
	}

	return errors.New("No operation for specified Activity type")
}

func (response *DefaultResponse) sendTextMessage(activity schema.Activity) error {
	url, err := response.prepareURL(activity)
	if err != nil {
		return err
	}

	activity = response.prepareActivity(activity)
	err = response.Client.Post(url, activity)
	if err != nil {
		return err
	}
	return nil
}

func (response *DefaultResponse) prepareActivity(activity schema.Activity) schema.Activity {
	return schema.Activity{
		Text:      activity.Text,
		Type:      activity.Type,
		From:      activity.Recipient,
		Recipient: activity.From,
		ID:        activity.ID,
	}
}

func (response *DefaultResponse) prepareURL(activity schema.Activity) (url.URL, error) {

	u, err := url.Parse(activity.ServiceURL)
	if err != nil {
		return *u, err
	}

	r := strings.NewReplacer("{conversationId}", activity.Conversation.ID,
		"{activityId}", activity.ID)

	u.Path = path.Join(u.Path, r.Replace(replyToAcitivityURL))

	return *u, nil
}

// NewActivityResponse provides a DefaultResponse implementaton of Response.
func NewActivityResponse(connectorClient client.Client) (Response, error) {
	if connectorClient == nil {
		return nil, errors.New("Invalid connector client for ActivityResponse")
	}

	return &DefaultResponse{connectorClient}, nil
}

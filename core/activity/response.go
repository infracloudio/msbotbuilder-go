package activity

import (
	"errors"
	"net/url"
	"strings"
	"path"

	"github.com/infracloudio/msbotbuilder-go/connector/client"
	"github.com/infracloudio/msbotbuilder-go/schema"
)

type Response interface {
	SendActivity(activity schema.Activity) error
}

const replyToAcitivityURL = "v3/conversations/{conversationId}/activities/{activityId}"


type ActivityResponse struct {
	Client client.Client
}

// SendActivity Send an activity to the BOT connector service
func (response *ActivityResponse) SendActivity(activity schema.Activity) error {
	if activity.Type == schema.MESSAGE {
		return response.sendTextMessage(activity)
	}

	return errors.New("No operation for specified Activity type")
}

func (response *ActivityResponse) sendTextMessage(activity schema.Activity) error {
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

func (response *ActivityResponse) prepareActivity(activity schema.Activity) schema.Activity {
	return schema.Activity{
		Text:      activity.Text,
		Type:      activity.Type,
		From:      activity.Recipient,
		Recipient: activity.From,
		Id:        activity.Id,
	}
}

func (response *ActivityResponse) prepareURL(activity schema.Activity) (url.URL, error) {

	u, err := url.Parse(activity.ServiceUrl)

	if err != nil {
		return *u, err
	}

	r := strings.NewReplacer("{conversationId}", activity.Conversation.Id,
		"{activityId}", activity.Id)

	u.Path = path.Join(u.Path, r.Replace(replyToAcitivityURL))

	return *u, nil
}

func NewActivityResponse(connectorClient client.Client) (Response, error) {
	if connectorClient == nil {
		return nil, errors.New("Invalid connector client for ActivityResponse")
	}

	return &ActivityResponse{connectorClient}, nil
}

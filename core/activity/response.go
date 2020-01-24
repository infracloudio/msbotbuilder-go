// Copyright (c) 2020 InfraCloud Technologies
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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

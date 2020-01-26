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
	"net/url"
	"path"
	"strings"

	"github.com/infracloudio/msbotbuilder-go/connector/client"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/pkg/errors"
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
	if activity.ReplyToID != "" {
		return response.ReplyToActivity(activity.Conversation.ID, activity.ReplyToID, activity)
	}
	return response.SendToConversation(activity.Conversation.ID, activity)
}

// ReplyToActivity sends reply to an activity.
func (response *DefaultResponse) ReplyToActivity(conversationID, activityID string, activity schema.Activity) error {
	url, err := response.prepareReplyToActivityURL(conversationID, activityID, activity.ServiceURL)
	if err != nil {
		return err
	}

	err = response.Client.Post(url, activity)
	if err != nil {
		return errors.Wrap(err, "Failed to send response.")
	}
	return nil
}

// SendToConversation sends an activity to the end of a conversation.
func (response *DefaultResponse) SendToConversation(conversationID string, activity schema.Activity) error {
	// TODO: yet to implement
	return nil
}

func (response *DefaultResponse) prepareReplyToActivityURL(conversationID, activityID, serviceURL string) (url.URL, error) {
	u, err := url.Parse(serviceURL)
	if err != nil {
		return *u, errors.Wrapf(err, "Failed to parse serviceURL %s.", serviceURL)
	}

	r := strings.NewReplacer("{conversationId}", conversationID,
		"{activityId}", activityID)

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

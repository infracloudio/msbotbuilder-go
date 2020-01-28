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
	"github.com/infracloudio/msbotbuilder-go/schema"
)

// GetCoversationReference returns conversation reference from the activity
func GetCoversationReference(activity schema.Activity) schema.ConversationReference {
	return schema.ConversationReference{
		ActivityID:   activity.ID,
		User:         activity.From,
		Bot:          activity.Recipient,
		Conversation: activity.Conversation,
		ChannelID:    activity.ChannelID,
		ServiceURL:   activity.ServiceURL,
	}
}

// ApplyConversationReference sets delivery information to the activity from conversation reference
func ApplyConversationReference(activity schema.Activity, reference schema.ConversationReference, isIncoming bool) schema.Activity {
	activity.ChannelID = reference.ChannelID
	activity.ServiceURL = reference.ServiceURL
	activity.Conversation = reference.Conversation
	if isIncoming {
		activity.From = reference.User
		activity.Recipient = reference.Bot
		if reference.ActivityID != "" {
			activity.ID = reference.ActivityID
		}
		return activity
	}
	activity.From = reference.Bot
	activity.Recipient = reference.User
	if reference.ActivityID != "" {
		activity.ReplyToID = reference.ActivityID
	}
	return activity
}

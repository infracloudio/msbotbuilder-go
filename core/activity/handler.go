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
	"fmt"

	"github.com/infracloudio/msbotbuilder-go/schema"
)

// Handler acts as the interface for the client program to define actions on various events from connector service.
type Handler interface {
	OnMessage(context *TurnContext) (schema.Activity, error)
	OnInvoke(context *TurnContext) (schema.Activity, error)
	OnConversationUpdate(context *TurnContext) (schema.Activity, error)
}

// HandlerFuncs is an adaptor to let client program specify as many or
// as few functions to handle events of the connector service while still implementing
// Handler.
type HandlerFuncs struct {
	OnMessageFunc            func(turn *TurnContext) (schema.Activity, error)
	OnInvokeFunc             func(turn *TurnContext) (schema.Activity, error)
	OnConversationUpdateFunc func(turn *TurnContext) (schema.Activity, error)
}

// OnMessage handles a 'message' event from connector service.
func (r HandlerFuncs) OnMessage(turn *TurnContext) (schema.Activity, error) {
	if r.OnMessageFunc != nil {
		return r.OnMessageFunc(turn)
	}
	return schema.Activity{}, errors.New("No handler found for this activity type")
}

// OnConversationUpdate handles a 'conversationUpdate' event from connector service.
func (r HandlerFuncs) OnConversationUpdate(turn *TurnContext) (schema.Activity, error) {
	if r.OnConversationUpdateFunc != nil {
		return r.OnConversationUpdateFunc(turn)
	}
	return schema.Activity{}, errors.New("No handler found for this activity type")
}

// OnInvoke handles a 'invoke' event from connector service.
func (r HandlerFuncs) OnInvoke(turn *TurnContext) (schema.Activity, error) {
	if r.OnInvokeFunc != nil {
		return r.OnInvokeFunc(turn)
	}
	return schema.Activity{}, errors.New("No handler found for this activity type")
}

// PrepareActivityContext routes the received Activity to respective handler function.
// Returns the result of the handler function.
func PrepareActivityContext(handler Handler, context *TurnContext) (schema.Activity, error) {
	switch context.Activity.Type {
	case schema.Message:
		return handler.OnMessage(context)
	case schema.Invoke:
		return handler.OnInvoke(context)
	case schema.ConversationUpdate:
		return handler.OnConversationUpdate(context)
	}
	return schema.Activity{}, fmt.Errorf("Activity type %s not supported yet", context.Activity.Type)
}

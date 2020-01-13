package activity

import (
	"errors"
	"github.com/infracloudio/msbotbuilder-go/schema"
)

// Handler acts as the interface for the client program to define actions on various events from connector service.
type Handler interface {
	OnMessage(context *TurnContext) (schema.Activity, error)
}

// HandlerFuncs is an adaptor to let client program specify as many or
// as few functions to handle events of the connector service while still implementing
// Handler.
type HandlerFuncs struct {
	OnMessageFunc func(turn *TurnContext) (schema.Activity, error)
}

// OnMessage handles a 'message' event from connector service.
func (r HandlerFuncs) OnMessage(turn *TurnContext) (schema.Activity, error) {
	if r.OnMessageFunc != nil {
		return r.OnMessageFunc(turn)
	}
	return schema.Activity{}, errors.New("No handler found for this activity type")
}

// PrepareActivityContext routes the received Activity to respective handler function. 
// Returns the result of the handler function.
func PrepareActivityContext(handler Handler, context *TurnContext) (schema.Activity, error) {
	if context.Activity.Type == schema.Message {
		activity, err := handler.OnMessage(context)
		if err != nil {
			return schema.Activity{}, err
		}
		return activity, nil
	}
	return schema.Activity{}, errors.New("Malformed Activity : Type not recognized")
}

package activity

import (
	"errors"
	"github.com/infracloudio/msbotbuilder-go/schema"
)

type Handler interface {
	OnMessage(context *TurnContext) (schema.Activity, error)
}

type HandlerFuncs struct {
	OnMessageFunc func(turn *TurnContext) (schema.Activity, error)
}

func (r HandlerFuncs) OnMessage(turn *TurnContext) (schema.Activity, error) {
	if r.OnMessageFunc != nil {
		return r.OnMessageFunc(turn)
	}
	return schema.Activity{}, errors.New("No handler found for this activity type")
}

func PrepareActivityContext(handler Handler, context *TurnContext) (schema.Activity, error) {
	if context.Activity.Type == schema.MESSAGE {
		activity, err := handler.OnMessage(context)
		if err != nil {
			return schema.Activity{}, err
		}
		return activity, nil
	}
	return schema.Activity{}, errors.New("Malformed Activity : Type not recognized")
}

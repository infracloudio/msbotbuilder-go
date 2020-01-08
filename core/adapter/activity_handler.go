package adapter


import (
	"github.com/infracloudio/msbotbuilder-go/schema"
)

type ActivityHandler interface {
	OnMessage(turn *TurnContext) interface{}
}

type ActivityHandlerFuncs struct {
	OnMessageFunc func(turn *TurnContext) interface{}
}

func (r ActivityHandlerFuncs) OnMessage(turn *TurnContext) interface{} {
	if r.OnMessageFunc != nil {
		return r.OnMessageFunc(turn)
	}
	return nil
}

func Activate(handler ActivityHandler, context *TurnContext) interface{} {
	if context.Activity.Type == schema.MESSAGE {
		return handler.OnMessage(context)
	}
	return nil
}

package adapter


import (
	"github.com/infracloudio/msbotbuilder-go/schema"
)

type ActivityHandler interface {
	OnMessage(turn *TurnContext) interface{}
}

type ActivityHandlerFuncs struct {
	MessageFuntion func(turn *TurnContext) interface{}
}

func (r ActivityHandlerFuncs) OnMessage(turn *TurnContext) interface{} {
	if r.MessageFuntion != nil {
		return r.MessageFuntion(turn)
	}
	return nil
}

func Activate(handler ActivityHandler, context *TurnContext) interface{} {
	if context.Activity.Type == schema.MESSAGE {
		return handler.OnMessage(context)
	}
	return nil
}

package activity

import (
	"github.com/infracloudio/msbotbuilder-go/schema"
)

type TurnContext struct {
	Activity schema.Activity
}

func (t *TurnContext) TextMessage(activity schema.Activity) schema.Activity {
	return activity
}

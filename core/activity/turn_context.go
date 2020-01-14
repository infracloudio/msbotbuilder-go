package activity

import (
	"github.com/infracloudio/msbotbuilder-go/schema"
)

// TurnContext wraps the Activity received and provides operations for the user
// program of this SDK.
//
// The return value is Activity as provided by the client program, to be send to the connector service.
type TurnContext struct {
	Activity schema.Activity
}

// TextMessage function to be used by the client program to send a plain text Activity.
func (t *TurnContext) TextMessage(activity schema.Activity) schema.Activity {
	return activity
}

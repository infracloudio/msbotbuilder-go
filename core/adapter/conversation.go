package adapter

import "github.com/infracloudio/msbotbuilder-go/schema"

type Conversation struct {
	Type schema.ActivityTypes
	Activity schema.Activity
}

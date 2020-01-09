package activity

import (
	"github.com/infracloudio/msbotbuilder-go/schema"
)

type Response interface {
	SendActivity(activity schema.Activity) error
}

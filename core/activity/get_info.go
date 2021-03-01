package activity

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/pkg/errors"
)

const getMemberInfo = "/%s/conversations/%s/members/%s"

// GetSenderInfo get conversation member info.
func (response *DefaultResponse) GetSenderInfo(activity schema.Activity) (*schema.ConversationMember, error) {
	u, err := url.Parse(activity.ServiceURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse ServiceURL %s.", activity.ServiceURL)
	}

	urlString := fmt.Sprintf(getMemberInfo, APIVersion, activity.Conversation.ID, activity.From.ID)

	u.Path = path.Join(u.Path, urlString)
	raw, err := response.Client.Get(*u)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get response.")
	}

	var convMember schema.ConversationMember

	err = json.Unmarshal(raw, &convMember)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse response.")
	}

	return &convMember, nil
}

package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/infracloudio/msbotbuilder-go/connector/auth"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/infracloudio/msbotbuilder-go/schema/customerror"
)

type ConnectorClient struct {
	Config ConnectorClientConfig
}

// Post a activity to given URL
func (client *ConnectorClient) Post(target url.URL, activity schema.Activity) error {

	token, err := client.getToken()

	if err != nil {
		return err
	}

	jsonStr, _ := json.Marshal(activity)
	req, err := http.NewRequest("POST", target.String(), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	replyClient := &http.Client{}
	resp, err := replyClient.Do(req)
	defer resp.Body.Close()

	if err != nil || !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
		return customerror.HttpError {
			StatusCode: resp.StatusCode,
			HtErr : err,
			Body : resp.Body,
		}
	}
	
	return nil
}

func (client *ConnectorClient) getToken() (string, error) {

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", client.Config.Credentials.GetAppId())
	data.Set("client_secret", client.Config.Credentials.GetAppPassword())
	data.Set("scope", auth.TO_CHANNEL_FROM_BOT_OAUTH_SCOPE)

	u, err := url.ParseRequestURI(client.Config.AuthURL.String())

	if err != nil {
		return "",err
	}

	authClient := &http.Client{}

	r, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	if err != nil {
		return "", err
	}

	resp, err := authClient.Do(r)
	defer resp.Body.Close()

	if err != nil {
		return "", customerror.HttpError {
			StatusCode: resp.StatusCode,
			HtErr : err,
			Body : resp.Body,
		}
	}

	a := &schema.AuthResponse{}
	json.NewDecoder(resp.Body).Decode(a)

	return a.AccessToken, nil
}

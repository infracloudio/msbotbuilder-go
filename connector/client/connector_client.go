package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/infracloudio/msbotbuilder-go/connector/auth"
)

type ConnectorClient struct {
	Config ConnectorClientConfig
}

func (client *ConnectorClient) Post(target url.URL, activity schema.Activity) (int, error) {

	token := client.getToken()

	jsonStr, _ := json.Marshal(activity)
	req, err := http.NewRequest("POST", target.String(), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	replyClient := &http.Client{}
	resp, err := replyClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, err
}

func (client *ConnectorClient) getToken() string {

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", client.Config.Credentials.GetAppId())
	data.Set("client_secret", client.Config.Credentials.GetAppPassword())
	data.Set("scope", auth.TO_CHANNEL_FROM_BOT_OAUTH_SCOPE)

	u, _ := url.ParseRequestURI(client.Config.AuthURL.String())

	authClient := &http.Client{}

	r, _ := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := authClient.Do(r)
	defer resp.Body.Close()

	a := &schema.AuthResponse{}
	json.NewDecoder(resp.Body).Decode(a)

	return a.AccessToken
}

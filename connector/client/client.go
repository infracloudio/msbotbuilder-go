package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"fmt"
	"errors"

	"github.com/infracloudio/msbotbuilder-go/connector/auth"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/infracloudio/msbotbuilder-go/schema/customerror"
)

// Client provides interface to send requests to the connector service.
type Client interface {
	Post(url url.URL, activity schema.Activity) error
}

// ConnectorClient implements Client to send HTTP requests to the connector service.
type ConnectorClient struct {
	Config
}

// NewClient constructs and returns a new ConnectorClient with provided configuration.
// Returns error if Config passed is nil.
func NewClient(config *Config) (Client, error) {
	
	if config == nil {
		return nil, errors.New("Invalid client configuration")
	}
	
	return &ConnectorClient{*config},nil
}

// Post a activity to given URL.
// 
// Creates a HTTP POST request with the provided activity as the body and a Bearer token in the header.
// Returns any error as received from the call to connector service.
func (client ConnectorClient) Post(target url.URL, activity schema.Activity) error {

	token, err := client.getToken()
	if err != nil {
		return err
	}

	jsonStr, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", target.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	replyClient := &http.Client{}
	
	resp, err := replyClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		return customerror.HTTPError{
			StatusCode: resp.StatusCode,
			HtErr:      err,
			Body:       resp.Body,
		}
	}

	defer resp.Body.Close()

	return nil
}

func (client *ConnectorClient) getToken() (string, error) {

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", client.Credentials.GetAppID())
	data.Set("client_secret", client.Credentials.GetAppPassword())
	data.Set("scope", auth.ToChannelFromBotOauthScope)

	u, err := url.ParseRequestURI(client.AuthURL.String())
	if err != nil {
		return "", err
	}

	authClient := &http.Client{}
	r, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := authClient.Do(r)
	if err != nil {
		return "", customerror.HTTPError{
			StatusCode: resp.StatusCode,
			HtErr:      err,
			Body:       resp.Body,
		}
	}

	defer resp.Body.Close()

	a := &schema.AuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(a)
	if err != nil {
		return "", fmt.Errorf("Invalid activity to send", err)
	}

	return a.AccessToken, nil
}

// Copyright (c) 2020 InfraCloud Technologies
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/infracloudio/msbotbuilder-go/connector/auth"
	"github.com/infracloudio/msbotbuilder-go/connector/cache"
	"github.com/infracloudio/msbotbuilder-go/schema"
	"github.com/infracloudio/msbotbuilder-go/schema/customerror"
)

// Client provides interface to send requests to the connector service.
type Client interface {
	Post(url url.URL, activity schema.Activity) error
	Delete(url url.URL, activity schema.Activity) error
	Get(url url.URL) ([]byte, error)
}

// ConnectorClient implements Client to send HTTP requests to the connector service.
type ConnectorClient struct {
	Config
	cache.AuthCache
}

// NewClient constructs and returns a new ConnectorClient with provided configuration and an empty cache.
// Returns error if Config passed is nil.
func NewClient(config *Config) (Client, error) {
	if config == nil {
		return nil, errors.New("Invalid client configuration")
	}

	return &ConnectorClient{*config, cache.AuthCache{}}, nil
}

// Get fetches data from given URL.
func (client *ConnectorClient) Get(target url.URL) ([]byte, error) {

	token, err := client.getToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", target.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Post an activity to given URL.
//
// Creates a HTTP POST request with the provided activity as the body and a Bearer token in the header.
// Returns any error as received from the call to connector service.
func (client *ConnectorClient) Post(target url.URL, activity schema.Activity) error {
	jsonStr, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, target.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	return client.sendRequest(req, activity)
}

// Delete an activity.
//
// Creates a HTTP DELETE request with the provided activity ID and a Bearer token in the header.
// Returns any error as received from the call to connector service.
func (client *ConnectorClient) Delete(target url.URL, activity schema.Activity) error {
	req, err := http.NewRequest(http.MethodDelete, target.String(), nil)
	if err != nil {
		return err
	}
	return client.sendRequest(req, activity)
}

func (client *ConnectorClient) sendRequest(req *http.Request, activity schema.Activity) error {
	token, err := client.getToken()
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	replyClient := &http.Client{}

	return client.checkRespError(replyClient.Do(req))
}

func (client *ConnectorClient) checkRespError(resp *http.Response, err error) error {
	allowedResp := []int{http.StatusOK, http.StatusCreated, http.StatusAccepted}
	if err != nil {
		return customerror.HTTPError{
			HtErr: err,
		}
	}
	defer resp.Body.Close()

	// Check if resp allowed
	for _, code := range allowedResp {
		if code == resp.StatusCode {
			return nil
		}
	}

	return customerror.HTTPError{
		HtErr:      errors.New("invalid response"),
		StatusCode: resp.StatusCode,
	}
}

func (client *ConnectorClient) getToken() (string, error) {

	// Return cached JWT
	if !client.AuthCache.IsExpired() {
		return client.AuthCache.Keys.(string), nil
	}

	// Get new JWT
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
		}
	}

	defer resp.Body.Close()

	a := &schema.AuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(a)
	if err != nil {
		return "", fmt.Errorf("Invalid activity to send %s", err)
	}

	// Update cache
	client.AuthCache = cache.AuthCache{
		Keys:   a.AccessToken,
		Expiry: time.Now().Add(time.Second * time.Duration(a.ExpireTime)),
	}

	return client.AuthCache.Keys.(string), nil
}

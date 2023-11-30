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
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	Post(ctx context.Context, url url.URL, activity schema.Activity) error
	Delete(ctx context.Context, url url.URL) error
	Get(ctx context.Context, url url.URL) (json.RawMessage, error)
	Put(ctx context.Context, url url.URL, activity schema.Activity) error
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

	if config.AuthClient == nil {
		config.AuthClient = &http.Client{}
	}

	if config.ReplyClient == nil {
		config.ReplyClient = &http.Client{}
	}

	return &ConnectorClient{*config, cache.AuthCache{}}, nil
}

// Post an activity to given URL.
//
// Creates a HTTP POST request with the provided activity as the body and a Bearer token in the header.
// Returns any error as received from the call to connector service.
func (client *ConnectorClient) Post(ctx context.Context, target url.URL, activity schema.Activity) error {
	jsonStr, err := json.Marshal(activity)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	return client.sendRequestWithRespErrCheck(req)
}

// Get a resource from given URL using authenticated request.
//
// This method is helpful for obtaining Teams context for your bot.
// Read more: https://learn.microsoft.com/en-us/microsoftteams/platform/bots/how-to/get-teams-context?tabs=json
func (client *ConnectorClient) Get(ctx context.Context, target url.URL) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := client.sendRequest(req)
	if err != nil {
		return nil, newHTTPError(err)
	}
	defer res.Body.Close()

	if wrappedErr := client.checkRespError(res); wrappedErr != nil {
		return nil, wrappedErr
	}

	var rawOutput json.RawMessage
	err = json.NewDecoder(res.Body).Decode(&rawOutput)
	if err != nil {
		return nil, err
	}

	return rawOutput, nil
}

// Delete an activity.
//
// Creates a HTTP DELETE request with the provided activity ID and a Bearer token in the header.
// Returns any error as received from the call to connector service.
func (client *ConnectorClient) Delete(ctx context.Context, target url.URL) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, target.String(), nil)
	if err != nil {
		return err
	}
	return client.sendRequestWithRespErrCheck(req)
}

// Put an activity.
//
// Creates a HTTP PUT request with the provided activity payload and a Bearer token in the header.
// Returns any error as received from the call to connector service.
func (client *ConnectorClient) Put(ctx context.Context, target url.URL, activity schema.Activity) error {
	jsonStr, err := json.Marshal(activity)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, target.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	return client.sendRequestWithRespErrCheck(req)
}

func (client *ConnectorClient) sendRequestWithRespErrCheck(req *http.Request) error {
	res, err := client.sendRequest(req)
	if err != nil {
		return newHTTPError(err)
	}

	defer res.Body.Close()
	return client.checkRespError(res)
}

func (client *ConnectorClient) sendRequest(req *http.Request) (*http.Response, error) {
	token, err := client.getToken(req.Context())
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return client.ReplyClient.Do(req)
}

func (client *ConnectorClient) checkRespError(resp *http.Response) error {
	allowedResp := []int{http.StatusOK, http.StatusCreated, http.StatusAccepted}
	// Check if resp allowed
	for _, code := range allowedResp {
		if code == resp.StatusCode {
			return nil
		}
	}

	return newHTTPErrorWithStatusCode(errors.New("invalid response"), resp.StatusCode)
}

func (client *ConnectorClient) getToken(ctx context.Context) (string, error) {

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

	r, err := http.NewRequestWithContext(ctx, "POST", u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.AuthClient.Do(r)
	if err != nil {
		return "", newHTTPErrorWithStatusCode(err, resp.StatusCode)
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

func newHTTPError(err error) error {
	return customerror.HTTPError{
		HtErr: err,
	}
}

func newHTTPErrorWithStatusCode(err error, statusCode int) error {
	return customerror.HTTPError{
		HtErr:      err,
		StatusCode: statusCode,
	}
}

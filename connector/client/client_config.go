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
	"errors"
	"github.com/infracloudio/msbotbuilder-go/connector/auth"
	"net/http"
	"net/url"
)

// Config represents the credentials for a user program and the URL for validating the credentials.
type Config struct {
	Credentials auth.CredentialProvider
	AuthURL     url.URL
	AuthClient  *http.Client
	ReplyClient *http.Client
}

// NewClientConfig creates configuration for ConnectorClient.
func NewClientConfig(credentials auth.CredentialProvider, tokenURL string) (*Config, error) {

	if len(credentials.GetAppID()) < 0 || len(credentials.GetAppPassword()) < 0 {
		return &Config{}, errors.New("Invalid client credentials")
	}

	parsedURL, err := url.Parse(tokenURL)
	if err != nil {
		return &Config{}, errors.New("Invalid token URL")
	}

	return &Config{
		Credentials: credentials,
		AuthURL:     *parsedURL,
	}, nil
}

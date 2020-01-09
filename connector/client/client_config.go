package client

import (
	"errors"
	"net/url"
	"github.com/infracloudio/msbotbuilder-go/connector/auth"
)

type ConnectorClientConfig struct {
	Credentials auth.CredentialProvider
	AuthURL     url.URL
}

// NewClientConfig creates configuration for ConnectorClient
func NewClientConfig(credentials auth.CredentialProvider, tokenURL string) (ConnectorClientConfig, error) {


	if len(credentials.GetAppId()) < 0 || len(credentials.GetAppPassword()) < 0 {
		errors.New("Invalid client credentials")
	}

	parsedURL, err := url.Parse(tokenURL)

	if err != nil {
		return ConnectorClientConfig{}, errors.New("Invalid token URL")
	}

	return ConnectorClientConfig{
		Credentials: credentials,
		AuthURL:     *parsedURL,
	}, nil
}

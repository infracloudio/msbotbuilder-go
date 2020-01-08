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

func NewClientConfig(credentials auth.CredentialProvider, tokenURL string) (ConnectorClientConfig, error) {

	parsedURL, err := url.Parse(tokenURL)

	if err != nil {
		return ConnectorClientConfig{}, errors.New("Invalid token URL")
	}

	return ConnectorClientConfig{
		Credentials: credentials,
		AuthURL:     *parsedURL,
	}, nil
}

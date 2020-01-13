package auth

// CredentialProvider represents and provides functionality for a type of Credential.
type CredentialProvider interface {
	IsValidAppID(appID string) bool
	GetAppPassword() string
	GetAppID() string
	IsAuthenticationDisabled() bool
}

// SimpleCredentialProvider can be used for authentication to the connector service using
// AppID and Password.
type SimpleCredentialProvider struct {
	AppID    string
	Password string
}

// IsValidAppID returns if the specified appID is valid.
func (sp SimpleCredentialProvider) IsValidAppID(appID string) bool {
	return sp.AppID == appID
}

// GetAppPassword returns the Password of the credential.
func (sp SimpleCredentialProvider) GetAppPassword() string {
	return sp.Password
}

// GetAppID returns the AppID of the credential.
func (sp SimpleCredentialProvider) GetAppID() string {
	return sp.AppID
}

// IsAuthenticationDisabled checks if no authentication is to be performed.
func (sp SimpleCredentialProvider) IsAuthenticationDisabled() bool {
	return sp.AppID == ""
}

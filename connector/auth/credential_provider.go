package auth

type CredentialProvider interface {
	IsValidAppID(appID string) bool
	GetAppPassword() string
	GetAppId() string
	IsAuthenticationDisabled() bool
}

type SimpleCredentialProvider struct {
	AppID    string
	Password string
}

func (sp SimpleCredentialProvider) IsValidAppID(appID string) bool {
	return sp.AppID == appID
}

func (sp SimpleCredentialProvider) GetAppPassword() string {
	return sp.Password
}

func (sp SimpleCredentialProvider) GetAppId() string {
	return sp.AppID
}

func (sp SimpleCredentialProvider) IsAuthenticationDisabled() bool {
	return sp.AppID == ""
}

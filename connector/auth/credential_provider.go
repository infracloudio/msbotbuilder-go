package auth

type CredentialProvider interface {
	IsValidAppID(appID string) bool
	GetAppPassword(appID string) string
	IsAuthenticationDisabled() bool
}

type SimpleCredentialProvider struct {
	AppID    string
	Password string
}

func (sp SimpleCredentialProvider) IsValidAppID(appID string) bool {
	return sp.AppID == appID
}

func (sp SimpleCredentialProvider) GetAppPassword(appID string) string {
	if sp.AppID != appID {
		return ""
	}
	return sp.Password
}

func (sp SimpleCredentialProvider) IsAuthenticationDisabled() bool {
	return sp.AppID == ""
}

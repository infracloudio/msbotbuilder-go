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

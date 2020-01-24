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

// Claim represents a claim in a JWT token.
type Claim interface{}

// ClaimsIdentity is the interface to process claims in a JWT token.
type ClaimsIdentity interface {
	GetClaimValue(string) string
	IsAuthenticated() bool
}

// DefaultClaim is the default implementation fo Claim.
type DefaultClaim struct {
	Type  string
	Value string
}

// NewClaim contructs and returns a new Claim value.
func NewClaim(tpe, val string) Claim {
	return &DefaultClaim{
		Type:  tpe,
		Value: val,
	}
}

// DefaultClaimIdentity implements ClaimsIdentity to create and process Claim values.
type DefaultClaimIdentity struct {
	claims          map[string]interface{}
	isAuthenticated bool
}

// NewClaimIdentity creates and returns a new ClaimsIdentity value.
func NewClaimIdentity(claims map[string]interface{}, isAuth bool) ClaimsIdentity {
	return &DefaultClaimIdentity{claims, isAuth}
}

// GetClaimValue returns value for a specified property of a claim.
func (ci DefaultClaimIdentity) GetClaimValue(cType string) string {
	return ci.claims[cType].(string)
}

// IsAuthenticated returns if the Claim is authenticated.
func (ci DefaultClaimIdentity) IsAuthenticated() bool {
	return ci.isAuthenticated
}

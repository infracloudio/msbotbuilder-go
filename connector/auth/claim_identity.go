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

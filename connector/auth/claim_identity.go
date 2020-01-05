package auth

type Claim interface{}
type ClaimsIdentity interface {
	GetClaimValue(string) string
	IsAuthenticated() bool
}

type DefaultClaim struct {
	Type  string
	Value string
}

func NewClaim(tpe, val string) Claim {
	return &DefaultClaim{
		Type:  tpe,
		Value: val,
	}
}

type DefaultClaimIdentity struct {
	claims          map[string]interface{}
	isAuthenticated bool
}

func NewClaimIdentity(claims map[string]interface{}, isAuth bool) ClaimsIdentity {
	return &DefaultClaimIdentity{claims, isAuth}
}

func (ci DefaultClaimIdentity) GetClaimValue(cType string) string {
	return ci.claims[cType].(string)
}

func (ci DefaultClaimIdentity) IsAuthenticated() bool {
	return ci.isAuthenticated
}

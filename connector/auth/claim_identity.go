package auth

type Claim interface{}
type ClaimsIdentity interface {
	GetClaimValue(string) string
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
	Claims          map[string]string
	IsAuthenticated bool
}

func NewClaimIdentity(claims map[string]string, isAuth bool) ClaimsIdentity {
	return &DefaultClaimIdentity{
		Claims:          claims,
		IsAuthenticated: isAuth,
	}
}

func (ci DefaultClaimIdentity) GetClaimValue(cType string) string {
	return ci.Claims[cType]
}

package auth

import (
	"time"

	"github.com/lestrrat-go/jwx/jwk"
)

type jwkCache struct {
	Keys   jwk.Set
	Expiry time.Time
}

func (cache *jwkCache) IsExpired() bool {
	if diff := time.Now().Sub(cache.Expiry).Hours(); diff > 0 {
		return true
	}
	return false
}

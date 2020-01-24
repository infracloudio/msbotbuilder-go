package cache

import "time"

// AuthCache is a general purpose cache
type AuthCache struct {
	Keys   interface{}
	Expiry time.Time
}

// IsExpired checks if the Keys have expired.
// Compares Expiry time with current time.
func (cache *AuthCache) IsExpired() bool {

	if diff := time.Now().Sub(cache.Expiry).Hours(); diff > 0 {
		return true
	}
	return false
}

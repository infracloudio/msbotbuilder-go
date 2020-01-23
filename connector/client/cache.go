package client

import "time"

type jwtCache struct {
	token  string
	Expiry time.Time
}

func (cache *jwtCache) IsExpired() bool {

	// if cache == nil {
	// 	return true
	// }

	if diff := time.Now().Sub(cache.Expiry).Hours(); diff > 0 {
		return true
	}
	return false
}

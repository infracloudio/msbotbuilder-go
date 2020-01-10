package schema

type AuthResponse struct {
	TokenType     string `json:"token_type"`
	ExpireTime    int    `json:"expires_in"`
	ExtExpireTime int    `json:"ext_expires_in"`
	AccessToken   string `json:"access_token"`
}

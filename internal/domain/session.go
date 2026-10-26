package domain

// Session represents an active authenticated user session.
type Session struct {
	AccessTokenHash  string `json:"access_token_hash"`
	RefreshTokenHash string `json:"refresh_token_hash"`
}

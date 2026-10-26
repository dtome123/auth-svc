package types

import "github.com/dtome123/auth-sdk/jwt"

type AuthUserType string

const (
	AuthUserTypeRSA  AuthUserType = "rsa"
	AuthUserTypeHMAC AuthUserType = "hmac"
)

type AuthM2MType string

const (
	AuthM2MTypeRSA  AuthM2MType = "rsa"
	AuthM2MTypeHMAC AuthM2MType = "hmac"
)

// ClientEntry defines authentication configuration for a client
type ClientEntry struct {
	Name                  string
	Type                  AuthM2MType
	AllowUserAuthenticate bool
	AllowAudiences        map[string]bool
	AccessTokenTTL        string
	RefreshTokenTTL       string

	Verifier     jwt.Verifier
	UserSigner   jwt.Signer
	UserVerifier jwt.Verifier
}

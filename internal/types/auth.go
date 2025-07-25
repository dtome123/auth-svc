package types

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

type AuthClientEntry struct {
	Name      string      `json:"name"`       // Service name
	Type      AuthM2MType `json:"type"`       // Auth method: "rsa" or "hmac"
	PublicKey string      `json:"public_key"` // RSA public key (for verification)
	SecretKey string      `json:"secret_key"` // HMAC secret key
}

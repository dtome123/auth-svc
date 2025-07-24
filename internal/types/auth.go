package types

type AuthClient string

const (
	AuthClientRSA  AuthClient = "rsa"
	AuthClientHMAC AuthClient = "hmac"
)

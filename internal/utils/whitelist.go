package utils

type WhitelistEntry struct {
	Name          string `json:"name"`
	PublicKeyPath string `json:"public_key"`
}

type WhitelistConfig struct {
	Services []WhitelistEntry `json:"services"`
}

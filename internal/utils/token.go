package utils

import (
	"strings"
)

// ParsedAuth represents the result of parsing the configured authentication header.
type ParsedAuth struct {
	Scheme string // e.g., "Bearer", "Custom", or empty if not required
	Token  string // the actual token value
}

type AuthConfig struct {
	Header string
	Scheme string
}

// ExtractExternalToken extracts a token from HTTP headers based on the provided AuthConfig.
// - If a scheme is defined, it verifies that the header value starts with that scheme (case-insensitive).
// - If no scheme is provided, it extracts the first word of the header value as the token.
func ExtractExternalToken(headers map[string]string, cfg AuthConfig) *ParsedAuth {
	headerValue := strings.TrimSpace(headers[strings.ToLower(cfg.Header)])
	if headerValue == "" {
		return nil
	}

	// If a scheme is specified (e.g., "Bearer"), ensure the header starts with it
	if cfg.Scheme != "" {
		expectedPrefix := strings.ToLower(cfg.Scheme) + " "
		if strings.HasPrefix(strings.ToLower(headerValue), expectedPrefix) {
			return &ParsedAuth{
				Scheme: cfg.Scheme,
				Token:  strings.TrimSpace(headerValue[len(expectedPrefix):]),
			}
		}
		// Scheme does not match
		return nil
	}

	// If no scheme is specified, extract the first token-like part (split by whitespace)
	parts := strings.Fields(headerValue)
	if len(parts) > 0 {
		return &ParsedAuth{
			Scheme: "",
			Token:  parts[0],
		}
	}

	return nil
}

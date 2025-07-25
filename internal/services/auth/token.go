package auth

import (
	"auth-svc/internal/types"
	"context"
	"fmt"
	"time"

	"github.com/dtome123/auth-sdk/client"
	"github.com/dtome123/auth-sdk/constants"
	"github.com/dtome123/auth-sdk/jwtutils"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

type TokenInput struct {
	GrantType           string
	ClientAssertionType string
	ClientAssertion     string
}

type TokenOutput struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
}

func (svc *AuthorizationService) Token(ctx context.Context, req TokenInput) (*TokenOutput, error) {

	// Step 2: Extract client ID from context and verify against whitelist
	clientID := client.ClientIDFromContext(ctx)
	issuer, ok := svc.clients[clientID]
	if !ok {
		return nil, fmt.Errorf("client %q is not whitelisted", clientID)
	}

	// Step 4: Extract JWT from gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}
	assertionTokens := md[constants.ClientAssertionKey]
	if len(assertionTokens) == 0 || assertionTokens[0] == "" {
		return nil, errors.New("missing client assertion")
	}
	assertionJWT := assertionTokens[0]

	// Step 5: Create JWT verifier based on auth type
	var verifier jwtutils.Verifier
	var err error
	switch issuer.Type {
	case types.AuthM2MTypeRSA:
		verifier, err = jwtutils.NewRS256VerifierFromString(issuer.PublicKey)
		if err != nil {
			return nil, errors.Wrap(err, "invalid public key")
		}
	case types.AuthM2MTypeHMAC:
		verifier = jwtutils.NewHMACVerifier([]byte(issuer.SecretKey))
	default:
		return nil, errors.Errorf("unsupported auth type %q for client %q", issuer.Type, issuer.Name)
	}

	// Step 6: Verify JWT and extract claims
	claims, err := verifier.Verify(assertionJWT)
	if err != nil {
		return nil, errors.Wrap(err, "JWT verification failed")
	}

	// Step 7: Validate standard claims
	if claims.Get("iss").AsString() == "" {
		return nil, errors.New("missing 'iss' claim")
	}
	if aud := claims.Get("aud").AsString(); aud != svc.cfg.Aud {
		return nil, errors.New("invalid or missing 'aud' claim")
	}
	exp := claims.Get("exp").AsInt64()
	if exp == 0 {
		return nil, errors.New("missing 'exp' claim")
	}
	now := time.Now().Unix()
	if exp < now {
		return nil, errors.New("JWT expired")
	}

	return &TokenOutput{
		AccessToken: assertionJWT,
		TokenType:   "Bearer",
		ExpiresIn:   exp - now,
	}, nil
}

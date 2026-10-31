package auth

import (
	"auth-svc/internal/utils"
	"context"
	"fmt"

	"github.com/dtome123/auth-sdk/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	// Parse and extract claims from client assertion (JWT)
	claims, err := jwt.NewClaimsFromTokenString(req.ClientAssertion)
	if err != nil {
		return nil, errors.Wrap(err, "invalid client assertion")
	}
	oauthClaims := claims.ToOauthClaims()

	// Get client ID from claims and check whitelist
	issuer, ok := svc.clients[oauthClaims.Iss]
	if !ok {
		return nil, fmt.Errorf("client %q is not whitelisted", oauthClaims.Iss)
	}

	// Verify signature of the client assertion
	if _, err := issuer.Verifier.Verify(req.ClientAssertion); err != nil {
		return nil, errors.Wrap(err, "failed to verify client assertion signature")
	}

	signer := svc.serverSigner
	// Convert to oauth claims and validate core fields (aud, exp, jti, etc.)
	if oauthClaims.Aud != svc.cfg.AuthConfig.Oauth.Aud {
		audience, ok := svc.clients[oauthClaims.Aud]
		if !ok || (!audience.AllowAudiences["*"] && !audience.AllowAudiences[oauthClaims.Iss]) {
			return nil, errors.New("unauthorized audience for client")
		}
		signer = audience.Signer
	}

	opts := []jwt.ValidateOption{}
	if svc.cfg.AuthConfig.EnableReplayCheck {
		opts = append(opts, jwt.WithReplayChecker(svc.oauthRelayChecker))
	}
	if err := oauthClaims.Validate(opts...); err != nil {
		return nil, err
	}

	// Generate access token with configured TTL
	accessTTL, err := utils.ParseFlexibleDuration(svc.cfg.AuthConfig.Oauth.AccessTokenTTL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid access token TTL")
	}
	claims.Set("jti", uuid.New().String())
	token, err := signer.Sign(claims, accessTTL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign access token")
	}

	// Return token response
	return &TokenOutput{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(accessTTL.Seconds()),
	}, nil
}

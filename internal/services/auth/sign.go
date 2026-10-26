package auth

import (
	"auth-svc/internal/domain"
	"auth-svc/internal/utils"
	"context"
	"fmt"
	"time"
	"uuid"

	"github.com/dtome123/auth-sdk/jwt"
	"github.com/dtome123/auth-sdk/session"
	"github.com/dtome123/auth-sdk/token"
)

type SignInput struct {
	UserID   string
	DeviceID string
	UserType string
	Claim    map[string]any
	ClientID string
}

type SignOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	ExpiresIn    int64
}

func (svc *AuthorizationService) Sign(ctx context.Context, req SignInput) (*SignOutput, error) {

	// Resolve client configuration (signer and TTLs)
	clientEntry, ok := svc.clients[req.ClientID]
	if !ok {
		return nil, fmt.Errorf("client %q is not whitelisted", req.ClientID)
	}
	if clientEntry.UserSigner == nil {
		return nil, fmt.Errorf("client %q does not have user signer configured", req.ClientID)
	}

	// Parse durations from config
	accessTTL, err := utils.ParseFlexibleDuration(clientEntry.AccessTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshTTL, err := utils.ParseFlexibleDuration(clientEntry.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}
	signer := clientEntry.UserSigner

	// Get user permissions
	userPermissions, err := svc.GetUserPermissions(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	permissions := formatUserPermissions(userPermissions)

	// Access token claims
	claims := jwt.NewClaims(req.Claim)
	claims.Set("aud", svc.cfg.AuthConfig.SystemAud)
	claims.Set("sub", req.UserID)
	claims.Set("exp", time.Now().Add(accessTTL).Unix())
	claims.Set("permissions", permissions)

	tokenBuilder := token.NewTokenBuilder(claims, accessTTL, refreshTTL, signer)

	idToken, err := tokenBuilder.Build()
	if err != nil {
		return nil, err
	}

	accessTokenHash := utils.HashSHA256(idToken.AccessToken)
	refreshTokenHash := utils.HashSHA256(idToken.RefreshToken)

	// Persist session in the DB
	go func() {
		if svc.sessionResource == nil {
			fmt.Println("warning: sessionResource is nil")
			return
		}

		err := svc.sessionResource.UpsertSession(ctx, session.Session[domain.Session]{
			JTI:      uuid.NewV7().String(),
			UserID:   req.UserID,
			DeviceID: req.DeviceID,
			Type:     req.UserType,
			TTL:      time.Now().Add(refreshTTL),
			Data: domain.Session{
				AccessTokenHash:  accessTokenHash,
				RefreshTokenHash: refreshTokenHash,
			},
		})
		if err != nil {
			fmt.Println("❌ Error persisting session in the DB:", err)
		}
	}()

	// Return tokens to client
	return &SignOutput{
		AccessToken:  idToken.AccessToken,
		RefreshToken: idToken.RefreshToken,
		ExpiresAt:    idToken.ExpiresAt,
		ExpiresIn:    int64(idToken.ExpiresIn),
	}, nil
}

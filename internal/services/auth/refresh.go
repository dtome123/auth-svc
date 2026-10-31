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

type RefreshInput struct {
	ClientID     string
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	ExpiresIn    int64
}

func (svc *AuthorizationService) Refresh(ctx context.Context, req RefreshInput) (*RefreshOutput, error) {

	// Parse metadata into a map
	claims, err := jwt.NewClaimsFromTokenString(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	refreshTokenHash := utils.HashSHA256(req.RefreshToken)

	userID := claims.Get("user_id").AsString()
	deviceID := claims.Get("device_id").AsString()
	userType := claims.Get("user_type").AsString()
	jti := claims.Get("jti").AsString()
	exp := claims.Get("exp").AsInt64()

	// Find s
	s, err := svc.sessionResource.GetSession(ctx, userID, deviceID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, fmt.Errorf("no session found for this refresh token")
	}

	if s.Data.RefreshTokenHash != refreshTokenHash {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if refresh token is expired
	if exp < time.Now().Unix() {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Resolve client configuration (signer and TTLs)
	// clientID := client.ClientIDFromContext(ctx)
	clientEntry, hasClient := svc.clients[req.ClientID]

	accessTTLStr := svc.cfg.AuthConfig.Oauth.UserTokenConfig.AccessTokenTTL
	refreshTTLStr := svc.cfg.AuthConfig.Oauth.UserTokenConfig.RefreshTokenTTL
	signer := svc.userSigner

	if hasClient {
		if clientEntry.AccessTokenTTL != "" {
			accessTTLStr = clientEntry.AccessTokenTTL
		}
		if clientEntry.RefreshTokenTTL != "" {
			refreshTTLStr = clientEntry.RefreshTokenTTL
		}
	}

	// Parse durations from config
	accessTTL, err := utils.ParseFlexibleDuration(accessTTLStr)
	if err != nil {
		return nil, err
	}
	refreshTTL, err := utils.ParseFlexibleDuration(refreshTTLStr)
	if err != nil {
		return nil, err
	}

	// Access token claims
	tokenBuilder := token.NewTokenBuilder(claims, accessTTL, refreshTTL, signer)
	idToken, err := tokenBuilder.Build()
	if err != nil {
		return nil, err
	}

	accessTokenHash := utils.HashSHA256(idToken.AccessToken)
	refreshTokenHash = utils.HashSHA256(idToken.RefreshToken)

	// Persist session in the DB
	go func(ctx context.Context) {

		if svc.sessionResource == nil {
			fmt.Println("warning: sessionResource is nil")
			return
		}

		err := svc.sessionResource.UpsertSession(ctx, session.Session[domain.Session]{
			JTI:      uuid.NewV7().String(),
			UserID:   userID,
			DeviceID: deviceID,
			Type:     userType,
			TTL:      time.Now().Add(refreshTTL),
			Data: domain.Session{
				AccessTokenHash:  accessTokenHash,
				RefreshTokenHash: refreshTokenHash,
			},
		})
		if err != nil {
			fmt.Println("❌ Error persisting session in the DB:", err)
		}
	}(context.Background())

	// Persist session in the DB
	err = svc.sessionResource.UpsertSession(ctx, session.Session[domain.Session]{
		JTI:      jti,
		UserID:   userID,
		DeviceID: deviceID,
		Type:     userType,
		TTL:      time.Now().Add(refreshTTL),
		Data: domain.Session{
			AccessTokenHash:  accessTokenHash,
			RefreshTokenHash: refreshTokenHash,
		},
	})
	if err != nil {
		return nil, err
	}

	// Return tokens to client
	return &RefreshOutput{
		AccessToken:  idToken.AccessToken,
		RefreshToken: idToken.RefreshToken,
		ExpiresAt:    idToken.ExpiresAt,
		ExpiresIn:    int64(idToken.ExpiresIn),
	}, nil
}

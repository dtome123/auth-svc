package auth

import (
	"auth-svc/internal/models"
	"auth-svc/internal/utils"
	"context"
	"fmt"
	"time"

	"github.com/dtome123/auth-sdk/jwtutils"
)

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshOutput struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	ExpiresIn    int64     `json:"expire_in"`
}

func (svc *AuthorizationService) Refresh(ctx context.Context, req RefreshInput) (*RefreshOutput, error) {

	// Parse metadata into a map
	claims, err := jwtutils.NewClaimsFromTokenString(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	refreshTokenHash := utils.HashSHA256(req.RefreshToken)

	userID := claims.Get("user_id").AsString()
	deviceID := claims.Get("device_id").AsString()
	userType := claims.Get("user_type").AsString()
	exp := claims.Get("exp").AsInt64()

	// Find session
	session, err := svc.authenticationRepo.GetSession(ctx, userID, deviceID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("no session found for this refresh token")
	}

	if session.RefreshTokenHash != refreshTokenHash {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if refresh token is expired
	if exp < time.Now().Unix() {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Parse durations from config
	accessTTL, err := utils.ParseFlexibleDuration(svc.cfg.Service.Session.AccessTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshTTL, err := utils.ParseFlexibleDuration(svc.cfg.Service.Session.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	// Access token claims
	token := jwtutils.NewTokenBuilder(claims, accessTTL, refreshTTL, svc.serverSigner)

	idToken, err := token.Build()
	if err != nil {
		return nil, err
	}

	accessTokenHash := utils.HashSHA256(idToken.AccessToken)
	refreshTokenHash = utils.HashSHA256(idToken.RefreshToken)

	// Persist session in the DB
	err = svc.authenticationRepo.UpsertSession(ctx, models.Session{
		UserID:           userID,
		DeviceID:         deviceID,
		Type:             userType,
		AccessTokenHash:  accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
		TTL:              time.Now().Add(refreshTTL),
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

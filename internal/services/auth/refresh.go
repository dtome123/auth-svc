package auth

import (
	"auth-svc/internal/models"
	"auth-svc/internal/utils"
	"context"
	"fmt"
	"maps"
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
	claimData, err := jwtutils.NewClaimsFromTokenString(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	refreshTokenHash := utils.HashSHA256(req.RefreshToken)

	userID := claimData.Get("user_id").AsString()
	deviceID := claimData.Get("device_id").AsString()
	userType := claimData.Get("user_type").AsString()
	exp := claimData.Get("exp").AsInt64()

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

	// Base claims for both tokens
	baseClaims := map[string]interface{}{
		"sub":       claimData.Get("sub").AsString(),
		"user_id":   userID,
		"device_id": deviceID,
		"user_type": userType,
	}

	// Merge metadata into base claims
	for k, v := range claimData {
		baseClaims[k] = v
	}

	// Access token claims
	accessClaims := maps.Clone(baseClaims)
	accessExp := time.Now().Add(accessTTL)
	accessClaims["exp"] = accessExp.Unix()
	accessToken, err := svc.serverSigner.Sign(accessClaims, accessTTL)
	if err != nil {
		return nil, err
	}
	accessTokenHash := utils.HashSHA256(accessToken)

	// Refresh token claims
	refreshClaims := maps.Clone(baseClaims)
	refreshClaims["exp"] = time.Now().Add(refreshTTL).Unix()
	refreshToken, err := svc.serverSigner.Sign(refreshClaims, refreshTTL)
	if err != nil {
		return nil, err
	}
	newRefreshTokenHash := utils.HashSHA256(refreshToken)

	// Persist session in the DB
	err = svc.authenticationRepo.UpsertSession(ctx, models.Session{
		UserID:           userID,
		DeviceID:         deviceID,
		Type:             userType,
		AccessTokenHash:  accessTokenHash,
		RefreshTokenHash: newRefreshTokenHash,
		TTL:              time.Now().Add(refreshTTL),
	})
	if err != nil {
		return nil, err
	}

	// Return tokens to client
	return &RefreshOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp,
		ExpiresIn:    int64(accessTTL.Seconds()),
	}, nil
}

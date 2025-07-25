package auth

import (
	"auth-svc/internal/models"
	"auth-svc/internal/utils"
	"context"
	"encoding/json"
	"maps"
	"time"
)

type SignInput struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
	UserType string `json:"user_type"`
	Metadata string `json:"metadata"`
}

type SignOutput struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	ExpiresIn    int64     `json:"expire_in"`
}

func (svc *AuthorizationService) Sign(ctx context.Context, req SignInput) (*SignOutput, error) {

	// Parse metadata into a map
	claimData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(req.Metadata), &claimData); err != nil {
		return nil, err
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
		"sub":       req.UserID,
		"user_id":   req.UserID,
		"device_id": req.DeviceID,
		"user_type": req.UserType,
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
	refreshTokenHash := utils.HashSHA256(refreshToken)

	// Persist session in the DB
	err = svc.authenticationRepo.UpsertSession(ctx, models.Session{
		UserID:           req.UserID,
		DeviceID:         req.DeviceID,
		Type:             req.UserType,
		AccessTokenHash:  accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
		TTL:              time.Now().Add(refreshTTL),
	})
	if err != nil {
		return nil, err
	}

	// Return tokens to client
	return &SignOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp,
		ExpiresIn:    int64(accessTTL.Seconds()),
	}, nil
}

package auth

import (
	"auth-svc/internal/utils"
	"context"
	"net/http"

	"github.com/dtome123/auth-sdk/jwt"
)

type CheckInput struct {
	AccessToken string
	FullMethod  string
}

type CheckOutput struct {
	Allowed    bool
	StatusCode int32
	Message    string
}

func (svc *AuthorizationService) Check(ctx context.Context, req CheckInput) (CheckOutput, error) {
	if req.AccessToken == "" {
		return CheckOutput{
			Allowed:    false,
			StatusCode: http.StatusUnauthorized,
			Message:    "token is required",
		}, nil
	}

	claims, err := jwt.Extract(req.AccessToken)
	if err != nil {
		return CheckOutput{
			Allowed:    false,
			StatusCode: http.StatusUnauthorized,
			Message:    "invalid token",
		}, nil
	}

	userID := claims.Get("sub").AsString()
	deviceID := claims.Get("device_id").AsString()

	if userID == "" {
		return CheckOutput{
			Allowed:    false,
			StatusCode: http.StatusUnauthorized,
			Message:    "user_id missing in token claims",
		}, nil
	}

	// Validate session in database if device_id is present
	if deviceID != "" {
		session, err := svc.sessionResource.GetSession(ctx, userID, deviceID)
		if err != nil || session == nil {
			return CheckOutput{
				Allowed:    false,
				StatusCode: http.StatusUnauthorized,
				Message:    "session not found or expired",
			}, nil
		}
		expectedHash := utils.HashSHA256(req.AccessToken)
		if session.Data.AccessTokenHash != expectedHash {
			return CheckOutput{
				Allowed:    false,
				StatusCode: http.StatusUnauthorized,
				Message:    "token mismatch with active session",
			}, nil
		}
	}

	return CheckOutput{
		Allowed:    true,
		StatusCode: http.StatusOK,
		Message:    "OK",
	}, nil
}

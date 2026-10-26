package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
	"github.com/dtome123/auth-sdk/middlewares"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GrpcServer) Sign(ctx context.Context, req *authPb.SignRequest) (*authPb.SignResponse, error) {

	claims, ok := middlewares.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated: missing claims in context")
	}
	identity := claims.ToOauthClaims()

	res, err := s.svc.GetAuthService().Sign(ctx, auth.SignInput{
		UserID:   req.UserId,
		DeviceID: req.DeviceId,
		UserType: req.UserType,
		ClientID: identity.Iss,
		Claim:    req.Claim.AsMap(),
	})

	if err != nil {
		return nil, err
	}

	return &authPb.SignResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
		ExpiresAt:    timestamppb.New(res.ExpiresAt),
	}, nil
}

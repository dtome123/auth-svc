package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GrpcServer) Refresh(ctx context.Context, req *authPb.RefreshRequest) (*authPb.RefreshResponse, error) {

	res, err := s.svc.GetAuthService().Refresh(ctx, auth.RefreshInput{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		return nil, err
	}

	return &authPb.RefreshResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
		ExpiresAt:    timestamppb.New(res.ExpiresAt),
	}, nil
}

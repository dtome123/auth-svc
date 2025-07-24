package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/gen/go/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GrpcServer) Sign(ctx context.Context, req *authPb.SignRequest) (*authPb.SignResponse, error) {

	res, err := s.svc.GetAuthService().Sign(ctx, auth.SignInput{
		UserID:   req.UserId,
		DeviceID: req.DeviceId,
		UserType: req.UserType,
		Metadata: req.Metadata,
	})

	if err != nil {
		return nil, err
	}

	return &authPb.SignResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    int64(res.ExpiresAt.Second()),
		ExpiresAt:    timestamppb.New(res.ExpiresAt),
	}, nil
}

package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) AssignRolesToUser(ctx context.Context, req *authPb.AssignRolesToUserRequest) (*authPb.AssignRolesToUserResponse, error) {

	err := s.svc.GetAuthService().AssignRolesToUser(ctx, auth.AssignRolesToUserInput{
		UserID:  req.UserId,
		RoleIDs: req.RoleIds,
	})
	if err != nil {
		return nil, err
	}

	return &authPb.AssignRolesToUserResponse{}, nil
}

package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) ListPermissions(ctx context.Context, req *authPb.ListPermissionsRequest) (*authPb.ListPermissionsResponse, error) {

	res, err := s.svc.GetAuthService().ListPermissions(ctx, auth.ListPermissionsInput{
		Domains: req.GetFilter().GetDomains(),
		Keyword: req.GetFilter().GetKeyword(),
	})
	if err != nil {
		return nil, err
	}

	return &authPb.ListPermissionsResponse{
		Permissions: ToPermissionsProto(res),
	}, nil
}

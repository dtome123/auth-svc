package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) ListPermissionPaths(ctx context.Context, req *authPb.ListPermissionPathsRequest) (*authPb.ListPermissionPathsResponse, error) {

	res, err := s.svc.GetAuthService().ListPermissionPaths(ctx, auth.ListPermissionPathsInput{
		Domains: req.GetFilter().GetDomains(),
		Keyword: req.GetFilter().GetKeyword(),
	})
	if err != nil {
		return nil, err
	}

	return &authPb.ListPermissionPathsResponse{
		Paths: ToPermissionPathsProto(res),
	}, nil
}

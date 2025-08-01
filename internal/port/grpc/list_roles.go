package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
)

func (s *GrpcServer) ListRole(ctx context.Context, req *authPb.ListRoleRequest) (*authPb.ListRoleResponse, error) {

	res, err := s.svc.GetAuthService().ListRole(ctx, auth.ListRoleInput{
		Keyword: req.Filter.Keyword,
		Scope: func(m map[string]*authPb.StringList) map[string][]string {
			res := make(map[string][]string)
			for k, v := range m {
				res[k] = v.Values
			}
			return res
		}(req.Filter.Scope),
	})
	if err != nil {
		return nil, err
	}

	return &authPb.ListRoleResponse{
		Roles: ToListRoleProto(res),
	}, nil
}

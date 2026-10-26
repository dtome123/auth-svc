package grpc

import (
	"auth-svc/internal/services/auth"
	"context"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

func (s *GrpcServer) ListRole(ctx context.Context, req *authPb.ListRoleRequest) (*authPb.ListRoleResponse, error) {
	var keyword string
	var scope map[string][]string
	if filter := req.GetFilter(); filter != nil {
		keyword = filter.GetKeyword()
		if filter.GetScope() != nil {
			scope = make(map[string][]string)
			for k, v := range filter.GetScope() {
				if v != nil {
					scope[k] = v.GetValues()
				}
			}
		}
	}

	var pageNum, pageSize int64
	if pag := req.GetPagination(); pag != nil {
		pageNum = pag.GetPageNumber()
		pageSize = pag.GetPageSize()
	}

	res, err := s.svc.GetAuthService().ListRole(ctx, auth.ListRoleInput{
		Keyword:    keyword,
		Scope:      scope,
		PageNumber: pageNum,
		PageLimit:  pageSize,
	})
	if err != nil {
		return nil, err
	}

	return &authPb.ListRoleResponse{
		Roles: ToListRoleProto(res),
	}, nil
}

package grpc

import (
	"auth-svc/internal/domain"
	"auth-svc/internal/types"

	authPb "github.com/dtome123/auth-sdk/api/auth/v1"
)

var RouteScopeMapping = map[types.RouteScope]authPb.RouteScope{
	types.RouteScopeAuthenticated: authPb.RouteScope_ROUTE_SCOPE_AUTHENTICATED,
	types.RouteScopeAuthorized:    authPb.RouteScope_ROUTE_SCOPE_AUTHORIZED,
}

var PathTypeMappingInv = map[authPb.RouteScope]types.RouteScope{
	authPb.RouteScope_ROUTE_SCOPE_AUTHENTICATED: types.RouteScopeAuthenticated,
	authPb.RouteScope_ROUTE_SCOPE_AUTHORIZED:    types.RouteScopeAuthorized,
}

func ToPermissionsProto(in []*domain.Permission) []*authPb.Permission {
	var permissions []*authPb.Permission
	for _, p := range in {
		if p == nil {
			continue
		}
		permissions = append(permissions, &authPb.Permission{
			Id:       p.ID,
			Name:     p.Name,
			Domain:   p.Domain,
			Resource: p.Resource,
			Actions:  p.Actions,
		})
	}

	return permissions
}

func ToResourceActionsPermissionsProto(in []domain.ResourceAction) []*authPb.Permission {
	var permissions []*authPb.Permission
	for _, ra := range in {
		permissions = append(permissions, &authPb.Permission{
			Id:      ra.PermissionID,
			Actions: ra.Actions,
		})
	}
	return permissions
}

func ToRoleProto(in *domain.Role) *authPb.Role {
	if in == nil {
		return nil
	}

	return &authPb.Role{
		Id:          in.ID,
		Name:        in.Name,
		Permissions: ToPermissionsProto(in.Permissions),
		Scope:       in.Scope,
	}
}

func ToListRoleProto(in []*domain.Role) []*authPb.Role {
	var roles []*authPb.Role
	for _, r := range in {
		roles = append(roles, ToRoleProto(r))
	}
	return roles
}

func FromPermissionProto(in *authPb.Permission) *domain.Permission {
	if in == nil {
		return nil
	}
	return &domain.Permission{
		ID:       in.Id,
		Name:     in.Name,
		Domain:   in.Domain,
		Resource: in.Resource,
		Actions:  in.Actions,
	}
}

func FromListPermissionProto(in []*authPb.Permission) []*domain.Permission {
	var permissions []*domain.Permission
	for _, p := range in {
		permissions = append(permissions, FromPermissionProto(p))
	}
	return permissions
}

func FromRoleProto(in *authPb.Role) *domain.Role {
	if in == nil {
		return nil
	}
	var perms []*domain.Permission
	var ras []domain.ResourceAction
	for _, p := range in.Permissions {
		dPerm := FromPermissionProto(p)
		perms = append(perms, dPerm)
		if p.Id != "" {
			ras = append(ras, domain.ResourceAction{
				PermissionID: p.Id,
				Actions:      p.Actions,
			})
		}
	}
	return &domain.Role{
		ID:              in.Id,
		Name:            in.Name,
		ResourceActions: ras,
		Permissions:     perms,
		Scope:           in.Scope,
	}
}

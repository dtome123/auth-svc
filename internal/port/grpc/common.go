package grpc

import (
	"auth-svc/internal/models"
	"auth-svc/internal/types"

	authPb "github.com/dtome123/auth-sdk/api/go/auth/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToPermissionsProto(in []*models.Permission) []*authPb.Permission {
	var permissions []*authPb.Permission
	for _, p := range in {

		actions := make([]*authPb.ActionResource, 0)
		for _, action := range p.ImpliedByActions {
			actions = append(actions, &authPb.ActionResource{
				Resource: action.Resource,
				Action:   action.Action,
			})
		}

		permissions = append(permissions, &authPb.Permission{
			Id:               p.ID.Hex(),
			Name:             p.Name,
			Domain:           p.Domain,
			Resource:         p.Resource,
			Action:           p.Action,
			ImpliedByActions: actions,
			Description:      p.Description,
		})
	}

	return permissions
}

func ToPermissionPathsProto(in []*models.PermissionPath) []*authPb.PermissionPath {
	var paths []*authPb.PermissionPath
	for _, p := range in {

		if p == nil {
			continue
		}

		paths = append(paths, &authPb.PermissionPath{
			Id:       p.ID.Hex(),
			Domain:   p.Domain,
			Path:     p.Path,
			Resource: p.Resource,
			Action:   p.Action,
			Type:     authPb.RouteScope(p.Type),
		})
	}
	return paths
}

func ToRoleProto(in *models.Role) *authPb.Role {

	if in == nil {
		return nil
	}

	permissionIds := make([]string, len(in.PermissionIDs))
	for i, permissionID := range in.PermissionIDs {
		permissionIds[i] = permissionID.Hex()
	}

	return &authPb.Role{
		Id:            in.ID.Hex(),
		Name:          in.Name,
		PermissionIds: permissionIds,
		Description:   in.Description,
		Permissions:   ToPermissionsProto(in.Permissions),
	}
}

func ToListRoleProto(in []*models.Role) []*authPb.Role {
	var roles []*authPb.Role
	for _, r := range in {
		roles = append(roles, ToRoleProto(r))
	}
	return roles
}

func FromActionResourceProto(in *authPb.ActionResource) models.ActionResource {
	return models.ActionResource{
		Resource: in.Resource,
		Action:   in.Action,
	}
}

func FromListActionResourceProto(in []*authPb.ActionResource) []models.ActionResource {
	var actions []models.ActionResource
	for _, a := range in {
		actions = append(actions, FromActionResourceProto(a))
	}
	return actions
}

func FromPermissionProto(in *authPb.Permission) *models.Permission {

	objectID, _ := primitive.ObjectIDFromHex(in.Id)

	return &models.Permission{
		ID:               objectID,
		Name:             in.Name,
		Domain:           in.Domain,
		Resource:         in.Resource,
		Action:           in.Action,
		ImpliedByActions: FromListActionResourceProto(in.ImpliedByActions),
		Description:      in.Description,
	}
}

func FromListPermissionProto(in []*authPb.Permission) []*models.Permission {
	var permissions []*models.Permission
	for _, p := range in {
		permissions = append(permissions, FromPermissionProto(p))
	}
	return permissions
}

func FromPermissionPathProto(in *authPb.PermissionPath) *models.PermissionPath {
	objectID, _ := primitive.ObjectIDFromHex(in.Id)
	return &models.PermissionPath{
		ID:       objectID,
		Domain:   in.Domain,
		Path:     in.Path,
		Resource: in.Resource,
		Action:   in.Action,
		Type:     types.RouteScope(in.Type),
	}
}

func FromListPermissionPathProto(in []*authPb.PermissionPath) []*models.PermissionPath {
	var paths []*models.PermissionPath
	for _, p := range in {
		paths = append(paths, FromPermissionPathProto(p))
	}
	return paths
}

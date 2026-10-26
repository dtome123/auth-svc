package models

import (
	"auth-svc/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResourceAction struct {
	PermissionID primitive.ObjectID `bson:"permission_id" json:"permission_id"`
	Actions      []string           `bson:"actions" json:"actions"`
}

type Role struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name            string               `bson:"name" json:"name"`
	ResourceActions []ResourceAction     `bson:"resource_actions,omitempty" json:"resource_actions,omitempty"`
	PermissionIDs   []primitive.ObjectID `bson:"permission_ids" json:"permission_ids"`
	Permissions     []*Permission        `bson:"-" json:"permissions"`
	Description     string               `bson:"description" json:"description"`
	Scope           map[string]string    `bson:"scope" json:"scope"`
}

func (Role) CollectionName() string {
	return "roles"
}

func (r *Role) ToDomain() *domain.Role {
	if r == nil {
		return nil
	}

	permIDs := make([]string, len(r.PermissionIDs))
	for i, pid := range r.PermissionIDs {
		permIDs[i] = pid.Hex()
	}

	var ras []domain.ResourceAction
	if len(r.ResourceActions) > 0 {
		ras = make([]domain.ResourceAction, len(r.ResourceActions))
		for i, ra := range r.ResourceActions {
			ras[i] = domain.ResourceAction{
				PermissionID: ra.PermissionID.Hex(),
				Actions:      ra.Actions,
			}
		}
	}

	var perms []*domain.Permission
	if len(r.Permissions) > 0 {
		perms = make([]*domain.Permission, len(r.Permissions))
		for i, p := range r.Permissions {
			perms[i] = p.ToDomain()
		}
	}

	return &domain.Role{
		ID:              r.ID.Hex(),
		Name:            r.Name,
		ResourceActions: ras,
		PermissionIDs:   permIDs,
		Permissions:     perms,
		Description:     r.Description,
		Scope:           r.Scope,
	}
}

func RoleFromDomain(r *domain.Role) *Role {
	if r == nil {
		return nil
	}

	objID, _ := primitive.ObjectIDFromHex(r.ID)
	permIDSet := make(map[string]struct{})
	for _, pid := range r.PermissionIDs {
		if pid != "" {
			permIDSet[pid] = struct{}{}
		}
	}

	var ras []ResourceAction
	for _, ra := range r.ResourceActions {
		oid, _ := primitive.ObjectIDFromHex(ra.PermissionID)
		ras = append(ras, ResourceAction{
			PermissionID: oid,
			Actions:      ra.Actions,
		})
		if ra.PermissionID != "" {
			permIDSet[ra.PermissionID] = struct{}{}
		}
	}

	var permIDs []primitive.ObjectID
	for pid := range permIDSet {
		oid, _ := primitive.ObjectIDFromHex(pid)
		permIDs = append(permIDs, oid)
	}

	var perms []*Permission
	if len(r.Permissions) > 0 {
		perms = make([]*Permission, len(r.Permissions))
		for i, p := range r.Permissions {
			perms[i] = PermissionFromDomain(p)
		}
	}

	return &Role{
		ID:              objID,
		Name:            r.Name,
		ResourceActions: ras,
		PermissionIDs:   permIDs,
		Permissions:     perms,
		Description:     r.Description,
		Scope:           r.Scope,
	}
}

package models

import (
	"auth-svc/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Assignment struct {
	UserID  string               `bson:"user_id" json:"user_id"`
	RoleIDs []primitive.ObjectID `bson:"role_id" json:"role_id"`
	Roles   []Role               `json:"role" bson:"-"`
}

func (Assignment) CollectionName() string {
	return "assignments"
}

func (a *Assignment) ToDomain() *domain.Assignment {
	if a == nil {
		return nil
	}

	roleIDs := make([]string, len(a.RoleIDs))
	for i, rid := range a.RoleIDs {
		roleIDs[i] = rid.Hex()
	}

	var roles []domain.Role
	if len(a.Roles) > 0 {
		roles = make([]domain.Role, len(a.Roles))
		for i, r := range a.Roles {
			dr := r.ToDomain()
			if dr != nil {
				roles[i] = *dr
			}
		}
	}

	return &domain.Assignment{
		UserID:  a.UserID,
		RoleIDs: roleIDs,
		Roles:   roles,
	}
}

func AssignmentFromDomain(a *domain.Assignment) *Assignment {
	if a == nil {
		return nil
	}

	roleIDs := make([]primitive.ObjectID, len(a.RoleIDs))
	for i, rid := range a.RoleIDs {
		oid, _ := primitive.ObjectIDFromHex(rid)
		roleIDs[i] = oid
	}

	var roles []Role
	if len(a.Roles) > 0 {
		roles = make([]Role, len(a.Roles))
		for i, r := range a.Roles {
			roles[i] = *RoleFromDomain(&r)
		}
	}

	return &Assignment{
		UserID:  a.UserID,
		RoleIDs: roleIDs,
		Roles:   roles,
	}
}

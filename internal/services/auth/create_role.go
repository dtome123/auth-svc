package auth

import (
	"auth-svc/internal/models"
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateRoleInput struct {
	Name          string
	PermissionIDs []string
	Description   string
}

func (svc *AuthorizationService) CreateRole(ctx context.Context, input CreateRoleInput) (string, error) {
	permissionIDs := make([]primitive.ObjectID, len(input.PermissionIDs))
	for i, permissionID := range input.PermissionIDs {
		primPermissionID, _ := primitive.ObjectIDFromHex(permissionID)
		permissionIDs[i] = primPermissionID
	}
	role := &models.Role{
		Name:          input.Name,
		PermissionIDs: permissionIDs,
		Description:   input.Description,
	}
	if err := svc.authorizationRepo.CreateRole(ctx, role); err != nil {
		return "", errors.Wrap(err, "failed to create role")
	}
	return role.ID.Hex(), nil
}

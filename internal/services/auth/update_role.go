package auth

import (
	"auth-svc/internal/models"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateRoleInput struct {
	ID            string
	Name          string
	PermissionIDs []string
	Description   string
}

func (svc *AuthorizationService) UpdateRole(ctx context.Context, input UpdateRoleInput) (string, error) {
	permissionIDs := make([]primitive.ObjectID, len(input.PermissionIDs))
	for i, permissionID := range input.PermissionIDs {
		primPermissionID, _ := primitive.ObjectIDFromHex(permissionID)
		permissionIDs[i] = primPermissionID
	}

	roleID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return "", fmt.Errorf("invalid role id: %w", err)
	}
	role := &models.Role{
		ID:            roleID,
		Name:          input.Name,
		PermissionIDs: permissionIDs,
		Description:   input.Description,
	}
	if err := svc.authorizationRepo.UpdateRole(ctx, role); err != nil {
		return "", errors.Wrap(err, "failed to update role")
	}
	return role.ID.Hex(), nil
}

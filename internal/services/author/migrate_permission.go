package author

import (
	"auth-svc/internal/models"
	"context"
)

type MigratePermissionInput struct {
	Permissions     []models.Permission
	PermissionPaths []models.PermissionPath
}

func (svc *AuthorizationService) MigratePermission(ctx context.Context, req MigratePermissionInput) error {

	return nil
}

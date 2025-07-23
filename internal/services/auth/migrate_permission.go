package auth

import (
	"auth-svc/internal/models"
	"context"
)

type MigratePermissionInput struct {
	Permissions     []*models.Permission
	PermissionPaths []*models.PermissionPath
}

func (svc *AuthorizationService) MigratePermission(ctx context.Context, req MigratePermissionInput) error {
	err := svc.authorizationRepo.BatchUpsertPermissions(ctx, req.Permissions)
	if err != nil {
		return err
	}

	err = svc.authorizationRepo.BatchUpsertPermissionPaths(ctx, req.PermissionPaths)
	if err != nil {
		return err
	}
	return nil
}

package auth

import (
	"auth-svc/internal/domain"
	"context"
	"fmt"
)

type MigratePermissionInput struct {
	Permissions []*domain.Permission
	Roles       []*domain.Role
}

func (svc *AuthorizationService) MigratePermission(ctx context.Context, req MigratePermissionInput) error {
	if len(req.Permissions) > 0 {
		err := svc.authorizationRepo.BatchUpsertPermissions(ctx, req.Permissions)
		if err != nil {
			return err
		}
	}

	if len(req.Roles) > 0 {
		for _, role := range req.Roles {
			if role == nil {
				continue
			}
			existing, err := svc.authorizationRepo.ListRoles(ctx, domain.ListRoleInput{Keyword: role.Name})
			var found *domain.Role
			if err == nil {
				for _, r := range existing {
					if r.Name == role.Name {
						found = r
						break
					}
				}
			}
			if found != nil {
				role.ID = found.ID
				_ = svc.authorizationRepo.UpdateRole(ctx, role)
			} else {
				_ = svc.authorizationRepo.CreateRole(ctx, role)
			}
		}
	}

	if err := svc.LoadAllPolicies(ctx); err != nil {
		fmt.Printf("failed to reload casbin policies on migration: %v\n", err)
	}

	svc.publishPolicySync(ctx)

	if svc.cfg.Caching.Enable {
		if err := svc.authorCache.ClearAllPermissionsAndChecks(ctx); err != nil {
			fmt.Printf("failed to clear all caches on migration: %v \n", err)
		}
	}

	return nil
}

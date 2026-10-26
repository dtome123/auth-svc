package auth

import (
	"auth-svc/internal/domain"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type CreateRoleInput struct {
	Name            string
	ResourceActions []domain.ResourceAction
	PermissionIDs   []string
	Description     string
	Scope           map[string]string
}

func (svc *AuthorizationService) CreateRole(ctx context.Context, input CreateRoleInput) (string, error) {
	permIDs := input.PermissionIDs
	if len(permIDs) == 0 && len(input.ResourceActions) > 0 {
		for _, ra := range input.ResourceActions {
			if ra.PermissionID != "" {
				permIDs = append(permIDs, ra.PermissionID)
			}
		}
	}

	role := &domain.Role{
		Name:            input.Name,
		ResourceActions: input.ResourceActions,
		PermissionIDs:   permIDs,
		Description:     input.Description,
		Scope:           input.Scope,
	}
	if err := svc.authorizationRepo.CreateRole(ctx, role); err != nil {
		return "", errors.Wrap(err, "failed to create role")
	}

	if err := svc.syncRolePoliciesInEnforcer(ctx, role.ID, role.Name, role.ResourceActions, permIDs); err != nil {
		fmt.Printf("failed to sync role policies in casbin enforcer: %v\n", err)
	}

	svc.publishPolicySync(ctx)

	return role.ID, nil
}

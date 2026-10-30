package auth

import (
	"auth-svc/internal/domain"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type CreateRoleInput struct {
	Code            string
	Name            string
	ResourceActions []domain.ResourceAction
	Description     string
	Scope           map[string]string
}

func (svc *AuthorizationService) CreateRole(ctx context.Context, input CreateRoleInput) (string, error) {
	role := &domain.Role{
		Code:            input.Code,
		Name:            input.Name,
		ResourceActions: input.ResourceActions,
		Description:     input.Description,
		Scope:           input.Scope,
	}
	if err := svc.authorizationRepo.CreateRole(ctx, role); err != nil {
		return "", errors.Wrap(err, "failed to create role")
	}

	if err := svc.syncRolePoliciesInEnforcer(ctx, role.ID, role.Name, role.ResourceActions); err != nil {
		fmt.Printf("failed to sync role policies in casbin enforcer: %v\n", err)
	}

	svc.publishPolicySync(ctx)

	return role.ID, nil
}

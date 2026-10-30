package auth

import (
	"auth-svc/internal/domain"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type UpdateRoleInput struct {
	ID              string
	Name            string
	ResourceActions []domain.ResourceAction
	Description     string
	Scope           map[string]string
}

func (svc *AuthorizationService) UpdateRole(ctx context.Context, input UpdateRoleInput) (string, error) {
	role := &domain.Role{
		ID:              input.ID,
		Name:            input.Name,
		ResourceActions: input.ResourceActions,
		Description:     input.Description,
		Scope:           input.Scope,
	}
	if err := svc.authorizationRepo.UpdateRole(ctx, role); err != nil {
		return "", errors.Wrap(err, "failed to update role")
	}

	if err := svc.syncRolePoliciesInEnforcer(ctx, role.ID, role.Name, role.ResourceActions); err != nil {
		fmt.Printf("failed to sync role policies in casbin enforcer: %v\n", err)
	}

	svc.publishPolicySync(ctx)

	// Invalidate cache for users assigned to this role
	assignments, err := svc.authorizationRepo.GetAssignmentsByRoleID(ctx, role.ID)
	if err == nil {
		for _, assignment := range assignments {
			_ = svc.clearCheckCache(ctx, assignment.UserID)
		}
	}

	return role.ID, nil
}

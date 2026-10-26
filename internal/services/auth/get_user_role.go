package auth

import (
	"auth-svc/internal/domain"
	"context"
)

func (svc *AuthorizationService) GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error) {

	assignments, err := svc.authorizationRepo.GetAssignmentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if assignments == nil || len(assignments.RoleIDs) == 0 {
		return nil, nil
	}

	roles, err := svc.authorizationRepo.ListRolesByIDs(ctx, assignments.RoleIDs)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

package auth

import (
	"auth-svc/internal/models"
	"context"
)

func (svc *AuthorizationService) GetUserRoles(ctx context.Context, userID string) ([]models.Role, error) {

	assignments, err := svc.authorizationRepo.GetAssignmentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles, err := svc.authorizationRepo.ListRolesByIDs(ctx, assignments.RoleIDs)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

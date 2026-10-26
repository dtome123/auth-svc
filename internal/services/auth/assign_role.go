package auth

import (
	"auth-svc/internal/domain"
	"context"
	"fmt"
)

type AssignRolesToUserInput struct {
	UserID  string
	RoleIDs []string
}

func (svc *AuthorizationService) AssignRolesToUser(ctx context.Context, req AssignRolesToUserInput) error {

	if err := svc.authorizationRepo.UpsertAssignment(ctx, &domain.Assignment{
		UserID:  req.UserID,
		RoleIDs: req.RoleIDs,
	}); err != nil {
		return err
	}

	if err := svc.clearCheckCache(ctx, req.UserID); err != nil {
		return err
	}
	return nil
}

func (svc *AuthorizationService) clearCheckCache(ctx context.Context, userID string) error {
	if !svc.cfg.Caching.Enable {
		return nil
	}
	if err := svc.authorCache.ClearUserPermissionCheck(ctx, userID); err != nil {
		fmt.Printf("failed to clear user cache: %v \n", err)
	}
	if err := svc.authorCache.ClearUserPermissions(ctx, userID); err != nil {
		fmt.Printf("failed to clear user cache: %v \n", err)
	}

	return nil
}

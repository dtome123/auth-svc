package auth

import (
	"auth-svc/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssignRolesToUserInput struct {
	UserID  string
	RoleIDs []string
}

func (svc *AuthorizationService) AssignRolesToUser(ctx context.Context, req AssignRolesToUserInput) error {

	roleIDs := make([]primitive.ObjectID, len(req.RoleIDs))
	for i, roleID := range req.RoleIDs {
		primRoleID, _ := primitive.ObjectIDFromHex(roleID)
		roleIDs[i] = primRoleID
	}

	if err := svc.authorizationRepo.UpsertAssignment(ctx, &models.Assignment{
		UserID:  req.UserID,
		RoleIDs: roleIDs,
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

package author

import (
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPermissionsByUserID retrieves all permissions assigned to a user via their roles.
func (repo *impl) GetPermissionsByUserID(ctx context.Context, userID string) ([]domain.ResourceAction, error) {
	// Get user's assignment (list of role IDs)
	assignment, err := repo.assignmentCol.FindOne(ctx, bson.M{"user_id": userID},
		options.FindOne().SetHint(IdxAssignmentUserId),
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	// No roles assigned
	if len(assignment.RoleIDs) == 0 {
		return nil, nil
	}

	// Find roles by IDs
	roles, err := repo.roleCol.Find(
		ctx,
		bson.M{"_id": bson.M{"$in": assignment.RoleIDs}},
		options.Find().SetHint("_id_"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Collect unique permission IDs from role.ResourceActions
	actions := []domain.ResourceAction{}
	for _, role := range roles {
		for _, ra := range role.ResourceActions {
			if ra.PermissionID != "" {
				actions = append(actions, domain.ResourceAction{
					PermissionID: ra.PermissionID,
					Actions:      ra.Actions,
				})
			}
		}
	}

	return actions, nil
}

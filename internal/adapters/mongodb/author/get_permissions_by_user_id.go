package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPermissionsByUserID retrieves all permissions assigned to a user via their roles.
func (repo *AuthorizationRepository) GetPermissionsByUserID(ctx context.Context, userID string) ([]*models.Permission, error) {
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

	// Collect unique permission IDs
	permIDSet := make(map[primitive.ObjectID]struct{})
	for _, role := range roles {
		for _, pid := range role.PermissionIDs {
			permIDSet[pid] = struct{}{}
		}
	}

	if len(permIDSet) == 0 {
		return nil, nil
	}

	// Convert set to slice
	permissionIDs := make([]primitive.ObjectID, 0, len(permIDSet))
	for pid := range permIDSet {
		permissionIDs = append(permissionIDs, pid)
	}

	// Find permissions
	permissions, err := repo.permissionCol.Find(
		ctx,
		bson.M{"_id": bson.M{"$in": permissionIDs}},
		options.Find().SetHint("_id_"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

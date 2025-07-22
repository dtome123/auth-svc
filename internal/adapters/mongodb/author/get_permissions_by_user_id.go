package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPermissionsByUserID retrieves all permissions assigned to a user via their roles.
// It returns all permissions including the implied by actions stored inside each permission.
func (repo *AuthorizationRepository) GetPermissionsByUserID(ctx context.Context, userID string) ([]models.Permission, error) {
	var permissions []models.Permission

	// Convert userID string to ObjectID
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Find all assignments for the user to get role IDs
	var assignments []models.Assignment
	cursor, err := repo.AssignmentCol.Find(ctx, bson.M{"user_id": uid}, &options.FindOptions{
		Hint: IdxAssignmentUserId,
	})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &assignments); err != nil {
		return nil, err
	}

	if len(assignments) == 0 {
		return permissions, nil
	}

	// Collect all role IDs from assignments
	roleIDs := make([]primitive.ObjectID, 0, len(assignments))
	for _, a := range assignments {
		roleIDs = append(roleIDs, a.RoleID)
	}

	// Find all roles by the collected role IDs
	var roles []models.Role
	cursor, err = repo.RoleCol.Find(ctx, bson.M{"_id": bson.M{"$in": roleIDs}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	// Collect unique permission IDs from roles
	permIDSet := make(map[primitive.ObjectID]struct{})
	for _, r := range roles {
		for _, pid := range r.PermissionIDs {
			permIDSet[pid] = struct{}{}
		}
	}

	permissionIDs := make([]primitive.ObjectID, 0, len(permIDSet))
	for pid := range permIDSet {
		permissionIDs = append(permissionIDs, pid)
	}

	// Query the permissions collection by permission IDs
	cursor, err = repo.PermissionCol.Find(ctx, bson.M{"_id": bson.M{"$in": permissionIDs}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

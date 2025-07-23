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
	var assignment models.Assignment
	assignmentRes := repo.AssignmentCol.FindOne(ctx, bson.M{"user_id": uid}, &options.FindOneOptions{
		Hint: IdxAssignmentUserId,
	})

	if err := assignmentRes.Decode(&assignment); err != nil {
		return nil, err
	}

	// Collect all role IDs from assignments
	roleIDs := assignment.RoleIDs

	// Find all roles by the collected role IDs
	var roles []models.Role
	cursor, err := repo.RoleCol.Find(ctx, bson.M{"_id": bson.M{"$in": roleIDs}})
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

package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAssignmentByUserID retrieves the assignment document associated with a given user ID.
func (repo *AuthorizationRepository) GetAssignmentByUserID(ctx context.Context, userID string) (*models.Assignment, error) {
	var assignment models.Assignment

	err := repo.AssignmentCol.FindOne(ctx, bson.M{"user_id": userID},
		options.FindOne().SetHint(IdxAssignmentUserId),
	).Decode(&assignment)

	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

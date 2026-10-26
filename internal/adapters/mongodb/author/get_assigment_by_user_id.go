package author

import (
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAssignmentByUserID retrieves the assignment document associated with a given user ID.
func (repo *AuthorizationRepository) GetAssignmentByUserID(ctx context.Context, userID string) (*domain.Assignment, error) {

	assignmentDoc, err := repo.assignmentCol.FindOne(ctx, bson.M{"user_id": userID},
		options.FindOne().SetHint(IdxAssignmentUserId),
	)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return assignmentDoc.ToDomain(), nil
}

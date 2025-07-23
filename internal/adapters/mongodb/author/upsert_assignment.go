package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (repo *AuthorizationRepository) UpsertAssignment(ctx context.Context, assignment *models.Assignment) error {
	_, err := repo.AssignmentCol.UpdateOne(
		ctx,
		bson.M{
			"user_id": assignment.UserID},
		bson.M{"$set": assignment},
	)
	return err
}

package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertAssignment updates an existing assignment for the user or inserts it if not exists.
// The operation matches the document by "user_id" field.
func (repo *AuthorizationRepository) UpsertAssignment(ctx context.Context, assignment *models.Assignment) error {
	updateOpts := options.Update().SetUpsert(true)

	err := repo.assignmentCol.UpdateSetOne(
		ctx,
		bson.M{"user_id": assignment.UserID},
		bson.M{"$set": assignment},
		updateOpts,
	)
	return err
}

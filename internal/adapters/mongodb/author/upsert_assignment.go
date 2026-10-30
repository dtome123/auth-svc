package author

import (
	mongoModels "auth-svc/internal/adapters/mongodb/models"
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertAssignment updates an existing assignment for the user or inserts it if not exists.
func (repo *impl) UpsertAssignment(ctx context.Context, assignment *domain.Assignment) error {
	doc := mongoModels.AssignmentFromDomain(assignment)
	updateOpts := options.Update().SetUpsert(true)

	err := repo.assignmentCol.UpdateOne(
		ctx,
		bson.M{"user_id": doc.UserID},
		bson.M{"$set": doc},
		updateOpts,
	)
	return err
}

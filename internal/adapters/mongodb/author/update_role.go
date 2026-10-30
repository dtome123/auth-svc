package author

import (
	mongoModels "auth-svc/internal/adapters/mongodb/models"
	"auth-svc/internal/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateRole updates an existing role document identified by role.ID.
func (repo *impl) UpdateRole(ctx context.Context, role *domain.Role) error {
	if role == nil {
		return errors.New("role is nil")
	}

	objID, err := primitive.ObjectIDFromHex(role.ID)
	if err != nil {
		return errors.New("invalid role ID")
	}

	doc := mongoModels.RoleFromDomain(role)
	doc.ID = objID

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": doc}

	err = repo.roleCol.UpdateOne(ctx, filter, update, &options.UpdateOptions{
		Hint: "_id",
	})
	return err
}

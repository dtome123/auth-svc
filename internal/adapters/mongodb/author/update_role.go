package author

import (
	"auth-svc/internal/models"
	"context"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateRole updates an existing role document identified by role.ID.
// It replaces the fields of the role document with the values from the provided role object.
func (repo *AuthorizationRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	if role == nil {
		return errors.New("role is nil")
	}

	// Ensure the role ID is valid ObjectID
	objID, err := primitive.ObjectIDFromHex(role.ID.Hex())
	if err != nil {
		return errors.New("invalid role ID")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": role}

	err = repo.roleCol.UpdateSetOne(ctx, filter, update, &options.UpdateOptions{
		Hint: "_id",
	})
	return err
}

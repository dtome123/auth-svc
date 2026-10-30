package author

import (
	"auth-svc/internal/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetRole retrieves a role by its ID (hex string).
func (repo *impl) GetRole(ctx context.Context, id string) (*domain.Role, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid role ID format")
	}

	roleDoc, err := repo.roleCol.FindOne(
		ctx,
		bson.M{"_id": objID},
		options.FindOne().SetHint("_id_"),
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return roleDoc.ToDomain(), nil
}

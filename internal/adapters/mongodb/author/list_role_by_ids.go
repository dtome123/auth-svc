package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *AuthorizationRepository) ListRolesByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Role, error) {

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	cursor, err := repo.RoleCol.Find(ctx, filter, &options.FindOptions{})
	if err != nil {
		return nil, err
	}

	var roles []models.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

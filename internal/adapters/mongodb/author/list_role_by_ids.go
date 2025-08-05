package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *AuthorizationRepository) ListRolesByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Role, error) {

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	roles, err := repo.roleCol.Find(ctx, filter, &options.FindOptions{
		Hint: "_id",
	}, nil)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

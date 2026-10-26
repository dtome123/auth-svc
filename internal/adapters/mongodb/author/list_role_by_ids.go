package author

import (
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *AuthorizationRepository) ListRolesByIDs(ctx context.Context, ids []string) ([]*domain.Role, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	if len(objIDs) == 0 {
		return nil, nil
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": objIDs,
		},
	}

	roleDocs, err := repo.roleCol.Find(ctx, filter, &options.FindOptions{
		Hint: "_id",
	}, nil)
	if err != nil {
		return nil, err
	}

	roles := make([]*domain.Role, len(roleDocs))
	for i, r := range roleDocs {
		roles[i] = r.ToDomain()
	}

	return roles, nil
}

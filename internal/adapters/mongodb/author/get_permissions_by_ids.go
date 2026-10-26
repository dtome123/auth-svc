package author

import (
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPermissionsByIDs returns the permission documents referenced by roles.
func (repo *AuthorizationRepository) GetPermissionsByIDs(ctx context.Context, ids []string) ([]*domain.Permission, error) {
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

	permDocs, err := repo.permissionCol.Find(
		ctx,
		bson.M{"_id": bson.M{"$in": objIDs}},
		options.Find().SetHint("_id_"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	permissions := make([]*domain.Permission, len(permDocs))
	for i, p := range permDocs {
		permissions[i] = p.ToDomain()
	}

	return permissions, nil
}

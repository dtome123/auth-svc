package author

import (
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPermissionsByIDs returns the permission documents referenced by roles.
func (repo *impl) GetPermissionsByIDs(ctx context.Context, ids []string) ([]*domain.Permission, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	permDocs, err := repo.permissionCol.Find(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}},
		options.Find(),
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

package author

import (
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetRoleByCode retrieves a role by its code.
func (repo *impl) GetRoleByCode(ctx context.Context, code string, scope map[string][]string) (*domain.Role, error) {
	filter := bson.M{"code": code}

	for k, v := range scope {
		if len(v) == 0 {
			continue
		}
		filter["scope."+k] = bson.M{"$in": v}
	}

	roleDoc, err := repo.roleCol.FindOne(
		ctx,
		filter,
		options.FindOne(),
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return roleDoc.ToDomain(), nil
}

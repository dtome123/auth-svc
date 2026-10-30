package author

import (
	"auth-svc/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetRoleByCode retrieves a role by its code.
func (repo *impl) GetRoleByCode(ctx context.Context, code string) (*domain.Role, error) {
	roleDoc, err := repo.roleCol.FindOne(
		ctx,
		bson.M{"code": code},
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

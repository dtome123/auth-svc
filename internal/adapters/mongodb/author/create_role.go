package author

import (
	mongoModels "auth-svc/internal/adapters/mongodb/models"
	"auth-svc/internal/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRole inserts a new role into the database.
func (repo *impl) CreateRole(ctx context.Context, role *domain.Role) error {
	if role == nil {
		return errors.New("role is nil")
	}

	doc := mongoModels.RoleFromDomain(role)
	if doc.ID.IsZero() {
		doc.ID = primitive.NewObjectID()
	}

	err := repo.roleCol.InsertOne(ctx, *doc)
	if err != nil {
		return err
	}

	role.ID = doc.ID.Hex()
	return nil
}

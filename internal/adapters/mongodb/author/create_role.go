package author

import (
	"auth-svc/internal/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRole inserts a new role into the database.
func (repo *AuthorizationRepository) CreateRole(ctx context.Context, role *models.Role) error {
	if role == nil {
		return errors.New("role is nil")
	}

	result, err := repo.RoleCol.InsertOne(ctx, &role)
	if err != nil {
		return err
	}

	role.ID, _ = result.InsertedID.(primitive.ObjectID)

	return nil
}

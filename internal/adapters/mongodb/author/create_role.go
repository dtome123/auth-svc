package author

import (
	"auth-svc/internal/models"
	"context"
	"errors"
)

// CreateRole inserts a new role into the database.
func (repo *AuthorizationRepository) CreateRole(ctx context.Context, role *models.Role) error {
	if role == nil {
		return errors.New("role is nil")
	}

	err := repo.roleCol.InsertOne(ctx, *role)
	if err != nil {
		return err
	}

	return nil
}

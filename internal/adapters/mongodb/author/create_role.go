package author

import (
	"auth-svc/internal/models"
	"context"
)

func (repo *AuthorizationRepository) CreateRole(ctx context.Context, role *models.Role) error {
	_, err := repo.RoleCol.InsertOne(ctx, role)

	return err
}

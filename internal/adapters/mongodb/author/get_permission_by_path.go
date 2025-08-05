package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPermissionByPath retrieves a permission path document by its path.
func (r *AuthorizationRepository) GetPermissionByPath(ctx context.Context, path string) (*models.PermissionPath, error) {

	permission, err := r.permissionPathCol.FindOne(ctx, bson.M{"path": path},
		options.FindOne().SetHint(IdxPermissionPath),
	)

	return permission, err
}

package author

import (
	"auth-svc/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *AuthorizationRepository) GetPermissionByPath(ctx context.Context, path string) (models.PermissionPath, error) {
	var permission models.PermissionPath
	err := r.PathPermissionCol.FindOne(ctx, bson.M{"path": path}, &options.FindOneOptions{
		Hint: IdxPermissionPath,
	}).Decode(&permission)
	return permission, err
}

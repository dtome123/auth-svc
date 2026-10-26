package author

import (
	"auth-svc/internal/domain"
	"auth-svc/internal/utils"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListPermissionsInput = domain.ListPermissionsInput

// ListPermissions retrieves permissions with optional domain filter and pagination.
func (r *AuthorizationRepository) ListPermissions(ctx context.Context, input domain.ListPermissionsInput) ([]*domain.Permission, error) {
	skip, limit := utils.PaginationToSkipLimit(input.PageNumber, input.PageLimit)

	findOpts := options.Find()
	if limit > 0 {
		findOpts.SetLimit(limit)
	}
	if skip > 0 {
		findOpts.SetSkip(skip)
	}

	filter := bson.M{}

	if input.Keyword != "" {
		filter["name"] = bson.M{
			"$regex":   input.Keyword,
			"$options": "i",
		}
	}

	if len(input.Domains) > 0 {
		filter["domain"] = bson.M{
			"$in": input.Domains,
		}
	}

	permDocs, err := r.permissionCol.Find(ctx, filter, findOpts, nil)
	if err != nil {
		return nil, err
	}

	permissions := make([]*domain.Permission, len(permDocs))
	for i, p := range permDocs {
		permissions[i] = p.ToDomain()
	}

	return permissions, nil
}

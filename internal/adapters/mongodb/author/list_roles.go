package author

import (
	"auth-svc/internal/domain"
	"auth-svc/internal/utils"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ListRoles retrieves roles from the database with pagination support by page number and page limit.
func (repo *impl) ListRoles(ctx context.Context, input domain.ListRoleInput) ([]*domain.Role, error) {
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
			"$options": "i", // case-insensitive
		}
	}

	if input.Scope != nil {
		for k, v := range input.Scope {
			if len(v) == 0 {
				continue
			}
			filter["scope."+k] = bson.M{"$in": v}
		}
	}

	roleDocs, err := repo.roleCol.Find(ctx, filter, findOpts, nil)
	if err != nil {
		return nil, err
	}

	roles := make([]*domain.Role, len(roleDocs))
	for i, r := range roleDocs {
		roles[i] = r.ToDomain()
	}

	return roles, nil
}

package auth

import (
	"auth-svc/internal/adapters/mongodb/author"
	"auth-svc/internal/models"
	"context"

	"github.com/pkg/errors"
)

type ListRoleInput struct {
	Keyword    string
	Scope      map[string][]string
	PageNumber int64
	PageLimit  int64
}

func (s *AuthorizationService) ListRole(ctx context.Context, input ListRoleInput) ([]models.Role, error) {

	roles, err := s.authorizationRepo.ListRoles(ctx, author.ListRoleInput{
		Keyword:    input.Keyword,
		Scope:      input.Scope,
		PageNumber: input.PageNumber,
		PageLimit:  input.PageLimit,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to list permission paths")
	}

	return roles, nil
}

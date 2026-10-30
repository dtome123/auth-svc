package auth

import (
	"auth-svc/internal/domain"
	"context"

	"github.com/pkg/errors"
)

type ListRoleInput = domain.ListRoleInput

func (s *AuthorizationService) ListRole(ctx context.Context, input ListRoleInput) ([]*domain.Role, error) {
	roles, err := s.authorizationRepo.ListRoles(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list roles")
	}

	allPermIDsSet := make(map[string]struct{})
	for _, r := range roles {
		for _, ra := range r.ResourceActions {
			if ra.PermissionID != "" {
				allPermIDsSet[ra.PermissionID] = struct{}{}
			}
		}
	}
	if len(allPermIDsSet) > 0 {
		allPermIDs := make([]string, 0, len(allPermIDsSet))
		for pid := range allPermIDsSet {
			allPermIDs = append(allPermIDs, pid)
		}
		perms, err := s.authorizationRepo.GetPermissionsByIDs(ctx, allPermIDs)
		if err == nil {
			permsByID := make(map[string]*domain.Permission, len(perms))
			for _, p := range perms {
				permsByID[p.ID] = p
			}
			for _, r := range roles {
				for _, ra := range r.ResourceActions {
					if p, ok := permsByID[ra.PermissionID]; ok {
						r.Permissions = append(r.Permissions, p)
					}
				}
			}
		}
	}

	return roles, nil
}

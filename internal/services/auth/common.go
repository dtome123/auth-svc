package auth

import (
	"auth-svc/internal/domain"
	casbinPkg "auth-svc/tools/casbin"
	"context"
	"fmt"
)

// LoadAllPolicies fetches all roles and permissions from repository
// and bulk-loads them into the in-memory Casbin Enforcer.
// Only Role -> Permission policies are loaded into RAM (no users).
func (svc *AuthorizationService) LoadAllPolicies(ctx context.Context) error {
	if svc.authorizationRepo == nil || svc.enforcer == nil {
		return nil
	}

	roles, err := svc.authorizationRepo.ListRoles(ctx, domain.ListRoleInput{})
	if err != nil {
		return fmt.Errorf("failed to list roles: %w", err)
	}

	permissions, err := svc.authorizationRepo.ListPermissions(ctx, domain.ListPermissionsInput{})
	if err != nil {
		return fmt.Errorf("failed to list permissions: %w", err)
	}

	policies := buildRolePolicies(roles, permissions)

	return svc.enforcer.LoadPolicies(nil, policies)
}

// getUserRoles retrieves the list of role names and IDs assigned to a user from repository.
func (svc *AuthorizationService) getUserRoles(ctx context.Context, userID string) ([]string, error) {
	if svc.authorizationRepo == nil || userID == "" {
		return nil, nil
	}
	assignment, err := svc.authorizationRepo.GetAssignmentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if assignment == nil || len(assignment.RoleIDs) == 0 {
		return nil, nil
	}
	roleList, err := svc.authorizationRepo.ListRolesByIDs(ctx, assignment.RoleIDs)
	if err != nil {
		return nil, err
	}
	roles := make([]string, 0, len(roleList)*2)
	for _, r := range roleList {
		if r.Name != "" {
			roles = append(roles, r.Name)
		}
		roles = append(roles, r.ID)
	}
	return roles, nil
}

// syncRolePoliciesInEnforcer synchronizes a role's permissions in the Casbin Enforcer.
func (svc *AuthorizationService) syncRolePoliciesInEnforcer(
	ctx context.Context,
	roleID, roleName string,
	resourceActions []domain.ResourceAction,
	permissionIDs []string,
) error {
	if svc.enforcer == nil {
		return nil
	}

	var allPermIDs []string
	if len(permissionIDs) > 0 {
		allPermIDs = append(allPermIDs, permissionIDs...)
	} else {
		for _, ra := range resourceActions {
			if ra.PermissionID != "" {
				allPermIDs = append(allPermIDs, ra.PermissionID)
			}
		}
	}

	permissions, err := svc.authorizationRepo.GetPermissionsByIDs(ctx, allPermIDs)
	if err != nil {
		return fmt.Errorf("failed to get role permissions: %w", err)
	}

	var policies []casbinPkg.Policy
	if len(resourceActions) > 0 {
		permsByID := make(map[string]*domain.Permission, len(permissions))
		for _, p := range permissions {
			permsByID[p.ID] = p
		}
		for _, ra := range resourceActions {
			perm, ok := permsByID[ra.PermissionID]
			if !ok {
				continue
			}
			actions := ra.Actions
			if len(actions) == 0 {
				actions = perm.Actions
			}
			for _, act := range actions {
				policies = append(policies, casbinPkg.Policy{Resource: perm.Resource, Action: act})
			}
		}
	} else {
		policies = permissionsToCasbinPolicies(permissions)
	}

	if err := svc.enforcer.SetRolePolicies(roleID, policies); err != nil {
		return err
	}
	if roleName != "" && roleName != roleID {
		return svc.enforcer.SetRolePolicies(roleName, policies)
	}
	return nil
}

// publishPolicySync sends a notification across all running service instances via Redis Pub/Sub.
func (svc *AuthorizationService) publishPolicySync(ctx context.Context) {
	if svc.enforcer == nil {
		return
	}
	svc.enforcer.PublishSync(ctx)
}

// ---------------------------------------------------------------------------
// Data transformation helpers
// ---------------------------------------------------------------------------

// buildRolePolicies converts domain Role & Permission models into raw Casbin policy slices.
// Supports matching both role ID and role name.
func buildRolePolicies(
	roles []*domain.Role,
	permissions []*domain.Permission,
) (policies [][]string) {

	permsByID := make(map[string]*domain.Permission, len(permissions))
	for _, perm := range permissions {
		permsByID[perm.ID] = perm
	}

	for _, role := range roles {
		roleID := role.ID
		if len(role.ResourceActions) > 0 {
			for _, ra := range role.ResourceActions {
				perm, ok := permsByID[ra.PermissionID]
				if !ok {
					continue
				}
				actions := ra.Actions
				if len(actions) == 0 {
					actions = perm.Actions
				}
				for _, act := range actions {
					policies = append(policies, []string{roleID, perm.Resource, act})
					if role.Name != "" && role.Name != roleID {
						policies = append(policies, []string{role.Name, perm.Resource, act})
					}
				}
			}
		} else {
			for _, permID := range role.PermissionIDs {
				perm, ok := permsByID[permID]
				if !ok {
					continue
				}
				for _, act := range perm.Actions {
					policies = append(policies, []string{roleID, perm.Resource, act})
					if role.Name != "" && role.Name != roleID {
						policies = append(policies, []string{role.Name, perm.Resource, act})
					}
				}
			}
		}
	}

	return policies
}

// permissionsToCasbinPolicies converts a slice of Permission models into casbin.Policy slice.
func permissionsToCasbinPolicies(permissions []*domain.Permission) []casbinPkg.Policy {
	var result []casbinPkg.Policy
	for _, perm := range permissions {
		for _, act := range perm.Actions {
			result = append(result, casbinPkg.Policy{Resource: perm.Resource, Action: act})
		}
	}
	return result
}

// formatUserPermissions converts []*domain.Permission into a slice of deduplicated "resource:action" strings.
func formatUserPermissions(permissions []*domain.Permission) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)

	add := func(res, act string) {
		if res == "" || act == "" {
			return
		}
		key := fmt.Sprintf("%s:%s", res, act)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, key)
		}
	}

	for _, perm := range permissions {
		if perm == nil {
			continue
		}
		for _, act := range perm.Actions {
			add(perm.Resource, act)
		}
	}

	return result
}

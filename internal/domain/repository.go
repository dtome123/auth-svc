package domain

import "context"

// ListRoleInput defines filter criteria for querying roles.
type ListRoleInput struct {
	Keyword    string
	Scope      map[string][]string
	PageNumber int64
	PageLimit  int64
}

// ListPermissionsInput defines filter criteria for querying permissions.
type ListPermissionsInput struct {
	Keyword    string
	Domains    []string
	PageNumber int64
	PageLimit  int64
}

// AuthorizationRepository is the secondary port interface for role and permission persistence.
type AuthorizationRepository interface {
	// Role management
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
	GetRole(ctx context.Context, id string) (*Role, error)
	ListRoles(ctx context.Context, input ListRoleInput) ([]*Role, error)
	ListRolesByIDs(ctx context.Context, ids []string) ([]*Role, error)

	// Permission management
	ListPermissions(ctx context.Context, input ListPermissionsInput) ([]*Permission, error)
	GetPermissionsByIDs(ctx context.Context, ids []string) ([]*Permission, error)
	GetPermissionsByUserID(ctx context.Context, userID string) ([]*Permission, error)
	BatchUpsertPermissions(ctx context.Context, permissions []*Permission) error

	// Assignment management
	GetAssignmentByUserID(ctx context.Context, userID string) (*Assignment, error)
	GetAssignmentsByRoleID(ctx context.Context, roleID string) ([]*Assignment, error)
	UpsertAssignment(ctx context.Context, assignment *Assignment) error
	ListAllAssignments(ctx context.Context) ([]*Assignment, error)
}

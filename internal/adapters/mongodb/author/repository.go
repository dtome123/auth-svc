package author

import (
	"auth-svc/internal/domain"
	"context"
)

// AuthorizationRepository is the secondary port interface for role and permission persistence.
type AuthorizationRepository interface {
	// Role management
	CreateRole(ctx context.Context, role *domain.Role) error
	UpdateRole(ctx context.Context, role *domain.Role) error
	GetRole(ctx context.Context, id string) (*domain.Role, error)
	GetRoleByCode(ctx context.Context, code string) (*domain.Role, error)
	ListRoles(ctx context.Context, input domain.ListRoleInput) ([]*domain.Role, error)
	ListRolesByIDs(ctx context.Context, ids []string) ([]*domain.Role, error)

	// Permission management
	ListPermissions(ctx context.Context, input domain.ListPermissionsInput) ([]*domain.Permission, error)
	GetPermissionsByIDs(ctx context.Context, ids []string) ([]*domain.Permission, error)
	GetPermissionsByUserID(ctx context.Context, userID string) ([]domain.ResourceAction, error)
	BatchUpsertPermissions(ctx context.Context, permissions []*domain.Permission) error

	// Assignment management
	GetAssignmentByUserID(ctx context.Context, userID string) (*domain.Assignment, error)
	GetAssignmentsByRoleID(ctx context.Context, roleID string) ([]*domain.Assignment, error)
	UpsertAssignment(ctx context.Context, assignment *domain.Assignment) error
	ListAllAssignments(ctx context.Context) ([]*domain.Assignment, error)
}

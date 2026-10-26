package casbin

import "context"

// IAM is the primary interface for Identity and Access Management policy enforcement,
// role-based access control (RBAC), and multi-instance synchronization.
type IAM interface {
	// Enforce evaluates if a subject/role is permitted to perform an action on a resource.
	// Matches are performed entirely in-memory with sub-microsecond latency.
	Enforce(subject, resource, action string) (bool, error)

	// EnforceRoles checks if ANY of the given roles is permitted to perform (resource, action).
	// Short-circuits immediately upon encountering the first allowed role.
	EnforceRoles(roles []string, resource, action string) (bool, error)

	// LoadPolicies atomically flushes and replaces all in-memory groupings and permission rules.
	//   - groupings: optional role hierarchy or user mappings, each entry is [sub, role]. Can be nil.
	//   - policies:  role-to-permission mappings, each entry is [role, resource, action].
	LoadPolicies(groupings, policies [][]string) error

	// SetUserRoles dynamically updates the assigned roles for a specific user.
	SetUserRoles(userID string, roleIDs []string) error

	// SetRolePolicies dynamically updates the permission rules for a specific role.
	SetRolePolicies(roleID string, policies []Policy) error

	// PublishSync broadcasts a policy-reload notification across all service instances via Redis Pub/Sub.
	PublishSync(ctx context.Context)

	// NotifySync is an alias for PublishSync.
	NotifySync(ctx context.Context)

	// OnSyncEvent registers a callback to execute whenever a policy-reload event is received
	// from Redis Pub/Sub. The subscription listener is started in a background goroutine.
	OnSyncEvent(fn func(ctx context.Context) error)
}

// Enforcer is an alias for IAM to maintain compatibility with Casbin naming conventions.
type Enforcer = IAM

// Policy specifies a resource and action permission tuple.
type Policy struct {
	Resource string `json:"resource" yaml:"resource"`
	Action   string `json:"action" yaml:"action"`
}

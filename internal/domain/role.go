package domain

// ResourceAction specifies actions permitted on a permission.
type ResourceAction struct {
	PermissionID string   `json:"permission_id"`
	Actions      []string `json:"actions"`
}

// Role represents an authorization role containing a set of permissions.
type Role struct {
	ID              string            `json:"id"`
	Code            string            `json:"code,omitempty"`
	Name            string            `json:"name"`
	ResourceActions []ResourceAction  `json:"resource_actions,omitempty"`
	Description     string            `json:"description,omitempty"`
	Scope           map[string]string `json:"scope,omitempty"`
	Permissions     []*Permission     `json:"permissions,omitempty"`
}

// ListRoleInput defines filter criteria for querying roles.
type ListRoleInput struct {
	Keyword    string
	Scope      map[string][]string
	PageNumber int64
	PageLimit  int64
}

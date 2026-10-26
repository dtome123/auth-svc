package domain

// ResourceAction specifies actions permitted on a permission.
type ResourceAction struct {
	PermissionID string   `json:"permission_id"`
	Actions      []string `json:"actions"`
}

// Role represents an authorization role containing a set of permissions.
type Role struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	ResourceActions []ResourceAction  `json:"resource_actions,omitempty"`
	PermissionIDs   []string          `json:"permission_ids,omitempty"`
	Description     string            `json:"description,omitempty"`
	Scope           map[string]string `json:"scope,omitempty"`
	Permissions     []*Permission     `json:"permissions,omitempty"`
}

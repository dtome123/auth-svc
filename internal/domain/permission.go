package domain

// Permission represents a single access rule targeting a resource and actions.
type Permission struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Domain      string   `json:"domain"`
	Description string   `json:"description,omitempty"`
	Resource    string   `json:"resource"`
	Actions     []string `json:"actions"`
}

// ListPermissionsInput defines filter criteria for querying permissions.
type ListPermissionsInput struct {
	Keyword    string
	Domains    []string
	PageNumber int64
	PageLimit  int64
}

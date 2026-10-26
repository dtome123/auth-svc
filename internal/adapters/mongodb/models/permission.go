package models

import (
	"auth-svc/internal/domain"
)

type Permission struct {
	ID          string   `bson:"_id,omitempty" json:"_id"`
	Name        string   `bson:"name" json:"name"`
	Domain      string   `bson:"domain" json:"domain"`
	Description string   `bson:"description,omitempty" json:"description,omitempty"`
	Resource    string   `bson:"resource" json:"resource"`
	Actions     []string `bson:"actions" json:"actions"`
}

func (Permission) CollectionName() string {
	return "permissions"
}

func (p *Permission) ToDomain() *domain.Permission {
	if p == nil {
		return nil
	}
	return &domain.Permission{
		ID:          p.ID,
		Name:        p.Name,
		Domain:      p.Domain,
		Description: p.Description,
		Resource:    p.Resource,
		Actions:     p.Actions,
	}
}

func PermissionFromDomain(p *domain.Permission) *Permission {
	if p == nil {
		return nil
	}
	return &Permission{
		ID:          p.ID,
		Name:        p.Name,
		Domain:      p.Domain,
		Description: p.Description,
		Resource:    p.Resource,
		Actions:     p.Actions,
	}
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	SchemaName string    `json:"schema_name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type TenantList struct {
	List []Tenant `json:"list"`
}

func NewTenant(tenant_request *Tenant) (*Tenant, error) {
	user := &Tenant{
		ID:         uuid.New(),
		Name:       tenant_request.Name,
		SchemaName: tenant_request.SchemaName,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return user, nil
}

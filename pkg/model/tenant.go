package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/potatowski/brazilcode"
)

type Tenant struct {
	ID         uuid.UUID  `json:"id"`
	GroupID    *uuid.UUID `json:"group_id,omitempty"` // optional, since not every tenant needs a group
	CNPJ       string     `json:"cnpj"`
	Name       string     `json:"name"`
	SchemaName string     `json:"schema_name"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
}

type TenantList struct {
	List []Tenant `json:"list"`
}

func NewTenant(tenant_request *Tenant) (*Tenant, error) {
	tenant := &Tenant{
		ID:         tenant_request.ID,
		CNPJ:       tenant_request.CNPJ,
		Name:       tenant_request.Name,
		SchemaName: tenant_request.ID.String(),
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return tenant, nil
}
func (t *Tenant) CheckCNPJ(cnpj string) bool {
	err := brazilcode.CNPJIsValid(cnpj)
	if err != nil {
		logger.Error("Error to check CNPJ", err)
		return false
	}

	return true
}

package dto

import (
	"time"

	"github.com/google/uuid"
)

type TenantRequestDtoInput struct {
	Name    string    `json:"name" binding:"required"`
	CNPJ    string    `json:"cnpj" binding:"required"`
	GroupID uuid.UUID `json:"group_id" binding:"required"`
}

type TenantRequestDtoOutPut struct {
	ID         uuid.UUID `json:"id"`
	GroupID    uuid.UUID `json:"group_id"`
	Name       string    `json:"name"`
	CNPJ       string    `json:"cnpj"`
	SchemaName string    `json:"schema_name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

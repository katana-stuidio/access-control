package dto

import (
	"time"

	"github.com/google/uuid"
)

type TenantRequestDtoInput struct {
	Name string `json:"name"`
	CNPJ string `json:"cnpj"`
}

type TenantRequestDtoOutPut struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CNPJ       string    `json:"cnpj"`
	SchemaName string    `json:"schema_name"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

package dto

import (
	"time"

	"github.com/google/uuid"
)

// TenantGroupCreateRequest represents the request to create a new tenant group
type TenantGroupCreateRequest struct {
	Name string `json:"name" binding:"required"`
	CNPJ string `json:"cnpj" binding:"required"`
}

// TenantGroupUpdateRequest represents the request to update a tenant group
type TenantGroupUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	CNPJ     string `json:"cnpj" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// TenantGroupResponse represents the response for tenant group operations
type TenantGroupResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TenantGroupListResponse represents the response for listing tenant groups
type TenantGroupListResponse struct {
	Total int64                 `json:"total"`
	Page  int64                 `json:"current_page"`
	Last  int64                 `json:"last_page"`
	Data  []TenantGroupResponse `json:"data"`
}

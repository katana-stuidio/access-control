package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/potatowski/brazilcode"
)

type TenantGroup struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type TenantGroupList struct {
	List []TenantGroup `json:"list"`
}

func NewTenantGroup(tenant_group_request *TenantGroup) (*TenantGroup, error) {
	tenant_group := &TenantGroup{
		ID:        uuid.New(),
		Name:      tenant_group_request.Name,
		CNPJ:      tenant_group_request.CNPJ,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return tenant_group, nil
}

func (tg *TenantGroup) CheckCNPJ(cnpj string) bool {
	err := brazilcode.CNPJIsValid(cnpj)
	if err != nil {
		logger.Error("Error to check CNPJ", err)
		return false
	}

	return true
}

func (tg *TenantGroup) PrepareToSave() {
	dt := time.Now()

	if tg.ID == uuid.Nil {
		tg.ID = uuid.New()
		tg.CreatedAt = dt
		tg.UpdatedAt = dt
	} else {
		tg.UpdatedAt = dt
	}
}

package tenant_group

import (
	"context"

	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/pkg/adapter/pgsql"
	"github.com/katana-stuidio/access-control/pkg/model"
)

type TenantGroupServiceInterface interface {
	GetAll(ctx context.Context, limit, page int64) (*model.Paginate, error)
	GetByID(ctx context.Context, ID uuid.UUID) *model.TenantGroup
	GetByCNPJ(ctx context.Context, CNPJ string) (*model.TenantGroup, error)
	Create(ctx context.Context, tenantGroup *model.TenantGroup) (*model.TenantGroup, error)
	Update(ctx context.Context, ID uuid.UUID, tenantGroup *model.TenantGroup) int64
	Delete(ctx context.Context, ID uuid.UUID) int64
	GetExistCNPJ(ctx context.Context, cnpj string) (bool, error)
}

type TenantGroup_service struct {
	dbp pgsql.DatabaseInterface
}

func NewTenantGroupService(database_pool pgsql.DatabaseInterface) *TenantGroup_service {
	return &TenantGroup_service{
		dbp: database_pool,
	}
}

func (tgs *TenantGroup_service) GetAll(ctx context.Context, limit, page int64) (*model.Paginate, error) {
	// Get total count
	var total int64
	err := tgs.dbp.GetDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM tb_tenant_group").Scan(&total)
	if err != nil {
		logger.Error("Error getting total count", err)
		return nil, err
	}

	// Create pagination object
	paginate := model.NewPaginate(limit, page, total)

	// Get paginated data
	offset := (paginate.Page - 1) * paginate.Limit
	rows, err := tgs.dbp.GetDB().QueryContext(ctx,
		"SELECT id, name, cnpj, is_active, created_at, updated_at FROM tb_tenant_group LIMIT $1 OFFSET $2",
		paginate.Limit, offset)
	if err != nil {
		logger.Error("Error querying tenant groups", err)
		return nil, err
	}
	defer rows.Close()

	tenant_group_list := &model.TenantGroupList{}
	for rows.Next() {
		tg := model.TenantGroup{}
		if err := rows.Scan(&tg.ID, &tg.Name, &tg.CNPJ, &tg.IsActive, &tg.CreatedAt, &tg.UpdatedAt); err != nil {
			logger.Error("Error scanning tenant group", err)
			return nil, err
		}
		tenant_group_list.List = append(tenant_group_list.List, tg)
	}

	paginate.Paginate(tenant_group_list)
	return paginate, nil
}

func (tgs *TenantGroup_service) GetByID(ctx context.Context, ID uuid.UUID) *model.TenantGroup {
	stmt, err := tgs.dbp.GetDB().PrepareContext(ctx, "SELECT id, name, cnpj, is_active, created_at, updated_at FROM tb_tenant_group WHERE id = $1")
	if err != nil {
		logger.Error(err.Error(), err)
	}

	defer stmt.Close()

	tg := model.TenantGroup{}

	if err := stmt.QueryRowContext(ctx, ID).Scan(&tg.ID, &tg.Name, &tg.CNPJ, &tg.IsActive, &tg.CreatedAt, &tg.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
	}

	return &tg
}

func (tgs *TenantGroup_service) Create(ctx context.Context, tenantGroup *model.TenantGroup) (*model.TenantGroup, error) {
	tx, err := tgs.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error(), err)
		return tenantGroup, err
	}

	query := "INSERT INTO tb_tenant_group (id, name, cnpj, is_active) VALUES ($1, $2, $3, $4)"

	_, err = tx.ExecContext(ctx, query, tenantGroup.ID, tenantGroup.Name, tenantGroup.CNPJ, tenantGroup.IsActive)
	if err != nil {
		logger.Error("Error executing SQL query insert tenant group", err)
		return tenantGroup, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error committing transaction", err)
		return tenantGroup, err
	} else {
		logger.Info("Insert Transaction committed")
	}

	return tenantGroup, nil
}

func (tgs *TenantGroup_service) Update(ctx context.Context, ID uuid.UUID, tenantGroup *model.TenantGroup) int64 {
	tx, err := tgs.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction", err)
	}

	query := "UPDATE tb_tenant_group SET name = $1, cnpj = $2, is_active = $3 WHERE id = $4"

	result, err := tx.ExecContext(ctx, query, tenantGroup.Name, tenantGroup.CNPJ, tenantGroup.IsActive, ID)
	if err != nil {
		logger.Error("Error updating tenant group", err)
		return 0
	}

	err = tx.Commit()
	if err != nil {
		logger.Error("Error committing transaction", err)
		tx.Rollback()
		return 0
	} else {
		logger.Info("Update Transaction committed")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected", err)
		return 0
	}

	return rowsAff
}

func (tgs *TenantGroup_service) Delete(ctx context.Context, ID uuid.UUID) int64 {
	tx, err := tgs.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction", err)
	}

	query := "DELETE FROM tb_tenant_group WHERE id = $1"

	result, err := tx.ExecContext(ctx, query, ID)
	if err != nil {
		logger.Error("Error deleting tenant group", err)
		return 0
	}

	err = tx.Commit()
	if err != nil {
		logger.Error("Error committing transaction", err)
		tx.Rollback()
		return 0
	} else {
		logger.Info("Delete Transaction committed")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected", err)
		return 0
	}

	return rowsAff
}

func (tgs *TenantGroup_service) GetExistCNPJ(ctx context.Context, cnpj string) (bool, error) {
	query := "SELECT COUNT(*) FROM tb_tenant_group WHERE cnpj = $1"
	var count int

	err := tgs.dbp.GetDB().QueryRowContext(ctx, query, cnpj).Scan(&count)
	if err != nil {
		logger.Error("Error checking existing CNPJ", err)
		return false, err
	}

	return count > 0, nil
}

func (tgs *TenantGroup_service) GetByCNPJ(ctx context.Context, CNPJ string) (*model.TenantGroup, error) {
	stmt, err := tgs.dbp.GetDB().PrepareContext(ctx, "SELECT id, name, cnpj, is_active, created_at, updated_at FROM tb_tenant_group WHERE cnpj = $1")
	tg := model.TenantGroup{}
	if err != nil {
		logger.Error(err.Error(), err)
		return &tg, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, CNPJ).Scan(&tg.ID, &tg.Name, &tg.CNPJ, &tg.IsActive, &tg.CreatedAt, &tg.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
		return &tg, err
	}

	return &tg, nil
}

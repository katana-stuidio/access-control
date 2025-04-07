package tenant

import (
	"context"

	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/pkg/adapter/pgsql"
	"github.com/katana-stuidio/access-control/pkg/model"
)

type TenantServiceInterface interface {
	GetAll(ctx context.Context) *model.TenantList
	GetByID(ctx context.Context, ID uuid.UUID) *model.Tenant
	Create(ctx context.Context, Tenant *model.Tenant) (*model.Tenant, error)
	Update(ctx context.Context, ID uuid.UUID, Tenant *model.Tenant) int64
	Delete(ctx context.Context, ID uuid.UUID) int64
	GetExistTenantName(ctx context.Context, TenantName string) (bool, error)
}

type Tenant_service struct {
	dbp pgsql.DatabaseInterface
}

func NewTenantService(database_pool pgsql.DatabaseInterface) *Tenant_service {
	return &Tenant_service{
		dbp: database_pool,
	}
}

func (us *Tenant_service) GetAll(ctx context.Context) *model.TenantList {

	rows, err := us.dbp.GetDB().QueryContext(ctx, "SELECT id, name, schema_name, is_active, created_at, updated_at  FROM tenants LIMIT 100")
	if err != nil {
		logger.Error(err.Error(), err)
	}

	defer rows.Close()

	Tenant_list := &model.TenantList{}

	for rows.Next() {
		tn := model.Tenant{}
		if err := rows.Scan(&tn.ID, &tn.Name, &tn.SchemaName, &tn.IsActive, &tn.CreatedAt, &tn.UpdatedAt); err != nil {
			logger.Error(err.Error(), err)
		} else {
			Tenant_list.List = append(Tenant_list.List, tn)
		}
	}

	return Tenant_list
}

func (us *Tenant_service) GetByID(ctx context.Context, ID uuid.UUID) *model.Tenant {
	stmt, err := us.dbp.GetDB().PrepareContext(ctx, "SELECT id, name, schema_name, is_active, created_at, updated_at  FROM tenants WHERE id = $1")
	if err != nil {
		logger.Error(err.Error(), err)
	}

	defer stmt.Close()

	tn := model.Tenant{}

	if err := stmt.QueryRowContext(ctx, ID).Scan(&tn.ID, &tn.Name, &tn.SchemaName, &tn.IsActive, &tn.CreatedAt, &tn.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
	}

	return &tn
}

func (us *Tenant_service) Create(ctx context.Context, Tenant *model.Tenant) (*model.Tenant, error) {
	tx, err := us.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error(), err)

		return Tenant, err
	}

	query := "INSERT INTO tenants (name, schema_name, is_active, created_at, updated_at ) VALUES ($1, $2, $3, $4, $5)"

	_, err = tx.ExecContext(ctx, query, Tenant.Name, Tenant.SchemaName, Tenant.IsActive, Tenant.CreatedAt, Tenant.UpdatedAt)
	if err != nil {
		logger.Error("Error executing SQL query insert Tenant", err)

		return Tenant, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error committing transaction", err)

		return Tenant, err
	} else {
		logger.Info("Insert Transaction committed")
	}

	return Tenant, nil
}

func (us *Tenant_service) Update(ctx context.Context, ID uuid.UUID, Tenant *model.Tenant) int64 {
	tx, err := us.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction", err)
	}

	query := "UPDATE tenants SET name = $1, schema_name = $2, is_active = $3, updated_at = $4 WHERE id = $5"

	result, err := tx.ExecContext(ctx, query, Tenant.Name, Tenant.SchemaName, Tenant.IsActive, Tenant.UpdatedAt, ID)
	if err != nil {
		logger.Error("Error updating Tenant", err)
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

func (us *Tenant_service) Delete(ctx context.Context, ID uuid.UUID) int64 {
	tx, err := us.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction", err)
	}

	query := "DELETE FROM tenants WHERE id = $1"

	result, err := tx.ExecContext(ctx, query, ID)
	if err != nil {
		logger.Error("Error deleting Tenant", err)
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

func (us *Tenant_service) GetExistTenantName(ctx context.Context, tenantName string) (bool, error) {
	query := "SELECT COUNT(*) FROM tenants WHERE name = $1"
	var count int

	err := us.dbp.GetDB().QueryRowContext(ctx, query, tenantName).Scan(&count)
	if err != nil {
		logger.Error("Error checking existing Tenantname", err)
		return false, err
	}

	return count > 0, nil // Return true if count is greater than 0
}

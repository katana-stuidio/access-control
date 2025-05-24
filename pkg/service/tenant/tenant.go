package tenant

import (
	"context"

	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/pkg/adapter/pgsql"
	"github.com/katana-stuidio/access-control/pkg/model"
)

type TenantServiceInterface interface {
	GetAll(ctx context.Context, limit, page int64) (*model.Paginate, error)
	GetByID(ctx context.Context, ID uuid.UUID) *model.Tenant
	GetByCNPJ(ctx context.Context, CNPJ string) (*model.Tenant, error)
	Create(ctx context.Context, tenant *model.Tenant) (*model.Tenant, error)
	Update(ctx context.Context, ID uuid.UUID, tenant *model.Tenant) int64
	Delete(ctx context.Context, ID uuid.UUID) int64
	GetExistCNPJ(ctx context.Context, cnpj string) (bool, error)
}

type Tenant_service struct {
	dbp pgsql.DatabaseInterface
}

func NewTenantService(database_pool pgsql.DatabaseInterface) *Tenant_service {
	return &Tenant_service{
		dbp: database_pool,
	}
}

func (ts *Tenant_service) GetAll(ctx context.Context, limit, page int64) (*model.Paginate, error) {
	// Get total count
	var total int64
	err := ts.dbp.GetDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM tb_tenant").Scan(&total)
	if err != nil {
		logger.Error("Error getting total count", err)
		return nil, err
	}

	// Create pagination object
	paginate := model.NewPaginate(limit, page, total)

	// Get paginated data
	offset := (paginate.Page - 1) * paginate.Limit
	rows, err := ts.dbp.GetDB().QueryContext(ctx,
		"SELECT id, cnpj, name, schema_name, is_active, created_at, updated_at FROM tb_tenant LIMIT $1 OFFSET $2",
		paginate.Limit, offset)
	if err != nil {
		logger.Error("Error querying tenants", err)
		return nil, err
	}
	defer rows.Close()

	tenant_list := &model.TenantList{}
	for rows.Next() {
		t := model.Tenant{}
		if err := rows.Scan(&t.ID, &t.CNPJ, &t.Name, &t.SchemaName, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
			logger.Error("Error scanning tenant", err)
			return nil, err
		}
		tenant_list.List = append(tenant_list.List, t)
	}

	paginate.Paginate(tenant_list)
	return paginate, nil
}

func (ts *Tenant_service) GetByID(ctx context.Context, ID uuid.UUID) *model.Tenant {
	stmt, err := ts.dbp.GetDB().PrepareContext(ctx, "SELECT id, cnpj, name, schema_name, is_active, created_at, updated_at FROM tb_tenant WHERE id = $1")
	if err != nil {
		logger.Error(err.Error(), err)
	}

	defer stmt.Close()

	t := model.Tenant{}

	if err := stmt.QueryRowContext(ctx, ID).Scan(&t.ID, &t.CNPJ, &t.Name, &t.SchemaName, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
	}

	return &t
}

func (ts *Tenant_service) Create(ctx context.Context, tenant *model.Tenant) (*model.Tenant, error) {
	tx, err := ts.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error(), err)
		return tenant, err
	}

	query := "INSERT INTO tb_tenant (id, cnpj, name, schema_name, is_active) VALUES ($1, $2, $3, $4, $5)"

	_, err = tx.ExecContext(ctx, query, tenant.ID, tenant.CNPJ, tenant.Name, tenant.SchemaName, tenant.IsActive)
	if err != nil {
		logger.Error("Error executing SQL query insert tenant", err)
		return tenant, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error committing transaction", err)
		return tenant, err
	} else {
		logger.Info("Insert Transaction committed")
	}

	return tenant, nil
}

func (ts *Tenant_service) Update(ctx context.Context, ID uuid.UUID, tenant *model.Tenant) int64 {
	tx, err := ts.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction", err)
	}

	query := "UPDATE tb_tenant SET cnpj = $1, name = $2, schema_name = $3, is_active = $4 WHERE id = $5"

	result, err := tx.ExecContext(ctx, query, tenant.CNPJ, tenant.Name, tenant.SchemaName, tenant.IsActive, ID)
	if err != nil {
		logger.Error("Error updating tenant", err)
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

func (ts *Tenant_service) Delete(ctx context.Context, ID uuid.UUID) int64 {
	tx, err := ts.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction", err)
	}

	query := "DELETE FROM tb_tenant WHERE id = $1"

	result, err := tx.ExecContext(ctx, query, ID)
	if err != nil {
		logger.Error("Error deleting tenant", err)
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

func (ts *Tenant_service) GetExistCNPJ(ctx context.Context, cnpj string) (bool, error) {
	query := "SELECT COUNT(*) FROM tb_tenant WHERE cnpj = $1"
	var count int

	err := ts.dbp.GetDB().QueryRowContext(ctx, query, cnpj).Scan(&count)
	if err != nil {
		logger.Error("Error checking existing CNPJ", err)
		return false, err
	}

	return count > 0, nil
}

func (ts *Tenant_service) GetByCNPJ(ctx context.Context, CNPJ string) (*model.Tenant, error) {
	stmt, err := ts.dbp.GetDB().PrepareContext(ctx, "SELECT id, cnpj, name, schema_name, is_active, created_at, updated_at FROM tb_tenant WHERE cnpj = $1")
	t := model.Tenant{}
	if err != nil {
		logger.Error(err.Error(), err)
		return &t, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, CNPJ).Scan(&t.ID, &t.CNPJ, &t.Name, &t.SchemaName, &t.IsActive, &t.CreatedAt, &t.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
		return &t, err
	}

	return &t, nil
}

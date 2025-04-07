package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/pkg/adapter/pgsql"
	"github.com/katana-stuidio/access-control/pkg/model"
)

type UserServiceInterface interface {
	GetAll(ctx context.Context) *model.UserList
	GetByID(ctx context.Context, ID uuid.UUID) *model.User
	GetByEmail(ctx context.Context, email string) (usr *model.User, err error)
	Create(ctx context.Context, User *model.User) (*model.User, error)
	Update(ctx context.Context, ID uuid.UUID, User *model.User) int64
	Delete(ctx context.Context, ID uuid.UUID) int64
	GetExistUserName(ctx context.Context, userName string) (bool, error)
	Authenticate(username, password string) (*model.User, error)
}

type User_service struct {
	dbp pgsql.DatabaseInterface
}

func NewUserService(database_pool pgsql.DatabaseInterface) *User_service {
	return &User_service{
		dbp: database_pool,
	}
}

func (us *User_service) GetAll(ctx context.Context) *model.UserList {
	rows, err := us.dbp.GetDB().QueryContext(ctx, "SELECT id, id_tanant, username, name_full, email, enabled, created_at, updated_at FROM tb_user LIMIT 100")
	if err != nil {
		logger.Error(err.Error(), err)
	}

	defer rows.Close()

	User_list := &model.UserList{}

	for rows.Next() {
		u := model.User{}
		if err := rows.Scan(&u.ID, &u.TanantID, &u.Username, &u.Name, &u.Email, &u.Enable, &u.CreatedAt, &u.UpdatedAt); err != nil {
			logger.Error(err.Error(), err)
		} else {
			User_list.List = append(User_list.List, &u)
		}
	}

	return User_list
}

func (us *User_service) GetByID(ctx context.Context, ID uuid.UUID) *model.User {
	stmt, err := us.dbp.GetDB().PrepareContext(ctx, "SELECT id, id_tanant, username, name_full, email, enabled, created_at, updated_at FROM tb_user WHERE id = $1")
	if err != nil {
		logger.Error(err.Error(), err)
	}

	defer stmt.Close()

	u := model.User{}

	if err := stmt.QueryRowContext(ctx, ID).Scan(&u.ID, &u.TanantID, &u.Username, &u.Name, &u.Email, &u.Enable, &u.CreatedAt, &u.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
	}

	return &u
}

func (us *User_service) Create(ctx context.Context, User *model.User) (*model.User, error) {
	tx, err := us.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error(), err)

		return User, err
	}
	User.CheckPassword(User.Password)
	User.PrepareToSave()

	query := "INSERT INTO tb_user (id_tanant, username, name_full, hashed_password, email, enabled) VALUES ($1, $2, $3, $4, $5)"

	_, err = tx.ExecContext(ctx, query, User.TanantID, User.Username, User.Name, User.HashedPassword, User.Email, User.Enable)
	if err != nil {
		logger.Error("Error executing SQL query insert user", err)

		return User, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error committing transaction", err)

		return User, err
	} else {
		logger.Info("Insert Transaction committed")
	}

	return User, nil
}

func (us *User_service) Update(ctx context.Context, ID uuid.UUID, User *model.User) int64 {
	tx, err := us.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction", err)
	}

	query := "UPDATE tb_user SET id_tanant = $1, username = $2, name_full = $3, password = $4, email = $5, enabled = $6 WHERE id = $7"

	result, err := tx.ExecContext(ctx, query, User.TanantID, User.Username, User.Name, User.Password, User.Email, User.Enable, ID)
	if err != nil {
		logger.Error("Error updating user", err)
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

func (us *User_service) Delete(ctx context.Context, ID uuid.UUID) int64 {
	tx, err := us.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction", err)
	}

	query := "DELETE FROM tb_user WHERE id = $1"

	result, err := tx.ExecContext(ctx, query, ID)
	if err != nil {
		logger.Error("Error deleting user", err)
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

func (us *User_service) GetExistUserName(ctx context.Context, userName string) (bool, error) {
	query := "SELECT COUNT(*) FROM tb_user WHERE username = $1"
	var count int

	err := us.dbp.GetDB().QueryRowContext(ctx, query, userName).Scan(&count)
	if err != nil {
		logger.Error("Error checking existing username", err)
		return false, err
	}

	return count > 0, nil // Return true if count is greater than 0
}

func (us *User_service) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	stmt, err := us.dbp.GetDB().PrepareContext(ctx, "SELECT id, id_tanant, username, name_full, email, enabled, created_at, updated_at FROM tb_user WHERE username = $1")
	u := model.User{}
	if err != nil {
		logger.Error(err.Error(), err)
		return &u, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, email).Scan(&u.ID, &u.TanantID, &u.Username, &u.Name, &u.Email, &u.Enable, &u.CreatedAt, &u.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
		return &u, err
	}

	return &u, nil
}

func (us *User_service) Authenticate(username, password string) (*model.User, error) {
	ctx := context.Background() // Ou use um contexto relevante

	stmt, err := us.dbp.GetDB().PrepareContext(ctx, "SELECT id, id_tanant, username, name_full , email, enabled, hashed_password, created_at, updated_at FROM tb_user WHERE username = $1")
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, err
	}
	defer stmt.Close()

	u := &model.User{}
	var hashedPassword string

	if err := stmt.QueryRowContext(ctx, username).Scan(&u.ID, &u.TanantID, &u.Username, &u.Name, &u.Email, &u.Enable, &hashedPassword, &u.CreatedAt, &u.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.New("invalid username or password")
	}

	u.HashedPassword = hashedPassword

	if !u.CheckPassword(password) {
		return nil, errors.New("invalid username or password")
	}

	return u, nil
}

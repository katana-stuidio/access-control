package user

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/katana-stuidio/access-control/pkg/adapter/pgsql"
	"github.com/katana-stuidio/access-control/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	GetAll(ctx context.Context, limit, page int64) (*model.Paginate, error)
	GetByID(ctx context.Context, ID uuid.UUID) *model.User
	GetByUserName(ctx context.Context, userName string) (usr *model.User, err error)
	Create(ctx context.Context, User *model.User) (*model.User, error)
	Update(ctx context.Context, ID uuid.UUID, User *model.User) int64
	Delete(ctx context.Context, ID uuid.UUID) int64
	GetExistUserName(ctx context.Context, userName string) (bool, error)
	Authenticate(username, password string) (*model.User, error)
	GetByCNPJ(ctx context.Context, CNPJ string) (tenant_id string, err error)
	ChangePassword(ctx context.Context, userName, currentPassword, newPassword string) error
	UpdatePassword(ctx context.Context, userName, newPassword string) int64
	EmailExists(ctx context.Context, email string) (bool, error)
}

type User_service struct {
	dbp pgsql.DatabaseInterface
}

func NewUserService(database_pool pgsql.DatabaseInterface) *User_service {
	return &User_service{
		dbp: database_pool,
	}
}

func (us *User_service) GetAll(ctx context.Context, limit, page int64) (*model.Paginate, error) {
	// Count total records
	var total int64
	countQuery := `SELECT COUNT(*) FROM user`
	err := us.dbp.GetDB().QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, err
	}

	// Create pagination
	paginate := model.NewPaginate(limit, page, total)

	// Get paginated results
	query := `
        SELECT id, cnpj, username, name_full, email, enabled, role_usr, created_at, updated_at 
        FROM user 
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2`

	offset := (paginate.Page - 1) * paginate.Limit
	rows, err := us.dbp.GetDB().QueryContext(ctx, query, paginate.Limit, offset)
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(
			&user.ID,
			&user.CNPJ,
			&user.Username,
			&user.Name,
			&user.Email,
			&user.Enable,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			logger.Error(err.Error(), err)
			return nil, err
		}
		users = append(users, user)
	}

	paginate.Paginate(users)
	return paginate, nil
}

func (us *User_service) GetByID(ctx context.Context, ID uuid.UUID) *model.User {
	stmt, err := us.dbp.GetDB().PrepareContext(ctx, "SELECT id, id_tanant, username, name_full, email, enabled, role_usr, created_at, updated_at FROM tb_user WHERE id = $1")
	if err != nil {
		logger.Error(err.Error(), err)
	}

	defer stmt.Close()

	u := model.User{}

	if err := stmt.QueryRowContext(ctx, ID).Scan(&u.ID, &u.TenantID, &u.Username, &u.Name, &u.Email, &u.Enable, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
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

	// Ensure password is hashed
	if User.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("Error hashing password", err)
			return User, err
		}
		User.HashedPassword = string(hashedPassword)
	}

	logger.Info("Creating user with role: " + User.Role)

	query := "INSERT INTO tb_user (id, id_tanant, username, name_full, hashed_password, email, enabled, role_usr) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	_, err = tx.ExecContext(ctx, query, User.ID, User.TenantID, User.Username, User.Name, User.HashedPassword, User.Email, User.Enable, User.Role)
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

	query := "UPDATE tb_user SET id_tanant = $1, username = $2, name_full = $3, password = $4, email = $5, enabled = $6, role_usr = $7 WHERE id = $8"

	result, err := tx.ExecContext(ctx, query, User.TenantID, User.Username, User.Name, User.Password, User.Email, User.Enable, User.Role, ID)
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

func (us *User_service) GetByUserName(ctx context.Context, email string) (*model.User, error) {
	stmt, err := us.dbp.GetDB().PrepareContext(ctx, "SELECT id, id_tanant, username, name_full, email, enabled, hashed_password, role_usr, created_at, updated_at FROM tb_user WHERE username = $1")
	u := model.User{}
	if err != nil {
		logger.Error(err.Error(), err)
		return &u, err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, email).Scan(&u.ID, &u.TenantID, &u.Username, &u.Name, &u.Email, &u.Enable, &u.HashedPassword, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
		return &u, err
	}

	return &u, nil
}

func (us *User_service) Authenticate(username, password string) (*model.User, error) {
	ctx := context.Background() // Ou use um contexto relevante

	stmt, err := us.dbp.GetDB().PrepareContext(ctx, "SELECT id, id_tanant, username, name_full, email, enabled, hashed_password, role_usr, created_at, updated_at FROM tb_user WHERE username = $1")
	if err != nil {
		logger.Error(err.Error(), err)
		return nil, err
	}
	defer stmt.Close()

	u := &model.User{}
	var hashedPassword string

	if err := stmt.QueryRowContext(ctx, username).Scan(&u.ID, &u.TenantID, &u.Username, &u.Name, &u.Email, &u.Enable, &hashedPassword, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		logger.Error(err.Error(), err)
		return nil, errors.New("invalid username or password")
	}

	u.HashedPassword = hashedPassword

	if !u.CheckPassword(password) {
		return nil, errors.New("invalid username or password")
	}

	return u, nil
}

func (us *User_service) GetByCNPJ(ctx context.Context, CNPJ string) (tenant_id string, err error) {
	query := "SELECT id FROM tb_tenant WHERE cnpj = $1"
	err = us.dbp.GetDB().QueryRowContext(ctx, query, CNPJ).Scan(&tenant_id)
	if err != nil {
		logger.Error("Error getting tenant id by cnpj", err)
		return "", err
	}
	return tenant_id, nil
}

func (us *User_service) UpdatePassword(ctx context.Context, userName, newPassword string) int64 {
	_, err := us.GetByUserName(ctx, userName)
	if err != nil {
		logger.Error("User not found for password update: "+userName, err)
		return 0
	}

	tx, err := us.dbp.GetDB().BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction for user: "+userName, err)
		return 0
	}
	defer tx.Rollback() // Rollback if not committed

	query := "UPDATE tb_user SET hashed_password = $1 WHERE username = $2"
	logger.Info("Executing query: " + query)

	result, err := tx.ExecContext(ctx, query, newPassword, userName)
	if err != nil {
		logger.Error("Error updating password for user: "+userName, err)
		return 0
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected for user: "+userName, err)
		return 0
	}

	if rowsAff == 0 {
		logger.Error("No rows affected for user: "+userName, nil)
		return 0
	}

	err = tx.Commit()
	if err != nil {
		logger.Error("Error committing transaction for user: "+userName, err)
		return 0
	}

	logger.Info("Password successfully updated for user: " + userName)
	return rowsAff
}

func (us *User_service) ChangePassword(ctx context.Context, userName, currentPassword, newPassword string) error {
	logger.Info("Starting password change process for user: " + userName)

	var user *model.User
	var err error

	user, err = us.GetByUserName(ctx, userName)
	if err != nil {
		logger.Error("Failed to find user: "+userName, err)
		return fmt.Errorf("failed to find user: %v", err)
	}

	logger.Info("User found, verifying current password")
	logger.Info(fmt.Sprintf("user: %v", user))

	// Check if hashed password is valid
	if user.HashedPassword == "" {
		logger.Error("No hashed password found for user: "+userName, nil)
		return fmt.Errorf("no password set for this user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(currentPassword)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			logger.Error("Current password does not match for user: "+userName, err)
			return fmt.Errorf("current password is incorrect")
		}
		logger.Error("Error comparing passwords for user: "+userName, err)
		return fmt.Errorf("error verifying current password: %v", err)
	}

	logger.Info("Current password verified, validating new password requirements")

	if len(newPassword) < 8 {
		logger.Info("New password length is too short for user: " + userName)
		return fmt.Errorf("new password length is too short")
	}

	const (
		uppercasePattern = "^(.*[A-Z]).*$"
		numberPattern    = "^(.*[0-9]).*$"
		symbolPattern    = "^(.*[!@#$%^&*()\\-_+=]).*$"
	)

	matchCapital, errCapital := regexp.MatchString(uppercasePattern, newPassword)
	if errCapital != nil || !matchCapital {
		logger.Info("New password must contain at least one uppercase letter for user: " + userName)
		return fmt.Errorf("new password must contain at least one uppercase letter")
	}

	matchNumber, errNumber := regexp.MatchString(numberPattern, newPassword)
	if errNumber != nil || !matchNumber {
		logger.Info("New password must contain at least one number for user: " + userName)
		return fmt.Errorf("new password must contain at least one number")
	}

	matchSymbol, errSymbol := regexp.MatchString(symbolPattern, newPassword)
	if errSymbol != nil || !matchSymbol {
		logger.Info("New password must contain at least one symbol for user: " + userName)
		return fmt.Errorf("new password must contain at least one symbol")
	}

	logger.Info("New password requirements met, generating hash")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error generating hashed password for user: "+userName, err)
		return fmt.Errorf("error generating password hash: %v", err)
	}

	pw := string(hashedPassword)
	user.Password = pw
	user.ChangePassword = false

	logger.Info("Updating password in database for user: " + userName)
	result := us.UpdatePassword(ctx, userName, pw)

	if result == 0 {
		logger.Error("Failed to update password in database for user: "+userName, nil)
		return fmt.Errorf("failed to update password in database")
	}

	logger.Info("Password successfully changed for user: " + userName)
	return nil
}

func (us *User_service) EmailExists(ctx context.Context, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM tb_user WHERE email = $1"
	var count int
	err := us.dbp.GetDB().QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		logger.Error("Error checking if email exists", err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

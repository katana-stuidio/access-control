package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/katana-stuidio/access-control/pkg/model"

	pgsql_mocks "github.com/katana-stuidio/access-control/pkg/adapter/pgsql/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceTestSuite struct {
	suite.Suite
	db      *sql.DB
	mock    sqlmock.Sqlmock
	service *User_service
	ctrl    *gomock.Controller                 // Adicione o controlador do gomock
	mockDB  *pgsql_mocks.MockDatabaseInterface // Adicione o mock da interface
}

func (suite *UserServiceTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.ctrl = gomock.NewController(suite.T())                    // Inicialize o controlador do gomock
	suite.mockDB = pgsql_mocks.NewMockDatabaseInterface(suite.ctrl) // Crie o mock da interface

	// Configure o comportamento do mock
	suite.mockDB.EXPECT().GetDB().Return(suite.db).AnyTimes()

	suite.service = &User_service{
		dbp: suite.mockDB, // Use o mock no lugar da implementação real
	}
}

func (suite *UserServiceTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	suite.db.Close()
	suite.ctrl.Finish() // Finalize o controlador do gomock
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestGetAll() {
	ctx := context.Background()
	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "enable", "created_at", "updated_at"}).
		AddRow(uuid.New(), "user1", "User One", "user1@example.com", true, time.Now(), time.Now()).
		AddRow(uuid.New(), "user2", "User Two", "user2@example.com", false, time.Now(), time.Now())

	suite.mock.ExpectQuery("SELECT id, username, name, email, enable, created_at, updated_at FROM tb_user LIMIT 100").WillReturnRows(rows)

	userList := suite.service.GetAll(ctx)

	assert.NotNil(suite.T(), userList)
	assert.Equal(suite.T(), 2, len(userList.List))
}
func (suite *UserServiceTestSuite) TestGetByID() {
	ctx := context.Background()
	userID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "enable", "created_at", "updated_at"}).
		AddRow(userID, "user1", "User One", "user1@example.com", true, time.Now(), time.Now())

	suite.mock.ExpectQuery("SELECT id, username, name, email, enable, created_at, updated_at FROM tb_user WHERE id = ?").
		WithArgs(userID).
		WillReturnRows(rows)

	user := suite.service.GetByID(ctx, userID)

	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), userID, user.ID)
}

func (suite *UserServiceTestSuite) TestGetByCpf() {
	ctx := context.Background()
	cpf := "12345678901"
	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "enable", "created_at", "updated_at"}).
		AddRow(uuid.New(), "user1", "User One", "user1@example.com", true, time.Now(), time.Now())

	suite.mock.ExpectQuery("SELECT id, username, name, email, enable, created_at, updated_at FROM tb_user WHERE username = ?").
		WithArgs(cpf).
		WillReturnRows(rows)

	user, err := suite.service.GetByCpf(ctx, cpf)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
}

func (suite *UserServiceTestSuite) TestCreate() {
	ctx := context.Background()
	user := &model.User{
		ID:       uuid.New(),
		Username: "user1",
		Name:     "User One",
		Email:    "user1@example.com",
		Enable:   true,
	}

	suite.mock.ExpectExec("INSERT INTO tb_user (username, name, hashed_password, email, enable) VALUES (?, ?, ?, ?, ?)").
		WithArgs(user.Username, user.Name, user.HashedPassword, user.Email, user.Enable).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := suite.service.Create(ctx, user)

	assert.NoError(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestUpdate() {
	ctx := context.Background()
	userID := uuid.New()
	user := &model.User{
		Username: "user1",
		Name:     "User One",
		Email:    "user1@example.com",
		Enable:   true,
	}

	suite.mock.ExpectExec("UPDATE tb_user SET username = ?, name = ?, password = ?, email = ?, enable = ? WHERE id = ?").
		WithArgs(user.Username, user.Name, user.Password, user.Email, user.Enable, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rowsAffected := suite.service.Update(ctx, userID, user)

	assert.Equal(suite.T(), int64(1), rowsAffected)
}
func (suite *UserServiceTestSuite) TestDelete() {
	ctx := context.Background()
	userID := uuid.New()

	suite.mock.ExpectExec("DELETE FROM tb_user WHERE id = ?").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rowsAffected := suite.service.Delete(ctx, userID)

	assert.Equal(suite.T(), int64(1), rowsAffected)
}

func (suite *UserServiceTestSuite) TestGetExistUserName() {
	ctx := context.Background()
	username := "user1"
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(1)

	suite.mock.ExpectQuery("SELECT COUNT(*) FROM tb_user WHERE username = ?").
		WithArgs(username).
		WillReturnRows(rows)

	exists, err := suite.service.GetExistUserName(ctx, username)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *UserServiceTestSuite) TestAuthenticate() {
	username := "user1"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "enable", "hashed_password", "created_at", "updated_at"}).
		AddRow(uuid.New(), username, "User One", "user1@example.com", true, string(hashedPassword), time.Now(), time.Now())

	suite.mock.ExpectQuery("SELECT id, username, name, email, enable, hashed_password, created_at, updated_at FROM tb_user WHERE username = ?").
		WithArgs(username).
		WillReturnRows(rows)

	user, err := suite.service.Authenticate(username, password)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), username, user.Username)
}

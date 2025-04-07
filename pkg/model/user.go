package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/katana-stuidio/access-control/internal/config/logger"
	"github.com/potatowski/brazilcode"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	Password       string    `json:"password"`
	HashedPassword string    `json:"hashed_password"`
	Email          string    `json:"email"`
	Enable         bool      `json:"enable"`
	Role           string    `bson:"role" json:"role"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

type UserList struct {
	List []*User `json:"list"`
}

func (u *User) passwordToHash() {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
		if err != nil {
			log.Println("Erro to SetPassWord", err.Error())
		}

		u.HashedPassword = string(hashedPassword)
	}
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err != nil {
		log.Println("Erro to CheckPassword", err.Error())
		return false
	}
	return true
}

func (u *User) CheckCpf(cpf string) bool {
	err := brazilcode.CPFIsValid(cpf)
	if err != nil {
		logger.Error("Erro to CheckCpf", err)
		return false
	}

	return true
}

func (u *User) PrepareToSave() {
	dt := time.Now()

	u.passwordToHash()
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
		u.CreatedAt = dt
		u.UpdatedAt = dt
	} else {
		u.UpdatedAt = dt
	}
}

func NewUser(user_request *User) (*User, error) {
	user := &User{
		ID:             uuid.New(),
		Username:       user_request.Username,
		Name:           user_request.Name,
		Password:       user_request.Password,
		HashedPassword: user_request.HashedPassword,
		Email:          user_request.Username + "@mpla.com.br",
		Enable:         true,
		Role:           "driver",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return user, nil
}

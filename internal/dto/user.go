package dto

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRequestDtoInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	CNPJ     string `json:"cnpj"`
	Email    string `json:"email"`
}

type UserRequestDtoOutPut struct {
	ID        uuid.UUID `json:"id"`
	CNPJ      string    `json:"cnpj,omitempty"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Enable    bool      `json:"enable"`
	Role      string    `bson:"role" json:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UsertoOutPut struct {
	CNPJ string `json:"cnpj"`
	CPF  string `json:"cpf"`
}

type UserChangePasswordOutPut struct {
	Username    string `json:"username"`
	NewPassowrd string `json:"new_password"`
	OldPassowrd string `json:"old_password"`
}

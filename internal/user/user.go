package user

import (
	"time"
	"github.com/google/uuid"
)


type User struct {
	ID           string
	Email        string
	PasswordHash string
	IsVerified	 bool
	CreatedAt    time.Time
}

func NewUser(email string, passwordHash string) User{
	return User{
		ID:uuid.NewString(),
		Email: email,
		PasswordHash: passwordHash,
		IsVerified: false,
		CreatedAt: time.Now().UTC(),
	}
}
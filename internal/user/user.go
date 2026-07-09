package user

import (
	"time"
	"github.com/google/uuid"
)


type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func NewUser(email string, passwordHash string) User{
	return User{
		ID:uuid.NewString(),
		Email: email,
		PasswordHash: passwordHash,
		CreatedAt: time.Now().UTC(),
	}
}
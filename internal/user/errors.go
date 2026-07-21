package user

import (
	"errors"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrEmailAlreadyExists = errors.New("email already exists")
var dberr = errors.New("database error")
var userExistsConflict = errors.New("email already exists")
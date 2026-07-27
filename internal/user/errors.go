package user

import (
	"errors"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrEmailAlreadyExists = errors.New("email already exists")
var Dberr = errors.New("database error")
var UserExistsConflict = errors.New("email already exists")
var UserDoesNotExistErr = errors.New("no user exists with email")
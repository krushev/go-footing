package util

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidRole        = errors.New("invalid role/s")
	ErrEmailAlreadyExists = errors.New("email %s already exists")
	ErrEmptyPassword      = errors.New("password can't be empty")
	ErrInvalidAuthToken   = errors.New("invalid auth-token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	SomethingGetWrong     = errors.New("something get wrong")
)

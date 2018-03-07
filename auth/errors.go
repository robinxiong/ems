package auth

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidAccount = errors.New("invalid account")
	ErrUnauthorized = errors.New("Unauthorized")
)
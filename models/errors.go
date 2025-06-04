package models

import "errors"

var (
	ErrUserExists        = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrMasterKeyNotFound = errors.New("master key not found")
)

package service

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNotFound          = errors.New("not found")
	ErrValid             = errors.New("invalid data")
)

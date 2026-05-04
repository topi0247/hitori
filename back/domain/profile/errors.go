package profile

import "errors"

var (
	ErrInvalidUserName = errors.New("invalid_user_name")
	ErrNotFound        = errors.New("not_found")
	ErrAlreadyExists   = errors.New("already_exists")
)

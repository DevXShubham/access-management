package apperrors

import "errors"

var (
	ErrTenantNotFound      = errors.New("tenant not found")
	ErrTenantAlreadyExists = errors.New("tenant already exists")
	ErrInvalidInput        = errors.New("invalid input")
	ErrDB                  = errors.New("database error")
	ErrParseTime           = errors.New("time parse error")
)

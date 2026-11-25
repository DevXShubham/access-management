package apperrors

import "errors"

var (
	ErrTenantNotFound      = errors.New("tenant not found")
	ErrTenantAlreadyExists = errors.New("tenant already exists")
	ErrInvalidInput        = errors.New("invalid input")
	ErrDBExecution         = errors.New("database query execution error")
	ErrDBConnection        = errors.New("database connection error")
	ErrParseTime           = errors.New("time parse error")
)

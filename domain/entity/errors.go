package entity

import (
	"errors"
)
//Errors of http
var (
	// ErrBadRequest is http.StatusBadRequest
	ErrBadRequest = errors.New("bad request")
	// ErrInternalServerError is http.StatusInternalServerError
	ErrInternalServerError = errors.New("internal server error")
)

//Errors of user
var (
	// ErrUserNotFound user not found error
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidUser invalid user request error(private)
	ErrInvalidUser = errors.New("invalid user")
)

//Errors of gorm
var (
	// ErrRecordNotFound record not found error
	ErrRecordNotFound = errors.New("record not found")
	// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = errors.New("no valid transaction")
	// ErrNotImplemented not implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrMissingWhereClause missing where clause
	ErrMissingWhereClause = errors.New("WHERE conditions required")
	// ErrUnsupportedRelation unsupported relations
	ErrUnsupportedRelation = errors.New("unsupported relations")
	// ErrPrimaryKeyRequired primary keys required
	ErrPrimaryKeyRequired = errors.New("primary key required")
	// ErrModelValueRequired model value required
	ErrModelValueRequired = errors.New("model value required")
	// ErrInvalidData unsupported data
	ErrInvalidData = errors.New("unsupported data")
	// ErrUnsupportedDriver unsupported driver
	ErrUnsupportedDriver = errors.New("unsupported driver")
	// ErrRegistered registered
	ErrRegistered = errors.New("registered")
	// ErrInvalidField invalid field
	ErrInvalidField = errors.New("invalid field")
	// ErrEmptySlice empty slice found
	ErrEmptySlice = errors.New("empty slice found")
	// ErrDryRunModeUnsupported dry run mode unsupported
	ErrDryRunModeUnsupported = errors.New("dry run mode unsupported")
)

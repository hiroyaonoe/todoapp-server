package entity

import (
	"errors"
	"github.com/jinzhu/gorm"
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

//Errors of jinzhu/gorm
var (
	// ErrRecordNotFound returns a "record not found error". Occurs only when attempting to query the database with a struct; querying with a slice won't return this error
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrInvalidSQL occurs when you attempt a query with invalid SQL
	ErrInvalidSQL = gorm.ErrInvalidSQL
	// ErrInvalidTransaction occurs when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
	// ErrCantStartTransaction can't start transaction when you are trying to start one with `Begin`
	ErrCantStartTransaction = gorm.ErrCantStartTransaction
	// ErrUnaddressable unaddressable value
	ErrUnaddressable = gorm.ErrUnaddressable
)

// //Errors of go-gorm/gorm
// var (
// 	// ErrRecordNotFound record not found error
// 	ErrRecordNotFound = gorm.ErrRecordNotFound
// 	// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
// 	ErrInvalidTransaction = gorm.ErrInvalidTransaction
// 	// ErrNotImplemented not implemented
// 	ErrNotImplemented = gorm.ErrNotImplemented
// 	// ErrMissingWhereClause missing where clause
// 	ErrMissingWhereClause = gorm.ErrMissingWhereClause
// 	// ErrUnsupportedRelation unsupported relations
// 	ErrUnsupportedRelation = gorm.ErrUnsupportedRelation
// 	// ErrPrimaryKeyRequired primary keys required
// 	ErrPrimaryKeyRequired = gorm.ErrPrimaryKeyRequired
// 	// ErrModelValueRequired model value required
// 	ErrModelValueRequired = gorm.ErrModelValueRequired
// 	// ErrInvalidData unsupported data
// 	ErrInvalidData = gorm.ErrInvalidData
// 	// ErrUnsupportedDriver unsupported driver
// 	ErrUnsupportedDriver = gorm.ErrUnsupportedDriver
// 	// ErrRegistered registered
// 	ErrRegistered = gorm.ErrRegistered
// 	// ErrInvalidField invalid field
// 	ErrInvalidField = gorm.ErrInvalidField
// 	// ErrEmptySlice empty slice found
// 	ErrEmptySlice = gorm.ErrEmptySlice
// 	// ErrDryRunModeUnsupported dry run mode unsupported
// 	ErrDryRunModeUnsupported = gorm.ErrDryRunModeUnsupported
// )

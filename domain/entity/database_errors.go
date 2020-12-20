package entity

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
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

// Errors of go-sql-driver/mysql. Various errors the driver might return. Can change between driver versions.
var (
	ErrInvalidConn       = mysql.ErrInvalidConn
	ErrMalformPkt        = mysql.ErrMalformPkt
	ErrNoTLS             = mysql.ErrNoTLS
	ErrCleartextPassword = mysql.ErrCleartextPassword
	ErrNativePassword    = mysql.ErrNativePassword
	ErrOldPassword       = mysql.ErrOldPassword
	ErrUnknownPlugin     = mysql.ErrUnknownPlugin
	ErrOldProtocol       = mysql.ErrOldProtocol
	ErrPktSync           = mysql.ErrPktSync
	ErrPktSyncMul        = mysql.ErrPktSyncMul
	ErrPktTooLarge       = mysql.ErrPktTooLarge
	ErrBusyBuffer        = mysql.ErrBusyBuffer

	// errBadConnNoWrite is used for connection errors where nothing was sent to the database yet.
	// If this happens first in a function starting a database interaction, it should be replaced by driver.ErrBadConn
	// to trigger a resend.
	// See https://github.com/go-sql-driver/mysql/pull/302
	errBadConnNoWrite = errors.New("bad connection")
)

// func NewError(err error) error {
// 	if nerr,ok := err.(*mysql.MySQLError); ok{
// 		return ErrMySQL(*nerr)
// 	} else {
// 		return err
// 	}
// }

// type ErrMySQL mysql.MySQLError

// MySQLError is an error type which represents a single MySQL error
func NewErrMySQL(num uint16, str string) (err *mysql.MySQLError) {
	err = &mysql.MySQLError{
		Number:  num,
		Message: str,
	}
	return
}

// func (me ErrMySQL) Error() string {
// 	return fmt.Sprintf("ErrMySQL %d: %s", me.Number, me.Message)
// }

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

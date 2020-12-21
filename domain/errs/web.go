package errs

import "errors"

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
	// ErrDuplicatedEmail email already exists error
	ErrDuplicatedEmail = errors.New("email already exists")
	// ErrInvalidUser invalid user request error(private)
	ErrInvalidUser = errors.New("invalid user")
)

//Errors of task
var (
	// ErrTaskNotFound task not found error
	ErrTaskNotFound = errors.New("task not found")
	// ErrInvalidTask invalid task request error(private)
	ErrInvalidTask = errors.New("invalid task")
)

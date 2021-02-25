package usecase

import "errors"

//Errors of user
var (
	// ErrInvalidUser invalid user request error
	ErrInvalidUser = errors.New("invalid user")
)

//Errors of task
var (
	// ErrInvalidTask invalid task request error
	ErrInvalidTask = errors.New("invalid task")
)

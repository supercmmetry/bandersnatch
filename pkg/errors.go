package pkg

import "errors"

var (
	ErrNilNode = errors.New("the current state is non-existent")
	ErrAlreadyExists = errors.New("the resource already exists")
	ErrDatabase = errors.New("error occured in database")
	ErrNotFound = errors.New("record not found")
	ErrInvalidOperation = errors.New("invalid operation")
	ErrNoPathFound = errors.New("no possible path found")
)

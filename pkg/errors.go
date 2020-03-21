package pkg

import "errors"

var (
	ErrNilChild = errors.New("the next state was not loaded")
	ErrAlreadyExists = errors.New("the resource already exists")
)

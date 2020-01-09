package pkg

import "errors"

var (
	ErrNilChild = errors.New("The next state was not loaded.")
)

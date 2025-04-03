package grpcerror

import "errors"

var (
	ErrBookExists   = errors.New("book is exist")
	ErrBookNotFound = errors.New("book not found")
)

package domain

import "errors"


var (
	ErrBookExist = errors.New("book exists")
	ErrBookNotFound = errors.New("book not found")
	ErrServerInternal = errors.New("server error")
)
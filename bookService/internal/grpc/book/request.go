package book

import (
	"github.com/go-playground/validator/v10"
)


type BookRequest struct {
	Id int64 `validate:"required"`
	Title string `validate:"required"`
	Author string `validate:"required"`
	CategoryId *int64 `validate:"required"`
}




func (b BookRequest) ValidateUpdateReq() error {
	validate := validator.New()
	return validate.Struct(b)
}
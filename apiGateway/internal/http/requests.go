package http

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CreateBookReq struct {
	Title string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	CategoryId int `json:"category_id" validate:"required"`
}

type BookIdRequest struct {
	Id int `json:"id" validate:"required"`
}

type BookRequest struct {
	Id int `json:"id" validate:"required"`
	Title string `json:"title,omitempty" validate:"required"`
	Author string `json:"author,omitempty"`
	CategoryId int `json:"category_id,omitempty"`
}


func validateStruct(s interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(s)

	if err != nil {
		return fmt.Errorf("validate error: %w", err)
	}
	
	return nil
}


func (b CreateBookReq) Validate() error {
	return validateStruct(b)
}


func (b BookRequest) Validate() error {
	return validateStruct(b)
}

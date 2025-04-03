package pgrepo

import (
	"github.com/serj213/bookService/internal/domain"
	"github.com/serj213/bookService/internal/repository/model"
)


func bookToDomain(book model.BookModel) domain.Book {
	return domain.Book{
		Id: book.Id,
		Title: book.Title,
		Author: *book.Author,
		CategoryId: *book.CategoryId,
	}
}
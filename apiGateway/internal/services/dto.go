package services

import (
	bsv1 "github.com/serj213/bookService/pb/grpc"
	"github.com/serj213/bookServiceApi/internal/domain"
)


func bookToDomain(book *bsv1.BookResponse) domain.Book {

	domainBook := domain.Book{
		ID: int(book.Id),
		Title: book.Title,
		Author: book.Author,
		CategoryId: int(book.CategoryId),
		CreatedAt: book.CreatedAt.AsTime(),
	}

	if book.UpdatedAt != nil {
		updatedAt := book.UpdatedAt.AsTime()
		domainBook.UpdatedAt = &updatedAt
	}

	return domainBook

}
package http

import "github.com/serj213/bookServiceApi/internal/domain"


func reqBookToDomain(req BookRequest) domain.Book{
	return domain.Book{
		ID: req.Id,
		Title: req.Title,
		Author: req.Author,
		CategoryId: req.CategoryId,
	}
}
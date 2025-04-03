package domain

import "time"

type Book struct {
	Id         int64
	Title      string
	Author     string
	CategoryId int
	CreateAt time.Time
	UpdatedAt *time.Time
}

func NewBookDomain(id int, title string, author string, categoryId int) Book {
	return Book{
		Id:         int64(id),
		Title:      title,
		Author:     author,
		CategoryId: categoryId,
	}
}

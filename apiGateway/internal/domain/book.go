package domain

import "time"


type Book struct {
	ID int
	Title string
	Author string
	CategoryId int
	CreatedAt time.Time
	UpdatedAt *time.Time
}
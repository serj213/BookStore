package model

import "time"

type BookModel struct {
	Id         int64
	Title      string
	Author     *string
	CategoryId *int
	UpdateAt *time.Time
	CreateAt time.Time
}
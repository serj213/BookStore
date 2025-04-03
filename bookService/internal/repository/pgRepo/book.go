package pgrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/serj213/bookService/internal/domain"
	"github.com/serj213/bookService/internal/repository"
	"github.com/serj213/bookService/internal/repository/model"
	grpcerror "github.com/serj213/bookService/pkg/grpcError"
	"github.com/serj213/bookService/pkg/pg"
)

type bookRepo struct {
	db *pg.PgDb
}

func NewBookRepo(db *pg.PgDb) bookRepo {
	return bookRepo{
		db: db,
	}
}

func (r bookRepo) Create(ctx context.Context, title string, author string, category_id int) (domain.Book, error) {
	var book domain.Book
	// insert new book
	const queryInsert = "INSERT INTO books(title, author, categories_id) SELECT $1,$2,$3 WHERE EXISTS (SELECT 1 FROM categories WHERE id = $3) RETURNING id, created_at"
	err := r.db.QueryRow(ctx, queryInsert, title, author, category_id).Scan(&book.Id, &book.CreateAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == repository.PgCodeDublicate {
				return domain.Book{}, grpcerror.ErrBookExists
			}
		}
		return domain.Book{}, fmt.Errorf("failed create book: %w", err)
	}

	book.Title = title
	book.Author = author
	book.CategoryId = category_id

	return book, nil
}

func (r bookRepo) Delete(ctx context.Context, id int) error {
	const query = "DELETE FROM books WHERE id=$1"
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return grpcerror.ErrBookNotFound
	}
	return nil
}

func (r bookRepo) Update(ctx context.Context, book domain.Book) (domain.Book, error){

	update_at := time.Now()

	const query = `UPDATE books SET 
		title = COALESCE($2, title), author = COALESCE($3, author), 
		categories_id = COALESCE($4, categories_id), updated_at = $5
		WHERE id=$1
		RETURNING id, title, author, categories_id, updated_at, created_at
		`

	var domainBook domain.Book

	err := r.db.Pool.QueryRow(ctx, query, book.Id, book.Title, book.Author, book.CategoryId, update_at).
		Scan(&domainBook.Id, &domainBook.Title, &domainBook.Author, &domainBook.CategoryId, &domainBook.CreateAt, &domainBook.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Book{}, grpcerror.ErrBookNotFound
		}
		return domain.Book{}, fmt.Errorf("failed update book: %w", err)
	}

	
	return domainBook, nil
}

func (r bookRepo) GetBookById(ctx context.Context, id int) (domain.Book, error) {
	var bookModel model.BookModel
	const query = `SELECT id, title, author, categories_id, created_at, updated_at FROM books WHERE id = $1`
	
	err := r.db.QueryRow(ctx, query, id).
		Scan(&bookModel.Id, &bookModel.Title, &bookModel.Author, &bookModel.CategoryId, &bookModel.CreateAt, &bookModel.UpdateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Book{}, grpcerror.ErrBookNotFound
		}
		return domain.Book{}, fmt.Errorf("failed get book by id: %w", err)
	}

	book := bookToDomain(bookModel)

	return book, nil
}

func (r bookRepo) GetBooks(ctx context.Context)([]domain.Book, error) {
	const query = `SELECT id, title, author, categories_id, created_at, updated_at FROM books`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []domain.Book{}, err
	}

	var books []domain.Book

	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.CategoryId, &book.CreateAt, &book.UpdatedAt)
		if err != nil {
			return []domain.Book{}, err
		}
		books = append(books, book)
	}
	return books, nil
}
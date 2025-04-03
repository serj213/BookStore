package book

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/serj213/bookService/internal/domain"
	grpcerror "github.com/serj213/bookService/pkg/grpcError"
)

type BookRepository interface {
	Create(ctx context.Context, title string, author string, category_id int) (domain.Book, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, book domain.Book) (domain.Book, error)
	GetBookById(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context)([]domain.Book, error)
}

type BookService struct {
	log       *slog.Logger
	repo      BookRepository
	trManager *manager.Manager
}

func NewBookService(log *slog.Logger, repo BookRepository, trManager *manager.Manager) BookService {
	return BookService{
		log:       log,
		repo:      repo,
		trManager: trManager,
	}
}

func (s BookService) Create(ctx context.Context, title string, author string, categoryId int64) (domain.Book, error) {
	log := s.log.With(slog.String("method", "Create"))
	log.Info("create book enabled")

	book, err := s.repo.Create(ctx, title, author, int(categoryId))

	if err != nil {
		log.Error("failed create book: %w", err)
		return domain.Book{}, err
	}

	log.Info("create book finish")
	return book, nil
}

func (s BookService) Delete(ctx context.Context, id int) error {
	log := s.log.With(slog.With("op", "service.book.Delete"))
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Error("failed delete: %w", err)
		if errors.Is(err, grpcerror.ErrBookNotFound) {
			return grpcerror.ErrBookNotFound
		}
		return fmt.Errorf("failed delete book")
	}
	return nil
}

func (s BookService) GetBookById(ctx context.Context, id int) (domain.Book, error) {

	log := s.log.With(slog.String("method", "GetBookById"), slog.Int("id", id))

	log.Info("GetBookById started")

	book, err := s.repo.GetBookById(ctx, id)
	if err != nil {
		log.Error("failed get book by id: ", err.Error())
		return domain.Book{}, err
	}

	return book, nil
}

func (s BookService) GetAllBooks(ctx context.Context) ([]domain.Book, error) {

	log := s.log.With(slog.String("method", "GetAllBooks"))

	log.Info("get books active")

	books, err := s.repo.GetBooks(ctx)
	if err != nil {
		log.Error("failed get books repo: %w", err)
		return []domain.Book{}, fmt.Errorf("failed get books")
	}

	return books, nil
}

func (s BookService) Update(ctx context.Context, book domain.Book) (domain.Book, error) {
	log := s.log.With(slog.String("method", "Update"))

	log.Info("book ", book)

	book, err := s.repo.Update(ctx, book)
	if err != nil {
		log.Error("failed update book: %w", err)
		return domain.Book{}, err
	}
	return book, nil
}

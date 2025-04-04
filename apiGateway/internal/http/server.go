package http

import (
	"context"

	"github.com/serj213/bookServiceApi/internal/domain"
	"github.com/serj213/bookServiceApi/internal/kafka"
	"go.uber.org/zap"
)

type BookService interface {
	Create(ctx context.Context, title string, author string, categoryId int)(domain.Book, error)
	GetBooks(ctx context.Context) ([]domain.Book, error)
	UpdateBook(ctx context.Context, book domain.Book) (domain.Book, error)
	GetBookById(ctx context.Context, id int)(domain.Book, error)
	DeleteBook(ctx context.Context, id int) error
}


type HTTPServer struct {
	log *zap.SugaredLogger
	BookService BookService
	Kafka *kafka.Kafka
}


func New(log *zap.SugaredLogger, bookService BookService, kafka *kafka.Kafka) *HTTPServer {
	return &HTTPServer{
		log: log,
		BookService: bookService,
		Kafka: kafka,
	}
}
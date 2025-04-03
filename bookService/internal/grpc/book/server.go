package book

import (
	"context"
	"errors"

	"github.com/serj213/bookService/internal/domain"
	bsv1 "github.com/serj213/bookService/pb/grpc"
	grpcerror "github.com/serj213/bookService/pkg/grpcError"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type serverApi struct {
	bsv1.UnimplementedBookServer
	book Book
}

type Book interface {
	Create(ctx context.Context, title string, author string, category_id int64) (domain.Book, error)
	Delete(ctx context.Context, id int) error
	GetBookById(ctx context.Context, id int) (domain.Book, error)
	GetAllBooks(ctx context.Context) ([]domain.Book, error)
	Update(ctx context.Context, book domain.Book) (domain.Book, error)
}

func RegisterGrpc(server *grpc.Server, book Book) {
	bsv1.RegisterBookServer(server, &serverApi{book: book})
}

func (s serverApi) Create(ctx context.Context, in *bsv1.BookCreateRequest) (*bsv1.BookResponse, error) {

	book, err := s.book.Create(ctx, in.GetTitle(), in.GetAuthor(), in.GetCategoryId())
	if err != nil {
		if errors.Is(err, grpcerror.ErrBookExists) {
			return nil, status.Error(codes.Internal, "book is exist")
		}
		return nil, status.Error(codes.Internal, "failed create book")
	}

	return &bsv1.BookResponse{
		Id:         int64(book.Id),
		Title:      book.Title,
		Author:     book.Author,
		CategoryId: int64(book.CategoryId),
		UpdatedAt: nil,
		CreatedAt: timestamppb.New(book.CreateAt),
	}, nil
}

func (s serverApi) Delete(ctx context.Context, in *bsv1.BookDeleteRequest) (*emptypb.Empty, error) {
    err := s.book.Delete(ctx, int(in.GetId()))
	if err != nil {
		if errors.Is(err, grpcerror.ErrBookNotFound){
			return nil, status.Error(codes.NotFound, "book not found")
		}
		return nil, status.Error(codes.Internal, "server error")
	}

    return &emptypb.Empty{}, nil
}

func (s serverApi) GetBookById(ctx context.Context, in *bsv1.BookGetBookByIdRequest) (*bsv1.BookResponse, error) {
	book, err := s.book.GetBookById(ctx, int(in.GetId()))
	if err != nil {
		if errors.Is(err, grpcerror.ErrBookNotFound){
			return nil, status.Error(codes.NotFound, "book not found")
		}
		return nil, status.Error(codes.Internal, "server error")
	}

	return &bsv1.BookResponse{
		Id:         book.Id,
		Title:      book.Title,
		Author:     book.Author,
		CategoryId: int64(book.CategoryId),
	}, nil

}

func (s serverApi) GetBooks(ctx context.Context, in *emptypb.Empty) (*bsv1.BookListResponse, error) {
	books, err := s.book.GetAllBooks(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resBooks := make([]*bsv1.BookResponse, len(books))

	for i, val := range books {
		resBooks[i] = domainToGrpcFormat(val)
	}

	return &bsv1.BookListResponse{
		Books: resBooks,
	}, nil
}

// Валидацию добавить
func (s serverApi) UpdateBook(ctx context.Context, in *bsv1.BookRequest) (*bsv1.BookResponse, error) {
	domainBook := domain.NewBookDomain(
		int(in.GetId()), 
		in.GetTitle(),
		in.GetAuthor(), 
		int(in.CategoryId), 
	)

	book, err := s.book.Update(ctx, domainBook)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed update book")
	}

	return &bsv1.BookResponse{
		Id:         int64(book.Id),
		Title:      book.Title,
		Author:     book.Author,
		CategoryId: int64(book.CategoryId),
		UpdatedAt: timestamppb.New(*book.UpdatedAt),
		CreatedAt: timestamppb.New(book.CreateAt),
	}, nil
}

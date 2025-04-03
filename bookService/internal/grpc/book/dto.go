package book

import (
	"github.com/serj213/bookService/internal/domain"
	bsv1 "github.com/serj213/bookService/pb/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)
	

func domainToGrpcFormat(book domain.Book) *bsv1.BookResponse{

	proto := &bsv1.BookResponse{
		Id: book.Id,
		Title: book.Title,
		Author: book.Author,
		CategoryId: int64(book.CategoryId),
		CreatedAt: timestamppb.New(book.CreateAt),
	}

	if book.UpdatedAt != nil {
		proto.UpdatedAt = timestamppb.New(*book.UpdatedAt)
	}

	return proto
}
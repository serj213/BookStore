package lib

import (
	"errors"

	grpcerror "github.com/serj213/bookService/pkg/grpcError"
	"google.golang.org/grpc/codes"
)


func GetResCode(err error) codes.Code {
	switch{
	case errors.Is(err, grpcerror.ErrBookNotFound):
			return codes.NotFound	
	default:
		return codes.Internal
	}
}
package lib

import (
	"errors"
	"net/http"

	"github.com/serj213/bookServiceApi/internal/domain"
)


func GetStatusError(err error) int {
	switch {
	case errors.Is(err, domain.ErrBookNotFound):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
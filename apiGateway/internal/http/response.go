package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/serj213/bookServiceApi/internal/domain"
)

const (
	StatusFailed = "failed"
)

type ResponseErr struct {
	Status string `json:"status"`
	Msg string `json:"msg"`
}

type BookResponse struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	CategoryId int `json:"category_id"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseOkBody struct {
	Status string `json:"status"`
	Data domain.Book `json:"data,omitempty"`
}

type GetBooksResponseOk struct {
	Status string `json:"status"`
	Books []BookResponse`json:"books"`
}


func ErrResponse(msg string, w http.ResponseWriter, r *http.Request, status int) {
	resp := ResponseErr{
		Status: StatusFailed,
		Msg: msg,
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(&resp)
}

func ResponseOk(data interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
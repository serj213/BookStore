package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/serj213/bookServiceApi/internal/domain"
	"github.com/serj213/bookServiceApi/internal/lib"
)

func (h HTTPServer) Create(w http.ResponseWriter, r *http.Request) {
	var bookReq CreateBookReq

	if err := json.NewDecoder(r.Body).Decode(&bookReq); err != nil {
		ErrResponse("invalid request", w, r, http.StatusBadRequest)
		return
	}

	if err := bookReq.Validate(); err != nil {
		ErrResponse("invalid request", w, r, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	book, err := h.BookService.Create(r.Context(), bookReq.Title, bookReq.Author, bookReq.CategoryId)
	if err != nil {
		h.log.Errorf("failed create serivce: %w", err)
		if errors.Is(err, domain.ErrBookExist) {
			ErrResponse("book is exists", w, r, http.StatusInternalServerError)
			return
		}
		ErrResponse("server error", w, r, http.StatusInternalServerError)
		return
	}

	ResponseOk(BookResponse{
		Id:         book.ID,
		Title:      book.Title,
		Author:     book.Author,
		CategoryId: book.CategoryId,
		CreatedAt:  book.CreatedAt,
	}, w)
}

func (h HTTPServer) GetBooks(w http.ResponseWriter, r *http.Request) {

	books, err := h.BookService.GetBooks(r.Context())
	if err != nil {
		ErrResponse("internal", w, r, http.StatusInternalServerError)
		return
	}

	resBooks := make([]BookResponse, len(books))

	for i, book := range books {

		if book.UpdatedAt != nil {
			resBooks[i].UpdatedAt = book.UpdatedAt
		}

		resBooks[i] = BookResponse{
			Id:         book.ID,
			Title:      book.Title,
			Author:     book.Author,
			CategoryId: book.CategoryId,
			CreatedAt:  book.CreatedAt,
		}
	}

	resOk := GetBooksResponseOk{
		Status: "success",
		Books:  resBooks,
	}

	ResponseOk(resOk, w)

}

func (h HTTPServer) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var req BookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrResponse("invalid request", w, r, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		ErrResponse("invalid request", w, r, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	bookDomain := reqBookToDomain(req)

	updateBook, err := h.BookService.UpdateBook(r.Context(), bookDomain)
	if err != nil {
		ErrResponse("server error", w, r, http.StatusInternalServerError)
		return
	}

	respData := BookResponse{
		Id:         updateBook.ID,
		Title:      updateBook.Title,
		Author:     updateBook.Author,
		CategoryId: updateBook.CategoryId,
		CreatedAt:  updateBook.CreatedAt,
		UpdatedAt:  updateBook.UpdatedAt,
	}

	ResponseOk(respData, w)
}

func (s HTTPServer) GetBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ErrResponse("book id is empty", w, r, http.StatusBadRequest)
		return
	}

	bookId, err := strconv.Atoi(id)
	if err != nil {
		ErrResponse("invalid request", w, r, http.StatusBadRequest)
		return
	}

	book, err := s.BookService.GetBookById(r.Context(), bookId)
	if err != nil {
		ErrResponse(err.Error(), w, r, lib.GetStatusError(err))
		return
	}

	data := ResponseOkBody{
		Status: "Ok",
		Data:   book,
	}

	ResponseOk(data, w)
}

func (s HTTPServer) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ErrResponse("book id is empty", w, r, http.StatusBadRequest)
		return
	}

	bookId, err := strconv.Atoi(id)
	if err != nil {
		ErrResponse("invalid request", w, r, http.StatusBadRequest)
		return
	}

	err = s.BookService.DeleteBook(r.Context(), bookId)
	if err != nil {
		ErrResponse(err.Error(), w, r, lib.GetStatusError(err))
		return
	}

	data := ResponseOkBody{
		Status: "Success",
	}

	ResponseOk(data, w)
}

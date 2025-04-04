package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	bsv1 "github.com/serj213/bookService/pb/grpc"
	_ "github.com/serj213/bookServiceApi/docs"
	"github.com/serj213/bookServiceApi/internal/config"
	HTTPServer "github.com/serj213/bookServiceApi/internal/http"
	"github.com/serj213/bookServiceApi/internal/services"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	local = "local"
	dev = "develop"
)

// @title Book Store API
// @version 1.0
// @description This is a sample server for a book store.
// @host localhost:8083
// @BasePath /api/v1

func main(){
	cfg, err := config.Deal()
	if err != nil {
		panic(err)
	}

	log := setupLogger(cfg.Env)
	logSugar := log.Sugar()
	logSugar = logSugar.With(zap.String("env", cfg.Env))

	logSugar.Info("logger is enabled")

	// инициализация grpc клиента


	conn, err := grpc.Dial(fmt.Sprintf("book-service:%d", cfg.GRPC.Port), 
		grpc.WithTransportCredentials(insecure.NewCredentials() ),
	)
	if err != nil {
		logSugar.Infof("failed start grpc client: %w", err)
		panic(err)
	}
	
	logSugar.Info("grpc client started...")

	defer conn.Close()

	bookClient := bsv1.NewBookClient(conn)

	bookService := services.New(logSugar, bookClient)

	httpServer := HTTPServer.New(logSugar, bookService)

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	router.HandleFunc("/book/create", httpServer.Create).Methods(http.MethodPost)
	router.HandleFunc("/books", httpServer.GetBooks).Methods(http.MethodGet)
	router.HandleFunc("/book/update", httpServer.UpdateBook).Methods(http.MethodPut)
	router.HandleFunc("/book/{id}", httpServer.GetBook).Methods(http.MethodGet)
	router.HandleFunc("/book/delete/{id}", httpServer.DeleteBook).Methods(http.MethodDelete)

	addr, _ := os.LookupEnv("HTTP_ADDR")

	srv := &http.Server{
		Handler: router,
		Addr: addr,
	}

	logSugar.Infof("http server started: %d...", addr)

	if err := srv.ListenAndServe(); err != nil {
		logSugar.Infof("failed http server: %w", err)
		panic(err)
	}
}



func setupLogger(env string) *zap.Logger{
	var log *zap.Logger

	switch(env){
	case local:
		log = zap.Must(zap.NewDevelopment())
	case dev:
		log = zap.Must(zap.NewProduction())
	default:
		log = zap.Must(zap.NewProduction())
	}
	return log
}
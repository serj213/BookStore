package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpcBook "github.com/serj213/bookService/internal/grpc/book"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	bookService grpcBook.Book,
	port int,
) *App {

	recoveryOpt := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("recovery from panic", slog.Any("panic", p))

			return status.Error(codes.Internal, "internal error")
		}),
	}

	loggingOpt := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived,
			logging.PayloadSent,
		),
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpt...),
		logging.UnaryServerInterceptor(loggingInterceptor(log), loggingOpt...),
	))

	grpcBook.RegisterGrpc(grpcServer, bookService)

	return &App{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func loggingInterceptor(l *slog.Logger) logging.Logger {

	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Создаём listener, который будет слушить TCP-сообщения, адресованные нашему gRPC-серверу
func (a *App) Run() error {

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", a.port))
	if err != nil {
		return fmt.Errorf("failed run grpcServer: %w", err)
	}

	a.log.Info("grpc server started...")

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("failed Serve grpc: %w", err)
	}

	return nil
}

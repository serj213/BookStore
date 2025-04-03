package lib

import "google.golang.org/grpc/status"

func GetDescGrpcErr(err error) string {
	if grpcErr, ok := status.FromError(err); ok {
		return grpcErr.Message()
	}
	return ""
}
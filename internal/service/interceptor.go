package service

import (
	"context"

	// userpb "github.com/fibonachyy/sternx/internal/api"

	"github.com/fibonachyy/sternx/pkg/logger"
	"google.golang.org/grpc"
)

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// Perform JWT authentication before calling the actual gRPC method
	// if err := authenticate(ctx); err != nil {
	// 	return nil, err
	// }

	// Perform data validation before calling the actual gRPC method
	// if req, ok := req.(*userpb.UserRequest); ok {
	// 	if err := validateUserRequest(req); err != nil {
	// 		return nil, err
	// 	}
	// }

	return logger.GrpcLogger(ctx, req, info, handler)
}

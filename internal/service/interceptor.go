package service

import (
	"context"
	"log"

	"github.com/fibonachyy/sternx/internal/api/user"
	"google.golang.org/grpc"
)

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Unary interceptor: %s called with request: %+v", info.FullMethod, req)

	// Perform JWT authentication before calling the actual gRPC method
	if err := authenticate(ctx); err != nil {
		return nil, err
	}

	// Perform data validation before calling the actual gRPC method
	if req, ok := req.(*user.UserRequest); ok {
		if err := validateUserRequest(req); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

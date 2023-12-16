package service

import (
	"context"

	"github.com/fibonachyy/sternx/internal/api/user"
	"google.golang.org/grpc"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	GrpcServer *grpc.Server
}

func NewUserServiceServer(grpcServer *grpc.Server) *UserServiceServer {

	server := &UserServiceServer{GrpcServer: grpcServer}
	user.RegisterUserServiceServer(grpcServer, server)
	return server
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	// Implement the logic to create a user
	// You can use your repository methods here
	// For now, let's return a placeholder response
	res := &user.UserResponse{
		UserId: "123",
		Name:   req.Name,
		Email:  req.Email,
	}
	return res, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.UserResponse, error) {
	// Implement the logic to get a user by ID
	// You can use your repository methods here
	// For now, let's return a placeholder response
	return &user.UserResponse{
		UserId: "123",
		Name:   "John Doe",
		Email:  "john@example.com",
	}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	// Implement the logic to update a user by ID
	// You can use your repository methods here
	// For now, let's return a placeholder response
	return &user.UserResponse{
		UserId: "123",
		Name:   req.Name,
		Email:  req.Email,
	}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.UserResponse, error) {
	// Implement the logic to delete a user by ID
	// You can use your repository methods here
	// For now, let's return a placeholder response
	return &user.UserResponse{
		UserId: "123",
		Name:   "John Doe",
		Email:  "john@example.com",
	}, nil
}

package service

import (
	"context"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/repository"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	UserRepo repository.IRepository
}

func NewUserServiceServer(repo repository.IRepository) *UserServiceServer {
	return &UserServiceServer{UserRepo: repo}
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	// Implement the logic to get a user by ID
	// You can use your repository methods here
	// For now, let's return a placeholder response
	return &userpb.UserResponse{
		User: &userpb.User{UserId: "123", Name: "John", Email: "a@gmail.com"},
	}, nil
}

package service

import (
	"context"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	// Implement the logic to update a user by ID
	// You can use your repository methods here
	// For now, let's return a placeholder response

	res := &userpb.UserResponse{
		User: &userpb.User{UserId: "123", Name: "john", Email: "a@GMAIL.COM", PasswordChangedAt: &timestamppb.Timestamp{}, CreatedAt: &timestamppb.Timestamp{}},
	}
	return res, nil
}

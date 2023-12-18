package service

import (
	"context"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.UpdateUserResponse, error) {
	authPayload, err := s.authorizeUser(ctx, []string{domain.AdminRole, domain.StandardRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	violations := validateDeleteUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if req.GetEmail() != authPayload.Email && authPayload.Role != domain.AdminRole {
		return nil, status.Errorf(codes.PermissionDenied, "cannot delete other user")
	}
	err = s.UserRepo.DeleteUserByEmail(ctx, authPayload.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return &userpb.UpdateUserResponse{
		Success: true,
	}, nil
}
func validateDeleteUserRequest(req *userpb.DeleteUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := domain.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}

package service

import (
	"context"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	authPayload, err := s.authorizeUser(ctx, []string{domain.AdminRole, domain.StandardRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Email != req.GetEmail() || authPayload.Role == domain.AdminRole {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}
	updatedUser, err := s.UserRepo.PartialUpdateUserByEmail(ctx, req.GetEmail(), domain.User{Name: req.GetName()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update user info")
	}
	res := ConvertToUserResponse(*updatedUser)
	return res, nil
}

func validateUpdateUserRequest(req *userpb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := domain.ValidateName(req.GetName()); err != nil {
		violations = append(violations, fieldViolation("name", err))
	}

	if err := domain.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}

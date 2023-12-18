package service

import (
	"context"
	"strconv"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {

	authPayload, err := s.authorizeUser(ctx, []string{domain.AdminRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	_ = authPayload
	violations := validateGetUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	id, _ := strconv.Atoi(req.GetUserId()) // its check in validation befor

	userData, err := s.UserRepo.GetUserByID(ctx, id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	return ConvertToUserResponse(*userData), nil
}

func validateGetUserRequest(req *userpb.GetUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := domain.ValidateUserIdString(req.GetUserId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}
	return violations
}

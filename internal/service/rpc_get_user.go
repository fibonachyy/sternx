package service

import (
	"context"
	"strconv"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/domain"
	"github.com/fibonachyy/sternx/internal/logger"
	"github.com/fibonachyy/sternx/pkg/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	log := logger.FromContext(ctx)

	authPayload, err := s.authorizeUser(ctx, []string{domain.AdminRole})
	if err != nil {
		log.Errorf(ctx, "Authorization failed for GetUser request: %v", err)
		return nil, unauthenticatedError(err)
	}
	_ = authPayload

	violations := validateGetUserRequest(req)
	if violations != nil {
		log.Error(ctx, "Validation failed for GetUser request", "violations", violations)
		return nil, invalidArgumentError(violations)
	}

	id, _ := strconv.Atoi(req.GetUserId()) // It's checked in validation before

	userData, err := s.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		log.Errorf(ctx, "Failed to find user by ID: %d, error: %v", id, err)
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	// Log user retrieval without sensitive details
	log.Infof(ctx, "User retrieved successfully: ID=%d, Email=%s, Role=%s", userData.ID, utils.MaskEmail(userData.Email), userData.Role)

	return ConvertToUserResponse(*userData), nil
}

func validateGetUserRequest(req *userpb.GetUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := domain.ValidateUserIdString(req.GetUserId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}
	return violations
}

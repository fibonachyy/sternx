package service

import (
	"context"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/domain"
	"github.com/fibonachyy/sternx/internal/logger"
	"github.com/fibonachyy/sternx/pkg/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	log := logger.FromContext(ctx)

	authPayload, err := s.authorizeUser(ctx, []string{domain.AdminRole, domain.StandardRole})
	if err != nil {
		log.Errorf(ctx, "Authorization failed for UpdateUser request: %v", err)
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		log.Errorf(ctx, "Validation failed for UpdateUser request: %v", violations)
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Email != req.GetEmail() || authPayload.Role == domain.AdminRole {
		log.Warn(ctx, "Permission denied: cannot update other user's info")
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	updatedUser, err := s.UserRepo.PartialUpdateUserByEmail(ctx, req.GetEmail(), domain.User{Name: req.GetName()})
	if err != nil {
		log.Errorf(ctx, "Failed to update user info: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update user info")
	}

	res := ConvertToUserResponse(*updatedUser)

	// Log successful user update without sensitive details
	log.Infof(ctx, "User info updated successfully: ID=%d, Email=%s, Role=%s", updatedUser.ID, utils.MaskEmail(updatedUser.Email), updatedUser.Role)

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

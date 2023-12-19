package service

import (
	"context"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/domain"
	"github.com/fibonachyy/sternx/internal/logger"
	"github.com/fibonachyy/sternx/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	log := logger.FromContext(ctx)

	tracer := otel.Tracer("grpc-server")
	ctx, span := tracer.Start(ctx, "UserService/CreateUser") // Use a standardized name
	defer span.End()

	span.SetAttributes(
		attribute.String("service.method.name", "updateUser"),
		attribute.String("user.email", req.GetEmail()),
	)
	ctx = trace.ContextWithSpan(ctx, span)

	authPayload, err := s.authorizeUser(ctx, []string{domain.AdminRole, domain.StandardRole})
	if err != nil {
		log.Errorf(ctx, "Authorization failed for UpdateUser request: %v", err)
		span.RecordError(err)
		return nil, unauthenticatedError(err)
	}
	span.SetAttributes(
		attribute.String("Applicant.email", authPayload.Email),
		attribute.String("Applicant.role", authPayload.Role),
	)

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		log.Errorf(ctx, "Validation failed for UpdateUser request: %v", violations)
		span.SetAttributes(domain.ConvertFieldViolationsToAttributes(violations)...)
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Email != req.GetEmail() || authPayload.Role == domain.AdminRole {
		log.Warn(ctx, "Permission denied: cannot update other user's info")
		err = status.Errorf(codes.PermissionDenied, "cannot update other user's info")
		span.RecordError(err)
		return nil, err
	}

	updatedUser, err := s.UserRepo.PartialUpdateUserByEmail(ctx, req.GetEmail(), domain.User{Name: req.GetName()})
	if err != nil {
		log.Errorf(ctx, "Failed to update user info: %v", err)
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to update user info")
	}
	span.SetAttributes(
		attribute.String("user.prevName", req.GetName()),
		attribute.String("user.newName", updatedUser.Name),
	)

	res := ConvertToUserResponse(*updatedUser)

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

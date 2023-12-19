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

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.UpdateUserResponse, error) {
	log := logger.FromContext(ctx)

	tracer := otel.Tracer("grpc-server")
	ctx, span := tracer.Start(ctx, "UserService/CreateUser") // Use a standardized name
	defer span.End()

	span.SetAttributes(
		attribute.String("service.method.name", "DeleteUser"),
		attribute.String("user.email", req.GetEmail()),
	)
	ctx = trace.ContextWithSpan(ctx, span)

	authPayload, err := s.authorizeUser(ctx, []string{domain.AdminRole, domain.StandardRole})
	if err != nil {
		log.Errorf(ctx, "Authorization failed for DeleteUser request: %v", err)
		span.RecordError(err)
		return nil, unauthenticatedError(err)
	}

	span.SetAttributes(
		attribute.String("Applicant.email", authPayload.Email),
		attribute.String("Applicant.role", authPayload.Role),
	)

	violations := validateDeleteUserRequest(req)
	if violations != nil {
		log.Error(ctx, "Validation failed for DeleteUser request", "violations", violations)
		span.SetAttributes(domain.ConvertFieldViolationsToAttributes(violations)...)
		return nil, invalidArgumentError(violations)
	}

	if req.GetEmail() != authPayload.Email && authPayload.Role != domain.AdminRole {
		log.Errorf(ctx, "Permission denied for deleting user with email: %s", req.GetEmail())
		span.RecordError(err)
		return nil, status.Errorf(codes.PermissionDenied, "cannot delete other user")
	}

	err = s.UserRepo.DeleteUserByEmail(ctx, authPayload.Email)
	if err != nil {
		log.Errorf(ctx, "Failed to delete user with email %s: %v", authPayload.Email, err)
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	// Log user deletion without sensitive details
	log.Infof(ctx, "User deleted successfully: Email=%s", utils.MaskEmail(authPayload.Email))

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

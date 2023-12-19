package service

import (
	"context"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/domain"
	"github.com/fibonachyy/sternx/internal/logger"
	"github.com/fibonachyy/sternx/internal/metrics"
	"github.com/fibonachyy/sternx/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *UserServiceServer) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.LoginUserResponse, error) {
	log := logger.FromContext(ctx)
	meter := metrics.FromContext(ctx)

	tracer := otel.Tracer("grpc-server")
	ctx, span := tracer.Start(ctx, "UserService/LoginUser") // Use a standardized name
	defer span.End()

	ctx = trace.ContextWithSpan(ctx, span)

	span.SetAttributes(
		attribute.String("service.method.name", "login"),
		attribute.String("user.email", req.GetEmail()),
	)

	violations := validateLoginUserRequest(req)
	if violations != nil {
		log.Error(ctx, "Validation failed for LoginUser request", "violations", violations)
		span.SetAttributes(domain.ConvertFieldViolationsToAttributes(violations)...)
		return nil, invalidArgumentError(violations)
	}

	user, err := server.UserRepo.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		log.Errorf(ctx, "Failed to find user by email: %s, error: %v", utils.MaskEmail(req.GetEmail()), err)
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}
	span.SetAttributes(
		attribute.String("user.role", user.Role),
	)

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		log.Errorf(ctx, "Incorrect password for user: %s", utils.MaskEmail(user.Email))
		span.RecordError(err)
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Email,
		user.Role,
		server.Config.JWTDuration,
	)
	if err != nil {
		log.Errorf(ctx, "Failed to create access token for user: %s, error: %v", utils.MaskEmail(user.Email), err)
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	rsp := &userpb.LoginUserResponse{
		User:                 ConvertToUserResponse(*user).User,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
	}
	loginCounter, _ := meter.Int64Counter("login")
	loginCounter.Add(ctx, 1)
	log.Infof(ctx, "User login successful: ID=%d, Email=%s, Role=%s", user.ID, utils.MaskEmail(user.Email), user.Role)

	return rsp, nil
}

func validateLoginUserRequest(req *userpb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := domain.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := domain.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}

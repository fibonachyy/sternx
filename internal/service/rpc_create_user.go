package service

import (
	"context"
	"fmt"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/domain"
	"github.com/fibonachyy/sternx/internal/logger"
	"github.com/fibonachyy/sternx/internal/repository"
	"github.com/fibonachyy/sternx/pkg/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	log := logger.FromContext(ctx)

	violations := validateCreateUserRequest(req)
	if violations != nil {
		log.Error(ctx, "Validation failed for CreateUser request", "violations", violations)
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		log.Errorf(ctx, "Failed to hash password for user creation: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	userParam := repository.CreateUserParams{
		Name:           req.GetName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
		Role:           domain.StandardRole,
	}

	user, err := s.UserRepo.CreateUser(ctx, userParam)
	if err != nil {
		log.Errorf(ctx, "Failed to create user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	log.Infof(ctx, "User created successfully: ID=%d, Email=%s, Role=%s", user.ID, utils.MaskEmail(user.Email), user.Role)

	return ConvertToUserResponse(*user), nil
}

func (s *UserServiceServer) CreateAdmin(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	log := logger.FromContext(ctx)

	authPayload, err := s.authorizeUser(ctx, []string{domain.AdminRole})
	if err != nil {
		log.Errorf(ctx, "Authorization failed for CreateAdmin request: %v", err)
		return nil, unauthenticatedError(err)
	}

	violations := validateCreateUserRequest(req)
	if violations != nil {
		log.Error(ctx, "Validation failed for CreateAdmin request", "violations", violations)
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		log.Errorf(ctx, "Failed to hash password for admin user creation: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	userParam := repository.CreateUserParams{
		Name:           req.GetName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
		Role:           authPayload.Role,
	}

	user, err := s.UserRepo.CreateUser(ctx, userParam)
	if err != nil {
		log.Errorf(ctx, "Failed to create admin user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create admin user: %v", err)
	}

	log.Infof(ctx, "Admin user created successfully: ID=%d, Email=%s, Role=%s", user.ID, utils.MaskEmail(user.Email), user.Role)

	return ConvertToUserResponse(*user), nil
}

func ConvertToUserResponse(user domain.User) *userpb.UserResponse {
	return &userpb.UserResponse{
		User: &userpb.User{
			UserId:            fmt.Sprint(user.ID),
			Name:              user.Name,
			Email:             user.Email,
			Role:              domain.StringToRole(user.Role),
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}
}

func validateCreateUserRequest(req *userpb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := domain.ValidateName(req.GetName()); err != nil {
		violations = append(violations, fieldViolation("name", err))
	}

	if err := domain.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := domain.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if req.GetRole() != userpb.Role_ADMIN && req.GetRole() != userpb.Role_STANDARD {
		violations = append(violations, fieldViolation("role", fmt.Errorf("invalid role")))
	}
	return violations
}

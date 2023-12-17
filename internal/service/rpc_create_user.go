package service

import (
	"context"
	"fmt"

	"github.com/fibonachyy/sternx/domain"
	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/repository"
	"github.com/fibonachyy/sternx/internal/validator"
	"github.com/fibonachyy/sternx/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	log.Info().Msg("new call")
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	userParam := repository.CreateUserParams{
		Name:           req.GetName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := s.UserRepo.CreateUser(ctx, userParam)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return ConvertToUserResponse(*user), nil
}

func ConvertToUserResponse(user domain.User) *userpb.UserResponse {
	return &userpb.UserResponse{
		User: &userpb.User{
			UserId:            fmt.Sprint(user.ID),
			Name:              user.Name,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}
}

func validateCreateUserRequest(req *userpb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateName(req.GetName()); err != nil {
		violations = append(violations, fieldViolation("name", err))
	}

	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}

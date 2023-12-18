package service

import (
	"context"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/domain"
	"github.com/fibonachyy/sternx/pkg/utils"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *UserServiceServer) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.UserRepo.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Email,
		user.Role,
		server.Config.JWTDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	rsp := &userpb.LoginUserResponse{
		User:                 ConvertToUserResponse(*user).User,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
	}
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

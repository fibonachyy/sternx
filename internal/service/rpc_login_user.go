package service

import (
	"context"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/validator"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *UserServiceServer) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.LoginUserResponse, error) {
	return &userpb.LoginUserResponse{}, nil
}

// {
// 	violations := validateLoginUserRequest(req)
// 	if violations != nil {
// 		return nil, invalidArgumentError(violations)
// 	}

// 	user, err := server.store.GetUser(ctx, req.GetUsername())
// 	if err != nil {
// 		if errors.Is(err, db.ErrRecordNotFound) {
// 			return nil, status.Errorf(codes.NotFound, "user not found")
// 		}
// 		return nil, status.Errorf(codes.Internal, "failed to find user")
// 	}

// 	err = util.CheckPassword(req.Password, user.HashedPassword)
// 	if err != nil {
// 		return nil, status.Errorf(codes.NotFound, "incorrect password")
// 	}

// 	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
// 		user.Username,
// 		user.Role,
// 		server.config.AccessTokenDuration,
// 	)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create access token")
// 	}

// 	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
// 		user.Username,
// 		user.Role,
// 		server.config.RefreshTokenDuration,
// 	)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
// 	}

// 	mtdt := server.extractMetadata(ctx)
// 	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
// 		ID:           refreshPayload.ID,
// 		Username:     user.Username,
// 		RefreshToken: refreshToken,
// 		UserAgent:    mtdt.UserAgent,
// 		ClientIp:     mtdt.ClientIP,
// 		IsBlocked:    false,
// 		ExpiresAt:    refreshPayload.ExpiredAt,
// 	})
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create session")
// 	}

// 	rsp := &pb.LoginUserResponse{
// 		User:                  convertUser(user),
// 		SessionId:             session.ID.String(),
// 		AccessToken:           accessToken,
// 		RefreshToken:          refreshToken,
// 		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
// 		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
// 	}
// 	return rsp, nil
// }

func validateLoginUserRequest(req *userpb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}

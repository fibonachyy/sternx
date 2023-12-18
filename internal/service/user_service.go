package service

import (
	"fmt"

	"github.com/fibonachyy/sternx/pkg/token"

	userpb "github.com/fibonachyy/sternx/internal/api"
	"github.com/fibonachyy/sternx/internal/repository"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	UserRepo repository.IRepository

	Config     Config
	tokenMaker token.Maker
}

func NewUserServiceServer(repo repository.IRepository, config Config) (*UserServiceServer, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	tokenMaker, err := createTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	if config.JWTDuration == 0 {
		defaultConfig := DefaultConfig()
		config.JWTDuration = defaultConfig.JWTDuration
	}

	return &UserServiceServer{UserRepo: repo, tokenMaker: tokenMaker, Config: config}, nil
}

func createTokenMaker(symmetricKey string) (token.Maker, error) {
	return token.NewPasetoMaker(symmetricKey)
}

package repository

import (
	"context"

	"github.com/fibonachyy/sternx/domain"
)

type IRepository interface {
	IMigrateTable
	IUserRepository
}
type IMigrateTable interface {
	Migrate(path string) error
}
type IUserRepository interface {
	CreateUser(ctx context.Context, params CreateUserParams) (*domain.User, error)
	FindUserByID(ctx context.Context, userID string) (*domain.User, error)
	PartialUpdateUser(ctx context.Context, userID string, updatedUser domain.User) (*domain.User, error)
	DeleteUserByID(ctx context.Context, userID string) error
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)
}

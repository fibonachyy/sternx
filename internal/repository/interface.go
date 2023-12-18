package repository

import (
	"context"

	"github.com/fibonachyy/sternx/internal/domain"
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
	GetUserByEmail(ctx context.Context, userEmail string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID int) (*domain.User, error)
	PartialUpdateUserByEmail(ctx context.Context, email string, updatedUser domain.User) (*domain.User, error)
	DeleteUserByEmail(ctx context.Context, email string) error
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)
}

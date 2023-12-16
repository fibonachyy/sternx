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
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
}

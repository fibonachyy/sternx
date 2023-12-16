// repository/user_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fibonachyy/sternx/domain"
)

type UserModel struct {
	id    string
	name  string
	email string
}

func (u UserModel) ToDomain() *domain.User {
	return &domain.User{
		ID:    u.id,
		Name:  u.name,
		Email: u.email,
	}
}

func (p *postgres) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := "SELECT id, name, email FROM users WHERE id = $1"
	row := p.conn.QueryRow(ctx, query, id)

	var user UserModel
	err := row.Scan(&user.id, &user.name, &user.email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return user.ToDomain(), nil
}

func (p *postgres) Save(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)"
	_, err := p.conn.Exec(ctx, query, user.ID, user.Name, user.Email)
	return err
}

// Additional methods for querying and manipulating user data can be added.

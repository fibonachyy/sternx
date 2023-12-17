// repository/user_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fibonachyy/sternx/domain"
	"golang.org/x/crypto/bcrypt"
)

type userModel struct {
	id                int
	name              string
	email             string
	hashedPassword    string
	passwordChangedAt time.Time
	createdAt         time.Time
}

func (u userModel) ToDomain() *domain.User {
	return &domain.User{
		ID:                u.id,
		Name:              u.name,
		Email:             u.email,
		HashedPassword:    u.hashedPassword,
		PasswordChangedAt: u.passwordChangedAt,
		CreatedAt:         u.createdAt,
	}
}

type CreateUserParams struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (p *postgres) CreateUser(ctx context.Context, params CreateUserParams) (*domain.User, error) {
	createdAt := time.Now()
	passwordChangedAt := time.Now()

	// Insert query
	insertQuery := "INSERT INTO users (name, email, hashed_password, password_changed_at, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var userID int
	err := p.conn.QueryRow(ctx, insertQuery, params.Name, params.Email, params.HashedPassword, passwordChangedAt, createdAt).Scan(&userID)
	fmt.Println(userID)
	if err != nil {
		p.logger.LogError(ctx, err, "failed to insert user into database", "email", params.Email)
		return nil, fmt.Errorf("failed to insert user into database")
	}

	user := &domain.User{
		ID:                userID,
		Name:              params.Name,
		Email:             params.Email,
		HashedPassword:    params.HashedPassword,
		PasswordChangedAt: passwordChangedAt,
		CreatedAt:         createdAt,
	}

	p.logger.LogInfo(ctx, "user created successfully: ID=%d, Email=%s", user.ID, user.Email)

	return user, nil
}
func (p *postgres) FindUserByID(ctx context.Context, userID string) (*domain.User, error) {
	query := "SELECT id, name, email, hashed_password, password_changed_at, created_at FROM users WHERE id = $1"
	var user userModel

	err := p.conn.QueryRow(ctx, query, userID).Scan(&user.id, &user.name, &user.email, &user.hashedPassword, &user.passwordChangedAt, &user.createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a specific error when the user is not found
			return nil, fmt.Errorf("user with ID %s not found: %w", userID, err)
		}
		// Handle other database errors
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return user.ToDomain(), nil
}

func (p *postgres) PartialUpdateUser(ctx context.Context, userID string, updatedUser domain.User) (*domain.User, error) {
	updatedUser.Email = ""
	updatedUser.PasswordChangedAt = time.Time{}
	updatedUser.CreatedAt = time.Time{}
	updatedUser.HashedPassword = ""

	query := "UPDATE users SET name = $1 WHERE id = $2 RETURNING id, password_changed_at, created_at, hashed_password"

	var newUserID int
	var newPasswordChangedAt, newCreatedAt time.Time
	var newHashedPassword string
	err := p.conn.QueryRow(ctx, query, updatedUser.Name, userID).Scan(&newUserID, &newPasswordChangedAt, &newCreatedAt, &newHashedPassword)
	if err != nil {

		return nil, fmt.Errorf("failed to partially update user by ID %s: %w", userID, err)
	}

	updatedUser.ID = newUserID
	updatedUser.PasswordChangedAt = newPasswordChangedAt
	updatedUser.CreatedAt = newCreatedAt

	return &updatedUser, nil
}

func (p *postgres) DeleteUserByID(ctx context.Context, userID string) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := p.conn.Exec(ctx, query, userID)
	if err != nil {
		// Handle the database deletion error
		return fmt.Errorf("failed to delete user by ID %s: %w", userID, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		// Return a specific error when no rows were affected (user not found)
		return fmt.Errorf("user with ID %s not found: %w", userID, sql.ErrNoRows)
	}

	return nil
}
func (p *postgres) AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error) {
	query := "SELECT id, name, email, hashed_password, password_changed_at, created_at FROM users WHERE email = $1"
	var user userModel

	err := p.conn.QueryRow(ctx, query, email).Scan(&user.id, &user.name, &user.email, &user.hashedPassword, &user.passwordChangedAt, &user.createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a specific error when the user is not found
			return nil, fmt.Errorf("user with email %s not found: %w", email, err)
		}
		// Handle other database errors
		return nil, fmt.Errorf("failed to retrieve user by email: %w", err)
	}

	// Compare hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.hashedPassword), []byte(password))
	if err != nil {
		// Passwords do not match
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Passwords match, user is authenticated
	return user.ToDomain(), nil
}

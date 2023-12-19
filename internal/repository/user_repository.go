// repository/user_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fibonachyy/sternx/internal/domain"
	"github.com/fibonachyy/sternx/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

type userModel struct {
	id                int
	name              string
	email             string
	hashedPassword    string
	role              string
	passwordChangedAt time.Time
	createdAt         time.Time
}

func (u userModel) ToDomain() *domain.User {
	return &domain.User{
		ID:                u.id,
		Name:              u.name,
		Email:             u.email,
		Role:              u.role,
		HashedPassword:    u.hashedPassword,
		PasswordChangedAt: u.passwordChangedAt,
		CreatedAt:         u.createdAt,
	}
}

type CreateUserParams struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	HashedPassword string `json:"hashed_password"`
}

func (p *postgres) CreateUser(ctx context.Context, params CreateUserParams) (*domain.User, error) {
	logFromCtx := logger.FromContext(ctx)

	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("repository.method.name", "CreateUser"),
		attribute.String("user.email", params.Email),
	)

	createdAt := time.Now()
	passwordChangedAt := time.Now()

	insertQuery := "INSERT INTO users (name, email, role, hashed_password, password_changed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	var userID int
	err := p.conn.QueryRow(ctx, insertQuery, params.Name, params.Email, params.Role, params.HashedPassword, passwordChangedAt, createdAt).Scan(&userID)
	if err != nil {
		logFromCtx.Errorf(ctx, "failed to insert user into database: %v", err)
		span.RecordError(err)
		return nil, fmt.Errorf("failed to insert user into database: %w", err)
	}

	user := &domain.User{
		ID:                userID,
		Name:              params.Name,
		Email:             params.Email,
		Role:              params.Role,
		HashedPassword:    params.HashedPassword,
		PasswordChangedAt: passwordChangedAt,
		CreatedAt:         createdAt,
	}

	logFromCtx.Infof(ctx, "user created successfully: ID=%d, Email=%s, Role=%s", user.ID, user.Email, user.Role)

	return user, nil
}

func (p *postgres) GetUserByEmail(ctx context.Context, userEmail string) (*domain.User, error) {
	logFromCtx := logger.FromContext(ctx)

	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "GetUserByEmail")
	defer span.End()

	span.SetAttributes(
		attribute.String("repository.method.name", "GetUserByEmail"),
		attribute.String("user.email", userEmail),
	)

	query := "SELECT id, name, email, role, hashed_password, password_changed_at, created_at FROM users WHERE email = $1"
	var user userModel

	err := p.conn.QueryRow(ctx, query, userEmail).Scan(
		&user.id, &user.name, &user.email, &user.role, &user.hashedPassword, &user.passwordChangedAt, &user.createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logFromCtx.Errorf(ctx, "user not found with the provided email: %v", err)
			span.RecordError(err)

			return nil, fmt.Errorf("user not found with the provided email: %w", err)
		}
		logFromCtx.Errorf(ctx, "failed to find user by email: %v", err)
		span.RecordError(err)

		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user.ToDomain(), nil
}

func (p *postgres) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	logFromCtx := logger.FromContext(ctx)

	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "GetUserByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("repository.method.name", "GetUserByID"),
		attribute.Int("user.id", userID),
	)

	query := "SELECT id, name, email, role, hashed_password, password_changed_at, created_at FROM users WHERE id = $1"
	var user userModel

	err := p.conn.QueryRow(ctx, query, userID).Scan(
		&user.id, &user.name, &user.email, &user.role, &user.hashedPassword, &user.passwordChangedAt, &user.createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logFromCtx.Errorf(ctx, "user not found with the provided ID: %v", err)
			span.RecordError(err)

			return nil, fmt.Errorf("user not found with the provided ID: %w", err)
		}
		logFromCtx.Errorf(ctx, "failed to find user by id: %v", err)
		span.RecordError(err)

		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user.ToDomain(), nil
}

func (p *postgres) PartialUpdateUserByEmail(ctx context.Context, email string, updatedUser domain.User) (*domain.User, error) {
	logFromCtx := logger.FromContext(ctx)

	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "PartialUpdateUserByEmail")
	defer span.End()

	span.SetAttributes(
		attribute.String("repository.method.name", "PartialUpdateUserByEmail"),
		attribute.String("user.email", email),
	)

	updatedUser.Email = ""
	updatedUser.Role = ""
	updatedUser.PasswordChangedAt = time.Time{}
	updatedUser.CreatedAt = time.Time{}
	updatedUser.HashedPassword = ""

	query := "UPDATE users SET name = $1 WHERE email = $2 RETURNING id, password_changed_at, created_at, hashed_password, email, role"

	var newUserID int
	var newPasswordChangedAt, newCreatedAt time.Time
	var newHashedPassword, userEmail, role string
	err := p.conn.QueryRow(ctx, query, updatedUser.Name, email).Scan(&newUserID, &newPasswordChangedAt, &newCreatedAt, &newHashedPassword, &userEmail, &role)
	if err != nil {
		logFromCtx.Errorf(ctx, "failed to update user info by email: %s: %v", email, err)
		span.RecordError(err)
		return nil, fmt.Errorf("failed to partially update user by Email %s: %w", email, err)
	}

	updatedUser.ID = newUserID
	updatedUser.PasswordChangedAt = newPasswordChangedAt
	updatedUser.CreatedAt = newCreatedAt
	updatedUser.Email = userEmail
	updatedUser.Role = role
	return &updatedUser, nil
}
func (p *postgres) DeleteUserByEmail(ctx context.Context, email string) error {

	logFromCtx := logger.FromContext(ctx)

	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "DeleteUserByEmail")
	defer span.End()

	// Add attributes specific to the repository method
	span.SetAttributes(
		attribute.String("repository.method.name", "DeleteUserByEmail"),
		attribute.String("user.email", email),
	)

	query := "DELETE FROM users WHERE email = $1"

	result, err := p.conn.Exec(ctx, query, email)
	if err != nil {
		logFromCtx.Errorf(ctx, "failed to delete user by email: %s: %v", email, err)
		span.RecordError(err)

		return fmt.Errorf("failed to delete user by Email %s: %w", email, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		logFromCtx.Errorf(ctx, "user with Email %s not found: %v", email, sql.ErrNoRows)
		span.SetAttributes(attribute.Bool("user.deleted", false))

		return fmt.Errorf("user with Email %s not found: %w", email, sql.ErrNoRows)
	}
	span.SetAttributes(attribute.Bool("user.deleted", true))

	return nil
}

func (p *postgres) AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error) {
	logFromCtx := logger.FromContext(ctx)

	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "AuthenticateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("repository.method.name", "AuthenticateUser"),
		attribute.String("user.email", email),
	)

	query := "SELECT id, name, email, role, hashed_password, password_changed_at, created_at FROM users WHERE email = $1"
	var user userModel

	err := p.conn.QueryRow(ctx, query, email).Scan(&user.id, &user.name, &user.email, &user.role, &user.hashedPassword, &user.passwordChangedAt, &user.createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logFromCtx.Errorf(ctx, "user with email %s not found: %v", email, err)
			return nil, fmt.Errorf("user with email %s not found: %w", email, err)
		}
		logFromCtx.Errorf(ctx, "failed to retrieve user by email: %v", err)
		return nil, fmt.Errorf("failed to retrieve user by email: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.hashedPassword), []byte(password))
	if err != nil {
		logFromCtx.Errorf(ctx, "authentication failed: %v", err)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	span.SetAttributes(
		attribute.Bool("user.authenticated", true),
		attribute.Int("user.id", user.id),
		attribute.String("user.role", user.role),
	)

	logFromCtx.Infof(ctx, "user authenticated successfully: ID=%d, Email=%s, Role=%s", user.id, user.email, user.role)

	return user.ToDomain(), nil
}

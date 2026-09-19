package repository

import (
	"context"
	"errors"
	"fmt"
	"sitewatch/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository provides methods for working with users.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user in the database.
func (ur *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users(email,pass_hash)
					VALUES($1,$2)
					RETURNING id,created_at`
	err := ur.db.QueryRow(ctx, query, user.Email, user.PassHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user:%w", err)
	}
	return nil
}

// GetByEmail returns a user by email.
func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id,email,pass_hash,created_at
					FROM users
					WHERE email=$1`
	err := ur.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PassHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// GetByID returns a user by id.
func (ur *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id,email,pass_hash,created_at
					FROM users
					WHERE id=$1`
	err := ur.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PassHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

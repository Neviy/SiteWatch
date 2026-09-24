package service

import (
	"context"
	"errors"
	"sitewatch/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserProfile contains public user profile information.
type UserProfile struct {
	ID        int64
	Email     string
	CreatedAt time.Time
}

// UserRepository defines database operations for users.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

// UserService contains business logic for users.
type UserService struct {
	repo UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Register creates a new user.
func (us *UserService) Register(ctx context.Context, email, pass string) error {
	if len(pass) < 8 {
		return errors.New("the password must be at least 8 characters long")
	}
	example, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
		} else {
			return err
		}
	}
	if example != nil {
		return errors.New("this email address is already taken")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user, err := model.NewUser(email, string(hash))
	if err != nil {
		return err
	}
	if err := us.repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

// Login authenticates a user.
func (us *UserService) Login(ctx context.Context, email, pass string) (string, error) {
	user, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(pass)); err != nil {
		return "", errors.New("invalid email or password")
	}
	//JWT
	return "", nil
}

// GetProfile returns public profile information for a user.
func (us *UserService) GetProfile(ctx context.Context, userID int64) (*UserProfile, error) {
	user, err := us.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &UserProfile{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

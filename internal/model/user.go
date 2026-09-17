// Package model contains domain models used by the application.
package model

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

// User represents an application user.
type User struct {
	ID        int64
	Email     string
	PassHash  string
	CreatedAt time.Time
}

// NewUser creates and returns a new User.
func NewUser(email, hash string) (*User, error) {
	var errs []error
	if _, err := mail.ParseAddress(email); err != nil {
		errs = append(errs, errors.New("invalid email"))
	}
	if strings.TrimSpace(hash) == "" {
		errs = append(errs, errors.New("password hash cannot be empty"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return &User{
		Email:    email,
		PassHash: hash,
	}, nil
}

// Package model contains domain models used by the application.
package model

import (
	"errors"
	"time"
)

// Check represents a single website check result.
type Check struct {
	ID           int64
	SiteID       int64
	StatusCode   int
	ResponseTime int64
	Error        string
	CheckedAt    time.Time
}

// NewCheck creates and returns a new Check.
func NewCheck(siteID int64, statusCode int, responseTime int64, err string) (*Check, error) {
	var errs []error
	if siteID <= 0 {
		errs = append(errs, errors.New("invalid siteID"))
	}
	if statusCode < 0 {
		errs = append(errs, errors.New("invalid status code"))
	}
	if responseTime < 0 {
		errs = append(errs, errors.New("invalid response time"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return &Check{
		SiteID:       siteID,
		StatusCode:   statusCode,
		ResponseTime: responseTime,
		Error:        err,
		CheckedAt:    time.Now(),
	}, nil
}

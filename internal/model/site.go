// Package model contains domain models used by the application.
package model

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

// Site represents a website monitored by the application.
type Site struct {
	ID          int64
	UserID      int64
	URL         string
	Status      bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	LastCheckAt *time.Time
}

// NewSite creates and returns a new Site.
func NewSite(userID int64, siteURL string) (*Site, error) {
	var errs []error
	if userID <= 0 {
		errs = append(errs, errors.New("invalid userID"))
	}
	if !strings.HasPrefix(siteURL, "http://") &&
		!strings.HasPrefix(siteURL, "https://") {
		siteURL = "https://" + siteURL
	}
	parsedURL, err := url.Parse(siteURL)
	if err != nil || parsedURL.Host == "" {
		errs = append(errs, errors.New("invalid URL"))
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		errs = append(errs, errors.New("invalid URL scheme"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return &Site{
		UserID: userID,
		URL:    siteURL,
	}, nil
}

// MarkUpdated updates the site's last update time.
func (s *Site) MarkUpdated() {
	now := time.Now()
	s.UpdatedAt = &now
}

// MarkChecked updates the site's last check time.
func (s *Site) MarkChecked() {
	now := time.Now()
	s.LastCheckAt = &now
}

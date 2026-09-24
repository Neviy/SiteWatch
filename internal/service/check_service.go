package service

import (
	"context"
	"sitewatch/internal/model"
)

// CheckRepository defines database operations for checks.
type CheckRepository interface {
	Create(ctx context.Context, check *model.Check) error
	GetByID(ctx context.Context, id int64) (*model.Check, error)
	GetBySiteID(ctx context.Context, siteID int64) ([]*model.Check, error)
}

// CheckService contains business logic for checks.
type CheckService struct {
	repo CheckRepository
}

// NewCheckService creates a new CheckService.
func NewCheckService(repo CheckRepository) *CheckService {
	return &CheckService{
		repo: repo,
	}
}

// Create creates a new check for a site.
func (cs *CheckService) Create(ctx context.Context, check *model.Check) error {
	return cs.repo.Create(ctx, check)
}

// GetByID returns a check by its ID.
func (cs *CheckService) GetByID(ctx context.Context, id int64) (*model.Check, error) {
	return cs.repo.GetByID(ctx, id)
}

// GetBySiteID returns all checks belonging to a site.
func (cs *CheckService) GetBySiteID(ctx context.Context, siteID int64) ([]*model.Check, error) {
	return cs.repo.GetBySiteID(ctx, siteID)
}

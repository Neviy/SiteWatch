package service

import (
	"context"
	"sitewatch/internal/model"
)

// SiteRepository defines database operations for sites.
type SiteRepository interface {
	Create(ctx context.Context, site *model.Site) error
	GetByID(ctx context.Context, id int64) (*model.Site, error)
	GetByUserID(ctx context.Context, userID int64) ([]*model.Site, error)
	UpdateStatus(ctx context.Context, id int64, status bool) error
	Delete(ctx context.Context, id int64) error
}

// SiteService contains business logic for sites.
type SiteService struct {
	repo SiteRepository
}

// NewSiteService creates a new SiteService.
func NewSiteService(repo SiteRepository) *SiteService {
	return &SiteService{
		repo: repo,
	}
}

// Create creates a new site for a user.
func (ss *SiteService) Create(ctx context.Context, userID int64, siteURL string) (*model.Site, error) {
	site, err := model.NewSite(userID, siteURL)
	if err != nil {
		return nil, err
	}

	if err := ss.repo.Create(ctx, site); err != nil {
		return nil, err
	}

	return site, nil
}

// GetByID returns a site by its ID.
func (ss *SiteService) GetByID(ctx context.Context, id int64) (*model.Site, error) {
	return ss.repo.GetByID(ctx, id)
}

// GetByUserID returns all sites belonging to a user.
func (ss *SiteService) GetByUserID(ctx context.Context, userID int64) ([]*model.Site, error) {
	return ss.repo.GetByUserID(ctx, userID)
}

// UpdateStatus updates the status of a site.
func (ss *SiteService) UpdateStatus(ctx context.Context, id int64, status bool) error {
	return ss.repo.UpdateStatus(ctx, id, status)
}

// Delete deletes a site by its ID.
func (ss *SiteService) Delete(ctx context.Context, id int64) error {
	return ss.repo.Delete(ctx, id)
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"sitewatch/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SiteRepository  provides methods for working with site.
type SiteRepository struct {
	db *pgxpool.Pool
}

// NewSiteRepository creates a new SiteRepository.
func NewSiteRepository(db *pgxpool.Pool) *SiteRepository {
	return &SiteRepository{
		db: db,
	}
}

// Create creates a new site in the database.
func (sr *SiteRepository) Create(ctx context.Context, site *model.Site) error {
	query := `INSERT INTO sites(user_id,url)
					VALUES($1,$2)
					RETURNING id,created_at,updated_at,last_check_at`
	err := sr.db.QueryRow(ctx, query, site.UserID, site.URL).Scan(&site.ID, &site.CreatedAt, &site.UpdatedAt, &site.LastCheckAt)
	if err != nil {
		return fmt.Errorf("failed to create site:%w", err)
	}
	return nil
}

// GetByID returns a site by id.
func (sr *SiteRepository) GetByID(ctx context.Context, id int64) (*model.Site, error) {
	site := &model.Site{}
	query := `SELECT id,user_id,url,status,created_at,updated_at,last_check_at
					FROM sites
					WHERE id=$1`
	err := sr.db.QueryRow(ctx, query, id).Scan(&site.ID, &site.UserID, &site.URL, &site.Status, &site.CreatedAt, &site.UpdatedAt, &site.LastCheckAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("site not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get site by id: %w", err)
	}
	return site, nil
}

// GetByUserID returns all of the  user`s sites.
func (sr *SiteRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Site, error) {
	var sites []*model.Site
	query := `SELECT id,user_id,url,status,created_at,updated_at,last_check_at
					FROM sites
					WHERE user_id=$1`
	rows, err := sr.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get by user id :%w", err)
	}
	defer rows.Close()
	for rows.Next() {
		site := &model.Site{}
		err = rows.Scan(&site.ID, &site.UserID, &site.URL, &site.Status, &site.CreatedAt, &site.UpdatedAt, &site.LastCheckAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan site: %w", err)
		}
		sites = append(sites, site)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sites: %w", err)
	}
	return sites, nil
}

// UpdateStatus updates the status of a site.
func (sr *SiteRepository) UpdateStatus(ctx context.Context, id int64, status bool) error {
	query := `UPDATE sites
					SET status=$1,
					    updated_at=NOW()
					WHERE id=$2`
	result, err := sr.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update site status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("site not found")
	}
	return nil
}

// Delete deletes a site from the database.
func (sr *SiteRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sites
					WHERE id=$1`
	result, err := sr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("site not found")
	}
	return nil
}

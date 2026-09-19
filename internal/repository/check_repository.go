package repository

import (
	"context"
	"errors"
	"fmt"
	"sitewatch/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CheckRepository  provides methods for working with site.
type CheckRepository struct {
	db *pgxpool.Pool
}

// NewCheckRepository creates a new CheckRepository.
func NewCheckRepository(db *pgxpool.Pool) *CheckRepository {
	return &CheckRepository{
		db: db,
	}
}

// Create creates a new check in the database.
func (cr *CheckRepository) Create(ctx context.Context, check *model.Check) error {
	query := `INSERT INTO checks (site_id,status_code,response_time,error,checked_at)
					VALUES($1,$2,$3,$4,$5)
					RETURNING id`
	err := cr.db.QueryRow(ctx, query, check.SiteID, check.StatusCode, check.ResponseTime, check.Error, check.CheckedAt).Scan(&check.ID)
	if err != nil {
		return fmt.Errorf("failed to create check:%w", err)
	}
	return nil
}

// GetBySiteID returns checks by site id.
func (cr *CheckRepository) GetBySiteID(ctx context.Context, siteID int64) ([]*model.Check, error) {
	var checks []*model.Check
	query := `SELECT id,site_id,status_code,response_time,error,checked_at
					FROM checks
					WHERE site_id=$1`
	rows, err := cr.db.Query(ctx, query, siteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get checks by site id: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		check := &model.Check{}
		err = rows.Scan(&check.ID, &check.SiteID, &check.StatusCode, &check.ResponseTime, &check.Error, &check.CheckedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan check: %w", err)
		}
		checks = append(checks, check)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate checks: %w", err)
	}
	return checks, nil
}

// GetByID returns a check by id.
func (cr *CheckRepository) GetByID(ctx context.Context, id int64) (*model.Check, error) {
	check := &model.Check{}
	query := `SELECT id,site_id,status_code,response_time,error,checked_at
					FROM checks
					WHERE id=$1`
	err := cr.db.QueryRow(ctx, query, id).Scan(&check.ID, &check.SiteID, &check.StatusCode, &check.ResponseTime, check.Error, &check.CheckedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("check not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get check by id: %w", err)
	}
	return check, nil
}

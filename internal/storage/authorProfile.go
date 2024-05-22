package storage

import (
	"time"

	"github.com/ruziba3vich/exam/internal/models"
)

func (s *Storage) UpdateBiography(req models.UpdateBiographyRequest) (*models.UpdateBiographyResponse, error) {
	query := `
		UPDATE Authors
		SET biography = $1,
			updated_at = $2
		WHERE id = $3
		RETURNING
			id,
			name,
			password,
			biography,
			created_at,
			updated_at;
	`
	ctx, cancel := getContext(s.ctx)
	defer cancel()
	row := s.db.QueryRowContext(ctx,
		query,
		req.Biography,
		time.Now(),
		req.Id)
	var response models.UpdateBiographyResponse
	if err := row.Scan(
		&response.Id,
		&response.Name,
		&response.Password,
		&response.Biography,
		&response.CreatedAt,
		&response.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Storage) UpdateBirthdate(req models.UpdateBirthdateRequest) (*models.UpdateBirthdateResponse, error) {
	query := `
		UPDATE Authors
		SET biography = $1,
			updated_at = $2
		WHERE id = $3
		RETURNING
			id,
			name,
			password,
			biography,
			created_at,
			updated_at;
	`
	ctx, cancel := getContext(s.ctx)
	defer cancel()
	row := s.db.QueryRowContext(ctx,
		query,
		req.Birthdate,
		time.Now(),
		req.Id)
	var response models.UpdateBirthdateResponse
	if err := row.Scan(
		&response.Id,
		&response.Name,
		&response.Password,
		&response.BirthDate,
		&response.CreatedAt,
		&response.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Storage) GetProfile(req models.GetProfileRequest) (*models.GetProfileResponse, error) {
	query := `
		SELECT
			id,
			name,
			password,
			biography,
			birthdate,
			created_at,
			updated_at
		FROM Authors
		WHERE id = $1;
	`
	ctx, cancel := getContext(s.ctx)
	defer cancel()
	row := s.db.QueryRowContext(ctx, query, req.Id)

	var response models.GetProfileResponse

	if err := row.Scan(
		&response.Id,
		&response.Name,
		&response.Password,
		&response.Biography,
		&response.BirthDate,
		&response.CreatedAt,
		&response.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/natapage/test_assignment/backend/internal/domain"
)

type LocationRepo struct {
	pool *pgxpool.Pool
}

func NewLocationRepo(pool *pgxpool.Pool) *LocationRepo {
	return &LocationRepo{pool: pool}
}

func (r *LocationRepo) List(ctx context.Context) ([]domain.Location, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, address, place_name, latitude, longitude, created_at, updated_at
		FROM locations ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Location
	for rows.Next() {
		var l domain.Location
		if err := rows.Scan(&l.ID, &l.Address, &l.PlaceName, &l.Latitude, &l.Longitude, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, rows.Err()
}

func (r *LocationRepo) GetByID(ctx context.Context, id int64) (*domain.Location, error) {
	var l domain.Location
	err := r.pool.QueryRow(ctx, `
		SELECT id, address, place_name, latitude, longitude, created_at, updated_at
		FROM locations WHERE id = $1`, id).
		Scan(&l.ID, &l.Address, &l.PlaceName, &l.Latitude, &l.Longitude, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LocationRepo) Create(ctx context.Context, l *domain.Location) (*domain.Location, error) {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO locations (address, place_name, latitude, longitude)
		VALUES ($1, $2, $3, $4)
		RETURNING id, address, place_name, latitude, longitude, created_at, updated_at`,
		l.Address, l.PlaceName, l.Latitude, l.Longitude).
		Scan(&l.ID, &l.Address, &l.PlaceName, &l.Latitude, &l.Longitude, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *LocationRepo) Update(ctx context.Context, l *domain.Location) (*domain.Location, error) {
	err := r.pool.QueryRow(ctx, `
		UPDATE locations SET address = $2, place_name = $3, latitude = $4, longitude = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING id, address, place_name, latitude, longitude, created_at, updated_at`,
		l.ID, l.Address, l.PlaceName, l.Latitude, l.Longitude).
		Scan(&l.ID, &l.Address, &l.PlaceName, &l.Latitude, &l.Longitude, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *LocationRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM locations WHERE id = $1`, id)
	return err
}

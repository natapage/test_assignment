package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/natapage/test_assignment/backend/internal/domain"
)

type MachineRepo struct {
	pool *pgxpool.Pool
}

func NewMachineRepo(pool *pgxpool.Pool) *MachineRepo {
	return &MachineRepo{pool: pool}
}

func (r *MachineRepo) List(ctx context.Context) ([]domain.Machine, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT m.id, m.name, m.serial_number, m.enabled, m.location_id,
		       l.id, l.address, l.place_name, l.latitude, l.longitude, l.created_at, l.updated_at,
		       m.created_at, m.updated_at
		FROM machines m
		LEFT JOIN locations l ON m.location_id = l.id
		ORDER BY m.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Machine
	for rows.Next() {
		var m domain.Machine
		var loc locationScan
		if err := rows.Scan(
			&m.ID, &m.Name, &m.SerialNumber, &m.Enabled, &m.LocationID,
			&loc.id, &loc.address, &loc.placeName, &loc.latitude, &loc.longitude, &loc.createdAt, &loc.updatedAt,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if loc.id != nil {
			m.Location = loc.toDomain()
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *MachineRepo) GetByID(ctx context.Context, id int64) (*domain.Machine, error) {
	var m domain.Machine
	var loc locationScan
	err := r.pool.QueryRow(ctx, `
		SELECT m.id, m.name, m.serial_number, m.enabled, m.location_id,
		       l.id, l.address, l.place_name, l.latitude, l.longitude, l.created_at, l.updated_at,
		       m.created_at, m.updated_at
		FROM machines m
		LEFT JOIN locations l ON m.location_id = l.id
		WHERE m.id = $1`, id).
		Scan(
			&m.ID, &m.Name, &m.SerialNumber, &m.Enabled, &m.LocationID,
			&loc.id, &loc.address, &loc.placeName, &loc.latitude, &loc.longitude, &loc.createdAt, &loc.updatedAt,
			&m.CreatedAt, &m.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}
	if loc.id != nil {
		m.Location = loc.toDomain()
	}
	return &m, nil
}

func (r *MachineRepo) Create(ctx context.Context, m *domain.Machine) (*domain.Machine, error) {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO machines (name, serial_number, enabled, location_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		m.Name, m.SerialNumber, m.Enabled, m.LocationID).
		Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, m.ID)
}

func (r *MachineRepo) Update(ctx context.Context, m *domain.Machine) (*domain.Machine, error) {
	_, err := r.pool.Exec(ctx, `
		UPDATE machines SET name = $2, serial_number = $3, enabled = $4, updated_at = NOW()
		WHERE id = $1`,
		m.ID, m.Name, m.SerialNumber, m.Enabled)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, m.ID)
}

func (r *MachineRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM machines WHERE id = $1`, id)
	return err
}

func (r *MachineRepo) UpdateLocationID(ctx context.Context, id int64, locationID *int64) error {
	db := getConn(ctx, r.pool)
	_, err := db.Exec(ctx, `
		UPDATE machines SET location_id = $2, updated_at = NOW()
		WHERE id = $1`, id, locationID)
	return err
}

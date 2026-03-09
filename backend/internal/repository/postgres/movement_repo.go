package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/natapage/test_assignment/backend/internal/domain"
)

type MovementRepo struct {
	pool *pgxpool.Pool
}

func NewMovementRepo(pool *pgxpool.Pool) *MovementRepo {
	return &MovementRepo{pool: pool}
}

func (r *MovementRepo) Create(ctx context.Context, m *domain.Movement) (*domain.Movement, error) {
	db := getConn(ctx, r.pool)
	err := db.QueryRow(ctx, `
		INSERT INTO movement_history (machine_id, from_location_id, to_location_id, moved_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, moved_at, created_at`,
		m.MachineID, m.FromLocationID, m.ToLocationID).
		Scan(&m.ID, &m.MovedAt, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MovementRepo) ListByMachineID(ctx context.Context, machineID int64) ([]domain.Movement, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT mh.id, mh.machine_id, mh.from_location_id, mh.to_location_id, mh.moved_at, mh.created_at,
		       fl.id, fl.address, fl.place_name, fl.latitude, fl.longitude, fl.created_at, fl.updated_at,
		       tl.id, tl.address, tl.place_name, tl.latitude, tl.longitude, tl.created_at, tl.updated_at
		FROM movement_history mh
		LEFT JOIN locations fl ON mh.from_location_id = fl.id
		JOIN locations tl ON mh.to_location_id = tl.id
		WHERE mh.machine_id = $1
		ORDER BY mh.moved_at DESC`, machineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Movement
	for rows.Next() {
		var m domain.Movement
		var fromLoc, toLoc locationScan
		if err := rows.Scan(
			&m.ID, &m.MachineID, &m.FromLocationID, &m.ToLocationID, &m.MovedAt, &m.CreatedAt,
			&fromLoc.id, &fromLoc.address, &fromLoc.placeName, &fromLoc.latitude, &fromLoc.longitude, &fromLoc.createdAt, &fromLoc.updatedAt,
			&toLoc.id, &toLoc.address, &toLoc.placeName, &toLoc.latitude, &toLoc.longitude, &toLoc.createdAt, &toLoc.updatedAt,
		); err != nil {
			return nil, err
		}
		if fromLoc.id != nil {
			m.FromLocation = fromLoc.toDomain()
		}
		m.ToLocation = toLoc.toDomain()
		list = append(list, m)
	}
	return list, rows.Err()
}

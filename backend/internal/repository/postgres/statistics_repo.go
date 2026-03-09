package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/natapage/test_assignment/backend/internal/domain"
)

type StatisticsRepo struct {
	pool *pgxpool.Pool
}

func NewStatisticsRepo(pool *pgxpool.Pool) *StatisticsRepo {
	return &StatisticsRepo{pool: pool}
}

func (r *StatisticsRepo) GetLocationDurations(ctx context.Context, machineID int64) ([]domain.LocationDuration, error) {
	rows, err := r.pool.Query(ctx, `
		WITH periods AS (
			SELECT
				to_location_id AS location_id,
				moved_at AS start_at,
				LEAD(moved_at) OVER (PARTITION BY machine_id ORDER BY moved_at) AS end_at
			FROM movement_history
			WHERE machine_id = $1
		)
		SELECT
			p.location_id,
			l.id, l.address, l.place_name, l.latitude, l.longitude, l.created_at, l.updated_at,
			COALESCE(SUM(EXTRACT(DAY FROM (COALESCE(p.end_at, NOW()) - p.start_at)))::INT, 0) AS days
		FROM periods p
		JOIN locations l ON p.location_id = l.id
		GROUP BY p.location_id, l.id, l.address, l.place_name, l.latitude, l.longitude, l.created_at, l.updated_at
		ORDER BY days DESC`, machineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.LocationDuration
	for rows.Next() {
		var d domain.LocationDuration
		var loc locationScan
		fields := append([]any{&d.LocationID}, loc.scanFields()...)
		fields = append(fields, &d.Days)
		if err := rows.Scan(fields...); err != nil {
			return nil, err
		}
		d.Location = loc.toDomain()
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *StatisticsRepo) GetMovementsCount(ctx context.Context, from, to time.Time) ([]domain.MovementsCount, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT m.id, m.name, COUNT(mh.id)::INT AS cnt
		FROM machines m
		LEFT JOIN movement_history mh ON m.id = mh.machine_id
			AND mh.moved_at BETWEEN $1 AND $2
		GROUP BY m.id, m.name
		ORDER BY cnt DESC, m.id`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.MovementsCount
	for rows.Next() {
		var mc domain.MovementsCount
		var name string
		if err := rows.Scan(&mc.MachineID, &name, &mc.Count); err != nil {
			return nil, err
		}
		mc.Machine = &domain.Machine{ID: mc.MachineID, Name: name}
		list = append(list, mc)
	}
	return list, rows.Err()
}

func (r *StatisticsRepo) GetMachineTimeline(ctx context.Context, machineID int64) ([]domain.TimelineEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT mh.moved_at,
		       fl.id, fl.address, fl.place_name, fl.latitude, fl.longitude, fl.created_at, fl.updated_at,
		       tl.id, tl.address, tl.place_name, tl.latitude, tl.longitude, tl.created_at, tl.updated_at
		FROM movement_history mh
		LEFT JOIN locations fl ON mh.from_location_id = fl.id
		JOIN locations tl ON mh.to_location_id = tl.id
		WHERE mh.machine_id = $1
		ORDER BY mh.moved_at ASC`, machineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.TimelineEntry
	for rows.Next() {
		var e domain.TimelineEntry
		var fromLoc, toLoc locationScan
		fields := append([]any{&e.MovedAt}, fromLoc.scanFields()...)
		fields = append(fields, toLoc.scanFields()...)
		if err := rows.Scan(fields...); err != nil {
			return nil, err
		}
		if fromLoc.id != nil {
			e.FromLocation = fromLoc.toDomain()
		}
		e.ToLocation = toLoc.toDomain()
		list = append(list, e)
	}
	return list, rows.Err()
}

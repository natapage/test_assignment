package repository

import (
	"context"
	"time"

	"github.com/natapage/test_assignment/backend/internal/domain"
)

type MachineRepository interface {
	List(ctx context.Context) ([]domain.Machine, error)
	GetByID(ctx context.Context, id int64) (*domain.Machine, error)
	Create(ctx context.Context, m *domain.Machine) (*domain.Machine, error)
	Update(ctx context.Context, m *domain.Machine) (*domain.Machine, error)
	Delete(ctx context.Context, id int64) error
	UpdateLocationID(ctx context.Context, id int64, locationID *int64) error
}

type LocationRepository interface {
	List(ctx context.Context) ([]domain.Location, error)
	GetByID(ctx context.Context, id int64) (*domain.Location, error)
	Create(ctx context.Context, l *domain.Location) (*domain.Location, error)
	Update(ctx context.Context, l *domain.Location) (*domain.Location, error)
	Delete(ctx context.Context, id int64) error
}

type MovementRepository interface {
	Create(ctx context.Context, m *domain.Movement) (*domain.Movement, error)
	ListByMachineID(ctx context.Context, machineID int64) ([]domain.Movement, error)
}

type StatisticsRepository interface {
	GetLocationDurations(ctx context.Context, machineID int64) ([]domain.LocationDuration, error)
	GetMovementsCount(ctx context.Context, from, to time.Time) ([]domain.MovementsCount, error)
	GetMachineTimeline(ctx context.Context, machineID int64) ([]domain.TimelineEntry, error)
}

// TxManager abstracts transaction handling.
type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

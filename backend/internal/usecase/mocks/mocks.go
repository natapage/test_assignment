package mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/natapage/test_assignment/backend/internal/domain"
)

type MachineRepo struct {
	Machines     map[int64]*domain.Machine
	ListErr      error
	GetErr       error
	CreateErr    error
	UpdateErr    error
	DeleteErr    error
	UpdateLocErr error
}

func NewMachineRepo() *MachineRepo {
	return &MachineRepo{Machines: make(map[int64]*domain.Machine)}
}

func (r *MachineRepo) List(_ context.Context) ([]domain.Machine, error) {
	if r.ListErr != nil {
		return nil, r.ListErr
	}
	var list []domain.Machine
	for _, m := range r.Machines {
		list = append(list, *m)
	}
	return list, nil
}

func (r *MachineRepo) GetByID(_ context.Context, id int64) (*domain.Machine, error) {
	if r.GetErr != nil {
		return nil, r.GetErr
	}
	m, ok := r.Machines[id]
	if !ok {
		return nil, fmt.Errorf("machine %d not found", id)
	}
	return m, nil
}

func (r *MachineRepo) Create(_ context.Context, m *domain.Machine) (*domain.Machine, error) {
	if r.CreateErr != nil {
		return nil, r.CreateErr
	}
	m.ID = int64(len(r.Machines) + 1)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	r.Machines[m.ID] = m
	return m, nil
}

func (r *MachineRepo) Update(_ context.Context, m *domain.Machine) (*domain.Machine, error) {
	if r.UpdateErr != nil {
		return nil, r.UpdateErr
	}
	existing, ok := r.Machines[m.ID]
	if !ok {
		return nil, fmt.Errorf("machine %d not found", m.ID)
	}
	existing.Name = m.Name
	existing.SerialNumber = m.SerialNumber
	existing.Enabled = m.Enabled
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (r *MachineRepo) Delete(_ context.Context, id int64) error {
	if r.DeleteErr != nil {
		return r.DeleteErr
	}
	delete(r.Machines, id)
	return nil
}

func (r *MachineRepo) UpdateLocationID(_ context.Context, id int64, locationID *int64) error {
	if r.UpdateLocErr != nil {
		return r.UpdateLocErr
	}
	m, ok := r.Machines[id]
	if !ok {
		return fmt.Errorf("machine %d not found", id)
	}
	m.LocationID = locationID
	return nil
}

type LocationRepo struct {
	Locations map[int64]*domain.Location
	ListErr   error
	GetErr    error
	CreateErr error
	UpdateErr error
	DeleteErr error
}

func NewLocationRepo() *LocationRepo {
	return &LocationRepo{Locations: make(map[int64]*domain.Location)}
}

func (r *LocationRepo) List(_ context.Context) ([]domain.Location, error) {
	if r.ListErr != nil {
		return nil, r.ListErr
	}
	var list []domain.Location
	for _, l := range r.Locations {
		list = append(list, *l)
	}
	return list, nil
}

func (r *LocationRepo) GetByID(_ context.Context, id int64) (*domain.Location, error) {
	if r.GetErr != nil {
		return nil, r.GetErr
	}
	l, ok := r.Locations[id]
	if !ok {
		return nil, fmt.Errorf("location %d not found", id)
	}
	return l, nil
}

func (r *LocationRepo) Create(_ context.Context, l *domain.Location) (*domain.Location, error) {
	if r.CreateErr != nil {
		return nil, r.CreateErr
	}
	l.ID = int64(len(r.Locations) + 1)
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()
	r.Locations[l.ID] = l
	return l, nil
}

func (r *LocationRepo) Update(_ context.Context, l *domain.Location) (*domain.Location, error) {
	if r.UpdateErr != nil {
		return nil, r.UpdateErr
	}
	existing, ok := r.Locations[l.ID]
	if !ok {
		return nil, fmt.Errorf("location %d not found", l.ID)
	}
	existing.Address = l.Address
	existing.PlaceName = l.PlaceName
	existing.Latitude = l.Latitude
	existing.Longitude = l.Longitude
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (r *LocationRepo) Delete(_ context.Context, id int64) error {
	if r.DeleteErr != nil {
		return r.DeleteErr
	}
	delete(r.Locations, id)
	return nil
}

type MovementRepo struct {
	Movements []*domain.Movement
	CreateErr error
	ListErr   error
}

func NewMovementRepo() *MovementRepo {
	return &MovementRepo{}
}

func (r *MovementRepo) Create(_ context.Context, m *domain.Movement) (*domain.Movement, error) {
	if r.CreateErr != nil {
		return nil, r.CreateErr
	}
	m.ID = int64(len(r.Movements) + 1)
	m.MovedAt = time.Now()
	m.CreatedAt = time.Now()
	r.Movements = append(r.Movements, m)
	return m, nil
}

func (r *MovementRepo) ListByMachineID(_ context.Context, machineID int64) ([]domain.Movement, error) {
	if r.ListErr != nil {
		return nil, r.ListErr
	}
	var list []domain.Movement
	for _, m := range r.Movements {
		if m.MachineID == machineID {
			list = append(list, *m)
		}
	}
	return list, nil
}

type TxManager struct{}

func (tm *TxManager) RunInTx(_ context.Context, fn func(ctx context.Context) error) error {
	return fn(context.Background())
}

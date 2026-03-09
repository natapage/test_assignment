package usecase

import (
	"context"
	"errors"

	"github.com/natapage/test_assignment/backend/internal/domain"
	"github.com/natapage/test_assignment/backend/internal/repository"
)

var (
	ErrSameLocation = errors.New("machine is already at this location")
	ErrNoLocation   = errors.New("machine has no current location")
)

type MovementUseCase struct {
	machineRepo  repository.MachineRepository
	locationRepo repository.LocationRepository
	movementRepo repository.MovementRepository
	txManager    repository.TxManager
}

func NewMovementUseCase(
	machineRepo repository.MachineRepository,
	locationRepo repository.LocationRepository,
	movementRepo repository.MovementRepository,
	txManager repository.TxManager,
) *MovementUseCase {
	return &MovementUseCase{
		machineRepo:  machineRepo,
		locationRepo: locationRepo,
		movementRepo: movementRepo,
		txManager:    txManager,
	}
}

func (uc *MovementUseCase) MoveMachine(ctx context.Context, machineID, toLocationID int64) (*domain.Movement, error) {
	var created *domain.Movement
	err := uc.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		machine, err := uc.machineRepo.GetByID(txCtx, machineID)
		if err != nil {
			return err
		}

		if _, err := uc.locationRepo.GetByID(txCtx, toLocationID); err != nil {
			return err
		}

		if machine.LocationID != nil && *machine.LocationID == toLocationID {
			return ErrSameLocation
		}

		movement := &domain.Movement{
			MachineID:      machineID,
			FromLocationID: machine.LocationID,
			ToLocationID:   toLocationID,
		}

		created, err = uc.movementRepo.Create(txCtx, movement)
		if err != nil {
			return err
		}
		return uc.machineRepo.UpdateLocationID(txCtx, machineID, &toLocationID)
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *MovementUseCase) GetHistory(ctx context.Context, machineID int64) ([]domain.Movement, error) {
	if _, err := uc.machineRepo.GetByID(ctx, machineID); err != nil {
		return nil, err
	}
	return uc.movementRepo.ListByMachineID(ctx, machineID)
}

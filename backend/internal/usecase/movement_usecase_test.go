package usecase_test

import (
	"context"
	"testing"

	"github.com/natapage/test_assignment/backend/internal/domain"
	"github.com/natapage/test_assignment/backend/internal/usecase"
	"github.com/natapage/test_assignment/backend/internal/usecase/mocks"
)

func setupMovement() (*usecase.MovementUseCase, *mocks.MachineRepo, *mocks.LocationRepo, *mocks.MovementRepo) {
	machineRepo := mocks.NewMachineRepo()
	locationRepo := mocks.NewLocationRepo()
	movementRepo := mocks.NewMovementRepo()
	txManager := &mocks.TxManager{}

	locID := int64(1)
	locationRepo.Locations[1] = &domain.Location{ID: 1, Address: "addr1", PlaceName: "place1"}
	locationRepo.Locations[2] = &domain.Location{ID: 2, Address: "addr2", PlaceName: "place2"}
	machineRepo.Machines[1] = &domain.Machine{ID: 1, Name: "M1", LocationID: &locID}

	uc := usecase.NewMovementUseCase(machineRepo, locationRepo, movementRepo, txManager)
	return uc, machineRepo, locationRepo, movementRepo
}

func TestMoveMachine_Success(t *testing.T) {
	uc, machineRepo, _, movementRepo := setupMovement()

	movement, err := uc.MoveMachine(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if movement == nil {
		t.Fatal("expected movement, got nil")
	}
	if movement.ToLocationID != 2 {
		t.Errorf("expected to_location_id=2, got %d", movement.ToLocationID)
	}
	if movement.FromLocationID == nil || *movement.FromLocationID != 1 {
		t.Error("expected from_location_id=1")
	}

	// Verify machine location updated
	m := machineRepo.Machines[1]
	if m.LocationID == nil || *m.LocationID != 2 {
		t.Errorf("expected machine location_id=2, got %v", m.LocationID)
	}

	// Verify movement recorded
	if len(movementRepo.Movements) != 1 {
		t.Errorf("expected 1 movement, got %d", len(movementRepo.Movements))
	}
}

func TestMoveMachine_SameLocation(t *testing.T) {
	uc, _, _, _ := setupMovement()

	_, err := uc.MoveMachine(context.Background(), 1, 1)
	if err == nil {
		t.Fatal("expected error for same location")
	}
	if err != usecase.ErrSameLocation {
		t.Errorf("expected ErrSameLocation, got: %v", err)
	}
}

func TestMoveMachine_MachineNotFound(t *testing.T) {
	uc, _, _, _ := setupMovement()

	_, err := uc.MoveMachine(context.Background(), 999, 2)
	if err == nil {
		t.Fatal("expected error for missing machine")
	}
}

func TestMoveMachine_LocationNotFound(t *testing.T) {
	uc, _, _, _ := setupMovement()

	_, err := uc.MoveMachine(context.Background(), 1, 999)
	if err == nil {
		t.Fatal("expected error for missing location")
	}
}

func TestGetHistory_Success(t *testing.T) {
	uc, _, _, movementRepo := setupMovement()

	fromID := int64(1)
	movementRepo.Movements = append(movementRepo.Movements, &domain.Movement{
		ID: 1, MachineID: 1, FromLocationID: &fromID, ToLocationID: 2,
	})

	movements, err := uc.GetHistory(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(movements) != 1 {
		t.Errorf("expected 1 movement, got %d", len(movements))
	}
}

func TestGetHistory_MachineNotFound(t *testing.T) {
	uc, _, _, _ := setupMovement()

	_, err := uc.GetHistory(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for missing machine")
	}
}

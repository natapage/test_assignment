package usecase_test

import (
	"context"
	"testing"

	"github.com/natapage/test_assignment/backend/internal/domain"
	"github.com/natapage/test_assignment/backend/internal/usecase"
	"github.com/natapage/test_assignment/backend/internal/usecase/mocks"
)

func TestMachineList(t *testing.T) {
	repo := mocks.NewMachineRepo()
	repo.Machines[1] = &domain.Machine{ID: 1, Name: "M1"}
	repo.Machines[2] = &domain.Machine{ID: 2, Name: "M2"}
	uc := usecase.NewMachineUseCase(repo)

	machines, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(machines) != 2 {
		t.Errorf("expected 2 machines, got %d", len(machines))
	}
}

func TestMachineGetByID(t *testing.T) {
	repo := mocks.NewMachineRepo()
	repo.Machines[1] = &domain.Machine{ID: 1, Name: "M1"}
	uc := usecase.NewMachineUseCase(repo)

	m, err := uc.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Name != "M1" {
		t.Errorf("expected M1, got %s", m.Name)
	}
}

func TestMachineGetByID_NotFound(t *testing.T) {
	repo := mocks.NewMachineRepo()
	uc := usecase.NewMachineUseCase(repo)

	_, err := uc.GetByID(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for missing machine")
	}
}

func TestMachineCreate(t *testing.T) {
	repo := mocks.NewMachineRepo()
	uc := usecase.NewMachineUseCase(repo)

	m, err := uc.Create(context.Background(), &domain.Machine{Name: "New", SerialNumber: "SN-001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if len(repo.Machines) != 1 {
		t.Errorf("expected 1 machine in repo, got %d", len(repo.Machines))
	}
}

func TestMachineDelete(t *testing.T) {
	repo := mocks.NewMachineRepo()
	repo.Machines[1] = &domain.Machine{ID: 1, Name: "M1"}
	uc := usecase.NewMachineUseCase(repo)

	err := uc.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.Machines) != 0 {
		t.Error("expected machine to be deleted")
	}
}

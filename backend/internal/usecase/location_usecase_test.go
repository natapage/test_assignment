package usecase_test

import (
	"context"
	"testing"

	"github.com/natapage/test_assignment/backend/internal/domain"
	"github.com/natapage/test_assignment/backend/internal/usecase"
	"github.com/natapage/test_assignment/backend/internal/usecase/mocks"
)

func TestLocationList(t *testing.T) {
	repo := mocks.NewLocationRepo()
	repo.Locations[1] = &domain.Location{ID: 1, PlaceName: "L1"}
	repo.Locations[2] = &domain.Location{ID: 2, PlaceName: "L2"}
	uc := usecase.NewLocationUseCase(repo)

	locations, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locations) != 2 {
		t.Errorf("expected 2 locations, got %d", len(locations))
	}
}

func TestLocationGetByID(t *testing.T) {
	repo := mocks.NewLocationRepo()
	repo.Locations[1] = &domain.Location{ID: 1, PlaceName: "L1"}
	uc := usecase.NewLocationUseCase(repo)

	l, err := uc.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.PlaceName != "L1" {
		t.Errorf("expected L1, got %s", l.PlaceName)
	}
}

func TestLocationCreate(t *testing.T) {
	repo := mocks.NewLocationRepo()
	uc := usecase.NewLocationUseCase(repo)

	l, err := uc.Create(context.Background(), &domain.Location{Address: "addr", PlaceName: "place"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

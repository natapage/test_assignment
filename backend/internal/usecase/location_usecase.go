package usecase

import (
	"context"

	"github.com/natapage/test_assignment/backend/internal/domain"
	"github.com/natapage/test_assignment/backend/internal/repository"
)

type LocationUseCase struct {
	repo repository.LocationRepository
}

func NewLocationUseCase(repo repository.LocationRepository) *LocationUseCase {
	return &LocationUseCase{repo: repo}
}

func (uc *LocationUseCase) List(ctx context.Context) ([]domain.Location, error) {
	return uc.repo.List(ctx)
}

func (uc *LocationUseCase) GetByID(ctx context.Context, id int64) (*domain.Location, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *LocationUseCase) Create(ctx context.Context, l *domain.Location) (*domain.Location, error) {
	return uc.repo.Create(ctx, l)
}

func (uc *LocationUseCase) Update(ctx context.Context, l *domain.Location) (*domain.Location, error) {
	return uc.repo.Update(ctx, l)
}

func (uc *LocationUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

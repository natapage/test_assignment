package usecase

import (
	"context"

	"github.com/natapage/test_assignment/backend/internal/domain"
	"github.com/natapage/test_assignment/backend/internal/repository"
)

type MachineUseCase struct {
	repo repository.MachineRepository
}

func NewMachineUseCase(repo repository.MachineRepository) *MachineUseCase {
	return &MachineUseCase{repo: repo}
}

func (uc *MachineUseCase) List(ctx context.Context) ([]domain.Machine, error) {
	return uc.repo.List(ctx)
}

func (uc *MachineUseCase) GetByID(ctx context.Context, id int64) (*domain.Machine, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *MachineUseCase) Create(ctx context.Context, m *domain.Machine) (*domain.Machine, error) {
	return uc.repo.Create(ctx, m)
}

func (uc *MachineUseCase) Update(ctx context.Context, m *domain.Machine) (*domain.Machine, error) {
	return uc.repo.Update(ctx, m)
}

func (uc *MachineUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

package usecase

import (
	"context"
	"time"

	"github.com/natapage/test_assignment/backend/internal/domain"
	"github.com/natapage/test_assignment/backend/internal/repository"
)

type StatisticsUseCase struct {
	repo repository.StatisticsRepository
}

func NewStatisticsUseCase(repo repository.StatisticsRepository) *StatisticsUseCase {
	return &StatisticsUseCase{repo: repo}
}

func (uc *StatisticsUseCase) GetLocationDurations(ctx context.Context, machineID int64) ([]domain.LocationDuration, error) {
	return uc.repo.GetLocationDurations(ctx, machineID)
}

func (uc *StatisticsUseCase) GetMovementsCount(ctx context.Context, from, to time.Time) ([]domain.MovementsCount, error) {
	return uc.repo.GetMovementsCount(ctx, from, to)
}

func (uc *StatisticsUseCase) GetMachineTimeline(ctx context.Context, machineID int64) ([]domain.TimelineEntry, error) {
	return uc.repo.GetMachineTimeline(ctx, machineID)
}

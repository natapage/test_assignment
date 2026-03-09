package postgres

import (
	"time"

	"github.com/natapage/test_assignment/backend/internal/domain"
)

// locationScan is a helper for scanning nullable location joins.
type locationScan struct {
	id        *int64
	address   *string
	placeName *string
	latitude  *float64
	longitude *float64
	createdAt *time.Time
	updatedAt *time.Time
}

func (ls *locationScan) scanFields() []any {
	return []any{&ls.id, &ls.address, &ls.placeName, &ls.latitude, &ls.longitude, &ls.createdAt, &ls.updatedAt}
}

func (ls *locationScan) toDomain() *domain.Location {
	if ls.id == nil {
		return nil
	}
	return &domain.Location{
		ID:        *ls.id,
		Address:   deref(ls.address),
		PlaceName: deref(ls.placeName),
		Latitude:  ls.latitude,
		Longitude: ls.longitude,
		CreatedAt: derefTime(ls.createdAt),
		UpdatedAt: derefTime(ls.updatedAt),
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

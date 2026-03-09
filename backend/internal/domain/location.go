package domain

import "time"

type Location struct {
	ID        int64
	Address   string
	PlaceName string
	Latitude  *float64
	Longitude *float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

package domain

import "time"

type Machine struct {
	ID           int64
	Name         string
	SerialNumber string
	Enabled      bool
	LocationID   *int64
	Location     *Location
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

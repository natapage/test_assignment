package domain

import "time"

type Movement struct {
	ID             int64
	MachineID      int64
	FromLocationID *int64
	ToLocationID   int64
	FromLocation   *Location
	ToLocation     *Location
	MovedAt        time.Time
	CreatedAt      time.Time
}

type LocationDuration struct {
	LocationID int64
	Location   *Location
	Days       int
}

type MovementsCount struct {
	MachineID int64
	Machine   *Machine
	Count     int
}

type TimelineEntry struct {
	MovedAt      time.Time
	FromLocation *Location
	ToLocation   *Location
}

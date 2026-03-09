export interface Location {
  id: string
  address: string
  placeName: string
  latitude?: number
  longitude?: number
  createdAt: string
  updatedAt: string
}

export interface Machine {
  id: string
  name: string
  serialNumber: string
  enabled: boolean
  locationId?: string
  location?: Location
  createdAt: string
  updatedAt: string
}

export interface Movement {
  id: string
  machineId: string
  fromLocationId?: string
  toLocationId: string
  fromLocation?: Location
  toLocation?: Location
  movedAt: string
  createdAt: string
}

export interface LocationDuration {
  locationId: string
  location?: Location
  days: number
}

export interface MovementsCount {
  machineId: string
  machineName: string
  count: number
}

export interface TimelineEntry {
  movedAt: string
  fromLocation?: Location
  toLocation?: Location
}

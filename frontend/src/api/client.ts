import type { Machine, Location, Movement, LocationDuration, MovementsCount, TimelineEntry } from '@/types'

const BASE = '/api/v1'

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(BASE + url, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.message || `HTTP ${res.status}`)
  }
  return res.json()
}

export const api = {
  // Machines
  async listMachines(): Promise<Machine[]> {
    const data = await request<{ machines: Machine[] }>('/machines')
    return data.machines || []
  },
  async getMachine(id: string): Promise<Machine> {
    const data = await request<{ machine: Machine }>(`/machines/${id}`)
    return data.machine
  },
  async createMachine(body: { name: string; serial_number: string; enabled: boolean; location_id?: string }): Promise<Machine> {
    const data = await request<{ machine: Machine }>('/machines', { method: 'POST', body: JSON.stringify(body) })
    return data.machine
  },
  async updateMachine(id: string, body: { name: string; serial_number: string; enabled: boolean }): Promise<Machine> {
    const data = await request<{ machine: Machine }>(`/machines/${id}`, { method: 'PUT', body: JSON.stringify(body) })
    return data.machine
  },
  async deleteMachine(id: string): Promise<void> {
    await request(`/machines/${id}`, { method: 'DELETE' })
  },

  // Locations
  async listLocations(): Promise<Location[]> {
    const data = await request<{ locations: Location[] }>('/locations')
    return data.locations || []
  },
  async createLocation(body: { address: string; place_name: string; latitude?: number; longitude?: number }): Promise<Location> {
    const data = await request<{ location: Location }>('/locations', { method: 'POST', body: JSON.stringify(body) })
    return data.location
  },
  async updateLocation(id: string, body: { address: string; place_name: string; latitude?: number; longitude?: number }): Promise<Location> {
    const data = await request<{ location: Location }>(`/locations/${id}`, { method: 'PUT', body: JSON.stringify(body) })
    return data.location
  },
  async deleteLocation(id: string): Promise<void> {
    await request(`/locations/${id}`, { method: 'DELETE' })
  },

  // Movements
  async moveMachine(machineId: string, toLocationId: string): Promise<Movement> {
    const data = await request<{ movement: Movement }>(`/machines/${machineId}/move`, {
      method: 'POST',
      body: JSON.stringify({ to_location_id: toLocationId }),
    })
    return data.movement
  },
  async getMovementHistory(machineId: string): Promise<Movement[]> {
    const data = await request<{ movements: Movement[] }>(`/machines/${machineId}/movements`)
    return data.movements || []
  },

  // Statistics
  async getLocationDurations(machineId: string): Promise<LocationDuration[]> {
    const data = await request<{ durations: LocationDuration[] }>(`/statistics/durations/${machineId}`)
    return data.durations || []
  },
  async getMovementsCount(from: string, to: string): Promise<MovementsCount[]> {
    const data = await request<{ counts: MovementsCount[] }>(`/statistics/movements-count?from=${from}&to=${to}`)
    return data.counts || []
  },
  async getMachineTimeline(machineId: string): Promise<TimelineEntry[]> {
    const data = await request<{ entries: TimelineEntry[] }>(`/statistics/timeline/${machineId}`)
    return data.entries || []
  },
}

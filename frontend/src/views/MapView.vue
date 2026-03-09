<template>
  <div>
    <h1 class="text-2xl font-semibold mb-6 text-foreground">Карта аппаратов</h1>
    <Card class="overflow-hidden">
      <div ref="mapContainer" class="h-[600px]"></div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import L from 'leaflet'
import { api } from '@/api/client'
import type { Machine } from '@/types'
import { Card } from '@/components/ui/card'

const mapContainer = ref<HTMLDivElement>()
let mapInstance: L.Map | null = null

onMounted(async () => {
  const machines: Machine[] = await api.listMachines()

  const map = L.map(mapContainer.value!, {
    center: [55.75, 37.62],
    zoom: 5,
  })

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors',
  }).addTo(map)

  const goldIcon = L.divIcon({
    className: 'custom-marker',
    html: '<div style="background:#e6b800;width:28px;height:28px;border-radius:50%;border:3px solid #1c1c2e;display:flex;align-items:center;justify-content:center;font-weight:bold;font-size:11px;color:#1c1c2e;">G</div>',
    iconSize: [28, 28],
    iconAnchor: [14, 14],
  })

  const bounds: L.LatLngExpression[] = []

  for (const m of machines) {
    if (!m.location?.latitude || !m.location?.longitude) continue
    const latlng: L.LatLngExpression = [m.location.latitude, m.location.longitude]
    bounds.push(latlng)

    L.marker(latlng, { icon: goldIcon })
      .addTo(map)
      .bindPopup(`
        <div style="font-size:14px;line-height:1.5;">
          <b>${m.name}</b> (${m.serialNumber})<br/>
          <span style="color:#555;">${m.location.placeName}</span><br/>
          <span style="font-size:12px;color:#888;">${m.location.address}</span>
        </div>
      `)
  }

  if (bounds.length > 0) {
    map.fitBounds(bounds as L.LatLngBoundsExpression, { padding: [40, 40] })
  }

  mapInstance = map
})

onUnmounted(() => {
  mapInstance?.remove()
  mapInstance = null
})
</script>

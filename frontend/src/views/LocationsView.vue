<template>
  <div>
    <h1>Локации</h1>
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Название</th>
          <th>Адрес</th>
          <th>Координаты</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="loc in locations" :key="loc.id">
          <td>{{ loc.id }}</td>
          <td>{{ loc.placeName }}</td>
          <td>{{ loc.address }}</td>
          <td>{{ loc.latitude?.toFixed(6) }}, {{ loc.longitude?.toFixed(6) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'
import type { Location } from '@/types'

const locations = ref<Location[]>([])

onMounted(async () => {
  locations.value = await api.listLocations()
})
</script>

<template>
  <div>
    <h1>Статистика</h1>

    <div class="machine-select">
      <label>Аппарат</label>
      <select v-model="selectedMachineId" @change="loadStats">
        <option value="" disabled>Выберите аппарат</option>
        <option v-for="m in machines" :key="m.id" :value="m.id">
          {{ m.name }} ({{ m.serialNumber }})
        </option>
      </select>
    </div>

    <!-- Duration on locations -->
    <div v-if="durations.length" class="section">
      <h2>Дни на локациях</h2>
      <table>
        <thead>
          <tr>
            <th>Локация</th>
            <th>Адрес</th>
            <th>Дней</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="d in durations" :key="d.locationId">
            <td>{{ d.location?.placeName }}</td>
            <td>{{ d.location?.address }}</td>
            <td>{{ d.days }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Timeline -->
    <div v-if="timeline.length" class="section">
      <h2>Timeline перемещений</h2>
      <table>
        <thead>
          <tr>
            <th>Дата</th>
            <th>Откуда</th>
            <th>Куда</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(entry, i) in timeline" :key="i">
            <td>{{ formatDate(entry.movedAt) }}</td>
            <td>{{ entry.fromLocation?.placeName || 'Первая установка' }}</td>
            <td>{{ entry.toLocation?.placeName }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Movements count -->
    <div class="section">
      <h2>Количество перемещений за период</h2>
      <div class="date-row">
        <div class="date-field">
          <label>С</label>
          <input type="date" v-model="dateFrom" />
        </div>
        <span class="date-separator">&mdash;</span>
        <div class="date-field">
          <label>По</label>
          <input type="date" v-model="dateTo" />
        </div>
        <button class="btn-primary" @click="loadCounts">Показать</button>
      </div>
      <table v-if="counts.length">
        <thead>
          <tr>
            <th>Аппарат</th>
            <th>Перемещений</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in counts" :key="c.machineId">
            <td>{{ c.machineName }}</td>
            <td>{{ c.count }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'
import type { Machine, LocationDuration, TimelineEntry, MovementsCount } from '@/types'

const machines = ref<Machine[]>([])
const selectedMachineId = ref('')
const durations = ref<LocationDuration[]>([])
const timeline = ref<TimelineEntry[]>([])
const counts = ref<MovementsCount[]>([])
const dateFrom = ref('2020-01-01')
const dateTo = ref(new Date().toISOString().slice(0, 10))

async function loadStats() {
  if (!selectedMachineId.value) return
  const [d, t] = await Promise.all([
    api.getLocationDurations(selectedMachineId.value),
    api.getMachineTimeline(selectedMachineId.value),
  ])
  durations.value = d
  timeline.value = t
}

async function loadCounts() {
  counts.value = await api.getMovementsCount(
    new Date(dateFrom.value).toISOString(),
    new Date(dateTo.value).toISOString(),
  )
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('ru-RU', {
    year: 'numeric', month: 'long', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
}

onMounted(async () => {
  machines.value = await api.listMachines()
})
</script>

<style scoped>
.machine-select {
  margin-bottom: 20px;
}
.machine-select label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 4px;
  color: #555;
}
.machine-select select {
  width: 320px;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}
.section {
  margin-top: 28px;
}
.section h2 {
  font-size: 18px;
  margin-bottom: 14px;
}
.date-row {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  margin-bottom: 16px;
}
.date-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.date-field label {
  font-size: 13px;
  font-weight: 600;
  color: #555;
}
.date-field input {
  padding: 7px 10px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  width: 160px;
}
.date-separator {
  color: #999;
  padding-bottom: 8px;
  font-size: 16px;
}
.date-row .btn-primary {
  height: 36px;
  margin-bottom: 1px;
}
</style>

<template>
  <div>
    <h1>Аппараты</h1>
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Название</th>
          <th>Серийный номер</th>
          <th>Статус</th>
          <th>Локация</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="m in machines" :key="m.id">
          <td>{{ m.id }}</td>
          <td>{{ m.name }}</td>
          <td>{{ m.serialNumber }}</td>
          <td>{{ m.enabled ? 'Активен' : 'Отключён' }}</td>
          <td>{{ m.location ? m.location.placeName + ' — ' + m.location.address : '—' }}</td>
          <td>
            <button class="btn-primary" @click="openMoveDialog(m)">Переместить</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="showMove" class="modal-overlay" @click.self="showMove = false">
      <div class="modal">
        <h2>Переместить {{ selectedMachine?.name }}</h2>
        <div class="form-group">
          <label>Текущая локация</label>
          <input type="text" disabled :value="selectedMachine?.location?.placeName || '—'" />
        </div>
        <div class="form-group">
          <label>Новая локация</label>
          <select v-model="targetLocationId">
            <option value="" disabled>Выберите локацию</option>
            <option
              v-for="loc in locations"
              :key="loc.id"
              :value="loc.id"
              :disabled="loc.id === selectedMachine?.locationId"
            >
              {{ loc.placeName }} — {{ loc.address }}
            </option>
          </select>
        </div>
        <p v-if="moveError" class="error">{{ moveError }}</p>
        <div class="modal-actions">
          <button class="btn-secondary" @click="showMove = false">Отмена</button>
          <button class="btn-primary" @click="doMove" :disabled="!targetLocationId">Переместить</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'
import type { Machine, Location } from '@/types'

const machines = ref<Machine[]>([])
const locations = ref<Location[]>([])
const showMove = ref(false)
const selectedMachine = ref<Machine | null>(null)
const targetLocationId = ref('')
const moveError = ref('')

async function load() {
  const [m, l] = await Promise.all([api.listMachines(), api.listLocations()])
  machines.value = m
  locations.value = l
}

function openMoveDialog(m: Machine) {
  selectedMachine.value = m
  targetLocationId.value = ''
  moveError.value = ''
  showMove.value = true
}

async function doMove() {
  if (!selectedMachine.value || !targetLocationId.value) return
  try {
    await api.moveMachine(selectedMachine.value.id, targetLocationId.value)
    showMove.value = false
    await load()
  } catch (e: any) {
    moveError.value = e.message
  }
}

onMounted(load)
</script>

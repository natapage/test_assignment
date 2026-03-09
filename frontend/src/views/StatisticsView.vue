<template>
  <div>
    <h1 class="text-2xl font-semibold mb-6 text-foreground">Статистика</h1>

    <div class="mb-8 space-y-2">
      <Label>Аппарат</Label>
      <Select v-model="selectedMachineId" @update:model-value="loadStats">
        <SelectTrigger class="w-96">
          <SelectValue placeholder="Выберите аппарат" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem v-for="m in machines" :key="m.id" :value="m.id">
            {{ m.name }} ({{ m.serialNumber }})
          </SelectItem>
        </SelectContent>
      </Select>
    </div>

    <Card v-if="durations.length" class="mb-6">
      <CardHeader>
        <CardTitle class="text-lg">Дни на локациях</CardTitle>
      </CardHeader>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Локация</TableHead>
              <TableHead>Адрес</TableHead>
              <TableHead class="text-right">Дней</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="d in durations" :key="d.locationId">
              <TableCell class="font-medium">{{ d.location?.placeName }}</TableCell>
              <TableCell>{{ d.location?.address }}</TableCell>
              <TableCell class="text-right font-mono">{{ d.days }}</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <Card v-if="timeline.length" class="mb-6">
      <CardHeader>
        <CardTitle class="text-lg">Timeline перемещений</CardTitle>
      </CardHeader>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Дата</TableHead>
              <TableHead>Откуда</TableHead>
              <TableHead>Куда</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="(entry, i) in timeline" :key="i">
              <TableCell class="text-muted-foreground whitespace-nowrap">{{ formatDate(entry.movedAt) }}</TableCell>
              <TableCell>{{ entry.fromLocation?.placeName || 'Первая установка' }}</TableCell>
              <TableCell class="font-medium">{{ entry.toLocation?.placeName }}</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="text-lg">Количество перемещений за период</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="flex items-end gap-4 mb-5">
          <div class="flex flex-col gap-1.5">
            <Label>С</Label>
            <Input type="date" v-model="dateFrom" class="w-44" />
          </div>
          <span class="text-muted-foreground pb-2">&mdash;</span>
          <div class="flex flex-col gap-1.5">
            <Label>По</Label>
            <Input type="date" v-model="dateTo" class="w-44" />
          </div>
          <Button @click="loadCounts">Показать</Button>
        </div>
        <Table v-if="counts.length">
          <TableHeader>
            <TableRow>
              <TableHead>Аппарат</TableHead>
              <TableHead class="text-right">Перемещений</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="c in counts" :key="c.machineId">
              <TableCell class="font-medium">{{ c.machineName }}</TableCell>
              <TableCell class="text-right font-mono">{{ c.count }}</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'
import type { Machine, LocationDuration, TimelineEntry, MovementsCount } from '@/types'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'

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

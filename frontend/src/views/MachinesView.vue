<template>
  <div>
    <h1 class="text-2xl font-semibold mb-6 text-foreground">Аппараты</h1>
    <Card>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Название</TableHead>
              <TableHead>Серийный номер</TableHead>
              <TableHead>Статус</TableHead>
              <TableHead>Локация</TableHead>
              <TableHead class="text-right">Действия</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="m in machines" :key="m.id">
              <TableCell class="text-muted-foreground">{{ m.id }}</TableCell>
              <TableCell class="font-medium">{{ m.name }}</TableCell>
              <TableCell class="text-muted-foreground font-mono text-[13px]">{{ m.serialNumber }}</TableCell>
              <TableCell>
                <Badge :variant="m.enabled ? 'success' : 'secondary'">
                  {{ m.enabled ? 'Активен' : 'Отключён' }}
                </Badge>
              </TableCell>
              <TableCell>{{ m.location ? m.location.placeName + ' — ' + m.location.address : '—' }}</TableCell>
              <TableCell class="text-right">
                <Button size="sm" @click="openMoveDialog(m)">Переместить</Button>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <Dialog :open="showMove" @update:open="showMove = $event">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Переместить {{ selectedMachine?.name }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-4 pt-2">
          <div class="space-y-2">
            <Label>Текущая локация</Label>
            <Input type="text" disabled :value="selectedMachine?.location?.placeName || '—'" />
          </div>
          <div class="space-y-2">
            <Label>Новая локация</Label>
            <Select v-model="targetLocationId">
              <SelectTrigger>
                <SelectValue placeholder="Выберите локацию" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="loc in locations"
                  :key="loc.id"
                  :value="loc.id"
                  :disabled="loc.id === selectedMachine?.locationId"
                >
                  {{ loc.placeName }} — {{ loc.address }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <p v-if="moveError" class="text-sm text-destructive">{{ moveError }}</p>
        </div>
        <DialogFooter class="pt-2">
          <Button variant="outline" @click="showMove = false">Отмена</Button>
          <Button @click="doMove" :disabled="!targetLocationId">Переместить</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'
import type { Machine, Location } from '@/types'
import { Card, CardContent } from '@/components/ui/card'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Select, SelectTrigger, SelectContent, SelectItem, SelectValue } from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

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

<template>
  <div>
    <h1 class="text-2xl font-semibold mb-6 text-foreground">Локации</h1>
    <Card>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Название</TableHead>
              <TableHead>Адрес</TableHead>
              <TableHead>Координаты</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="loc in locations" :key="loc.id">
              <TableCell class="text-muted-foreground">{{ loc.id }}</TableCell>
              <TableCell class="font-medium">{{ loc.placeName }}</TableCell>
              <TableCell>{{ loc.address }}</TableCell>
              <TableCell class="text-muted-foreground font-mono text-[13px]">{{ loc.latitude?.toFixed(6) }}, {{ loc.longitude?.toFixed(6) }}</TableCell>
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
import type { Location } from '@/types'
import { Card, CardContent } from '@/components/ui/card'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'

const locations = ref<Location[]>([])

onMounted(async () => {
  locations.value = await api.listLocations()
})
</script>

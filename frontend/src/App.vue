<template>
  <div id="app" class="min-h-screen bg-background">
    <nav class="bg-[#1c1c2e] shadow-md">
      <div class="max-w-[1200px] mx-auto px-6 h-14 flex items-center gap-8">
        <div class="text-xl font-bold text-[#e6b800] tracking-wide select-none">Goldex</div>
        <div class="flex gap-1">
          <router-link
            v-for="link in links"
            :key="link.to"
            :to="link.to"
            class="text-[13px] text-slate-400 no-underline px-3.5 py-1.5 rounded-md transition-all duration-150 hover:bg-white/[0.08] hover:text-white [&.router-link-active]:bg-white/[0.12] [&.router-link-active]:text-white [&.router-link-active]:font-medium"
          >
            {{ link.label }}
          </router-link>
        </div>
      </div>
    </nav>
    <main class="max-w-[1200px] mx-auto py-8 px-6">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
const links = [
  { to: '/machines', label: 'Аппараты' },
  { to: '/locations', label: 'Локации' },
  { to: '/map', label: 'Карта' },
  { to: '/statistics', label: 'Статистика' },
]
</script>

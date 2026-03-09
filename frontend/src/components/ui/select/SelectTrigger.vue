<template>
  <SelectTriggerRaw :class="cn('flex h-9 w-full items-center justify-between whitespace-nowrap rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1', $attrs.class as string)">
    <slot />
    <SelectIcon as-child>
      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="opacity-50"><path d="m6 9 6 6 6-6"/></svg>
    </SelectIcon>
  </SelectTriggerRaw>
</template>

<script setup lang="ts">
import { SelectIcon, SelectTrigger as SelectTriggerRaw } from 'radix-vue'
import { cn } from '@/lib/utils'
</script>

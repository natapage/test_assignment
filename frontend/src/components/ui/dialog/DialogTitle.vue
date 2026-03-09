<template>
  <DialogTitleRaw :class="cn('text-lg font-semibold leading-none tracking-tight', $attrs.class as string)">
    <slot />
  </DialogTitleRaw>
</template>

<script setup lang="ts">
import { DialogTitle as DialogTitleRaw } from 'radix-vue'
import { cn } from '@/lib/utils'
</script>

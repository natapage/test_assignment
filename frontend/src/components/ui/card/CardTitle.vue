<template>
  <h3 :class="cn('font-semibold leading-none tracking-tight', $attrs.class as string)">
    <slot />
  </h3>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

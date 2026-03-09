<template>
  <thead :class="cn('bg-muted/50 [&_tr]:border-b', $attrs.class as string)">
    <slot />
  </thead>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

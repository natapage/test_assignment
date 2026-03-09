<template>
  <tr :class="cn('border-b border-border/60 transition-colors hover:bg-muted/40 data-[state=selected]:bg-muted', $attrs.class as string)">
    <slot />
  </tr>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

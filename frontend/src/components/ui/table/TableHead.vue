<template>
  <th :class="cn('h-12 px-4 text-left align-middle text-xs font-semibold uppercase tracking-wider text-muted-foreground [&:has([role=checkbox])]:pr-0', $attrs.class as string)">
    <slot />
  </th>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

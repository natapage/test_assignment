<template>
  <td :class="cn('px-4 py-3 align-middle text-[14px] [&:has([role=checkbox])]:pr-0', $attrs.class as string)">
    <slot />
  </td>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

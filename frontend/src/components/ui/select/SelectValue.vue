<template>
  <SelectValueRaw v-bind="$attrs">
    <slot />
  </SelectValueRaw>
</template>

<script setup lang="ts">
import { SelectValue as SelectValueRaw } from 'radix-vue'
</script>

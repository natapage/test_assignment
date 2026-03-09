<template>
  <button :class="cn(buttonVariants({ variant, size }), $attrs.class as string)" v-bind="filteredAttrs">
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed, useAttrs } from 'vue'
import { cn } from '@/lib/utils'
import { buttonVariants } from './variants'
import type { VariantProps } from 'class-variance-authority'

type ButtonVariants = VariantProps<typeof buttonVariants>

defineProps<{
  variant?: NonNullable<ButtonVariants['variant']>
  size?: NonNullable<ButtonVariants['size']>
}>()

const attrs = useAttrs()
const filteredAttrs = computed(() => {
  const { class: _, ...rest } = attrs
  return rest
})
</script>

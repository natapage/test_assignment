<template>
  <div :class="cn(badgeVariants({ variant }), $attrs.class as string)">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
import { badgeVariants } from './variants'
import type { VariantProps } from 'class-variance-authority'

type BadgeVariants = VariantProps<typeof badgeVariants>

defineProps<{
  variant?: NonNullable<BadgeVariants['variant']>
}>()
</script>

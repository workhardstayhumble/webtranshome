<script setup>
import { computed } from 'vue'
const props = defineProps({
  variant: { type: String, default: 'neutral', validator: v => ['neutral', 'primary', 'success', 'warning', 'info', 'danger', 'accent'].includes(v) },
})
const VARIAN = {
  neutral: 'bg-subtle text-foreground border-border',
  primary: 'bg-primary-soft text-primary border-primary/15',
  success: 'bg-success-soft text-success-ink border-success/15',
  warning: 'bg-warning-soft text-warning-ink border-warning/15',
  info: 'bg-info-soft text-info-ink border-info/15',
  danger: 'bg-danger-soft text-danger-ink border-danger/15',
  accent: 'bg-accent text-accent-foreground border-accent',
}
const kelas = computed(() => VARIAN[props.variant])
</script>
<template>
  <span class="inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs font-semibold whitespace-nowrap" :class="kelas"><slot /></span>
</template>
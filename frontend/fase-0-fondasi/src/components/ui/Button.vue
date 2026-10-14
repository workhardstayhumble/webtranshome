<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { LoaderCircle } from 'lucide-vue-next'

const props = defineProps({
  variant: { type: String, default: 'primary', validator: v => ['primary', 'secondary', 'outline', 'ghost', 'whatsapp', 'danger'].includes(v) },
  size: { type: String, default: 'md', validator: v => ['sm', 'md', 'lg'].includes(v) },
  type: { type: String, default: 'button' },
  to: { type: [String, Object], default: null },
  href: { type: String, default: null },
  loading: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  block: { type: Boolean, default: false },
  wrap: { type: Boolean, default: false },
})
const VARIAN = {
  primary: 'border border-primary bg-primary text-white hover:bg-primary-hover hover:border-primary-hover',
  secondary: 'border border-secondary bg-secondary text-white hover:bg-secondary-hover',
  outline: 'border border-primary bg-white text-primary hover:bg-primary-soft',
  ghost: 'border border-transparent bg-transparent text-primary hover:bg-primary-soft',
  whatsapp: 'border border-accent bg-accent text-accent-foreground hover:bg-whatsapp-hover',
  danger: 'border border-danger bg-danger text-white hover:bg-danger-ink',
}
const UKURAN = {
  sm: 'min-h-10 px-3.5 text-sm',
  md: 'min-h-11 px-4 text-sm',
  lg: 'min-h-13 px-6 text-base',
}
const nonaktif = computed(() => props.disabled || props.loading)
const tag = computed(() => nonaktif.value ? 'button' : props.to ? RouterLink : props.href ? 'a' : 'button')
const atribut = computed(() => {
  if (tag.value === RouterLink) return { to: props.to }
  if (tag.value === 'a') return { href: props.href, target: '_blank', rel: 'noopener noreferrer' }
  return { type: props.type, disabled: nonaktif.value }
})
const kelas = computed(() => [
  'inline-flex shrink-0 select-none items-center justify-center gap-2 rounded-md font-semibold shadow-[0_1px_2px_rgba(28,36,48,0.05)]',
  'transition-[background-color,border-color,box-shadow,transform] duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
  'disabled:cursor-not-allowed disabled:opacity-55',
  props.wrap ? 'py-2 text-center leading-snug' : 'whitespace-nowrap',
  UKURAN[props.size], VARIAN[props.variant], props.block && 'w-full',
])
</script>

<template>
  <component :is="tag" v-bind="atribut" :class="kelas" :data-variant="variant" :aria-busy="loading ? 'true' : undefined" data-tombol>
    <LoaderCircle v-if="loading" class="size-4 animate-spin" aria-hidden="true" />
    <slot v-else name="icon" />
    <slot />
  </component>
</template>
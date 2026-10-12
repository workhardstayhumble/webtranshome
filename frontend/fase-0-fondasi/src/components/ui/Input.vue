<script setup>
import { computed, useId, useSlots } from 'vue'

// Atribut lain (placeholder, autocomplete, inputmode, ...) diteruskan ke <input>.
defineOptions({ inheritAttrs: false })

const model = defineModel({ type: [String, Number], default: '' })

const props = defineProps({
  label: { type: String, required: true },
  type: { type: String, default: 'text' },
  error: { type: String, default: '' },
  hint: { type: String, default: '' },
  required: { type: Boolean, default: false },
  hideLabel: { type: Boolean, default: false },
  // Teks kecil di kiri isian, mis. "Rp".
  prefix: { type: String, default: '' },
})

const slots = useSlots()
const id = useId()
const keterangan = computed(() => {
  if (props.error) return `${id}-error`
  if (props.hint) return `${id}-hint`
  return undefined
})
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <label :for="id" :class="hideLabel ? 'sr-only' : 'text-sm font-semibold'">
      {{ label }}<span v-if="required" class="text-danger" aria-hidden="true"> *</span>
    </label>
    <div class="relative">
      <span v-if="prefix" class="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-sm text-muted" aria-hidden="true">{{ prefix }}</span>
      <input
        :id="id"
        v-model="model"
        :type="type"
        :required="required"
        :aria-invalid="error ? 'true' : undefined"
        :aria-describedby="keterangan"
        v-bind="$attrs"
        class="h-11 w-full min-w-0 rounded-md border bg-surface text-base text-foreground shadow-[0_1px_2px_rgba(28,36,48,0.02)] transition-colors duration-150 placeholder:text-faint focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/15 disabled:cursor-not-allowed disabled:bg-subtle"
        :class="[error ? 'border-danger' : 'border-border-strong focus-visible:border-primary', prefix ? (prefix.length > 3 ? 'pl-19' : 'pl-10') : 'pl-3', slots.akhir ? 'pr-12' : 'pr-3']"
      />
      <!-- Slot `akhir`: tombol/ikon di sisi kanan isian, mis. tampilkan password atau tanda SKU tersedia. -->
      <div v-if="slots.akhir" class="absolute inset-y-0 right-0 flex items-center pr-1">
        <slot name="akhir" />
      </div>
    </div>
    <p v-if="error" :id="`${id}-error`" class="text-sm font-medium text-danger-ink">{{ error }}</p>
    <p v-else-if="hint" :id="`${id}-hint`" class="text-sm text-muted">{{ hint }}</p>
    <slot name="bawah" />
  </div>
</template>
<script setup>
import { Inbox } from 'lucide-vue-next'
defineProps({ title: { type: String, required: true }, description: { type: String, default: '' } })
</script>
<template>
  <div class="flex flex-col items-center gap-3 rounded-xl border border-border bg-surface px-6 py-12 text-center shadow-[0_2px_16px_rgba(28,36,48,0.03)]">
    <div class="flex size-14 items-center justify-center rounded-xl bg-primary-soft text-primary"><slot name="icon"><Inbox class="size-7" aria-hidden="true" /></slot></div>
    <p class="text-lg font-bold">{{ title }}</p>
    <p v-if="description" class="max-w-md text-sm leading-6 text-muted">{{ description }}</p>
    <div v-if="$slots.actions" class="mt-2 flex flex-wrap justify-center gap-2"><slot name="actions" /></div>
  </div>
</template>
<script setup>
import { RouterLink } from 'vue-router'
import { ArrowUpRight } from 'lucide-vue-next'
defineProps({
  id: { type: String, required: true }, eyebrow: { type: String, default: '' },
  title: { type: String, required: true }, description: { type: String, default: '' },
  to: { type: String, default: '' }, linkLabel: { type: String, default: 'Lihat semua' },
})
</script>
<template>
  <div class="flex flex-wrap items-end justify-between gap-4">
    <div class="flex min-w-0 flex-col gap-2">
      <span v-if="eyebrow" class="th-section-label">{{ eyebrow }}</span>
      <h2 :id="id" class="th-section-heading">{{ title }}</h2>
      <p v-if="description" class="th-section-copy">{{ description }}</p>
    </div>
    <RouterLink v-if="to" :to="to" class="inline-flex min-h-11 items-center gap-1.5 rounded-sm text-sm font-bold text-primary hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
      {{ linkLabel }} <ArrowUpRight class="size-4" aria-hidden="true" />
    </RouterLink>
  </div>
</template>
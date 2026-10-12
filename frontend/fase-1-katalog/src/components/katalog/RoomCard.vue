<script setup>
import { RouterLink } from 'vue-router'
import { ArrowUpRight, ImageOff } from 'lucide-vue-next'
defineProps({ ruangan: { type: Object, required: true } })
</script>
<template>
  <RouterLink :to="`/ruangan/${ruangan.slug}`" class="group th-surface-card th-hover-lift flex h-full flex-col overflow-hidden focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
    <span class="block aspect-[4/3] overflow-hidden bg-[linear-gradient(145deg,#e9e4d9,#d9d1c1)]">
      <img v-if="ruangan.image_cover" :src="ruangan.image_cover" :alt="`Inspirasi ${ruangan.name}`" loading="lazy" width="400" height="300" class="size-full object-cover transition-transform duration-300 group-hover:scale-[1.035]" />
      <span v-else class="flex size-full flex-col items-center justify-center gap-2 text-muted"><ImageOff class="size-7" aria-hidden="true" /><span class="text-xs">Foto ruangan belum tersedia</span></span>
    </span>
    <span class="flex items-center justify-between gap-3 p-4 sm:p-5">
      <span class="flex min-w-0 flex-col gap-1"><span class="text-base font-bold">{{ ruangan.name }}</span><span v-if="ruangan.product_count != null" class="text-sm text-muted">{{ ruangan.product_count }} produk terkait</span></span>
      <span class="flex size-9 shrink-0 items-center justify-center rounded-full bg-primary-soft text-primary"><ArrowUpRight class="size-4" aria-hidden="true" /></span>
    </span>
  </RouterLink>
</template>
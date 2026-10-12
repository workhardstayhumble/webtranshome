<script setup>
import { RouterLink } from 'vue-router'
import { Newspaper } from 'lucide-vue-next'
import { formatTanggal } from '@/utils/format'
import { TIPE_ARTIKEL, tautanArtikel } from '@/utils/artikel'
defineProps({ artikel: { type: Object, required: true }, ringkas: { type: Boolean, default: false } })
</script>
<template>
  <article class="group relative th-surface-card th-hover-lift flex h-full flex-col overflow-hidden">
    <span class="flex aspect-video items-center justify-center overflow-hidden bg-[linear-gradient(140deg,#f0ece4,#dcd4c5)] text-muted">
      <img v-if="artikel.thumbnail_url" :src="artikel.thumbnail_url" :alt="artikel.title" width="800" height="450" loading="lazy" class="size-full object-cover transition-transform duration-300 group-hover:scale-[1.035]" />
      <span v-else class="flex flex-col items-center gap-2"><Newspaper class="size-7" aria-hidden="true" /><span class="text-xs">Gambar belum tersedia</span></span>
    </span>
    <div class="flex flex-1 flex-col gap-2 p-4 sm:p-5">
      <span class="self-start rounded-full px-2.5 py-1 text-xs font-bold" :class="TIPE_ARTIKEL[artikel.type]?.kelas">{{ TIPE_ARTIKEL[artikel.type]?.label ?? artikel.type }}</span>
      <h3 class="text-base leading-snug font-bold"><RouterLink :to="tautanArtikel(artikel.slug)" class="rounded-sm after:absolute after:inset-0 after:content-[''] hover:text-primary focus-visible:outline-none focus-visible:after:rounded-xl focus-visible:after:ring-2 focus-visible:after:ring-ring">{{ artikel.title }}</RouterLink></h3>
      <time v-if="artikel.published_at" :datetime="artikel.published_at" class="text-sm text-muted">{{ formatTanggal(artikel.published_at) }}</time>
      <p v-if="!ringkas && artikel.excerpt" class="line-clamp-3 text-sm leading-6 text-muted">{{ artikel.excerpt }}</p>
    </div>
  </article>
</template>
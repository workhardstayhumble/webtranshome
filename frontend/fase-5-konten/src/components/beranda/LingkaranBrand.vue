<script setup>
import { RouterLink } from 'vue-router'

// Satu lingkaran di "Brand pilihan" Beranda. { brand, warna } → tautan ke /brand/:slug berisi logo dari admin
// (atau nama brand bila logo belum diunggah). { hiasan: [warna1, warna2, sudut] } → lingkaran dua warna datar,
// pengisi tempat kosong bila brand lebih sedikit dari jumlah lingkaran. Diameter dari variabel CSS --d induknya.
// Logo memakai mix-blend-multiply: latar putih logo melebur ke warna lingkaran, jadi tidak tampak sebagai kotak putih.
defineProps({
  item: { type: Object, required: true },
})
</script>

<template>
  <RouterLink
    v-if="item.brand"
    :to="`/brand/${item.brand.slug}`"
    class="flex size-[var(--d)] shrink-0 items-center justify-center overflow-hidden rounded-full p-[16%] text-center ring-2 ring-transparent ring-offset-2 ring-offset-surface transition-shadow duration-150 hover:ring-primary focus-visible:ring-ring focus-visible:outline-none"
    :style="{ background: item.warna }"
    data-lingkaran-brand
  >
    <img v-if="item.brand.logo_url" :src="item.brand.logo_url" :alt="item.brand.name" width="96" height="96" loading="lazy" class="max-h-full max-w-full object-contain mix-blend-multiply" />
    <span v-else class="text-xs leading-tight font-extrabold break-words text-foreground sm:text-sm xl:text-base">{{ item.brand.name }}</span>
  </RouterLink>
  <span
    v-else
    class="block size-[var(--d)] shrink-0 rounded-full"
    :style="{ background: `linear-gradient(${item.hiasan[2]}deg, ${item.hiasan[0]} 50%, ${item.hiasan[1]} 50%)` }"
    aria-hidden="true"
    data-lingkaran-hiasan
  />
</template>

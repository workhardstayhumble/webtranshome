<script setup>
import lambang from '@/assets/logo-transhome.png'

// Lambang resmi (rumah + H, dari docs/wireframe-low-fidelity-transhome.xlsx sheet "Logo Transhome") + tulisan Transhome.
// Lambang dekoratif (alt kosong): nama toko sudah terbaca dari tulisan di sebelahnya.
defineProps({
  // true di latar gelap (footer, sidebar admin): "home" memakai warna oranye terang.
  gelap: { type: Boolean, default: false },
  ukuran: { type: String, default: 'size-9' },
})
</script>

<template>
  <span class="inline-flex items-center gap-2">
    <img :src="lambang" alt="" width="36" height="36" :class="[ukuran, 'shrink-0 rounded-md']" />
    <span>Trans<span :class="gelap ? 'text-primary-on-dark' : 'text-primary'">home</span></span>
  </span>
</template>

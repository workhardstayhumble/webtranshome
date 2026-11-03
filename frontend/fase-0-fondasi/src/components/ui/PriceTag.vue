<script setup>
import { computed } from 'vue'
import { formatRupiah } from '@/utils/format'

const props = defineProps({
  // products.price_general (DECIMAL dari API berupa string). null = harga tidak ditampilkan.
  price: { type: [Number, String], default: null },
  // products.show_price: false = admin menyembunyikan harga produk ini.
  show: { type: Boolean, default: true },
  // products.unit_sale, contoh "dus", "sak"
  unit: { type: String, default: '' },
  // `card`: kecil di kartu 2 kolom HP, normal mulai 640px.
  size: { type: String, default: 'md', validator: (v) => ['sm', 'md', 'lg', 'card'].includes(v) },
})

const UKURAN = { sm: 'text-lg', md: 'text-xl', lg: 'text-3xl', card: 'text-base sm:text-xl' }
// Pengganti harga: lebih kecil dari angka harga supaya tidak terlihat seperti harga.
const UKURAN_GANTI = { sm: 'text-base', md: 'text-base', lg: 'text-xl', card: 'text-sm sm:text-base' }
const tampil = computed(() => props.show && props.price != null && props.price !== '')
</script>

<template>
  <span v-if="!tampil" class="inline-flex items-baseline font-bold text-muted" :class="UKURAN_GANTI[size]" data-harga-tersembunyi>Harga via WhatsApp</span>
  <span v-else class="inline-flex flex-wrap items-baseline gap-x-1">
    <span class="font-extrabold tabular-nums" :class="UKURAN[size]">{{ formatRupiah(price) }}</span>
    <span v-if="unit" class="text-sm text-muted"><span class="sr-only">per </span><span aria-hidden="true">/</span>{{ unit }}</span>
  </span>
</template>

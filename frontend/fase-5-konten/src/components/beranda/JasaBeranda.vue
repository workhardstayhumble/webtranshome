<script setup>
import SkeletonCard from '@/components/ui/SkeletonCard.vue'
import KartuJasa from '@/components/jasa/KartuJasa.vue'
import { useMuat } from '@/composables/useMuat'
import { listServices } from '@/services/kontenService'
import JudulSeksi from './JudulSeksi.vue'
import PesanGagal from './PesanGagal.vue'

// Seksi "Jasa" di Beranda (KEPUTUSAN #20): 3 jasa teratas menurut urutan di admin + tautan ke /jasa.
// Tanpa jasa aktif, seksi disembunyikan.
const jasa = useMuat(() => listServices({ limit: 3 }))
</script>

<template>
  <section v-if="jasa.status.value !== 'ready' || jasa.data.value.length" aria-labelledby="judul-jasa-beranda" class="flex flex-col gap-4" data-jasa-beranda>
    <JudulSeksi id="judul-jasa-beranda" judul="Butuh tukang? Pakai jasa Transhome" tautan="/jasa" label-tautan="Lihat semua jasa" />
    <p class="-mt-2 max-w-2xl text-muted">Beli bahannya di sini, pemasangannya juga bisa. Konsultasikan dulu kebutuhanmu lewat WhatsApp.</p>
    <div v-if="jasa.status.value === 'loading'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3" aria-busy="true">
      <span class="sr-only">Memuat jasa…</span>
      <SkeletonCard v-for="n in 3" :key="n" />
    </div>
    <PesanGagal v-else-if="jasa.status.value === 'error'" pesan="Daftar jasa belum bisa dimuat." @coba-lagi="jasa.muat" />
    <ul v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <li v-for="j in jasa.data.value" :key="j.id"><KartuJasa :jasa="j" /></li>
    </ul>
  </section>
</template>

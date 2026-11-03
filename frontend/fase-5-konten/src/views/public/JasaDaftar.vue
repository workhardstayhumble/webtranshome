<script setup>
import { Wrench } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import SkeletonCard from '@/components/ui/SkeletonCard.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import KartuJasa from '@/components/jasa/KartuJasa.vue'
import { useMuat } from '@/composables/useMuat'
import { remahSeo, urlSitus, useSeo } from '@/composables/useSeo'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { listServices } from '@/services/kontenService'

// /jasa (KEPUTUSAN #20): daftar jasa aktif dari admin › Jasa. Konsultasi hanya lewat WhatsApp (KEPUTUSAN #12).
const CARA = [
  { judul: 'Pilih jasa', isi: 'Lihat jasa yang sesuai dengan pekerjaanmu.' },
  { judul: 'Konsultasi lewat WhatsApp', isi: 'Kirim foto, ukuran, dan lokasi pekerjaan ke tim kami.' },
  { judul: 'Atur jadwal dan biaya', isi: 'Tim kami membalas dengan biaya dan jadwal pengerjaan.' },
]

const jasa = useMuat(() => listServices())
const { generalLink } = useWhatsApp()
const REMAH = [{ label: 'Beranda', to: '/' }, { label: 'Jasa', to: '/jasa' }]

useSeo({
  title: 'Jasa pasang & renovasi rumah',
  description: 'Jasa Transhome di Sleman: pasang keramik, pengecatan, perlengkapan kamar mandi, AC & elektronik, service, pintu & jendela, rakit perabot. Konsultasi lewat WhatsApp.',
  canonical: '/jasa',
  jsonLd: () =>
    jasa.status.value === 'ready' && [
      remahSeo(REMAH),
      {
        '@type': 'ItemList',
        name: 'Jasa Transhome',
        itemListElement: jasa.data.value.map((j, i) => ({ '@type': 'ListItem', position: i + 1, url: urlSitus(`/jasa/${j.slug}`), name: j.name })),
      },
    ],
})
</script>

<template>
  <div class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:py-6">
    <Breadcrumb :items="REMAH.slice(0, 1).concat({ label: 'Jasa' })" />
    <h1 class="mt-2 text-2xl font-extrabold sm:text-3xl">Jasa Transhome</h1>
    <p class="mt-2 max-w-2xl text-muted">Beli bahannya di Transhome, pemasangannya juga bisa. Pilih jasa yang kamu butuhkan, lalu konsultasikan lewat WhatsApp.</p>

    <section aria-labelledby="judul-cara-jasa" class="mt-6 rounded-xl bg-primary-soft p-4 sm:p-5">
      <h2 id="judul-cara-jasa" class="text-lg font-extrabold">Cara pesan jasa</h2>
      <ol class="mt-3 grid gap-3 sm:grid-cols-3">
        <li v-for="(c, i) in CARA" :key="c.judul" class="flex gap-3">
          <span class="flex size-8 shrink-0 items-center justify-center rounded-full bg-primary text-sm font-extrabold text-primary-foreground" aria-hidden="true">{{ i + 1 }}</span>
          <span class="flex flex-col">
            <strong>{{ c.judul }}</strong>
            <span class="text-sm text-muted">{{ c.isi }}</span>
          </span>
        </li>
      </ol>
    </section>

    <div v-if="jasa.status.value === 'loading'" class="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3" aria-busy="true">
      <span class="sr-only">Memuat jasa…</span>
      <SkeletonCard v-for="n in 6" :key="n" />
    </div>
    <ErrorState v-else-if="jasa.status.value === 'error'" class="mt-6" :message="jasa.pesan.value" @retry="jasa.muat" />
    <EmptyState v-else-if="!jasa.data.value.length" class="mt-6" title="Daftar jasa sedang disiapkan" description="Tanyakan jasa yang kamu butuhkan langsung ke tim kami.">
      <template #icon><Wrench class="size-10" aria-hidden="true" /></template>
      <template #actions><WhatsAppButton :href="generalLink()" label="Tanya via WhatsApp" /></template>
    </EmptyState>
    <ul v-else class="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3" data-daftar-jasa>
      <li v-for="(j, i) in jasa.data.value" :key="j.id" v-reveal="i % 3"><KartuJasa :jasa="j" tingkat="h2" /></li>
    </ul>

    <div v-if="jasa.status.value === 'ready' && jasa.data.value.length" class="mt-8 flex flex-col items-start gap-3 rounded-xl border border-border bg-surface p-4 sm:flex-row sm:items-center sm:justify-between sm:p-5">
      <p><strong>Jasa yang kamu cari belum ada?</strong> <span class="text-muted">Tanyakan saja, tim kami bantu carikan.</span></p>
      <WhatsAppButton class="shrink-0" :href="generalLink()" label="Tanya via WhatsApp" />
    </div>
  </div>
</template>

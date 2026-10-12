<script setup>
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import NotFoundState from '@/components/ui/NotFoundState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import ChipNav from '@/components/katalog/ChipNav.vue'
import KatalogDaftar from '@/components/katalog/KatalogDaftar.vue'
import { useSeo } from '@/composables/useSeo'
import { cariKategori, getCategoryTree, pohonTersimpan } from '@/services/categoryService'

// /katalog (semua produk) dan /katalog/:kategori (induk termasuk produk anak-anaknya).
const route = useRoute()
// Pohon yang sudah ada dipakai langsung agar render pertama lengkap (posisi scroll saat Back tidak bergeser).
const pohon = ref(pohonTersimpan())
const statusPohon = ref(pohon.value ? 'ready' : 'loading') // loading | ready | error

async function muatPohon() {
  statusPohon.value = 'loading'
  try {
    pohon.value = await getCategoryTree()
    statusPohon.value = 'ready'
  } catch {
    statusPohon.value = 'error'
  }
}
if (!pohon.value) muatPohon()

const slug = computed(() => route.params.kategori || null)
const posisi = computed(() => (slug.value && pohon.value ? cariKategori(pohon.value, slug.value) : null))
const tidakAda = computed(() => Boolean(slug.value) && statusPohon.value === 'ready' && !posisi.value)

// Pohon gagal dimuat: judul diambil dari slug agar halaman tetap bisa dipakai.
const judul = computed(() => {
  if (!slug.value) return 'Semua produk'
  if (posisi.value) return posisi.value.kategori.name
  if (statusPohon.value === 'error') return slug.value.replace(/-/g, ' ').replace(/^./, (h) => h.toUpperCase())
  return null
})

const kunci = computed(() => (slug.value ? { category: slug.value } : {}))

const remah = computed(() => {
  const items = [{ label: 'Beranda', to: '/' }]
  if (!slug.value) return [...items, { label: 'Katalog' }]
  items.push({ label: 'Katalog', to: '/katalog' })
  if (posisi.value?.induk) items.push({ label: posisi.value.induk.name, to: `/katalog/${posisi.value.induk.slug}` })
  items.push({ label: judul.value ?? 'Kategori' })
  return items
})

// Semua produk → kategori induk. Kategori induk/anak → "Semua <induk>" + anak-anaknya.
const navKategori = computed(() => {
  if (!pohon.value) return []
  if (!slug.value) return pohon.value.map((k) => ({ key: k.slug, label: k.name, to: `/katalog/${k.slug}` }))
  const induk = posisi.value?.induk ?? posisi.value?.kategori
  if (!induk?.children?.length) return []
  return [
    { key: induk.slug, label: `Semua ${induk.name}`, to: `/katalog/${induk.slug}`, aktif: slug.value === induk.slug },
    ...induk.children.map((c) => ({ key: c.slug, label: c.name, to: `/katalog/${c.slug}`, aktif: slug.value === c.slug })),
  ]
})

useSeo({
  title: () => {
    if (tidakAda.value) return 'Kategori tidak ditemukan'
    if (!slug.value) return 'Katalog produk'
    return judul.value
  },
  description: () =>
    judul.value && !tidakAda.value
      ? `${slug.value ? judul.value : 'Bahan bangunan dan elektronik'} di Transhome Sleman: cek harga, status stok, dan tanya langsung lewat WhatsApp.`
      : null,
})
</script>

<template>
  <NotFoundState
    v-if="tidakAda"
    title="Kategori tidak ditemukan"
    description="Kategori yang kamu buka tidak ada atau sudah tidak aktif."
    action-label="Lihat semua produk"
    action-to="/katalog"
  />

  <div v-else class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:py-9">
    <Breadcrumb :items="remah" />
    <p class="th-section-label mt-6">Transhome / Katalog produk</p>
    <h1 class="mt-2 text-3xl font-extrabold tracking-tight sm:text-4xl">
      <template v-if="judul">{{ judul }}</template>
      <Skeleton v-else class="inline-block h-8 w-56 align-middle" />
    </h1>
    <p class="mt-2 text-sm text-muted">Jelajahi produk berdasarkan kategori, spesifikasi, dan kebutuhan rumahmu.</p>
    <ChipNav v-if="navKategori.length" class="mt-4" :label="slug ? 'Subkategori' : 'Kategori'" :items="navKategori" />
    <KatalogDaftar class="mt-6" :kunci="kunci" :topik="slug && judul ? judul.toLowerCase() : ''" />
  </div>
</template>
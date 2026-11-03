<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import NotFoundState from '@/components/ui/NotFoundState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import KatalogDaftar from '@/components/katalog/KatalogDaftar.vue'
import { remahSeo, urlSitus, useSeo } from '@/composables/useSeo'
import { getBrand } from '@/services/catalogService'
import { pesanError } from '@/services/errors'

// /brand/:slug : info brand + katalog dengan filter brand terkunci. Daftar produk dimuat bersamaan.
const route = useRoute()
const status = ref('loading') // loading | ready | tidak-ada | error
const brand = ref(null)
const pesan = ref('')

async function muat() {
  status.value = 'loading'
  try {
    brand.value = await getBrand(route.params.slug)
    status.value = 'ready'
  } catch (error) {
    if (error?.response?.status === 404) {
      status.value = 'tidak-ada'
      return
    }
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

watch(
  () => route.params.slug,
  (slug) => {
    if (slug) muat()
  },
  { immediate: true },
)

const kunci = computed(() => ({ brand: route.params.slug }))
const inisial = computed(() => brand.value?.name?.trim().charAt(0).toUpperCase() ?? '')

const remah = computed(() => [{ label: 'Beranda', to: '/' }, { label: 'Brand', to: '/brand' }, { label: brand.value?.name ?? 'Brand' }])

// SEO: halaman brand bersih diindeks (orang mencari "<nama brand> Sleman"). Filter/urutan katalog di URL = noindex
// + canonical ke halaman brand; brand tanpa produk dan brand tidak ditemukan juga noindex (halaman tipis).
const difilter = computed(() => Object.keys(route.query).some((k) => k !== 'page'))
const halaman = computed(() => Math.max(Number(route.query.page) || 1, 1))
useSeo({
  title: () => {
    if (status.value === 'tidak-ada') return 'Brand tidak ditemukan'
    return status.value === 'ready' ? `${brand.value.name}: produk, harga & stok di Sleman` : null
  },
  description: () => {
    if (status.value !== 'ready') return null
    const { name, product_count: jumlah } = brand.value
    return jumlah
      ? `${jumlah} produk ${name} di Transhome Sleman, Yogyakarta. Cek harga dan stok, lalu tanya langsung lewat WhatsApp.`
      : `Produk ${name} di Transhome Sleman, Yogyakarta. Tanyakan ketersediaannya lewat WhatsApp.`
  },
  canonical: () => {
    const dasar = `/brand/${route.params.slug}`
    return !difilter.value && halaman.value > 1 ? `${dasar}?page=${halaman.value}` : dasar
  },
  robots: () => {
    if (status.value === 'tidak-ada') return 'noindex'
    return difilter.value || (status.value === 'ready' && !brand.value.product_count) ? 'noindex, follow' : null
  },
  image: () => brand.value?.logo_url,
  jsonLd: () =>
    status.value === 'ready' && [
      remahSeo(remah.value.map((r, i) => (i === 2 ? { ...r, to: `/brand/${brand.value.slug}` } : r))),
      {
        '@type': 'CollectionPage',
        name: `Produk ${brand.value.name}`,
        url: urlSitus(`/brand/${brand.value.slug}`),
        about: { '@type': 'Brand', name: brand.value.name, ...(brand.value.logo_url ? { logo: urlSitus(brand.value.logo_url) } : {}) },
      },
    ],
})
</script>

<template>
  <NotFoundState
    v-if="status === 'tidak-ada'"
    title="Brand tidak ditemukan"
    description="Brand yang kamu buka tidak ada atau sudah tidak aktif."
    action-label="Lihat katalog"
    action-to="/katalog"
  />

  <div v-else class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:py-6">
    <Breadcrumb :items="remah" />

    <div class="mt-3 flex items-center gap-4">
      <div class="flex size-16 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-border bg-surface">
        <img v-if="brand?.logo_url" :src="brand.logo_url" :alt="`Logo ${brand.name}`" width="64" height="64" class="size-full object-contain p-2" />
        <span v-else-if="brand" class="text-2xl font-extrabold text-muted" aria-hidden="true">{{ inisial }}</span>
      </div>
      <div class="flex min-w-0 flex-col gap-1">
        <h1 class="text-2xl font-extrabold sm:text-3xl">
          <template v-if="brand">{{ brand.name }}</template>
          <Skeleton v-else class="inline-block h-8 w-40 align-middle" />
        </h1>
        <p v-if="brand" class="text-sm text-muted">{{ brand.product_count }} produk di Transhome</p>
        <Skeleton v-else class="h-5 w-36" />
      </div>
    </div>

    <ErrorState v-if="status === 'error'" class="mt-6" :message="pesan" @retry="muat" />
    <KatalogDaftar v-else class="mt-6" :kunci="kunci" :sembunyikan="['brand']" :topik="brand ? `produk ${brand.name}` : ''" />
  </div>
</template>

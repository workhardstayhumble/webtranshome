<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import NotFoundState from '@/components/ui/NotFoundState.vue'
import PriceTag from '@/components/ui/PriceTag.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import ProductGallery from '@/components/produk/ProductGallery.vue'
import ProductSummary from '@/components/produk/ProductSummary.vue'
import ProductTabs from '@/components/produk/ProductTabs.vue'
import RelatedProducts from '@/components/produk/RelatedProducts.vue'
import { useSeo } from '@/composables/useSeo'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { getProduct } from '@/services/productService'
import { pesanError } from '@/services/errors'
import { formatRupiah } from '@/utils/format'
import { hargaDisembunyikan, labelTanyaProduk } from '@/utils/whatsapp'

const route = useRoute()
const { productLink } = useWhatsApp()

const status = ref('loading') // loading | ready | tidak-ada | error
const produk = ref(null)
const pesan = ref('')
let pengendali = null

async function muat() {
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    produk.value = await getProduct(route.params.slug, { signal: pengendali.signal })
    status.value = 'ready'
  } catch (error) {
    if (error?.code === 'ERR_CANCELED') return
    // 404 juga untuk produk nonaktif (FR-A).
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
onBeforeUnmount(() => pengendali?.abort())


const remah = computed(() => {
  const items = [
    { label: 'Beranda', to: '/' },
    { label: 'Katalog', to: '/katalog' },
  ]
  const kategori = produk.value?.category
  if (kategori?.parent) items.push({ label: kategori.parent.name, to: `/katalog/${kategori.parent.slug}` })
  if (kategori) items.push({ label: kategori.name, to: `/katalog/${kategori.slug}` })
  items.push({ label: produk.value?.name ?? 'Produk' })
  return items
})

function deskripsiOtomatis(p) {
  const merek = p.brand ? ` dari ${p.brand.name}` : ''
  const harga = hargaDisembunyikan(p) ? 'Tanya harga' : `Harga ${formatRupiah(p.price_general)} per ${p.unit_sale}`
  return `${p.name}${merek}. ${harga}. Cek stok dan tanya langsung lewat WhatsApp Transhome.`
}

// Selama memuat, judul dibiarkan dari router (tidak memakai produk sebelumnya).
useSeo({
  title: () => {
    if (status.value === 'tidak-ada') return 'Produk tidak ditemukan'
    if (status.value !== 'ready') return null
    return produk.value.meta_title || produk.value.name
  },
  description: () => (status.value === 'ready' ? produk.value.meta_description || deskripsiOtomatis(produk.value) : null),
})
</script>

<template>
  <NotFoundState
    v-if="status === 'tidak-ada'"
    title="Produk tidak ditemukan"
    description="Produk yang kamu cari tidak ada atau sudah tidak dijual lagi."
    action-label="Kembali ke katalog"
    action-to="/katalog"
  />

  <div v-else class="mx-auto max-w-7xl px-4 pt-4 pb-12 sm:px-6 lg:pt-6 lg:pb-16">
    <div v-if="status === 'loading'" aria-busy="true">
      <span class="sr-only">Memuat produk…</span>
      <Skeleton class="h-5 w-64 max-w-full" />
      <div class="mt-4 grid gap-8 lg:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)] lg:gap-12">
        <Skeleton class="aspect-square rounded-xl" />
        <div class="flex flex-col gap-4">
          <Skeleton class="h-4 w-24" />
          <Skeleton class="h-9 w-4/5" />
          <Skeleton class="h-5 w-40" />
          <Skeleton class="h-10 w-48" />
          <Skeleton class="h-7 w-32" />
          <Skeleton class="h-14 w-full" />
          <Skeleton class="h-24 w-full rounded-xl" />
        </div>
      </div>
    </div>

    <ErrorState v-else-if="status === 'error'" class="mt-6" :message="pesan" @retry="muat" />

    <template v-else-if="produk">
      <Breadcrumb :items="remah" />
      <div class="mt-4 grid items-start gap-8 lg:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)] lg:gap-12">
        <ProductGallery :images="produk.images ?? []" :nama="produk.name" />
        <ProductSummary :produk="produk" />
      </div>
      <ProductTabs class="mt-10" :produk="produk" />
      <RelatedProducts class="mt-12" :produk="produk" />

      <!-- HP: harga + WhatsApp selalu terjangkau. Tombol WhatsApp melayang disembunyikan di halaman ini (meta.barAksiHp). -->
      <div class="fixed inset-x-0 bottom-0 z-40 border-t border-border bg-surface px-4 pt-3 pb-[calc(0.75rem+env(safe-area-inset-bottom))] lg:hidden">
        <div class="mx-auto flex max-w-7xl items-center gap-3">
          <PriceTag class="min-w-0 flex-1" size="sm" :price="produk.price_general" :show="produk.show_price !== false" :unit="produk.unit_sale" />
          <WhatsAppButton class="shrink-0" :href="productLink(produk)" :label="labelTanyaProduk(produk)" />
        </div>
      </div>
    </template>
  </div>
</template>

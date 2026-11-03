<script setup>
import { computed, defineAsyncComponent, ref, useId, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { SearchX } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import SkeletonCard from '@/components/ui/SkeletonCard.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import ActiveFilterChips from './ActiveFilterChips.vue'
import FilterPanel from './FilterPanel.vue'
import KatalogToolbar from './KatalogToolbar.vue'
import ProductCard from '@/components/katalog/ProductCard.vue'
import ProductListRow from '@/components/katalog/ProductListRow.vue'
import { useKatalog } from '@/composables/useKatalog'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { perilakuGulir } from '@/utils/gerak'

// Bottom sheet filter HP baru diunduh saat pertama kali dibuka.
const FilterSheet = defineAsyncComponent(() => import('./FilterSheet.vue'))

const props = defineProps({
  // Parameter API yang dikunci halaman: { category } | { q } | { brand } | { room }.
  kunci: { type: Object, default: () => ({}) },
  // Kunci URL yang tetap ada saat filter direset, mis. ['q'] di /cari.
  kunciUrl: { type: Array, default: () => [] },
  // Grup filter yang disembunyikan karena dikunci halaman, mis. ['brand'] di /brand/:slug.
  sembunyikan: { type: Array, default: () => [] },
  // Untuk pesan WhatsApp saat hasil kosong, mis. kata kunci atau nama kategori.
  topik: { type: String, default: '' },
  judulKosong: { type: String, default: 'Produk tidak ditemukan' },
  deskripsiKosong: { type: String, default: 'Belum ada produk di sini. Tanya tim kami, mungkin barangnya tersedia di toko.' },
})

const { status, produk, meta, facets, pesan, chip, filter, halaman, urutan, toggle, terapkan, gantiFilter, resetFilter, muatUlang } =
  useKatalog({ kunci: () => props.kunci, kunciUrl: props.kunciUrl })
const { searchLink } = useWhatsApp()

const idFilter = useId()
const idHasil = useId()
const atasDaftar = ref(null)
const judulHasil = ref(null)
const sheetBuka = ref(false)
const sheetDimuat = ref(false)

// Pilihan Grid/Daftar disimpan per browser; localStorage bisa diblokir, jadi selalu lewat try/catch.
const KUNCI_TAMPILAN = 'transhome.tampilanKatalog'
function bacaTampilan() {
  try {
    return localStorage.getItem(KUNCI_TAMPILAN) === 'daftar' ? 'daftar' : 'grid'
  } catch {
    return 'grid'
  }
}
const tampilan = ref(bacaTampilan())
watch(tampilan, (nilai) => {
  try {
    localStorage.setItem(KUNCI_TAMPILAN, nilai)
  } catch {
    // Tetap berganti tampilan, hanya tidak diingat.
  }
})

const total = computed(() => meta.value?.total ?? null)
const totalHalaman = computed(() => meta.value?.total_pages ?? 1)
// Jumlah kerangka mengikuti daftar sebelumnya supaya tinggi halaman tidak melompat.
const jumlahKerangka = computed(() => (produk.value.length ? Math.min(Math.max(produk.value.length, 4), 12) : 8))

function ubahUrutan(nilai) {
  terapkan({ sort: nilai === 'terbaru' ? null : nilai })
}

function bukaFilter() {
  sheetDimuat.value = true
  sheetBuka.value = true
}

function fokusHasil() {
  judulHasil.value?.focus({ preventScroll: true })
}

// Nomor halaman berupa <a href> (bisa diikuti mesin pencari); klik biasa tetap lewat gantiHalaman.
const route = useRoute()
const router = useRouter()
const tautanHalaman = (n) => router.resolve({ query: { ...route.query, page: n > 1 ? n : undefined } }).href

async function gantiHalaman(nomor) {
  await terapkan({ page: nomor > 1 ? nomor : null })
  atasDaftar.value?.scrollIntoView({ behavior: perilakuGulir(), block: 'start' })
  fokusHasil()
}
</script>

<template>
  <div class="grid gap-8 lg:grid-cols-[16rem_minmax(0,1fr)]">
    <aside class="hidden lg:block" :aria-labelledby="idFilter">
      <div class="mb-4 flex min-h-9 items-center justify-between gap-2">
        <h2 :id="idFilter" class="text-lg font-bold">Filter</h2>
        <Button v-if="chip.length" variant="ghost" size="sm" @click="resetFilter">Reset</Button>
      </div>
      <p v-if="!facets && status === 'error'" class="text-sm text-muted">Filter tampil setelah produk berhasil dimuat.</p>
      <FilterPanel v-else :nilai="filter" :facets="facets" :sembunyikan="sembunyikan" @toggle="toggle" @harga="terapkan" />
    </aside>

    <div ref="atasDaftar" class="flex min-w-0 scroll-mt-4 flex-col gap-4">
      <KatalogToolbar
        :urutan="urutan"
        :total="total"
        :memuat="status === 'loading'"
        :jumlah-filter="chip.length"
        v-model:tampilan="tampilan"
        @update:urutan="ubahUrutan"
        @buka-filter="bukaFilter"
      />
      <ActiveFilterChips :chip="chip" @hapus-semua="resetFilter" @kosong="fokusHasil" />

      <section :aria-labelledby="idHasil" :aria-busy="status === 'loading' ? 'true' : 'false'">
        <h2 :id="idHasil" ref="judulHasil" tabindex="-1" class="sr-only">Daftar produk</h2>

        <template v-if="status === 'loading'">
          <span class="sr-only">Memuat produk…</span>
          <div v-if="tampilan === 'grid'" class="grid grid-cols-2 gap-3 sm:gap-4 md:grid-cols-3 xl:grid-cols-4">
            <SkeletonCard v-for="n in jumlahKerangka" :key="n" />
          </div>
          <div v-else class="flex flex-col gap-3" aria-hidden="true">
            <div v-for="n in jumlahKerangka" :key="n" class="flex gap-4 rounded-xl border border-border bg-surface p-3">
              <Skeleton class="size-28 shrink-0 rounded-lg" />
              <div class="flex flex-1 flex-col gap-2 py-1">
                <Skeleton class="h-3 w-16" />
                <Skeleton class="h-4 w-3/4" />
                <Skeleton class="h-6 w-28" />
              </div>
            </div>
          </div>
        </template>

        <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muatUlang" />

        <EmptyState
          v-else-if="produk.length === 0"
          :title="judulKosong"
          :description="chip.length ? 'Tidak ada produk yang cocok dengan semua filter ini. Coba kurangi filter.' : deskripsiKosong"
        >
          <template #icon><SearchX class="size-10" aria-hidden="true" /></template>
          <template #actions>
            <Button v-if="chip.length" variant="outline" @click="resetFilter">Reset filter</Button>
            <WhatsAppButton :href="searchLink(topik)" label="Tanya via WhatsApp" />
          </template>
        </EmptyState>

        <ul v-else-if="tampilan === 'grid'" class="grid grid-cols-2 gap-3 sm:gap-4 md:grid-cols-3 xl:grid-cols-4">
          <!-- Kartu di bawah layar muncul bergantian saat digulir (v-reveal, per baris 4 kartu). -->
          <li v-for="(p, i) in produk" :key="p.id" v-reveal="i % 4" class="flex">
            <ProductCard class="w-full" :product="p" />
          </li>
        </ul>
        <ul v-else class="flex flex-col gap-3">
          <li v-for="(p, i) in produk" :key="p.id" v-reveal="i % 4">
            <ProductListRow :product="p" />
          </li>
        </ul>
      </section>

      <PaginationBar v-if="status === 'ready' && totalHalaman > 1" class="mt-4" :page="halaman" :total-pages="totalHalaman" :tautan="tautanHalaman" @update:page="gantiHalaman" />
    </div>

    <FilterSheet
      v-if="sheetDimuat"
      v-model:open="sheetBuka"
      :nilai="filter"
      :facets="facets"
      :total="total"
      :kunci="kunci"
      :sembunyikan="sembunyikan"
      @terapkan="gantiFilter"
    />
  </div>
</template>

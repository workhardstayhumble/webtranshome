<script setup>
import { computed, onBeforeUnmount, ref, useId, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search } from 'lucide-vue-next'
import { suggestProducts } from '@/services/productService'

defineProps({
  placeholder: { type: String, default: 'Cari SKU, nama produk, atau merek' },
})

const route = useRoute()
const router = useRouter()
const teks = ref('')
const terbuka = ref(false)
const status = ref('idle') // idle | loading | ready | error
const saran = ref({ products: [], categories: [], brands: [] })
const aktif = ref(-1)
const idInput = useId()
const idDaftar = useId()
let timer = null
let pengendali = null

const GRUP = [
  { kunci: 'products', nama: 'Produk', jenis: 'produk', ke: (x) => `/produk/${x.slug}`, detail: (x) => x.sku },
  { kunci: 'categories', nama: 'Kategori', jenis: 'kategori', ke: (x) => `/katalog/${x.slug}` },
  { kunci: 'brands', nama: 'Brand', jenis: 'brand', ke: (x) => `/brand/${x.slug}` },
]

// Satu daftar datar untuk navigasi ↑↓, tetap dikelompokkan saat ditampilkan.
const opsi = computed(() =>
  GRUP.flatMap((g) =>
    (saran.value[g.kunci] ?? []).map((x) => ({ id: `${g.kunci}-${x.id}`, grup: g.nama, jenis: g.jenis, label: x.name, detail: g.detail?.(x), ke: g.ke(x) })),
  ).map((o, indeks) => ({ ...o, indeks })),
)
const kelompok = computed(() => GRUP.map((g) => ({ nama: g.nama, item: opsi.value.filter((o) => o.grup === g.nama) })).filter((g) => g.item.length))

let dariUrl = false

watch(teks, (nilai) => {
  if (dariUrl) {
    dariUrl = false
    return
  }
  clearTimeout(timer)
  pengendali?.abort()
  aktif.value = -1
  const q = nilai.trim()
  if (q.length < 2) {
    terbuka.value = false
    status.value = 'idle'
    return
  }
  timer = setTimeout(async () => {
    pengendali = new AbortController()
    status.value = 'loading'
    terbuka.value = true
    try {
      saran.value = await suggestProducts(q, { signal: pengendali.signal })
      status.value = 'ready'
    } catch (error) {
      if (error?.code !== 'ERR_CANCELED') status.value = 'error'
    }
  }, 300)
})

// Di /cari kotak ini menampilkan kata kunci dari URL (juga setelah refresh) tanpa membuka saran.
// Dipasang setelah watcher `teks` supaya penanda `dariUrl` langsung dipakai dan direset olehnya.
watch(
  () => (route.path === '/cari' ? String(route.query.q ?? '') : null),
  (q) => {
    if (q === null || q === teks.value) return
    dariUrl = true
    teks.value = q
  },
  { immediate: true },
)

function pilih(item) {
  terbuka.value = false
  aktif.value = -1
  router.push(item.ke)
}

function kirim() {
  const terpilih = opsi.value[aktif.value]
  if (terbuka.value && terpilih) return pilih(terpilih)
  const q = teks.value.trim()
  if (!q) return
  terbuka.value = false
  router.push({ path: '/cari', query: { q } })
}

function onKeydown(event) {
  if (event.key === 'ArrowDown') {
    event.preventDefault()
    if (!opsi.value.length) return
    terbuka.value = true
    aktif.value = Math.min(aktif.value + 1, opsi.value.length - 1)
  } else if (event.key === 'ArrowUp') {
    event.preventDefault()
    aktif.value = Math.max(aktif.value - 1, -1)
  } else if (event.key === 'Escape' && terbuka.value) {
    event.preventDefault()
    terbuka.value = false
    aktif.value = -1
  }
}

function onFocusOut(event) {
  if (!event.currentTarget.contains(event.relatedTarget)) terbuka.value = false
}

onBeforeUnmount(() => {
  clearTimeout(timer)
  pengendali?.abort()
})
</script>

<template>
  <form role="search" class="relative min-w-0" @submit.prevent="kirim" @focusout="onFocusOut">
    <label :for="idInput" class="sr-only">Cari produk</label>
    <div class="flex h-11 items-center gap-1 rounded-md border border-border-strong bg-background pr-1 pl-3 transition-colors duration-150 focus-within:border-primary focus-within:ring-2 focus-within:ring-ring">
      <input
        :id="idInput"
        v-model="teks"
        type="search"
        role="combobox"
        aria-autocomplete="list"
        :aria-expanded="terbuka ? 'true' : 'false'"
        :aria-controls="idDaftar"
        :aria-activedescendant="aktif >= 0 && opsi[aktif] ? `${idDaftar}-${opsi[aktif].id}` : undefined"
        :placeholder="placeholder"
        autocomplete="off"
        enterkeyhint="search"
        class="h-full min-w-0 flex-1 bg-transparent text-base outline-none placeholder:text-faint"
        @keydown="onKeydown"
        @focus="teks.trim().length >= 2 && (terbuka = true)"
      />
      <button
        type="submit"
        aria-label="Cari"
        class="inline-flex size-9 shrink-0 items-center justify-center rounded-md text-muted transition-colors duration-150 hover:bg-subtle hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        <Search class="size-4.5" aria-hidden="true" />
      </button>
    </div>

    <div v-show="terbuka" class="absolute inset-x-0 top-full z-40 mt-1 overflow-hidden rounded-xl border border-border bg-surface shadow-lg">
      <ul :id="idDaftar" role="listbox" aria-label="Saran pencarian" class="max-h-80 overflow-y-auto py-2">
        <li v-if="status === 'loading'" role="presentation" class="px-4 py-3 text-sm text-muted">Mencari…</li>
        <li v-else-if="status === 'error'" role="presentation" class="px-4 py-3 text-sm text-muted">Saran gagal dimuat. Tekan Enter untuk mencari.</li>
        <li v-else-if="!opsi.length" role="presentation" class="px-4 py-3 text-sm text-muted">Tidak ada saran</li>
        <template v-else>
          <template v-for="grup in kelompok" :key="grup.nama">
            <li role="presentation" class="px-4 pt-2 pb-1 text-xs font-semibold tracking-wide text-faint uppercase">{{ grup.nama }}</li>
            <li
              v-for="item in grup.item"
              :id="`${idDaftar}-${item.id}`"
              :key="item.id"
              role="option"
              :aria-selected="item.indeks === aktif ? 'true' : 'false'"
              class="flex cursor-pointer items-center justify-between gap-3 px-4 py-2.5 text-sm"
              :class="item.indeks === aktif ? 'bg-primary-soft' : 'hover:bg-subtle'"
              @mousedown.prevent="pilih(item)"
              @mousemove="aktif = item.indeks"
            >
              <span class="truncate">{{ item.label }}<span class="sr-only">, {{ item.jenis }}</span></span>
              <span v-if="item.detail" class="shrink-0 font-mono text-xs text-muted">{{ item.detail }}</span>
            </li>
          </template>
        </template>
      </ul>
    </div>
  </form>
</template>
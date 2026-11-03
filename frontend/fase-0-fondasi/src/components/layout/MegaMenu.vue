<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, useId, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ChevronDown, Menu } from 'lucide-vue-next'
import { getCategoryTree } from '@/services/categoryService'
import Button from '@/components/ui/Button.vue'
import Skeleton from '@/components/ui/Skeleton.vue'

// Tautan tambahan yang di HP tidak muat di header.
const NAVIGASI_HP = [
  { label: 'Ruangan', to: '/ruangan' },
  { label: 'Jasa', to: '/jasa' },
  { label: 'Promo', to: '/promo' },
  { label: 'Artikel', to: '/artikel' },
  { label: 'Info toko', to: '/info-toko' },
]

const route = useRoute()
const terbuka = ref(false)
const status = ref('idle') // idle | loading | ready | error
const kategori = ref([])
const tombol = ref(null)
const panel = ref(null)
const idPanel = useId()

async function muat() {
  status.value = 'loading'
  try {
    kategori.value = await getCategoryTree()
    status.value = 'ready'
  } catch {
    status.value = 'error'
  }
}

async function buka() {
  terbuka.value = true
  if (status.value !== 'ready') await muat()
  await nextTick()
  panel.value?.querySelector('a, button')?.focus()
}

function tutup({ kembalikanFokus = false } = {}) {
  terbuka.value = false
  if (kembalikanFokus) tombol.value?.focus()
}

function onKeydown(event) {
  if (event.key === 'Escape' && terbuka.value) {
    event.preventDefault()
    tutup({ kembalikanFokus: true })
  }
}

function klikDiLuar(event) {
  if (terbuka.value && !event.target.closest?.('[data-mega-menu]')) tutup()
}

watch(() => route.fullPath, () => tutup())
onMounted(() => document.addEventListener('mousedown', klikDiLuar))
onBeforeUnmount(() => document.removeEventListener('mousedown', klikDiLuar))
</script>

<template>
  <!-- Panel diposisikan terhadap <header> (sticky), sehingga selebar layar. Tingginya dibatasi
       sisa layar di bawah header supaya isi panel tetap bisa digulir saat header menempel. -->
  <div data-mega-menu @keydown="onKeydown">
    <button
      ref="tombol"
      type="button"
      :aria-expanded="terbuka ? 'true' : 'false'"
      :aria-controls="idPanel"
      class="inline-flex h-11 items-center gap-2 rounded-md border border-foreground bg-surface px-3.5 text-sm font-semibold transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:px-4"
      @click="terbuka ? tutup() : buka()"
    >
      <Menu class="size-4.5" aria-hidden="true" />
      Kategori
      <ChevronDown class="size-4 transition-transform duration-150" :class="terbuka && 'rotate-180'" aria-hidden="true" />
    </button>

    <div v-show="terbuka" :id="idPanel" ref="panel" class="absolute inset-x-0 top-full z-30 border-t border-border bg-surface shadow-lg">
      <div class="mx-auto max-h-[min(70dvh,calc(100dvh-var(--tinggi-header-penuh,0px)-1rem))] max-w-7xl overflow-y-auto px-4 py-6 sm:px-6">
        <div v-if="status === 'loading'" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4" aria-busy="true">
          <span class="sr-only">Memuat kategori…</span>
          <div v-for="n in 8" :key="n" class="flex flex-col gap-2">
            <Skeleton class="h-5 w-32" />
            <Skeleton class="h-4 w-24" />
            <Skeleton class="h-4 w-28" />
          </div>
        </div>

        <div v-else-if="status === 'error'" class="flex flex-col items-start gap-3">
          <p class="font-semibold">Kategori gagal dimuat.</p>
          <Button variant="outline" size="sm" @click="muat">Coba lagi</Button>
        </div>

        <nav v-else aria-label="Kategori produk" class="grid gap-x-6 gap-y-6 sm:grid-cols-2 lg:grid-cols-4">
          <div v-for="induk in kategori" :key="induk.id" class="flex flex-col gap-1">
            <RouterLink
              :to="`/katalog/${induk.slug}`"
              class="rounded-md py-1 font-bold transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              {{ induk.name }}
            </RouterLink>
            <ul class="flex flex-col">
              <li v-for="anak in induk.children" :key="anak.id">
                <RouterLink
                  :to="`/katalog/${anak.slug}`"
                  class="block rounded-md py-1.5 text-sm text-muted transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                >
                  {{ anak.name }}
                </RouterLink>
              </li>
            </ul>
          </div>
        </nav>

        <nav aria-label="Navigasi lain" class="mt-6 flex flex-wrap gap-2 border-t border-border pt-4 xl:hidden">
          <RouterLink
            v-for="item in NAVIGASI_HP"
            :key="item.to"
            :to="item.to"
            class="inline-flex h-11 items-center rounded-full bg-subtle px-4 text-sm font-semibold transition-colors duration-150 hover:bg-border focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            {{ item.label }}
          </RouterLink>
        </nav>
      </div>
    </div>
  </div>
</template>

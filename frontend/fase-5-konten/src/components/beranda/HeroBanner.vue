<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { ArrowRight, ChevronDown, Pause, Play } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { useMuat } from '@/composables/useMuat'
import { listBanners } from '@/services/kontenService'
import { perilakuGulir } from '@/utils/gerak'

// Hero banner Beranda (wireframe 01, blok A11:L18) = landing page: setinggi layar di bawah header, media selebar
// layar (gambar atau video) + judul di kiri, tombol di kanan bawah. Isi dari banner posisi HOME_HERO di admin;
// tanpa banner tayang → teks bawaan tanpa media. "Lihat selengkapnya" selalu ada. Tanpa tombol WhatsApp di hero
// (keputusan pemilik 2 Nov 2026); WhatsApp tetap di header dan tombol melayang.
const BAWAAN = {
  title: 'Dari pondasi sampai atap, semua kebutuhan rumah ada di sini.',
  subtitle: 'Cari bahan untuk tiap tahap rumahmu, cek harga dan stok.',
  cta_label: 'Jelajahi katalog',
  link_url: '/katalog',
  cta2_label: null,
  cta2_url: null,
  image_url: null,
  video_url: null,
}

const muat = useMuat(async () => (await listBanners('HOME_HERO'))[0] ?? null)
// Kolom kosong di banner memakai nilai bawaan (mis. tombol tanpa label → "Jelajahi katalog").
const hero = computed(() => {
  const b = muat.data.value ?? {}
  return Object.fromEntries(Object.entries(BAWAAN).map(([k, v]) => [k, b[k] || v]))
})
// Tautan diawali / → RouterLink; https:// → tautan biasa.
const tujuan = (url) => (url.startsWith('/') ? { to: url } : { href: url })
const tombolKedua = computed(() => (hero.value.cta2_label && hero.value.cta2_url ? hero.value : null))

// Video berjalan otomatis tanpa suara. Tidak diputar otomatis bila pengguna meminta animasi dikurangi,
// dan tidak dimuat sama sekali bila mode hemat data aktif; gambar tetap tampil sebagai cadangan.
const kurangiGerak = window.matchMedia?.('(prefers-reduced-motion: reduce)').matches === true
const hematData = navigator.connection?.saveData === true
const videoEl = ref(null)
const videoGagal = ref(false)
const berjalan = ref(!kurangiGerak)
const pakaiVideo = computed(() => Boolean(hero.value.video_url) && !hematData && !videoGagal.value)

watch(
  () => hero.value.video_url,
  async () => {
    videoGagal.value = false
    await nextTick()
    if (videoEl.value && berjalan.value) videoEl.value.play().catch(() => (berjalan.value = false))
  },
)

function gantiPutar() {
  const v = videoEl.value
  if (!v) return
  if (berjalan.value) v.pause()
  else v.play().catch(() => {})
  berjalan.value = !berjalan.value
}

// "Lihat selengkapnya": gulir ke seksi sesudah hero dan pindahkan fokus ke sana (pengguna keyboard lanjut dari situ).
const bagianEl = ref(null)
function gulirKeBawah() {
  const berikut = bagianEl.value?.nextElementSibling
  if (!berikut) return
  berikut.scrollIntoView({ behavior: perilakuGulir(), block: 'start' })
  if (!berikut.hasAttribute('tabindex')) {
    berikut.setAttribute('tabindex', '-1')
    berikut.classList.add('outline-none')
  }
  berikut.focus({ preventScroll: true })
}
</script>

<template>
  <section ref="bagianEl" aria-labelledby="judul-beranda" class="relative isolate overflow-hidden bg-secondary text-secondary-foreground" data-hero-banner>
    <!-- Media sebagai latar; lapisan navy 70% menjaga kontras teks putih di atas foto/video apa pun. -->
    <div v-if="hero.image_url || pakaiVideo" class="absolute inset-0 -z-10" aria-hidden="true">
      <img v-if="hero.image_url" :src="hero.image_url" alt="" width="1920" height="800" fetchpriority="high" class="size-full object-cover" data-hero-gambar />
      <video
        v-if="pakaiVideo"
        ref="videoEl"
        :src="hero.video_url"
        :poster="hero.image_url ?? undefined"
        :autoplay="berjalan"
        muted
        loop
        playsinline
        preload="metadata"
        tabindex="-1"
        class="absolute inset-0 size-full object-cover"
        data-hero-video
        @error="videoGagal = true"
      />
      <div class="absolute inset-0 bg-secondary/70" />
    </div>

    <!-- Tinggi = layar dikurangi header (--tinggi-header-penuh dari AppHeader), minimal 26rem untuk layar pendek. -->
    <div class="mx-auto flex min-h-[max(26rem,calc(100svh-var(--tinggi-header-penuh,0px)))] max-w-7xl flex-col gap-8 px-4 pt-10 pb-4 sm:px-6 lg:pt-16">
      <div v-if="muat.status.value === 'loading'" class="flex max-w-2xl flex-1 flex-col justify-center gap-4" aria-busy="true">
        <span class="sr-only">Memuat hero…</span>
        <div v-for="w in ['w-full', 'w-3/4', 'w-2/3']" :key="w" class="h-10 animate-pulse rounded-md bg-secondary-hover motion-reduce:animate-none" :class="w" aria-hidden="true" />
      </div>
      <div v-else class="flex max-w-2xl flex-1 flex-col justify-center gap-4">
        <h1 id="judul-beranda" class="text-3xl leading-tight font-extrabold text-balance sm:text-4xl lg:text-5xl">{{ hero.title }}</h1>
        <p v-if="hero.subtitle" class="max-w-xl text-lg text-on-dark">{{ hero.subtitle }}</p>
      </div>

      <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <button
          v-if="pakaiVideo"
          type="button"
          class="inline-flex min-h-11 items-center gap-2 self-start rounded-md border border-on-dark-muted/60 px-3.5 text-sm font-semibold text-on-dark transition-colors duration-150 hover:bg-secondary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark sm:self-auto"
          data-jeda-hero
          @click="gantiPutar"
        >
          <component :is="berjalan ? Pause : Play" class="size-4" aria-hidden="true" />
          {{ berjalan ? 'Jeda video' : 'Putar video' }}
        </button>
        <!-- CTA kanan bawah (wireframe): tombol utama, lalu tombol kedua bila diisi admin. Di HP bertumpuk penuh. -->
        <div v-if="muat.status.value !== 'loading'" class="flex flex-col gap-3 sm:ml-auto sm:flex-row sm:flex-wrap sm:justify-end" data-hero-tombol>
          <Button size="lg" v-bind="tujuan(hero.link_url)">
            {{ hero.cta_label }}
            <ArrowRight class="size-5" aria-hidden="true" />
          </Button>
          <Button v-if="tombolKedua" variant="outline" size="lg" v-bind="tujuan(tombolKedua.cta2_url)">{{ tombolKedua.cta2_label }}</Button>
        </div>
      </div>

      <button
        type="button"
        class="mx-auto inline-flex min-h-11 flex-col items-center justify-center rounded-md px-3 text-sm font-semibold text-on-dark transition-colors duration-150 hover:text-secondary-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark"
        data-gulir-hero
        @click="gulirKeBawah"
      >
        Lihat selengkapnya
        <ChevronDown class="size-5" aria-hidden="true" />
      </button>
    </div>
  </section>
</template>

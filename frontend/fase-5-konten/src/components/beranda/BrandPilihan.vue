<script setup>
import { computed } from 'vue'
import { ArrowRight } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useMuat } from '@/composables/useMuat'
import { listBrands } from '@/services/catalogService'
import { formatAngka } from '@/utils/format'
import LingkaranBrand from './LingkaranBrand.vue'
import PesanGagal from './PesanGagal.vue'

// "Brand pilihan" (desain pemilik 2 Nov 2026): judul di tengah, diapit lingkaran berisi logo brand yang diunggah
// di admin › Katalog › Brand. Desktop: 9 lingkaran di kiri + 9 di kanan, baris tengah bergeser setengah langkah
// dan lingkaran terluar terpotong tepi panel. HP: 5 lingkaran di atas + 5 di bawah judul.
// Urutan: brand berlogo dulu, lalu produk terbanyak, lalu nama; posisi terdekat judul diisi lebih dulu.
// Semua brand tetap bisa dilihat di /brand.
const JUMLAH_LINGKARAN = 18
const WARNA = [
  'var(--color-primary-soft)',
  'var(--color-info-soft)',
  'var(--color-accent)',
  'var(--color-success-soft)',
  'var(--color-danger-soft)',
  'var(--color-on-dark)',
]
const HIASAN = [
  ['var(--color-primary)', 'var(--color-primary-hover)', 115],
  ['var(--color-accent)', 'var(--color-warning)', 65],
  ['var(--color-secondary)', 'var(--color-secondary-hover)', 150],
  ['var(--color-primary-on-dark)', 'var(--color-primary)', 30],
]

const merek = useMuat(listBrands)
const jumlah = computed(() => merek.data.value?.length ?? 0)

const isi = computed(() => {
  const urut = [...(merek.data.value ?? [])].sort(
    (a, b) =>
      Number(Boolean(b.logo_url)) - Number(Boolean(a.logo_url)) ||
      (b.product_count ?? 0) - (a.product_count ?? 0) ||
      a.name.localeCompare(b.name, 'id'),
  )
  return Array.from({ length: JUMLAH_LINGKARAN }, (_, i) =>
    urut[i] ? { kunci: `b${urut[i].id}`, brand: urut[i] } : { kunci: `h${i}`, hiasan: HIASAN[i % HIASAN.length] },
  )
})

// Desktop: peringkat genap ke kiri, ganjil ke kanan. Per sisi 3 baris × 3 kolom, kolom terdalam (dekat judul) dulu.
// Warna menurut posisi: tetangga kiri-kanan dan atas-bawah selalu berbeda.
const berwarna = (item, n) => (item.brand ? { ...item, warna: WARNA[n % WARNA.length] } : item)
function sisi(awal, kiri) {
  const kisi = [[], [], []]
  for (let k = 0; k < 9; k++) kisi[k % 3][Math.floor(k / 3)] = isi.value[awal + k * 2]
  return kisi.flatMap((baris, r) => (kiri ? [...baris].reverse() : baris).map((item, c) => berwarna(item, r * 2 + c + (kiri ? 0 : 3))))
}
const kiri = computed(() => sisi(0, true))
const kanan = computed(() => sisi(1, false))

// HP: peringkat teratas di tengah strip, lalu menyebar ke kiri-kanan.
const POSISI_STRIP = [3, 1, 0, 2, 4]
const strip = (awal) => POSISI_STRIP.map((p, i) => berwarna(isi.value[awal + p], i * 2 + awal))
const atas = computed(() => strip(0))
const bawah = computed(() => strip(5))

// Baris tengah bergeser setengah langkah (diameter + jarak) ke arah luar, seperti susunan bata.
const geser = (arah) => ({ transform: `translateX(calc(${arah} * (var(--d) + 1rem) / 2))` })
</script>

<template>
  <section
    v-if="merek.status.value !== 'ready' || jumlah"
    aria-labelledby="judul-brand-pilihan"
    class="relative overflow-hidden rounded-2xl border border-border bg-surface py-8 [--d:5rem] sm:[--d:6rem] lg:py-6 lg:[--d:7rem] xl:[--d:8rem]"
    data-brand-pilihan
  >
    <div v-if="merek.status.value === 'loading'" class="flex flex-col items-center gap-6" aria-busy="true">
      <span class="sr-only">Memuat brand…</span>
      <div class="flex gap-3"><Skeleton v-for="n in 5" :key="n" class="size-[var(--d)] rounded-full" /></div>
      <Skeleton class="h-10 w-72 max-w-full" />
    </div>
    <PesanGagal v-else-if="merek.status.value === 'error'" class="mx-4" pesan="Brand belum bisa dimuat." @coba-lagi="merek.muat" />

    <div v-else class="grid grid-cols-[minmax(0,1fr)] items-center gap-8 lg:grid-cols-[minmax(0,1fr)_minmax(16rem,24rem)_minmax(0,1fr)] lg:gap-6">
      <ul class="hidden grid-cols-[repeat(3,var(--d))] gap-4 justify-self-end lg:grid" data-lingkaran-kiri>
        <li v-for="(item, i) in kiri" :key="item.kunci" :style="i >= 3 && i < 6 ? geser(-1) : null">
          <LingkaranBrand :item="item" />
        </li>
      </ul>
      <ul class="flex justify-center gap-4 lg:hidden" data-lingkaran-atas>
        <li v-for="item in atas" :key="item.kunci"><LingkaranBrand :item="item" /></li>
      </ul>

      <div class="flex flex-col items-center gap-6 px-4 text-center">
        <h2 id="judul-brand-pilihan" class="text-3xl leading-tight font-extrabold text-balance sm:text-4xl xl:text-5xl">
          Siap untuk menjelajah lebih banyak Brand?
        </h2>
        <Button to="/brand" size="lg">
          Telusuri {{ formatAngka(jumlah) }} brand
          <ArrowRight class="size-5" aria-hidden="true" />
        </Button>
      </div>

      <ul class="flex justify-center gap-4 lg:hidden" :style="geser(1)" data-lingkaran-bawah>
        <li v-for="item in bawah" :key="item.kunci"><LingkaranBrand :item="item" /></li>
      </ul>
      <ul class="hidden grid-cols-[repeat(3,var(--d))] gap-4 justify-self-start lg:grid" data-lingkaran-kanan>
        <li v-for="(item, i) in kanan" :key="item.kunci" :style="i >= 3 && i < 6 ? geser(1) : null">
          <LingkaranBrand :item="item" />
        </li>
      </ul>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowRight } from 'lucide-vue-next'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useMuat } from '@/composables/useMuat'
import { listBanners } from '@/services/kontenService'
import { useAuthStore } from '@/stores/auth'
import BannerSlider from './BannerSlider.vue'
import PesanGagal from './PesanGagal.vue'

// Sketsa "Main": slider banner (HOME_SLIDER) + kolom samping (banner HOME_SIDE dan kartu Trans Family).
// Banner hanya yang sedang tayang (rentang tanggal & aktif), sudah difilter API.
const auth = useAuthStore()
// Hanya member yang punya kartu; admin yang sedang masuk diperlakukan seperti pengunjung.
const member = computed(() => auth.hasRole('CUSTOMER'))
const slider = useMuat(() => listBanners('HOME_SLIDER'))
const samping = useMuat(() => listBanners('HOME_SIDE'))

// Tanpa banner tayang, kolom slider disembunyikan dan kartu samping memakai lebar penuh.
const tampilSlider = computed(() => slider.status.value !== 'ready' || slider.data.value.length > 0)
const bannerSamping = computed(() => (samping.status.value === 'ready' ? (samping.data.value[0] ?? null) : null))
// Hanya atribut yang ada yang dipasang: `href: undefined` akan menimpa href buatan RouterLink.
const tautanSamping = computed(() => {
  const link = bannerSamping.value?.link_url
  if (!link) return { is: 'div', atribut: {} }
  return link.startsWith('/') ? { is: RouterLink, atribut: { to: link } } : { is: 'a', atribut: { href: link, target: '_blank', rel: 'noopener' } }
})
</script>

<template>
  <section aria-labelledby="judul-promo-pilihan" class="grid gap-4" :class="tampilSlider && 'lg:grid-cols-3'">
    <h2 id="judul-promo-pilihan" class="sr-only">Promo pilihan</h2>
    <div v-if="tampilSlider" class="lg:col-span-2">
      <Skeleton v-if="slider.status.value === 'loading'" class="h-90 rounded-xl" />
      <PesanGagal v-else-if="slider.status.value === 'error'" pesan="Banner promo belum bisa dimuat." @coba-lagi="slider.muat" />
      <BannerSlider v-else :banner="slider.data.value" />
    </div>

    <div class="flex flex-col gap-4 sm:flex-row" :class="tampilSlider && 'lg:flex-col'">
      <component
        :is="tautanSamping.is"
        v-if="bannerSamping"
        v-bind="tautanSamping.atribut"
        class="group relative flex min-h-40 flex-1 items-end overflow-hidden rounded-xl border border-border bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        data-banner-samping
      >
        <!-- Gambar jadi latar (absolute) supaya tinggi aslinya tidak mendorong label keluar kotak. -->
        <img :src="bannerSamping.image_url" alt="" width="800" height="520" loading="lazy" class="absolute inset-0 size-full object-cover" />
        <span class="relative flex p-3">
          <span class="rounded-md bg-surface px-3 py-2 text-base font-bold text-foreground group-hover:underline">
            {{ bannerSamping.title }}<span v-if="tautanSamping.is === 'a'" class="sr-only"> (tab baru)</span>
          </span>
        </span>
      </component>

      <RouterLink
        :to="member ? '/akun' : '/daftar'"
        class="flex min-h-40 flex-1 flex-col justify-end gap-1.5 rounded-xl bg-accent p-5 text-accent-foreground transition-colors duration-150 hover:bg-accent/85 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        data-kartu-trans-family
      >
        <span class="text-xs font-bold tracking-wide uppercase">Trans Family</span>
        <strong class="text-lg leading-snug">Belanja di toko, kumpulkan poin</strong>
        <span class="inline-flex items-center gap-1 text-sm font-bold">
          {{ member ? 'Lihat kartu member' : 'Daftar gratis' }}
          <ArrowRight class="size-4" aria-hidden="true" />
        </span>
      </RouterLink>
    </div>
  </section>
</template>

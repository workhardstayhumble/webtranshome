<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowUpRight, MapPin, MessageCircle } from 'lucide-vue-next'
import WhatsAppButton from './WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { useSettingsStore } from '@/stores/settings'
const BELANJA = [
  { label: 'Katalog', to: '/katalog' }, { label: 'Ruangan', to: '/ruangan' },
  { label: 'Brand', to: '/brand' }, { label: 'Promo', to: '/promo' },
  { label: 'Artikel', to: '/artikel' }, { label: 'Info toko', to: '/info-toko' },
]
const settings = useSettingsStore()
const { generalLink } = useWhatsApp()
const linkWhatsApp = computed(() => generalLink())
const tahun = new Date().getFullYear()
</script>
<template>
  <footer class="bg-secondary text-sm text-on-dark">
    <div class="mx-auto grid max-w-7xl gap-10 px-4 py-14 sm:grid-cols-2 sm:px-6 lg:grid-cols-4 lg:py-16">
      <div class="flex flex-col items-start gap-4">
        <RouterLink to="/" class="rounded-sm text-2xl font-extrabold tracking-tight text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark">Trans<span class="text-primary-on-dark">home</span></RouterLink>
        <p class="max-w-[220px] text-sm leading-6 text-on-dark-muted">Pilihan material bangunan dan elektronik dengan panduan yang membantu.</p>
        <span class="border-l-2 border-accent pl-3 text-xs font-semibold tracking-widest text-white uppercase">The Warm Guide</span>
      </div>
      <div class="flex flex-col items-start gap-3">
        <p class="text-base font-bold text-white">Kunjungi toko</p>
        <template v-if="settings.data">
          <address class="flex gap-2 leading-6 not-italic"><MapPin class="mt-0.5 size-4 shrink-0 text-accent" aria-hidden="true" />{{ settings.data.address }}</address>
          <p>Buka {{ settings.data.opening_hours }}</p>
          <a v-if="settings.data.maps_url" :href="settings.data.maps_url" target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-1 font-semibold text-white underline-offset-4 hover:underline">Buka di Google Maps <ArrowUpRight class="size-4" aria-hidden="true" /></a>
        </template>
        <p v-else class="text-on-dark-muted">Informasi toko belum tersedia.</p>
      </div>
      <nav aria-label="Belanja" class="flex flex-col items-start gap-3">
        <p class="text-base font-bold text-white">Jelajahi</p>
        <RouterLink v-for="item in BELANJA" :key="item.to" :to="item.to" class="rounded-sm text-on-dark-muted transition-colors hover:text-white hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark">{{ item.label }}</RouterLink>
      </nav>
      <div class="flex flex-col items-start gap-3">
        <p class="text-base font-bold text-white">Butuh bantuan?</p>
        <p class="max-w-xs leading-6 text-on-dark-muted">Tanya stok, harga, dan kebutuhan proyek secara langsung kepada tim Transhome.</p>
        <WhatsAppButton :href="linkWhatsApp" label="Chat WhatsApp CS" />
        <p class="flex items-center gap-2 text-xs text-on-dark-muted"><MessageCircle class="size-4" aria-hidden="true" /> Konsultasi tanpa checkout online.</p>
      </div>
    </div>
    <div class="border-t border-white/15"><div class="mx-auto flex max-w-7xl flex-wrap justify-between gap-2 px-4 py-5 text-xs text-on-dark-muted sm:px-6"><span>© {{ tahun }} CV Trans Home · Sleman, Yogyakarta</span><span>Material tepat. Rumah terasa dekat.</span></div></div>
  </footer>
</template>
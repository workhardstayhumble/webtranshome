<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import WhatsAppButton from './WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { useSettingsStore } from '@/stores/settings'

const BELANJA = [
  { label: 'Katalog', to: '/katalog' },
  { label: 'Ruangan', to: '/ruangan' },
  { label: 'Brand', to: '/brand' },
  { label: 'Promo', to: '/promo' },
  { label: 'Artikel', to: '/artikel' },
  { label: 'Info toko', to: '/info-toko' },
]

const settings = useSettingsStore()
const { generalLink } = useWhatsApp()
const linkWhatsApp = computed(() => generalLink())
const tahun = new Date().getFullYear()

const kelasTautan =
  'rounded-sm transition-colors duration-150 hover:text-secondary-foreground hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark'
</script>

<template>
  <footer class="bg-secondary text-sm text-on-dark">
    <div class="mx-auto grid max-w-7xl gap-8 px-4 py-12 sm:grid-cols-2 sm:px-6 lg:grid-cols-4">
      <div class="flex flex-col gap-2">
        <p class="text-2xl font-extrabold tracking-tight text-secondary-foreground">Trans<span class="text-primary-on-dark">home</span></p>
        <p>Supermarket bahan bangunan dan elektronik.</p>
      </div>

      <div class="flex flex-col gap-2">
        <p class="font-bold text-secondary-foreground">Kunjungi toko</p>
        <template v-if="settings.data">
          <address class="not-italic">{{ settings.data.address }}</address>
          <p>Buka {{ settings.data.opening_hours }}</p>
          <a v-if="settings.data.maps_url" :href="settings.data.maps_url" target="_blank" rel="noopener" :class="[kelasTautan, 'self-start font-semibold']">
            Buka di Google Maps
          </a>
        </template>
        <p v-else class="text-on-dark-muted">Info toko belum tersedia.</p>
      </div>

      <nav aria-label="Belanja" class="flex flex-col gap-2">
        <p class="font-bold text-secondary-foreground">Belanja</p>
        <RouterLink v-for="item in BELANJA" :key="item.to" :to="item.to" :class="[kelasTautan, 'self-start']">{{ item.label }}</RouterLink>
      </nav>

      <div class="flex flex-col items-start gap-3">
        <p class="font-bold text-secondary-foreground">Butuh bantuan?</p>
        <p>Tanya stok, harga, atau kebutuhan proyek langsung ke tim kami.</p>
        <WhatsAppButton :href="linkWhatsApp" label="Chat WhatsApp CS" />
      </div>
    </div>
    <div class="border-t border-secondary-hover">
      <p class="mx-auto max-w-7xl px-4 py-5 text-on-dark-muted sm:px-6">© {{ tahun }} CV Trans Home · Sleman, Yogyakarta</p>
    </div>
  </footer>
</template>
<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { Instagram } from 'lucide-vue-next'
import WhatsAppButton from './WhatsAppButton.vue'
import LogoTranshome from '@/components/layout/LogoTranshome.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { useSettingsStore } from '@/stores/settings'

const BELANJA = [
  { label: 'Katalog', to: '/katalog' },
  { label: 'Ruangan', to: '/ruangan' },
  { label: 'Brand', to: '/brand' },
  { label: 'Jasa', to: '/jasa' },
  { label: 'Promo', to: '/promo' },
  { label: 'Artikel', to: '/artikel' },
  { label: 'Info toko', to: '/info-toko' },
]

// Akun media sosial resmi (dari pemilik, 2 Nov 2026). Ikon TikTok tidak ada di lucide → SVG di template.
const SOSIAL = {
  instagram: 'https://www.instagram.com/transhome/',
  tiktok: 'https://www.tiktok.com/@transhome_id',
}

const settings = useSettingsStore()
const { generalLink } = useWhatsApp()
const linkWhatsApp = computed(() => generalLink())
const tahun = new Date().getFullYear()

// Tombol media sosial di latar navy: garis tepi tipis, teks terang, area sentuh 44px.
const kelasSosial =
  'inline-flex min-h-11 items-center gap-2 rounded-md border border-on-dark-muted/60 px-3.5 font-semibold text-secondary-foreground transition-colors duration-150 hover:border-on-dark hover:bg-secondary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark'

const kelasTautan =
  'rounded-sm transition-colors duration-150 hover:text-secondary-foreground hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark'
</script>

<template>
  <footer class="bg-secondary text-sm text-on-dark">
    <div class="mx-auto grid max-w-7xl gap-8 px-4 py-12 sm:grid-cols-2 sm:px-6 lg:grid-cols-4">
      <div class="flex flex-col gap-2">
        <p class="text-2xl font-extrabold tracking-tight text-secondary-foreground"><LogoTranshome gelap /></p>
        <p>Supermarket bahan bangunan dan elektronik.</p>
        <p class="mt-3 font-bold text-secondary-foreground">Ikuti kami</p>
        <ul class="flex flex-wrap gap-2" data-sosial>
          <li>
            <a :href="SOSIAL.instagram" target="_blank" rel="noopener noreferrer" :class="kelasSosial">
              <Instagram class="size-5" aria-hidden="true" />
              Instagram<span class="sr-only"> Transhome (tab baru)</span>
            </a>
          </li>
          <li>
            <a :href="SOSIAL.tiktok" target="_blank" rel="noopener noreferrer" :class="kelasSosial">
              <svg viewBox="0 0 24 24" fill="currentColor" class="size-5" aria-hidden="true">
                <path
                  d="M12.525.02c1.31-.02 2.61-.01 3.91-.02.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z"
                />
              </svg>
              TikTok<span class="sr-only"> Transhome (tab baru)</span>
            </a>
          </li>
        </ul>
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

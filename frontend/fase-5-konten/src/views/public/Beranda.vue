<script setup>
import AjakanMember from '@/components/beranda/AjakanMember.vue'
import BannerPromo from '@/components/beranda/BannerPromo.vue'
import BrandPilihan from '@/components/beranda/BrandPilihan.vue'
import HeroRumah from '@/components/beranda/HeroRumah.vue'
import IntroHero from '@/components/beranda/IntroHero.vue'
import InspirasiRuangan from '@/components/beranda/InspirasiRuangan.vue'
import KategoriTahap from '@/components/beranda/KategoriTahap.vue'
import Keunggulan from '@/components/beranda/Keunggulan.vue'
import ProdukUnggulan from '@/components/beranda/ProdukUnggulan.vue'
import PromoEvent from '@/components/beranda/PromoEvent.vue'
import { useSeo } from '@/composables/useSeo'
useSeo({
  title: 'Material bangunan & elektronik pilihan di Sleman',
  description: 'Jelajahi katalog, inspirasi ruang, spesifikasi dan konsultasi WhatsApp Transhome di Sleman, Yogyakarta.',
})
</script>
<template>
  <div>
    <IntroHero />
    <div class="mx-auto flex max-w-7xl flex-col gap-15 px-4 py-12 sm:gap-18 sm:px-6 sm:py-18">
      <KategoriTahap v-reveal />
      <ProdukUnggulan v-reveal />
      <BannerPromo v-reveal />
      <InspirasiRuangan v-reveal />
      <PromoEvent v-reveal />
      <BrandPilihan v-reveal />
    </div>
    <section class="mt-4" aria-label="Eksplorasi interaktif bagian rumah"><div class="mx-auto max-w-7xl px-4 pb-5 sm:px-6"><p class="th-section-label">Eksplorasi interaktif</p><h2 class="th-section-heading mt-2">Kenali material di setiap bagian rumah</h2></div><HeroRumah /></section>
    <div class="mx-auto flex max-w-7xl flex-col gap-12 px-4 py-14 sm:px-6 sm:py-18">
      <Keunggulan v-reveal />
      <AjakanMember v-reveal />
    </div>
  </div>
</template>
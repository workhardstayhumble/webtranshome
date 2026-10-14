<script setup>
import AjakanMember from '@/components/beranda/AjakanMember.vue'
import BannerPromo from '@/components/beranda/BannerPromo.vue'
import BrandPilihan from '@/components/beranda/BrandPilihan.vue'
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
    <div class="mx-auto flex max-w-7xl flex-col gap-12 px-4 py-10 sm:px-6 lg:gap-16 lg:py-14">
      <KategoriTahap v-reveal />
      <ProdukUnggulan v-reveal />
      <BannerPromo v-reveal />
      <InspirasiRuangan v-reveal />
      <PromoEvent v-reveal />
      <BrandPilihan v-reveal />
    </div>
    <div class="mx-auto flex max-w-7xl flex-col gap-12 px-4 pb-14 sm:px-6 sm:pb-18">
      <Keunggulan v-reveal />
      <AjakanMember v-reveal />
    </div>
  </div>
</template>
<script setup>
import { onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppFooter from '@/components/layout/AppFooter.vue'
import AppHeader from '@/components/layout/AppHeader.vue'
import WhatsAppFloat from '@/components/layout/WhatsAppFloat.vue'
import { useSettingsStore } from '@/stores/settings'

const route = useRoute()
const settings = useSettingsStore()
onMounted(() => settings.load())
</script>

<template>
  <!-- meta.barAksiHp: halaman punya bar aksi tetap di bawah layar HP, jadi footer diberi ruang setinggi bar itu. -->
  <div class="flex min-h-dvh flex-col" :class="route.meta.barAksiHp && 'pb-[calc(6.75rem+env(safe-area-inset-bottom))] sm:pb-[calc(4.5rem+env(safe-area-inset-bottom))] lg:pb-0'">
    <a
      href="#konten"
      class="sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-60 focus:rounded-md focus:bg-surface focus:px-4 focus:py-2 focus:font-semibold focus:shadow-lg"
    >
      Langsung ke konten
    </a>
    <AppHeader />
    <main id="konten" tabindex="-1" class="flex-1 focus:outline-none">
      <RouterView />
    </main>
    <AppFooter />
    <WhatsAppFloat />
  </div>
</template>

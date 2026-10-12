<script setup>
import { RouterLink } from 'vue-router'
import { ArrowRight, BadgeCheck, ChevronRight, Layers3, Sparkles } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
const { generalLink } = useWhatsApp()
</script>

<template>
  <section aria-labelledby="intro-hero" class="overflow-hidden border-b border-border bg-subtle">
    <div class="mx-auto grid max-w-7xl items-center gap-8 px-4 py-12 sm:px-6 lg:grid-cols-[0.9fr_1.1fr] lg:gap-12 lg:py-20">
      <div class="relative z-10 flex flex-col items-start gap-6">
        <div class="inline-flex items-center gap-2 rounded-full border border-primary/15 bg-white px-3 py-1.5 text-xs font-bold tracking-wider text-primary uppercase">
          <Sparkles class="size-3.5" aria-hidden="true" /> Material & elektronik untuk rumah Anda
        </div>
        <h1 id="intro-hero" class="max-w-xl text-[clamp(2.5rem,4.7vw,4.5rem)] leading-[1.09] font-extrabold tracking-[-0.045em] text-foreground">
          Wujudkan <span class="text-primary">rumah impian</span> dengan pilihan yang tepat.
        </h1>
        <p class="max-w-lg text-base leading-8 text-muted sm:text-lg">Temukan material bangunan dan elektronik pilihan. Bandingkan spesifikasi, dapatkan inspirasi ruangan, dan konsultasikan kebutuhan dengan tim Transhome.</p>
        <div class="flex w-full flex-wrap items-center gap-3">
          <Button to="/katalog" size="lg">Jelajahi katalog <ArrowRight class="size-4.5" aria-hidden="true" /></Button>
          <WhatsAppButton :href="generalLink()" label="Konsultasi WhatsApp" size="lg" />
        </div>
        <div class="grid w-full max-w-xl grid-cols-2 gap-4 border-t border-border pt-6 text-sm sm:grid-cols-3">
          <div class="flex items-start gap-2"><BadgeCheck class="mt-0.5 size-4 shrink-0 text-primary" aria-hidden="true" /><span class="font-medium">Informasi jelas & terstruktur</span></div>
          <div class="flex items-start gap-2"><Layers3 class="mt-0.5 size-4 shrink-0 text-primary" aria-hidden="true" /><span class="font-medium">Pilihan lintas kategori</span></div>
          <div class="col-span-2 flex items-start gap-2 sm:col-span-1"><Sparkles class="mt-0.5 size-4 shrink-0 text-primary" aria-hidden="true" /><span class="font-medium">Konsultasi lebih mudah</span></div>
        </div>
      </div>
      <!-- Ilustrasi arsitektural dekoratif, bukan foto atau gambar produk asli. Ganti dengan foto toko berlisensi saat tersedia. -->
      <div class="relative min-h-[340px] overflow-hidden rounded-[10px] border border-[#d9d4ca] bg-[#eae5d8] shadow-[0_22px_60px_rgba(28,36,48,0.09)] sm:min-h-[490px] lg:min-h-[560px]" role="group" aria-label="Ilustrasi interior arsitektural hangat, bukan foto produk">
        <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_70%_20%,#fffaf1_0%,#e7ddcc_42%,#d6cfbf_90%)]"></div>
        <div class="absolute top-0 right-[13%] h-[76%] w-[61%] border-x-[12px] border-b-[15px] border-[#bdb5a6] bg-[linear-gradient(120deg,#eee7d9_0%,#fff9ed_48%,#dbcfbb_100%)] shadow-[-14px_7px_25px_rgba(64,51,39,0.08)] sm:border-x-[18px]"></div>
        <div class="absolute top-[8%] right-[20%] h-[52%] w-[47%] border-[14px] border-[#d1c5b1] bg-[linear-gradient(160deg,#fdfcf6,#e8e2d3)] shadow-inner sm:border-[22px]"></div>
        <div class="absolute top-[14%] right-[25%] h-[41%] w-[37%] border-l-[2px] border-r-[2px] border-[#c2baab] bg-[linear-gradient(142deg,#cfdbdd_0%,#f6f1e4_54%,#e8d9c3_100%)]"></div>
        <div class="absolute bottom-0 left-0 h-[26%] w-full border-t border-[#c4bbaa] bg-[repeating-linear-gradient(110deg,#cbbca5_0px,#cbbca5_2px,#d9cbb4_3px,#eadfcf_95px)]"></div>
        <div class="absolute right-[5%] bottom-[17%] h-[27%] w-[60%] rounded-[8px_8px_4px_4px] border border-[#c7bbae] bg-[linear-gradient(155deg,#fffdf9,#f5eee5)] shadow-[-16px_18px_18px_rgba(76,58,45,0.17)]"></div>
        <div class="absolute right-[16%] bottom-[17%] h-[10%] w-[33%] border-t-[8px] border-[#d9cdbc] bg-[#f2e7d5]"></div>
        <div class="absolute right-[10%] bottom-[17%] h-[3%] w-[51%] bg-[#b8aa99]"></div>
        <div class="absolute bottom-[18%] left-[12%] h-[36%] w-[14%] rounded-t-full bg-[#71846b] shadow-[14px_8px_0_#849479,-14px_18px_0_#677961]"></div>
        <div class="absolute bottom-[13%] left-[9%] h-[13%] w-[22%] rounded-[3px_3px_12px_12px] bg-[#bfa78e] shadow-md"></div>
        <div class="absolute right-5 bottom-5 left-5 flex items-end justify-between gap-3 rounded-lg border border-white/60 bg-white/90 p-4 shadow-lg backdrop-blur-sm sm:right-7 sm:bottom-7 sm:left-7 sm:p-5">
          <div class="flex flex-col gap-1"><span class="text-xs font-bold tracking-widest text-primary uppercase">Inspirasi ruang</span><strong class="text-base sm:text-lg">Material yang menyatu dengan ruang</strong><span class="text-xs text-muted">Visual konseptual · bukan foto produk</span></div>
          <RouterLink to="/ruangan" class="flex size-11 shrink-0 items-center justify-center rounded-md bg-primary text-white transition-colors hover:bg-primary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" aria-label="Jelajahi inspirasi ruangan"><ChevronRight class="size-5" aria-hidden="true" /></RouterLink>
        </div>
      </div>
    </div>
  </section>
</template>
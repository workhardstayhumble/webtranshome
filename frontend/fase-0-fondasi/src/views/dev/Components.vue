<script setup>
import DesignFoundations from './DesignFoundations.vue'
import DemoAkun from './bagian/DemoAkun.vue'
import DemoData from './bagian/DemoData.vue'
import DemoForm from './bagian/DemoForm.vue'
import DemoKatalog from './bagian/DemoKatalog.vue'
import DemoKeadaan from './bagian/DemoKeadaan.vue'
import DemoStatus from './bagian/DemoStatus.vue'
import DemoTombol from './bagian/DemoTombol.vue'
import DemoUmpanBalik from './bagian/DemoUmpanBalik.vue'
const BAGIAN = [
  { id:'foundations',judul:'Design foundations',komponen:DesignFoundations },
  { id:'tombol',judul:'Button & WhatsApp CTA',komponen:DemoTombol },
  { id:'form',judul:'Input, field & filter',komponen:DemoForm },
  { id:'status',judul:'Badge, stock, tabs & harga',komponen:DemoStatus },
  { id:'umpan-balik',judul:'Toast & dialog',komponen:DemoUmpanBalik },
  { id:'keadaan',judul:'Empty, error, loading',komponen:DemoKeadaan },
  { id:'data',judul:'Tabel & paginasi',komponen:DemoData },
  { id:'katalog',judul:'Kartu produk & katalog',komponen:DemoKatalog },
  { id:'akun',judul:'Akun & hak akses',komponen:DemoAkun },
]
</script>
<template>
  <div class="mx-auto max-w-7xl space-y-9 px-4 py-12 sm:px-6">
    <header class="space-y-4 border-b border-border pb-8">
      <p class="th-section-label">Transhome / Design System v1.0</p>
      <h1 class="text-4xl font-extrabold tracking-tight sm:text-5xl">Components <span class="text-primary">& Foundations</span></h1>
      <p class="max-w-3xl text-base leading-7 text-muted">Galeri nyata komponen WebTranshome, diselaraskan dengan referensi Warm Modern Specialist Retail. Hanya tersedia saat mode pengembangan.</p>
      <nav aria-label="Bagian galeri" class="flex flex-wrap gap-2"><a v-for="b in BAGIAN" :key="b.id" :href="`#${b.id}`" class="inline-flex min-h-10 items-center rounded-md border border-border bg-surface px-3.5 text-xs font-semibold text-muted transition-colors hover:border-primary hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">{{ b.judul }}</a></nav>
    </header>
    <section v-for="b in BAGIAN" :id="b.id" :key="b.id" :aria-labelledby="`judul-${b.id}`" class="scroll-mt-6 rounded-xl border border-border bg-surface p-5 shadow-[0_3px_20px_rgba(28,36,48,0.04)] sm:p-8">
      <h2 v-if="b.id !== 'foundations'" :id="`judul-${b.id}`" class="mb-6 border-b border-border pb-3 text-xl font-bold">{{ b.judul }}</h2>
      <component :is="b.komponen" />
    </section>
  </div>
</template>
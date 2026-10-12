<script setup>
import { computed } from 'vue'
const COLORS = [
  { label: 'Primary Blue', hex: '#0033CC', token: 'color/brand/primary' },
  { label: 'Brand Red', hex: '#FF3300', token: 'color/brand/red' },
  { label: 'Brand Yellow', hex: '#FFCC00', token: 'color/brand/yellow' },
  { label: 'Brand White / Card', hex: '#FFFFFF', token: 'color/surface/card' },
  { label: 'Warm Canvas', hex: '#F7F5F0', token: 'color/surface/canvas' },
  { label: 'Ink', hex: '#1C2430', token: 'color/text/ink' },
  { label: 'Secondary Text', hex: '#647084', token: 'color/text/secondary' },
  { label: 'Border', hex: '#E3E6EC', token: 'color/border/default' },
  { label: 'Soft Red Surface', hex: '#FFF0EB', token: 'color/brand/red-soft' },
  { label: 'Primary Soft', hex: '#EAEFFF', token: 'color/brand/primary-soft' },
  { label: 'Red for Text', hex: '#B82500', token: 'color/brand/red-text' },
  { label: 'Success', hex: '#1F7A4D', token: 'color/status/success' },
]
const TYPO = [
  { token: 'Display/Hero', specimen: 'Wujudkan rumah impian', css: 'text-[clamp(2.2rem,4vw,3.5rem)] font-extrabold tracking-tight' },
  { token: 'Heading/H1', specimen: 'Wujudkan rumah impian', css: 'text-4xl font-bold tracking-tight' },
  { token: 'Heading/H2', specimen: 'Material bangunan & elektronik pilihan', css: 'text-3xl font-bold tracking-tight' },
  { token: 'Heading/H3', specimen: 'Material bangunan & elektronik pilihan', css: 'text-2xl font-semibold tracking-tight' },
  { token: 'Product/Title', specimen: 'Keramik lantai 60×60 abu natural', css: 'text-lg font-semibold' },
  { token: 'Body/Regular', specimen: 'Material bangunan dan elektronik pilihan', css: 'text-base font-normal' },
  { token: 'Body/Medium', specimen: 'Material bangunan dan elektronik pilihan', css: 'text-base font-medium' },
  { token: 'Meta/Regular', specimen: 'Informasi produk yang relevan dan mudah dipahami', css: 'text-sm font-normal' },
  { token: 'Meta/Medium', specimen: 'Informasi produk yang relevan dan mudah dipahami', css: 'text-sm font-medium' },
  { token: 'Label/Semibold', specimen: 'KATEGORI MATERIAL', css: 'text-xs font-semibold tracking-wide' },
]
const SPACING = [4, 8, 12, 16, 24, 32, 48, 64]
const RADIUS = [6, 10, 16, 999]
function luminance(hex) {
  const rgb = hex.match(/[a-f0-9]{2}/gi).map(part => parseInt(part, 16) / 255)
  const values = rgb.map(v => v <= .04045 ? v / 12.92 : ((v + .055) / 1.055) ** 2.4)
  return values[0] * .2126 + values[1] * .7152 + values[2] * .0722
}
const CHECKS = [
  ['White on Primary Blue', '#FFFFFF', '#0033CC'],
  ['Ink on Brand Yellow', '#1C2430', '#FFCC00'],
  ['White on Brand Red', '#FFFFFF', '#FF3300'],
  ['Ink on Brand Red', '#1C2430', '#FF3300'],
  ['Red-text on White', '#B82500', '#FFFFFF'],
  ['Red-text on Soft Red', '#B82500', '#FFF0EB'],
  ['Secondary text on White', '#647084', '#FFFFFF'],
  ['Secondary text on Warm Canvas', '#647084', '#F7F5F0'],
  ['Ink on Warm Canvas', '#1C2430', '#F7F5F0'],
  ['Primary Blue on White', '#0033CC', '#FFFFFF'],
  ['Primary Blue on Primary Soft', '#0033CC', '#EAEFFF'],
  ['Success on White', '#1F7A4D', '#FFFFFF'],
  ['Primary Blue focus ring on White', '#0033CC', '#FFFFFF'],
  ['Border on White (decorative)', '#E3E6EC', '#FFFFFF'],
]
const KONTRAS = computed(() => CHECKS.map(([name, fg, bg]) => {
  const a = luminance(fg), b = luminance(bg)
  return { name, fg, bg, ratio: ((Math.max(a,b) + .05) / (Math.min(a,b) + .05)).toFixed(2) }
}))
</script>
<template>
  <section class="space-y-10" aria-label="Fondasi desain Transhome">
    <div class="space-y-2"><p class="th-section-label">01 · Foundations</p><h2 class="th-section-heading">Design System · Warm Modern Specialist Retail</h2><p class="text-sm text-muted">Color · Typography (Inter) · Spacing · Radius. Token menjadi rujukan bersama semua halaman.</p></div>
    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
      <div v-for="c in COLORS" :key="c.token" class="min-w-0">
        <div class="mb-2 h-20 rounded-lg border border-border shadow-sm" :style="{ backgroundColor: c.hex }" />
        <p class="text-xs font-semibold">{{ c.label }}</p><p class="mt-0.5 break-words text-[11px] text-muted">{{ c.hex }} · {{ c.token }}</p>
      </div>
    </div>
    <section class="space-y-4">
      <h3 class="text-lg font-bold">Contrast check <span class="font-normal text-muted">(WCAG 2.2 AA)</span></h3>
      <div class="overflow-x-auto rounded-lg border border-border"><table class="w-full min-w-[520px] text-left text-xs sm:text-sm"><thead class="bg-subtle text-xs font-bold uppercase text-muted"><tr><th class="p-3">Warna</th><th class="p-3">Contoh</th><th class="p-3">Kontras</th><th class="p-3">Normal text AA</th></tr></thead><tbody><tr v-for="item in KONTRAS" :key="item.name" class="border-t border-border"><td class="p-3">{{ item.name }}</td><td class="p-3"><span class="inline-block rounded px-3 py-1 text-xs font-bold" :style="{ backgroundColor: item.bg, color: item.fg }">Contoh teks</span></td><td class="p-3 font-mono">{{ item.ratio }} : 1</td><td class="p-3 font-semibold" :class="Number(item.ratio) >= 4.5 ? 'text-success' : 'text-danger'"><span v-if="item.name.includes('(decorative)')">Dekoratif saja</span><span v-else>{{ Number(item.ratio) >= 4.5 ? 'PASS' : 'FAIL — hindari teks kecil' }}</span></td></tr></tbody></table></div>
    </section>
    <section class="space-y-5"><h3 class="text-lg font-bold">Typography — Inter only</h3><div class="space-y-5"><div v-for="item in TYPO" :key="item.token" class="grid items-baseline gap-3 border-b border-border pb-4 md:grid-cols-[200px_1fr]"><span class="text-xs text-muted">{{ item.token }}</span><p class="min-w-0" :class="item.css">{{ item.specimen }}</p></div></div></section>
    <section class="space-y-5"><h3 class="text-lg font-bold">Spacing (4/8 scale) & radius</h3><div class="flex flex-wrap items-end gap-6"><div v-for="n in SPACING" :key="n" class="flex flex-col items-center gap-2"><span class="w-5 bg-primary" :style="{ height: `${n}px` }" /><span class="text-[11px] text-muted">{{ n }}px</span></div><div v-for="n in RADIUS" :key="n" class="flex flex-col items-center gap-2"><span class="size-13 border border-primary bg-primary-soft" :style="{ borderRadius: n === 999 ? '50%' : `${n}px` }" /><span class="text-[11px] text-muted">{{ n === 999 ? 'full' : `${n}px` }}</span></div></div><p class="text-xs leading-6 text-muted">Grid desktop maksimal 1280px · 12 kolom · gutter 24px. Grid mobile 4 kolom · gutter 16px. Target interaksi minimal 44×44px.</p></section>
  </section>
</template>
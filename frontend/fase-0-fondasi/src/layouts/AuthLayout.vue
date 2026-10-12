<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { ArrowLeft, ShieldCheck } from 'lucide-vue-next'
</script>
<template>
  <div class="flex min-h-dvh flex-col bg-background-admin">
    <header class="border-b border-border bg-surface">
      <div class="mx-auto flex h-18 max-w-7xl items-center justify-between gap-4 px-4 sm:px-6">
        <RouterLink to="/" class="rounded-md text-2xl font-extrabold tracking-tight focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" aria-label="Transhome, ke beranda">Trans<span class="text-primary">home</span></RouterLink>
        <RouterLink to="/" class="inline-flex min-h-11 items-center gap-2 rounded-sm text-sm font-semibold text-muted hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"><ArrowLeft class="size-4" aria-hidden="true" /> Kembali ke beranda</RouterLink>
      </div>
    </header>
    <main id="konten" class="mx-auto grid w-full max-w-6xl flex-1 items-center gap-12 px-4 py-8 sm:px-6 lg:grid-cols-2 lg:py-14">
      <section class="hidden flex-col gap-7 lg:flex">
        <span class="th-section-label">Transhome · The Warm Guide</span>
        <h1 class="max-w-lg text-5xl leading-[1.12] font-extrabold tracking-tight">Temukan kebutuhan rumah dengan lebih <span class="text-primary">mudah.</span></h1>
        <p class="max-w-lg text-lg leading-8 text-muted">Jelajahi material, pahami spesifikasi, dan dapatkan bantuan yang tepat dari tim Transhome.</p>
        <div class="max-w-md rounded-xl border border-border bg-surface p-6 shadow-[0_12px_40px_rgba(28,36,48,0.05)]">
          <div class="flex items-center gap-3"><span class="flex size-11 items-center justify-center rounded-lg bg-primary-soft text-primary"><ShieldCheck class="size-5" aria-hidden="true" /></span><div><strong class="block text-sm">Akses aman dan terarah</strong><span class="text-xs text-muted">Pengelolaan sesuai role dan permission.</span></div></div>
        </div>
      </section>
      <div class="mx-auto w-full max-w-md"><RouterView /></div>
    </main>
  </div>
</template>
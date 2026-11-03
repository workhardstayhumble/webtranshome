<script setup>
import Button from '@/components/ui/Button.vue'
import TombolGoogle from '@/components/akun/TombolGoogle.vue'
import { useAuthStore } from '@/stores/auth'

// Ajakan jadi member (sketsa "Landing3D"). Disembunyikan hanya bila yang masuk member (CUSTOMER);
// akun admin yang sedang membuka situs publik tetap melihatnya.
const auth = useAuthStore()
</script>

<template>
  <section
    v-if="!auth.hasRole('CUSTOMER')"
    aria-labelledby="judul-ajakan-member"
    class="grid items-center gap-8 overflow-hidden rounded-2xl bg-primary-soft p-6 sm:p-10 lg:grid-cols-2"
    data-ajakan-member
  >
    <div class="flex flex-col gap-4">
      <h2 id="judul-ajakan-member" class="text-3xl font-extrabold">Jadi member Trans Family, gratis.</h2>
      <p class="text-lg text-primary-hover">Daftar pakai email, dapat kode member digital. Sebutkan kodenya ke kasir, poin belanja langsung tercatat.</p>
      <div class="flex flex-wrap items-center gap-3">
        <Button to="/daftar" size="lg">Daftar sekarang</Button>
        <div class="w-full sm:w-64"><TombolGoogle label="Masuk dengan Google" /></div>
      </div>
    </div>
    <!-- Gambaran kartu member digital (dekorasi). -->
    <div class="flex justify-center py-4 [perspective:1000px]" aria-hidden="true">
      <div
        class="flex h-52 w-full max-w-90 flex-col justify-between rounded-2xl bg-secondary p-6 text-secondary-foreground shadow-xl"
        style="transform: rotateY(-18deg) rotateX(10deg) rotateZ(-2deg)"
      >
        <span class="flex items-center justify-between">
          <strong>Trans Family</strong>
          <span class="rounded bg-accent px-2.5 py-0.5 text-xs font-bold text-accent-foreground">Tahap 1</span>
        </span>
        <span class="font-mono text-2xl font-semibold tracking-widest">TF-000000</span>
        <span class="flex justify-between text-sm text-on-dark"><span>Nama kamu</span><span>0 poin</span></span>
      </div>
    </div>
  </section>
</template>

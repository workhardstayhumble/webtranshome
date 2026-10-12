<script setup>
import { CircleCheck, CircleX, Info, TriangleAlert, X } from 'lucide-vue-next'
import { useUiStore } from '@/stores/ui'

// Satu wadah toast untuk seluruh aplikasi, dipasang di App.vue. Isi dari store `ui`.
const ui = useUiStore()

const JENIS = {
  info: { ikon: Info, kelas: 'text-info' },
  success: { ikon: CircleCheck, kelas: 'text-success' },
  warning: { ikon: TriangleAlert, kelas: 'text-warning' },
  danger: { ikon: CircleX, kelas: 'text-danger' },
}
</script>

<template>
  <!-- Wadah live region selalu ada agar pembaca layar mengumumkan toast baru. -->
  <div
    role="status"
    aria-live="polite"
    class="pointer-events-none fixed inset-x-4 bottom-24 z-60 flex flex-col items-center gap-2 sm:inset-x-auto sm:right-6 sm:bottom-24 sm:items-end"
  >
    <TransitionGroup
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="translate-y-2 opacity-0"
      leave-active-class="transition duration-150 ease-in"
      leave-to-class="opacity-0"
    >
      <div
        v-for="toast in ui.toasts"
        :key="toast.id"
        class="pointer-events-auto flex w-full max-w-sm items-start gap-3 rounded-lg border border-border bg-surface p-4 shadow-[0_12px_32px_rgba(28,36,48,0.14)]"
      >
        <component :is="(JENIS[toast.jenis] ?? JENIS.info).ikon" class="mt-0.5 size-5 shrink-0" :class="(JENIS[toast.jenis] ?? JENIS.info).kelas" aria-hidden="true" />
        <p class="flex-1 text-sm font-medium">{{ toast.pesan }}</p>
        <button
          type="button"
          class="-m-1.5 inline-flex size-8 shrink-0 items-center justify-center rounded-md text-muted transition-colors duration-150 hover:bg-subtle hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          aria-label="Tutup notifikasi"
          @click="ui.hapusToast(toast.id)"
        >
          <X class="size-4" aria-hidden="true" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>
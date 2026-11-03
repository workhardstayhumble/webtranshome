<script setup>
import { computed } from 'vue'

const page = defineModel('page', { type: Number, default: 1 })

const props = defineProps({
  totalPages: { type: Number, required: true },
  // Opsional: (nomor) => href. Bila diisi, nomor halaman jadi tautan <a href> yang bisa diikuti mesin pencari;
  // klik biasa tetap ditangani lewat v-model (tanpa memuat ulang), Ctrl/⌘-klik membuka tab baru.
  tautan: { type: Function, default: null },
})

// Selalu tampilkan halaman pertama, terakhir, dan tetangga halaman aktif; sisanya "…".
const daftar = computed(() => {
  const total = props.totalPages
  const kini = page.value
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)
  const terpilih = [...new Set([1, total, kini - 1, kini, kini + 1])].filter((n) => n >= 1 && n <= total).sort((a, b) => a - b)
  const hasil = []
  terpilih.forEach((n, i) => {
    if (i > 0 && n - terpilih[i - 1] > 1) hasil.push(`jeda-${n}`)
    hasil.push(n)
  })
  return hasil
})

const atribut = (n) => (props.tautan ? { href: props.tautan(n) } : { type: 'button' })
function pilih(event, n) {
  if (props.tautan && (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey || event.button !== 0)) return
  event.preventDefault()
  page.value = n
}

const kelasTombol =
  'inline-flex h-11 min-w-11 items-center justify-center rounded-md px-3 text-sm font-semibold shadow-[0_1px_2px_rgba(28,36,48,0.04)] transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:h-10 md:min-w-10'
</script>

<template>
  <!-- Tombol di ujung disembunyikan, bukan dinonaktifkan (rencana build Fase 1). -->
  <nav v-if="totalPages > 1" aria-label="Halaman" class="flex flex-wrap items-center justify-center gap-1.5">
    <component
      :is="tautan ? 'a' : 'button'"
      v-if="page > 1"
      v-bind="atribut(page - 1)"
      :rel="tautan ? 'prev' : undefined"
      :class="[kelasTombol, 'border border-border-strong bg-surface hover:border-foreground']"
      @click="pilih($event, page - 1)"
    >
      Sebelumnya
    </component>
    <template v-for="item in daftar" :key="item">
      <span v-if="typeof item === 'string'" class="px-1 text-muted" aria-hidden="true">…</span>
      <component
        :is="tautan ? 'a' : 'button'"
        v-else
        v-bind="atribut(item)"
        :aria-current="item === page ? 'page' : undefined"
        :aria-label="`Halaman ${item}`"
        :class="[kelasTombol, item === page ? 'bg-secondary text-secondary-foreground' : 'border border-border-strong bg-surface hover:border-foreground']"
        @click="pilih($event, item)"
      >
        {{ item }}
      </component>
    </template>
    <component
      :is="tautan ? 'a' : 'button'"
      v-if="page < totalPages"
      v-bind="atribut(page + 1)"
      :rel="tautan ? 'next' : undefined"
      :class="[kelasTombol, 'border border-border-strong bg-surface hover:border-foreground']"
      @click="pilih($event, page + 1)"
    >
      Berikutnya
    </component>
  </nav>
</template>
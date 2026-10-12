<script setup>
import { computed } from 'vue'
import { CircleAlert, CircleCheck, CircleX, Clock } from 'lucide-vue-next'
import { LABEL_STOK } from '@/utils/produk'
const props = defineProps({
  status: { type: String, required: true },
  keterangan: { type: String, default: '' },
})
const STATUS = {
  TERSEDIA: { ikon: CircleCheck, kelas: 'bg-success-soft text-success-ink border-success/15' },
  SISA_STOK: { ikon: CircleAlert, kelas: 'bg-warning-soft text-warning-ink border-warning/15' },
  PRE_ORDER: { ikon: Clock, kelas: 'bg-info-soft text-info-ink border-info/15' },
  HABIS: { ikon: CircleX, kelas: 'bg-subtle text-foreground border-border' },
}
const info = computed(() => ({ ...(STATUS[props.status] ?? STATUS.HABIS), label: LABEL_STOK[props.status] ?? props.status }))
</script>
<template>
  <span class="inline-flex items-center gap-1.5 self-start rounded-full border px-2.5 py-1 text-xs font-semibold" :class="info.kelas">
    <component :is="info.ikon" class="size-3.5 shrink-0" aria-hidden="true" />
    <span>{{ info.label }}<template v-if="keterangan"> · {{ keterangan }}</template></span>
  </span>
</template>
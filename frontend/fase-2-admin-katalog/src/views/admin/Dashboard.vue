<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowRight, Plus } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import AktivitasTerbaru from '@/components/admin/AktivitasTerbaru.vue'
import PerhatianList from '@/components/admin/PerhatianList.vue'
import { usePermission } from '@/composables/usePermission'
import { getDashboard } from '@/services/adminService'
import { pesanError } from '@/services/errors'

// Aksi cepat tampil sesuai izin; isi ringkasan & "perlu perhatian" sudah difilter backend.
const AKSI_CEPAT = [
  { label: 'Tambah produk', to: '/admin/produk/baru', izin: 'catalog.product.create', varian: 'primary', ikon: Plus },
  { label: 'Input poin member', to: '/admin/poin/input', izin: 'point.create', varian: 'outline' },
  { label: 'Buat banner', to: '/admin/banner/baru', izin: 'banner.create', varian: 'outline' },
  { label: 'Tulis artikel', to: '/admin/artikel/baru', izin: 'article.create', varian: 'outline' },
]

const { can } = usePermission()
const status = ref('loading') // loading | ready | error
const data = ref(null)
const pesan = ref('')

const aksi = computed(() => AKSI_CEPAT.filter((a) => can(a.izin)))
const angka = new Intl.NumberFormat('id-ID')

async function muat() {
  status.value = 'loading'
  try {
    data.value = await getDashboard()
    status.value = 'ready'
  } catch (error) {
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

onMounted(muat)
</script>

<template>
  <div class="flex flex-col gap-6">
    <section v-if="aksi.length" aria-label="Aksi cepat" class="flex flex-wrap gap-2.5">
      <Button v-for="a in aksi" :key="a.to" :variant="a.varian" :to="a.to">
        <template v-if="a.ikon" #icon><component :is="a.ikon" class="size-4.5" aria-hidden="true" /></template>
        {{ a.label }}
      </Button>
    </section>

    <div v-if="status === 'loading'" class="flex flex-col gap-6" aria-busy="true">
      <span class="sr-only">Memuat dashboard…</span>
      <div class="grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-3">
        <Skeleton v-for="n in 3" :key="n" class="h-28 rounded-xl" />
      </div>
      <div class="grid gap-4 xl:grid-cols-2">
        <Skeleton class="h-72 rounded-xl" />
        <Skeleton class="h-72 rounded-xl" />
      </div>
    </div>

    <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

    <template v-else>
      <section v-if="data.summary.length" aria-label="Ringkasan" class="grid grid-cols-2 gap-3 sm:gap-4 lg:grid-cols-3 2xl:grid-cols-6">
        <component
          :is="k.link ? RouterLink : 'div'"
          v-for="k in data.summary"
          :key="k.key"
          :to="k.link ?? undefined"
          class="group flex flex-col gap-1.5 rounded-xl border border-border bg-surface p-5"
          :class="k.link && 'transition-colors duration-150 hover:border-border-strong focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'"
        >
          <span class="text-sm font-semibold text-muted">{{ k.label }}</span>
          <span class="text-3xl leading-tight font-extrabold tabular-nums">{{ angka.format(k.value) }}</span>
          <span v-if="k.link" class="inline-flex items-center gap-1 text-xs font-semibold text-primary">
            Lihat daftar
            <ArrowRight class="size-3.5 transition-transform duration-150 group-hover:translate-x-0.5" aria-hidden="true" />
          </span>
        </component>
      </section>

      <div class="grid items-start gap-4 xl:grid-cols-2">
        <PerhatianList :items="data.attention" :total="data.attention_total" />
        <AktivitasTerbaru :logs="data.logs" :scope="data.logs_scope" />
      </div>
    </template>
  </div>
</template>
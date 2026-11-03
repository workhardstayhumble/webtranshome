<script setup>
import { computed } from 'vue'
import { Film, ImageOff, Images, Plus } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Switch from '@/components/ui/Switch.vue'
import Table from '@/components/ui/Table.vue'
import TombolUrut from '@/components/admin/TombolUrut.vue'
import StatusBanner from '@/components/admin/konten/StatusBanner.vue'
import { useDaftarBanner } from '@/composables/useDaftarBanner'
import { usePermission } from '@/composables/usePermission'
import { ingatDaftar } from '@/composables/useTujuanKembali'
import { POSISI_BANNER } from '@/utils/banner'
import { formatTanggal } from '@/utils/format'

const d = useDaftarBanner()
const { can } = usePermission()
ingatDaftar('banner')

const KOLOM = [
  { key: 'banner', label: 'Banner', class: 'min-w-64' },
  { key: 'periode', label: 'Periode', class: 'whitespace-nowrap' },
  { key: 'status', label: 'Status' },
  { key: 'urutan', label: 'Urutan', align: 'center' },
  { key: 'aktif', label: 'Aktif', align: 'center' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]

const pesanNonaktif = computed(() => {
  const b = d.konfirmasi.value
  return b ? `Banner “${b.title}” tidak tampil lagi di Beranda. Data tetap tersimpan dan bisa diaktifkan lagi kapan saja.` : ''
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="text-sm text-muted">Banner tampil di Beranda selama periodenya. Urutan berlaku di dalam posisi yang sama.</p>
      <Button v-if="can('banner.create')" to="/admin/banner/baru">
        <template #icon><Plus class="size-4.5" aria-hidden="true" /></template>
        Buat banner
      </Button>
    </div>

    <div v-if="d.status.value === 'loading'" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4" aria-busy="true">
      <span class="sr-only">Memuat banner…</span>
      <Skeleton v-for="n in 5" :key="n" class="h-14" />
    </div>
    <ErrorState v-else-if="d.status.value === 'error'" :message="d.pesan.value" @retry="d.muat" />
    <EmptyState v-else-if="!d.daftar.value.length" title="Belum ada banner" description="Buat banner promo untuk slider di Beranda.">
      <template #icon><Images class="size-10" aria-hidden="true" /></template>
      <template v-if="can('banner.create')" #actions><Button to="/admin/banner/baru">Buat banner</Button></template>
    </EmptyState>

    <template v-else>
      <section
        v-for="(info, posisi) in POSISI_BANNER"
        :key="posisi"
        :aria-labelledby="`judul-${posisi}`"
        class="flex flex-col gap-3"
        :data-posisi-banner="posisi"
      >
        <h2 :id="`judul-${posisi}`" class="text-base font-bold">
          {{ info.label }} <span class="font-normal text-muted">({{ d.perPosisi.value[posisi].length }})</span>
        </h2>
        <p v-if="posisi === 'HOME_HERO'" class="-mt-2 text-sm text-muted">Paling atas Beranda: gambar atau video. Bila ada beberapa yang tayang, urutan pertama yang tampil.</p>
        <p v-if="!d.perPosisi.value[posisi].length" class="rounded-xl border border-dashed border-border-strong bg-surface px-4 py-6 text-center text-sm text-muted">
          Belum ada banner di posisi ini.
        </p>
        <Table v-else :columns="KOLOM" :rows="d.perPosisi.value[posisi]" :caption="`Banner ${info.label.toLowerCase()}`" min-width="min-w-200">
          <template #cell-banner="{ row }">
            <div class="flex items-center gap-3">
              <span class="flex h-12 w-24 shrink-0 items-center justify-center overflow-hidden rounded-md bg-subtle text-faint">
                <img v-if="row.image_url" :src="row.image_url" alt="" width="96" height="48" loading="lazy" class="size-full object-cover" />
                <ImageOff v-else class="size-4.5" aria-hidden="true" />
              </span>
              <span class="flex min-w-0 flex-col">
                <strong class="font-semibold">{{ row.title }}</strong>
                <span v-if="row.video_url" class="inline-flex items-center gap-1 text-xs font-semibold text-muted"><Film class="size-3.5" aria-hidden="true" />Dengan video</span>
                <span class="truncate font-mono text-xs text-faint">{{ row.link_url || 'Tanpa tautan' }}</span>
              </span>
            </div>
          </template>
          <template #cell-periode="{ row }">
            <span class="flex flex-col">
              <span>{{ formatTanggal(row.start_at) }} – {{ row.end_at ? formatTanggal(row.end_at) : '…' }}</span>
              <span v-if="!row.end_at" class="text-xs text-muted">Tanpa tanggal selesai</span>
            </span>
          </template>
          <template #cell-status="{ row }"><StatusBanner :status="row.status" /></template>
          <template #cell-urutan="{ row }">
            <TombolUrut
              :id="row.id"
              :nama="row.title"
              :pertama="d.perPosisi.value[posisi][0] === row"
              :terakhir="d.perPosisi.value[posisi].at(-1) === row"
              :disabled="!can('banner.update')"
              @geser="(arah) => d.geser(row, arah)"
            />
          </template>
          <template #cell-aktif="{ row }">
            <Switch
              :model-value="row.is_active"
              :label="`Aktif: ${row.title}`"
              hide-label
              :disabled="!can('banner.update') || row.sibuk"
              @update:model-value="(v) => d.ubahAktif(row, v)"
            />
          </template>
          <template #cell-aksi="{ row }">
            <div v-if="can('banner.update')" class="flex items-center justify-end gap-1.5 whitespace-nowrap">
              <Button v-if="row.status === 'BERAKHIR'" variant="secondary" size="sm" :to="`/admin/banner/${row.id}#periode`">
                Perpanjang <span class="sr-only">{{ row.title }}</span>
              </Button>
              <Button variant="outline" size="sm" :to="`/admin/banner/${row.id}`">Ubah <span class="sr-only">{{ row.title }}</span></Button>
            </div>
            <span v-else class="text-muted">–</span>
          </template>
        </Table>
      </section>
    </template>

    <p class="text-sm text-muted">Banner tidak dihapus permanen. Matikan "Aktif" agar tidak tampil. Setiap perubahan tercatat di log aktivitas.</p>

    <ConfirmModal
      :open="Boolean(d.konfirmasi.value)"
      title="Nonaktifkan banner?"
      :message="pesanNonaktif"
      confirm-label="Nonaktifkan"
      :loading="d.sedangNonaktif.value"
      @update:open="(buka) => !buka && (d.konfirmasi.value = null)"
      @confirm="d.nonaktifkan"
    />
  </div>
</template>

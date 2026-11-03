<script setup>
import { computed } from 'vue'
import { ImageOff, Plus, Wrench } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Switch from '@/components/ui/Switch.vue'
import Table from '@/components/ui/Table.vue'
import TombolUrut from '@/components/admin/TombolUrut.vue'
import { useDaftarJasa } from '@/composables/useDaftarJasa'
import { usePermission } from '@/composables/usePermission'
import { ingatDaftar } from '@/composables/useTujuanKembali'
import { tautanJasa } from '@/utils/jasa'

// /admin/jasa (KEPUTUSAN #20): daftar jasa, urutan (3 teratas tampil di Beranda), aktif/nonaktif. Jasa tidak dihapus permanen.
const d = useDaftarJasa()
const { can } = usePermission()
ingatDaftar('jasa')

const KOLOM = [
  { key: 'jasa', label: 'Jasa', class: 'min-w-64' },
  { key: 'kategori', label: 'Tampil di produk', class: 'min-w-48' },
  { key: 'urutan', label: 'Urutan', align: 'center' },
  { key: 'aktif', label: 'Aktif', align: 'center' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]

const pesanNonaktif = computed(() => {
  const j = d.konfirmasi.value
  return j ? `“${j.name}” tidak tampil lagi di website (halaman Jasa, Beranda, dan halaman produk). Data tetap tersimpan dan bisa diaktifkan lagi.` : ''
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="max-w-2xl text-sm text-muted">
        Jasa tampil di halaman Jasa, 3 teratas di Beranda, dan di halaman produk dari kategori yang dipilih. Harga tidak ditampilkan; semua jasa bertuliskan “Harga pasti murah”, dan pelanggan berkonsultasi lewat WhatsApp.
      </p>
      <Button v-if="can('service.create')" to="/admin/jasa/baru">
        <template #icon><Plus class="size-4.5" aria-hidden="true" /></template>
        Tambah jasa
      </Button>
    </div>

    <div v-if="d.status.value === 'loading'" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4" aria-busy="true">
      <span class="sr-only">Memuat jasa…</span>
      <Skeleton v-for="n in 5" :key="n" class="h-14" />
    </div>
    <ErrorState v-else-if="d.status.value === 'error'" :message="d.pesan.value" @retry="d.muat" />
    <EmptyState v-else-if="!d.daftar.value.length" title="Belum ada jasa" description="Tambahkan jasa yang bisa dipesan pelanggan, mis. pasang keramik atau antar material.">
      <template #icon><Wrench class="size-10" aria-hidden="true" /></template>
      <template v-if="can('service.create')" #actions><Button to="/admin/jasa/baru">Tambah jasa</Button></template>
    </EmptyState>

    <Table v-else :columns="KOLOM" :rows="d.urut.value" caption="Daftar jasa" min-width="min-w-200">
      <template #cell-jasa="{ row }">
        <div class="flex items-center gap-3">
          <span class="flex h-12 w-20 shrink-0 items-center justify-center overflow-hidden rounded-md bg-subtle text-faint">
            <img v-if="row.image_url" :src="row.image_url" alt="" width="80" height="48" loading="lazy" class="size-full object-cover" />
            <ImageOff v-else class="size-4.5" aria-hidden="true" />
          </span>
          <span class="flex min-w-0 flex-col">
            <strong class="font-semibold">{{ row.name }}</strong>
            <span class="truncate font-mono text-xs text-faint">{{ tautanJasa(row.slug) }}</span>
          </span>
        </div>
      </template>
      <template #cell-kategori="{ row }">
        <span v-if="row.categories?.length" class="text-sm">{{ row.categories.map((k) => k.name).join(', ') }}</span>
        <span v-else class="text-sm text-muted">Tidak dihubungkan</span>
      </template>
      <template #cell-urutan="{ row }">
        <TombolUrut
          :id="row.id"
          :nama="row.name"
          :pertama="d.urut.value[0] === row"
          :terakhir="d.urut.value.at(-1) === row"
          :disabled="!can('service.update')"
          @geser="(arah) => d.geser(row, arah)"
        />
      </template>
      <template #cell-aktif="{ row }">
        <Switch
          :model-value="row.is_active"
          :label="`Aktif: ${row.name}`"
          hide-label
          :disabled="!can('service.update') || row.sibuk"
          @update:model-value="(v) => d.ubahAktif(row, v)"
        />
      </template>
      <template #cell-aksi="{ row }">
        <div class="flex items-center justify-end gap-1.5 whitespace-nowrap">
          <Button v-if="row.is_active" variant="ghost" size="sm" :href="tautanJasa(row.slug)" target="_blank" rel="noopener">
            Lihat <span class="sr-only">{{ row.name }} di website (tab baru)</span>
          </Button>
          <Button v-if="can('service.update')" variant="outline" size="sm" :to="`/admin/jasa/${row.id}`">Ubah <span class="sr-only">{{ row.name }}</span></Button>
        </div>
      </template>
    </Table>

    <p class="text-sm text-muted">Jasa tidak dihapus permanen. Matikan "Aktif" agar tidak tampil. Setiap perubahan tercatat di log aktivitas.</p>

    <ConfirmModal
      :open="Boolean(d.konfirmasi.value)"
      title="Nonaktifkan jasa?"
      :message="pesanNonaktif"
      confirm-label="Nonaktifkan"
      :loading="d.sedangNonaktif.value"
      @update:open="(buka) => !buka && (d.konfirmasi.value = null)"
      @confirm="d.nonaktifkan"
    />
  </div>
</template>

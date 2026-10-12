<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { PackageSearch, Plus } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import ProdukFilterBar from '@/components/admin/produk/ProdukFilterBar.vue'
import ProdukTabel from '@/components/admin/produk/ProdukTabel.vue'
import { useFilter } from '@/composables/useFilter'
import { usePermission } from '@/composables/usePermission'
import { ingatDaftar } from '@/composables/useTujuanKembali'
import { listMaster } from '@/services/adminMasterService'
import { listAdminProducts, patchProduct } from '@/services/adminProductService'
import { pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'
import { perilakuGulir } from '@/utils/gerak'

// Filter disimpan di URL (kembali dari form produk → filter & halaman sama).
const KUNCI_FILTER = ['q', 'category', 'brand', 'stock', 'active', 'foto', 'room']

const { nilai, terapkan, reset } = useFilter()
const { can } = usePermission()
const ui = useUiStore()
ingatDaftar('produk')

const status = ref('loading') // loading | ready | error
const produk = ref([])
const meta = ref(null)
const pesan = ref('')
const master = ref({ kategori: [], brand: [], ruangan: [] })
const atasTabel = ref(null)
const konfirmasi = ref(null) // produk yang akan dinonaktifkan
const sedangNonaktif = ref(false)
let pengendali = null

const params = computed(() => {
  const hasil = { per_page: 20, page: nilai.value.page || undefined }
  for (const k of KUNCI_FILTER) if (nilai.value[k]) hasil[k] = nilai.value[k]
  return hasil
})
const adaFilter = computed(() => KUNCI_FILTER.some((k) => nilai.value[k]))
const rentang = computed(() => {
  if (!meta.value?.total) return ''
  const awal = (meta.value.page - 1) * meta.value.per_page + 1
  return `Menampilkan ${awal}–${awal + produk.value.length - 1} dari ${meta.value.total} produk`
})

async function muat() {
  pengendali?.abort()
  pengendali = new AbortController()
  status.value = 'loading'
  try {
    const hasil = await listAdminProducts(params.value, { signal: pengendali.signal })
    produk.value = hasil.data
    meta.value = hasil.meta
    status.value = 'ready'
  } catch (error) {
    if (error?.code === 'ERR_CANCELED') return
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

watch(() => JSON.stringify(params.value), muat, { immediate: true })
onBeforeUnmount(() => pengendali?.abort())

// Pilihan filter; kalau salah satu gagal dimuat, filter itu saja yang kosong.
onMounted(async () => {
  const [kategori, brand, ruangan] = await Promise.allSettled([listMaster('categories'), listMaster('brands'), listMaster('rooms')])
  master.value = { kategori: kategori.value ?? [], brand: brand.value ?? [], ruangan: ruangan.value ?? [] }
})

async function gantiHalaman(nomor) {
  await terapkan({ page: nomor > 1 ? nomor : null })
  atasTabel.value?.scrollIntoView({ behavior: perilakuGulir(), block: 'start' })
}

// Toggle cepat: UI langsung berubah, dikembalikan + toast bila gagal (Aturan Workflow Tombol).
async function ubahUnggulan(p) {
  const lama = p.is_featured
  p.is_featured = !lama
  p.sibuk = true
  try {
    await patchProduct(p.id, { is_featured: p.is_featured })
  } catch (error) {
    p.is_featured = lama
    ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    p.sibuk = false
  }
}

async function aktifkan(p) {
  p.is_active = true
  p.sibuk = true
  try {
    await patchProduct(p.id, { is_active: true })
    ui.tampilkanToast({ pesan: `${p.name} tampil di website.`, jenis: 'success' })
  } catch (error) {
    p.is_active = false
    ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    p.sibuk = false
  }
}

function ubahAktif(p, nilaiBaru) {
  if (nilaiBaru) aktifkan(p)
  else konfirmasi.value = p
}

async function nonaktifkan() {
  const p = konfirmasi.value
  sedangNonaktif.value = true
  try {
    await patchProduct(p.id, { is_active: false })
    p.is_active = false
    konfirmasi.value = null
    ui.tampilkanToast({ pesan: `${p.name} tidak tampil lagi di website.`, jenis: 'success' })
  } catch (error) {
    ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    sedangNonaktif.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="text-sm text-muted">Produk yang aktif tampil di katalog website. Cari, saring, lalu ubah dari sini.</p>
      <Button v-if="can('product.create')" to="/admin/produk/baru">
        <template #icon><Plus class="size-4.5" aria-hidden="true" /></template>
        Tambah produk
      </Button>
    </div>

    <ProdukFilterBar
      :nilai="nilai"
      :kategori="master.kategori"
      :brand="master.brand"
      :ruangan="master.ruangan"
      @ubah="terapkan"
      @reset="reset()"
    />

    <section ref="atasTabel" aria-label="Daftar produk" class="scroll-mt-24" :aria-busy="status === 'loading' ? 'true' : 'false'">
      <div v-if="status === 'loading'" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-4">
        <span class="sr-only">Memuat produk…</span>
        <div v-for="n in 8" :key="n" class="flex items-center gap-3">
          <Skeleton class="size-11 shrink-0" />
          <Skeleton class="h-4 flex-1" />
          <Skeleton class="h-4 w-24" />
          <Skeleton class="h-4 w-20" />
        </div>
      </div>

      <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

      <EmptyState
        v-else-if="!produk.length"
        :title="adaFilter ? 'Tidak ada produk yang cocok' : 'Belum ada produk'"
        :description="adaFilter ? 'Coba ubah kata kunci atau kurangi filter.' : 'Produk yang ditambahkan akan muncul di sini.'"
      >
        <template #icon><PackageSearch class="size-10" aria-hidden="true" /></template>
        <template #actions>
          <Button v-if="adaFilter" variant="outline" @click="reset()">Reset filter</Button>
          <Button v-else-if="can('product.create')" to="/admin/produk/baru">Tambah produk</Button>
        </template>
      </EmptyState>

      <div v-else class="overflow-hidden rounded-xl border border-border bg-surface shadow-[0_3px_20px_rgba(28,36,48,0.04)]">
        <ProdukTabel :produk="produk" @unggulan="ubahUnggulan" @aktif="ubahAktif" />
        <div class="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
          <span class="text-sm text-muted" aria-live="polite">{{ rentang }}</span>
          <PaginationBar v-if="meta.total_pages > 1" :page="meta.page" :total-pages="meta.total_pages" @update:page="gantiHalaman" />
        </div>
      </div>
    </section>

    <p class="text-sm text-muted">
      Produk tidak dihapus permanen. Matikan "Aktif" agar tidak tampil di website. Setiap perubahan tercatat di log aktivitas.
    </p>

    <ConfirmModal
      :open="Boolean(konfirmasi)"
      title="Nonaktifkan produk?"
      :message="konfirmasi ? `“${konfirmasi.name}” tidak akan tampil di website. Data tetap tersimpan dan bisa diaktifkan lagi kapan saja.` : ''"
      confirm-label="Nonaktifkan"
      :loading="sedangNonaktif"
      @update:open="(buka) => !buka && (konfirmasi = null)"
      @confirm="nonaktifkan"
    />
  </div>
</template>
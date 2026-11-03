<script setup>
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CircleCheck, LoaderCircle, Plus } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Combobox from '@/components/ui/Combobox.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Switch from '@/components/ui/Switch.vue'
import Textarea from '@/components/ui/Textarea.vue'
import BrandCepatModal from '@/components/admin/BrandCepatModal.vue'
import FormSection from '@/components/admin/FormSection.vue'
import DatasheetField from '@/components/admin/produk/DatasheetField.vue'
import FotoProduk from '@/components/admin/produk/FotoProduk.vue'
import RuanganPicker from '@/components/admin/produk/RuanganPicker.vue'
import SpesifikasiEditor from '@/components/admin/produk/SpesifikasiEditor.vue'
import { useFormProduk } from '@/composables/useFormProduk'
import { useKonfirmasiKeluar } from '@/composables/useKonfirmasiKeluar'
import { usePermission } from '@/composables/usePermission'
import { tujuanKembali } from '@/composables/useTujuanKembali'
import { listMaster } from '@/services/adminMasterService'
import { useUiStore } from '@/stores/ui'
import { formatRupiah } from '@/utils/format'
import { LABEL_STOK } from '@/utils/produk'

const route = useRoute()
const router = useRouter()
const ui = useUiStore()
const { can } = usePermission()
const f = useFormProduk()
const { form, errors } = f
const keluar = useKonfirmasiKeluar(f.kotor, { saatSesiHabis: f.simpanDraf })

const master = ref({ kategori: [], brand: [], ruangan: [] })
const brandBaruBuka = ref(false)
const formEl = ref(null)

const OPSI_STOK = Object.entries(LABEL_STOK).map(([value, label]) => ({ value, label }))
const opsiKategori = computed(() =>
  master.value.kategori
    .filter((k) => k.parent_id == null)
    .map((p) => {
      const anak = master.value.kategori.filter((k) => k.parent_id === p.id)
      const label = (k) => (k.is_active ? k.name : `${k.name} (nonaktif)`)
      return { label: p.name, options: (anak.length ? anak : [p]).map((k) => ({ value: k.id, label: label(k) })) }
    }),
)
const opsiBrand = computed(() => [
  { value: '', label: 'Tanpa brand' },
  ...master.value.brand.map((b) => ({ value: b.id, label: b.is_active ? b.name : `${b.name} (nonaktif)` })),
])
const pratinjauHarga = computed(() => (form.price_general ? `${formatRupiah(form.price_general)} per ${form.unit_sale || 'satuan'}` : ''))

async function fokusGalatPertama() {
  await nextTick()
  const el = formEl.value?.querySelector('[aria-invalid="true"], [data-galat]')
  el?.scrollIntoView({ block: 'center' })
  if (el?.matches('input, select, textarea')) el.focus({ preventScroll: true })
}

async function simpan(mode = 'simpan') {
  const berhasil = await f.kirim(mode)
  if (!berhasil) {
    fokusGalatPertama()
    return
  }
  ui.tampilkanToast({ pesan: 'Produk disimpan', jenis: 'success' })
  if (mode === 'tambah-lagi') {
    f.kosongkan()
    window.scrollTo({ top: 0 })
    await nextTick()
    formEl.value?.querySelector('input')?.focus()
    return
  }
  f.tandaiTersimpan()
  router.push(tujuanKembali('produk', '/admin/produk'))
}

function batal() {
  router.push(tujuanKembali('produk', '/admin/produk'))
}

function brandDibuat(brand) {
  master.value.brand = [...master.value.brand, brand].sort((a, b) => a.name.localeCompare(b.name, 'id'))
  form.brand_id = brand.id
  ui.tampilkanToast({ pesan: `Brand ${brand.name} ditambahkan.`, jenis: 'success' })
}

onMounted(async () => {
  const [kategori, brand, ruangan] = await Promise.allSettled([listMaster('categories'), listMaster('brands'), listMaster('rooms'), f.muat()])
  master.value = { kategori: kategori.value ?? [], brand: brand.value ?? [], ruangan: ruangan.value ?? [] }
  // Tautan dashboard: /admin/produk/:id#stok atau #foto langsung ke bagiannya.
  if (route.hash && f.status.value === 'ready') {
    await nextTick()
    const bagian = document.getElementById(route.hash.slice(1))
    bagian?.scrollIntoView({ block: 'start' })
    bagian?.querySelector('input, select, button')?.focus({ preventScroll: true })
  }
})
</script>

<template>
  <div v-if="f.status.value === 'loading'" class="grid gap-6 xl:grid-cols-[minmax(0,1fr)_20rem]" aria-busy="true">
    <span class="sr-only">Memuat produk…</span>
    <div class="flex flex-col gap-6">
      <Skeleton class="h-72 rounded-xl" />
      <Skeleton class="h-56 rounded-xl" />
    </div>
    <Skeleton class="h-80 rounded-xl" />
  </div>

  <section v-else-if="f.status.value === 'tidak-ada'" class="mx-auto flex max-w-md flex-col items-center gap-4 py-16 text-center">
    <h2 class="text-2xl font-bold">Produk tidak ditemukan</h2>
    <p class="text-muted">Produk ini tidak ada atau alamatnya salah.</p>
    <Button to="/admin/produk">Kembali ke daftar produk</Button>
  </section>

  <ErrorState v-else-if="f.status.value === 'error'" :message="f.pesan.value" @retry="f.muat" />

  <form v-else ref="formEl" novalidate class="flex flex-col gap-6" @submit.prevent="simpan('simpan')">
    <p v-if="f.dipulihkan.value" role="status" class="rounded-md bg-info-soft px-4 py-3 text-sm font-medium text-info-ink">
      Isian sebelum sesi habis sudah dipulihkan. Periksa lagi, lalu simpan.
    </p>
    <p v-if="f.galatUmum.value" role="alert" class="rounded-md bg-danger-soft px-4 py-3 text-sm font-semibold text-danger-ink">{{ f.galatUmum.value }}</p>

    <div class="grid items-start gap-6 xl:grid-cols-[minmax(0,1fr)_20rem]">
      <div class="flex min-w-0 flex-col gap-6">
        <FormSection judul="Informasi produk">
          <div class="grid gap-4 sm:grid-cols-2">
            <Input
              v-model="form.sku"
              label="SKU"
              required
              autocomplete="off"
              class="font-mono"
              :error="errors.sku ?? ''"
              :hint="f.statusSku.value === 'tersedia' ? 'SKU tersedia' : 'Kode unik dari Master Barang'"
              @blur="f.periksaSku"
            >
              <template #akhir>
                <LoaderCircle v-if="f.statusSku.value === 'memeriksa'" class="mr-3 size-4.5 animate-spin text-muted" aria-hidden="true" />
                <CircleCheck v-else-if="f.statusSku.value === 'tersedia'" class="mr-3 size-4.5 text-success" aria-hidden="true" />
              </template>
            </Input>
            <Input :model-value="form.name" label="Nama produk" required autocomplete="off" :error="errors.name ?? ''" @update:model-value="f.ubahNama" />
          </div>
          <Input
            :model-value="form.slug"
            label="Slug (alamat halaman)"
            required
            autocomplete="off"
            prefix="/produk/"
            :error="errors.slug ?? ''"
            hint="Dibuat otomatis dari nama. Boleh diubah: huruf kecil, angka, tanda hubung."
            @update:model-value="f.ubahSlug"
          />
          <Textarea v-model="form.description" label="Deskripsi" :rows="5" hint="Ceritakan kegunaan dan keunggulan produk untuk pelanggan." />
        </FormSection>

        <FormSection id="stok" judul="Harga dan stok" deskripsi="Satu harga umum untuk semua pelanggan (Muka 1).">
          <div class="flex flex-col gap-1">
            <Switch v-model="form.show_price" label="Tampilkan harga di website" />
            <p class="text-sm text-muted">
              {{ form.show_price ? 'Pelanggan melihat harga di katalog dan halaman produk.' : 'Pelanggan melihat “Harga via WhatsApp” dan tombol “Tanya harga”. Harga boleh dikosongkan.' }}
            </p>
          </div>
          <div class="grid gap-4 sm:grid-cols-3">
            <Input
              v-model="form.price_general"
              label="Harga umum"
              :required="form.show_price"
              prefix="Rp"
              inputmode="numeric"
              autocomplete="off"
              :error="errors.price_general ?? ''"
              :hint="pratinjauHarga || (form.show_price ? '' : 'Opsional, hanya untuk catatan admin.')"
            />
            <Input v-model="form.unit_sale" label="Satuan jual" required placeholder="dus, sak, pcs" autocomplete="off" :error="errors.unit_sale ?? ''" />
            <Input v-model="form.min_order" label="Minimal pembelian" type="number" min="1" inputmode="numeric" :error="errors.min_order ?? ''" />
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <Select v-model="form.stock_status" label="Status stok" :options="OPSI_STOK" />
            <Input v-model="form.stock_qty_label" label="Keterangan stok" placeholder="Mis. 3 dus" autocomplete="off" hint="Opsional, tampil di samping status (mis. Sisa stok · 3 dus)." />
          </div>
        </FormSection>

        <FormSection id="foto" judul="Foto produk" deskripsi="Minimal 1 foto utama agar produk bisa tampil di website. Foto utama tampil di kartu katalog.">
          <FotoProduk v-model="form.images" :error="errors.images ?? ''" />
        </FormSection>

        <FormSection judul="Spesifikasi" deskripsi="Tampil di tab Spesifikasi. Ukuran, finishing, warna, bahan, dan daya juga menjadi filter katalog.">
          <SpesifikasiEditor v-model="form.spesifikasi" :errors="errors.spesifikasi ?? []" />
        </FormSection>

        <FormSection judul="Dokumen">
          <DatasheetField v-model="form.datasheet_pdf_url" />
        </FormSection>

        <FormSection judul="SEO" deskripsi="Opsional. Kosongkan untuk memakai nama dan ringkasan produk.">
          <Input v-model="form.meta_title" label="Judul halaman (meta title)" maxlength="70" autocomplete="off" />
          <Textarea v-model="form.meta_description" label="Deskripsi singkat (meta description)" :rows="3" :maxlength="160" />
        </FormSection>
      </div>

      <div class="flex flex-col gap-6 xl:sticky xl:top-24">
        <FormSection judul="Status">
          <div class="flex flex-col">
            <Switch v-model="form.is_active" label="Tampil di website" />
            <Switch v-model="form.is_featured" label="Produk unggulan" />
          </div>
        </FormSection>

        <FormSection judul="Kategori dan brand">
          <Select v-model="form.category_id" label="Kategori" required placeholder="Pilih kategori" :groups="opsiKategori" :error="errors.category_id ?? ''" />
          <div class="flex flex-col gap-2">
            <Combobox v-model="form.brand_id" label="Brand" :options="opsiBrand" :error="errors.brand_id ?? ''" kosong="Brand tidak ditemukan. Tambahkan di bawah." />
            <Button v-if="can('brand.create')" variant="ghost" size="sm" class="self-start" @click="brandBaruBuka = true">
              <template #icon><Plus class="size-4" aria-hidden="true" /></template>
              Tambah brand
            </Button>
          </div>
        </FormSection>

        <FormSection judul="Ruangan" deskripsi="Produk tampil di halaman ruangan yang dipilih.">
          <RuanganPicker v-model="form.room_ids" :ruangan="master.ruangan" />
        </FormSection>
      </div>
    </div>

    <!-- HP: [Simpan & tambah lagi] selebar layar, lalu [Batal][Simpan]. Mulai 640px: satu baris di kanan. -->
    <div class="sticky bottom-0 z-20 -mx-4 grid grid-cols-2 gap-2 border-t border-border bg-surface px-4 py-3 sm:-mx-7 sm:flex sm:flex-wrap sm:items-center sm:justify-end sm:px-7">
      <span v-if="f.kotor.value" class="hidden text-sm text-muted sm:mr-auto sm:inline">Ada perubahan yang belum disimpan</span>
      <Button
        v-if="!f.id"
        variant="secondary"
        class="order-first col-span-2 sm:order-none sm:col-auto"
        :loading="f.menyimpan.value === 'tambah-lagi'"
        :disabled="Boolean(f.menyimpan.value)"
        @click="simpan('tambah-lagi')"
      >
        Simpan & tambah lagi
      </Button>
      <Button variant="outline" class="sm:order-first" :disabled="Boolean(f.menyimpan.value)" @click="batal">Batal</Button>
      <Button type="submit" :loading="f.menyimpan.value === 'simpan'" :disabled="Boolean(f.menyimpan.value)">Simpan</Button>
    </div>

    <BrandCepatModal v-model:open="brandBaruBuka" @dibuat="brandDibuat" />
    <ConfirmModal
      :open="keluar.terbuka.value"
      title="Buang perubahan?"
      message="Perubahan yang belum disimpan akan hilang."
      confirm-label="Buang perubahan"
      cancel-label="Tetap di sini"
      danger
      @confirm="keluar.buang"
      @cancel="keluar.tetap"
    />
  </form>
</template>

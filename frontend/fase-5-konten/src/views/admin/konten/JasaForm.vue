<script setup>
import { nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Checkbox from '@/components/ui/Checkbox.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Switch from '@/components/ui/Switch.vue'
import Textarea from '@/components/ui/Textarea.vue'
import FormSection from '@/components/admin/FormSection.vue'
import GambarField from '@/components/admin/GambarField.vue'
import { useFormJasa } from '@/composables/useFormJasa'
import { useKonfirmasiKeluar } from '@/composables/useKonfirmasiKeluar'
import { tujuanKembali } from '@/composables/useTujuanKembali'
import { getCategoryTree } from '@/services/categoryService'
import { useUiStore } from '@/stores/ui'
import { LABEL_HARGA_JASA } from '@/utils/jasa'

// /admin/jasa/baru dan /admin/jasa/:id (KEPUTUSAN #20).
const router = useRouter()
const ui = useUiStore()
const f = useFormJasa()
const { form, errors } = f
const keluar = useKonfirmasiKeluar(f.kotor, { saatSesiHabis: f.simpanDraf })
const formEl = ref(null)

// Kategori produk: jasa tampil di halaman produk dari kategori yang dicentang (induk = semua anaknya).
const pohon = ref([])
const galatKategori = ref('')
const dipilih = (id) => form.category_ids.includes(id)

async function fokusGalatPertama() {
  await nextTick()
  const el = formEl.value?.querySelector('[aria-invalid="true"], [data-galat]')
  el?.scrollIntoView({ block: 'center' })
  if (el?.matches('input, select, textarea')) el.focus({ preventScroll: true })
}

async function simpan() {
  const hasil = await f.kirim()
  if (!hasil) {
    fokusGalatPertama()
    return
  }
  ui.tampilkanToast({ pesan: 'Jasa disimpan', jenis: 'success' })
  f.tandaiTersimpan()
  router.push(tujuanKembali('jasa', '/admin/jasa'))
}

const batal = () => router.push(tujuanKembali('jasa', '/admin/jasa'))

onMounted(async () => {
  f.muat()
  try {
    pohon.value = await getCategoryTree()
  } catch {
    galatKategori.value = 'Daftar kategori belum bisa dimuat. Muat ulang halaman untuk mencoba lagi.'
  }
})
</script>

<template>
  <div v-if="f.status.value === 'loading'" class="grid gap-6 xl:grid-cols-[minmax(0,1fr)_20rem]" aria-busy="true">
    <span class="sr-only">Memuat jasa…</span>
    <Skeleton class="h-96 rounded-xl" />
    <Skeleton class="h-56 rounded-xl" />
  </div>

  <section v-else-if="f.status.value === 'tidak-ada'" class="mx-auto flex max-w-md flex-col items-center gap-4 py-16 text-center">
    <h2 class="text-2xl font-bold">Jasa tidak ditemukan</h2>
    <p class="text-muted">Jasa ini tidak ada atau alamatnya salah.</p>
    <Button to="/admin/jasa">Kembali ke daftar jasa</Button>
  </section>

  <ErrorState v-else-if="f.status.value === 'error'" :message="f.pesan.value" @retry="f.muat" />

  <form v-else ref="formEl" novalidate class="flex flex-col gap-6" @submit.prevent="simpan">
    <p v-if="f.dipulihkan.value" role="status" class="rounded-md bg-info-soft px-4 py-3 text-sm font-medium text-info-ink">
      Isian sebelum sesi habis sudah dipulihkan. Periksa lagi, lalu simpan.
    </p>
    <p v-if="f.galatUmum.value" role="alert" class="rounded-md bg-danger-soft px-4 py-3 text-sm font-semibold text-danger-ink">{{ f.galatUmum.value }}</p>

    <div class="grid items-start gap-6 xl:grid-cols-[minmax(0,1fr)_20rem]">
      <div class="flex min-w-0 flex-col gap-6">
        <FormSection judul="Isi jasa">
          <Input :model-value="form.name" label="Nama jasa" required maxlength="100" autocomplete="off" placeholder="Jasa pasang keramik & granit" :error="errors.name ?? ''" @update:model-value="f.ubahNama" />
          <Input
            :model-value="form.slug"
            label="Slug"
            required
            autocomplete="off"
            :error="errors.slug ?? ''"
            hint="Alamat halaman: /jasa/slug. Dibuat otomatis dari nama; huruf kecil, angka, tanda hubung."
            @update:model-value="f.ubahSlug"
          />
          <Textarea v-model="form.summary" label="Ringkasan" required :rows="2" :maxlength="160" :error="errors.summary ?? ''" hint="Satu kalimat; tampil di kartu jasa dan hasil pencarian Google." />
          <GambarField v-model="form.image_url" label="Gambar jasa" rasio="aspect-video" lebar="w-64" :error="errors.image_url ?? ''" hint="JPG, PNG, atau WebP, maksimal 2 MB. Ukuran disarankan 1200 × 675 px. Opsional." />
          <Textarea
            v-model="form.description"
            label="Penjelasan"
            :rows="12"
            :error="errors.description ?? ''"
            hint="Baris kosong = paragraf baru. Awali baris dengan ## untuk subjudul dan - untuk daftar, mis. ## Yang dikerjakan."
          />
        </FormSection>

        <FormSection judul="Tampil di halaman produk" deskripsi="Centang kategori produk yang cocok dengan jasa ini. Kategori induk berlaku untuk semua kategori di bawahnya.">
          <p v-if="galatKategori" role="alert" class="text-sm font-medium text-danger-ink">{{ galatKategori }}</p>
          <div v-else class="grid gap-4 sm:grid-cols-2" data-pilih-kategori>
            <fieldset v-for="induk in pohon" :key="induk.id" class="flex flex-col gap-1.5 rounded-lg border border-border p-3">
              <legend class="sr-only">{{ induk.name }}</legend>
              <Checkbox :model-value="dipilih(induk.id)" :label="`${induk.name} (semua)`" @update:model-value="(v) => f.ubahKategori(induk.id, v)" />
              <div v-if="induk.children?.length" class="flex flex-col gap-1.5 pl-6">
                <Checkbox
                  v-for="anak in induk.children"
                  :key="anak.id"
                  :model-value="dipilih(induk.id) || dipilih(anak.id)"
                  :disabled="dipilih(induk.id)"
                  :label="anak.name"
                  @update:model-value="(v) => f.ubahKategori(anak.id, v)"
                />
              </div>
            </fieldset>
          </div>
          <p v-if="errors.category_ids" data-galat class="text-sm font-medium text-danger-ink">{{ errors.category_ids }}</p>
        </FormSection>
      </div>

      <div class="flex flex-col gap-6 xl:sticky xl:top-24">
        <FormSection judul="Area & harga">
          <Input v-model="form.service_area" label="Area layanan" maxlength="100" autocomplete="off" placeholder="Sleman dan sekitarnya" :error="errors.service_area ?? ''" />
          <p class="rounded-md bg-background px-3 py-2 text-sm" data-label-harga>
            <span class="text-muted">Harga tidak ditampilkan. Semua jasa bertuliskan</span> <strong>“{{ LABEL_HARGA_JASA }}”</strong><span class="text-muted">; biaya dibahas lewat WhatsApp.</span>
          </p>
        </FormSection>
        <FormSection judul="Tampil di website">
          <Switch v-model="form.is_active" label="Aktif" />
          <p class="text-sm text-muted">Jasa nonaktif tidak tampil di halaman Jasa, Beranda, maupun halaman produk.</p>
        </FormSection>
      </div>
    </div>

    <div class="sticky bottom-0 z-20 -mx-4 flex flex-wrap items-center justify-end gap-2 border-t border-border bg-surface px-4 py-3 sm:-mx-7 sm:px-7">
      <span v-if="f.kotor.value" class="mr-auto hidden text-sm text-muted sm:inline">Ada perubahan yang belum disimpan</span>
      <Button variant="outline" :disabled="f.menyimpan.value" @click="batal">Batal</Button>
      <Button type="submit" :loading="f.menyimpan.value">Simpan</Button>
    </div>

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

<script setup>
import { computed, nextTick, reactive, ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Modal from '@/components/ui/Modal.vue'
import Select from '@/components/ui/Select.vue'
import Switch from '@/components/ui/Switch.vue'
import GambarField from './GambarField.vue'
import { createMaster, updateMaster } from '@/services/adminMasterService'
import { errorField, pesanError } from '@/services/errors'
import { buatSlug, POLA_SLUG } from '@/utils/slug'

// Form tambah/ubah kategori, brand, dan ruangan dalam modal.
const open = defineModel('open', { type: Boolean, default: false })

const props = defineProps({
  jenis: { type: String, required: true, validator: (v) => ['categories', 'brands', 'rooms'].includes(v) },
  // null = tambah baru
  baris: { type: Object, default: null },
  // Kategori utama untuk pilihan induk (khusus kategori).
  induk: { type: Array, default: () => [] },
})

const emit = defineEmits(['tersimpan'])

const KONFIG = {
  categories: { nama: 'kategori', awalan: '/katalog/', gambar: 'image_url', labelGambar: 'Gambar kategori', rasio: 'aspect-square', wajib: false },
  brands: {
    nama: 'brand', awalan: '/brand/', gambar: 'logo_url', labelGambar: 'Logo brand', rasio: 'aspect-square', wajib: false,
    hint: 'JPG, PNG, atau WebP, maksimal 2 MB. Tampil di lingkaran "Brand pilihan" Beranda dan halaman Brand; pakai logo persegi berlatar putih atau transparan. Brand berlogo tampil lebih dulu di Beranda.',
  },
  rooms: { nama: 'ruangan', awalan: '/ruangan/', gambar: 'image_cover', labelGambar: 'Cover ruangan', rasio: 'aspect-4/3', wajib: true },
}
const k = computed(() => KONFIG[props.jenis])

const form = reactive({ name: '', slug: '', parent_id: '', gambar: '', is_active: true })
const errors = ref({})
const galatUmum = ref('')
const menyimpan = ref(false)
const slugManual = ref(false)
const formEl = ref(null)
const idForm = `form-master-${props.jenis}`

watch(open, (buka) => {
  if (!buka) return
  const b = props.baris
  Object.assign(form, {
    name: b?.name ?? '',
    slug: b?.slug ?? '',
    parent_id: b?.parent_id ?? '',
    gambar: b?.[k.value.gambar] ?? '',
    is_active: b?.is_active ?? true,
  })
  slugManual.value = Boolean(b)
  errors.value = {}
  galatUmum.value = ''
})

const opsiInduk = computed(() => [
  { value: '', label: 'Tidak ada (kategori utama)' },
  ...props.induk.filter((x) => x.id !== props.baris?.id).map((x) => ({ value: x.id, label: x.name })),
])

function ubahNama(nilai) {
  form.name = nilai
  if (!slugManual.value) form.slug = buatSlug(nilai)
}

function ubahSlug(nilai) {
  slugManual.value = true
  form.slug = nilai
}

function validasi() {
  const e = {}
  if (!form.name.trim()) e.name = 'Nama wajib diisi'
  if (!form.slug) e.slug = 'Slug wajib diisi'
  else if (!POLA_SLUG.test(form.slug)) e.slug = 'Slug hanya boleh huruf kecil, angka, dan tanda hubung'
  if (k.value.wajib && !form.gambar) e[k.value.gambar] = `${k.value.labelGambar} wajib diunggah`
  errors.value = e
  return !Object.keys(e).length
}

async function simpan() {
  if (menyimpan.value) return
  galatUmum.value = ''
  if (!validasi()) {
    await nextTick()
    formEl.value?.querySelector('[aria-invalid="true"]')?.focus()
    return
  }
  const payload = { name: form.name.trim(), slug: form.slug, [k.value.gambar]: form.gambar || null, is_active: form.is_active }
  if (props.jenis === 'categories') payload.parent_id = form.parent_id === '' ? null : Number(form.parent_id)
  menyimpan.value = true
  try {
    const hasil = props.baris ? await updateMaster(props.jenis, props.baris.id, payload) : await createMaster(props.jenis, payload)
    emit('tersimpan', hasil, { baru: !props.baris })
    open.value = false
  } catch (error) {
    if ([409, 422].includes(error?.response?.status)) errors.value = errorField(error)
    else if (error?.response?.status !== 401) galatUmum.value = pesanError(error)
  } finally {
    menyimpan.value = false
  }
}
</script>

<template>
  <Modal v-model:open="open" :title="`${baris ? 'Ubah' : 'Tambah'} ${k.nama}`" :dismissible="!menyimpan">
    <form :id="idForm" ref="formEl" novalidate class="flex flex-col gap-4" @submit.prevent="simpan">
      <p v-if="galatUmum" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galatUmum }}</p>
      <Input :model-value="form.name" label="Nama" required autocomplete="off" :error="errors.name ?? ''" @update:model-value="ubahNama" />
      <Input
        :model-value="form.slug"
        label="Slug"
        required
        autocomplete="off"
        :prefix="k.awalan"
        :error="errors.slug ?? ''"
        hint="Dibuat otomatis dari nama. Huruf kecil, angka, tanda hubung."
        @update:model-value="ubahSlug"
      />
      <Select v-if="jenis === 'categories'" v-model="form.parent_id" label="Induk kategori" :options="opsiInduk" :error="errors.parent_id ?? ''" />
      <GambarField v-model="form.gambar" :label="k.labelGambar" :required="k.wajib" :rasio="k.rasio" :hint="k.hint" :error="errors[k.gambar] ?? ''" />
      <Switch v-model="form.is_active" label="Tampil di website" />
    </form>
    <template #footer>
      <Button variant="outline" :disabled="menyimpan" @click="open = false">Batal</Button>
      <Button type="submit" :form="idForm" :loading="menyimpan">Simpan</Button>
    </template>
  </Modal>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Modal from '@/components/ui/Modal.vue'
import Switch from '@/components/ui/Switch.vue'
import { patchProduct } from '@/services/adminProductService'
import { errorField, pesanError } from '@/services/errors'
import { formatRupiah } from '@/utils/format'

// Ubah cepat harga dari daftar produk, tanpa membuka form lengkap: harga, satuan jual, dan tampil/sembunyikan harga.
const open = defineModel('open', { type: Boolean, default: false })

const props = defineProps({
  // Baris GET /admin/products: { id, name, sku, price_general, unit_sale, show_price }
  produk: { type: Object, default: null },
})

const emit = defineEmits(['tersimpan'])

const angka = (teks) => String(teks ?? '').replace(/\D/g, '')
const isian = reactive({ price_general: '', unit_sale: '', show_price: true })
const errors = ref({})
const pesan = ref('')
const menyimpan = ref(false)

watch(open, (buka) => {
  if (!buka || !props.produk) return
  const p = props.produk
  Object.assign(isian, {
    price_general: p.price_general == null ? '' : Number(p.price_general).toFixed(0),
    unit_sale: p.unit_sale ?? '',
    show_price: p.show_price !== false,
  })
  errors.value = {}
  pesan.value = ''
})

// Error satu field hilang begitu isiannya diubah.
watch(
  () => ({ ...isian }),
  (baru, lama) => {
    const sisa = { ...errors.value }
    for (const k of Object.keys(baru)) if (baru[k] !== lama[k]) delete sisa[k]
    // Harga disembunyikan: "harga wajib" tidak berlaku lagi.
    if (!baru.show_price) delete sisa.price_general
    errors.value = sisa
  },
)

const pratinjau = computed(() => {
  if (!isian.show_price) return 'Pelanggan melihat “Harga via WhatsApp” dan tombol “Tanya harga”.'
  return angka(isian.price_general) ? `Tampil: ${formatRupiah(angka(isian.price_general))} per ${isian.unit_sale.trim() || 'satuan'}` : ''
})

function periksa() {
  const e = {}
  if (isian.show_price && !angka(isian.price_general)) e.price_general = 'Harga wajib diisi, atau matikan "Tampilkan harga di website"'
  if (!isian.unit_sale.trim()) e.unit_sale = 'Satuan jual wajib diisi, mis. dus atau sak'
  errors.value = e
  return !Object.keys(e).length
}

async function simpan() {
  if (menyimpan.value) return
  pesan.value = ''
  if (!periksa()) return
  menyimpan.value = true
  try {
    const baru = await patchProduct(props.produk.id, {
      price_general: angka(isian.price_general) ? Number(angka(isian.price_general)) : null,
      unit_sale: isian.unit_sale.trim(),
      show_price: isian.show_price,
    })
    emit('tersimpan', baru)
    open.value = false
  } catch (error) {
    errors.value = errorField(error)
    if (!Object.keys(errors.value).length && error?.response?.status !== 401) pesan.value = pesanError(error)
  } finally {
    menyimpan.value = false
  }
}
</script>

<template>
  <Modal v-model:open="open" :title="`Ubah harga ${produk?.name ?? ''}`" :description="produk ? `SKU ${produk.sku}` : ''" size="sm" :dismissible="!menyimpan">
    <form id="form-ubah-harga" class="flex flex-col gap-4" novalidate @submit.prevent="simpan">
      <p v-if="pesan" role="alert" class="rounded-md bg-danger-soft px-3 py-2 text-sm font-medium text-danger-ink">{{ pesan }}</p>
      <Switch v-model="isian.show_price" label="Tampilkan harga di website" />
      <div class="grid gap-4 sm:grid-cols-[minmax(0,1fr)_8rem]">
        <Input
          v-model="isian.price_general"
          label="Harga umum"
          :required="isian.show_price"
          prefix="Rp"
          inputmode="numeric"
          autocomplete="off"
          :error="errors.price_general ?? ''"
          :hint="isian.show_price ? '' : 'Opsional, hanya untuk catatan admin.'"
        />
        <Input v-model="isian.unit_sale" label="Satuan jual" required autocomplete="off" :error="errors.unit_sale ?? ''" />
      </div>
      <p v-if="pratinjau" class="text-sm text-muted" aria-live="polite" data-pratinjau-harga>{{ pratinjau }}</p>
    </form>
    <template #footer>
      <Button variant="outline" :disabled="menyimpan" @click="open = false">Batal</Button>
      <Button type="submit" form="form-ubah-harga" :loading="menyimpan">Simpan harga</Button>
    </template>
  </Modal>
</template>

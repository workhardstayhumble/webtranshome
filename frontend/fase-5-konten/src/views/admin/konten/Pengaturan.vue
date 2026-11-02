<script setup>
import { computed, nextTick, onMounted, reactive, ref } from 'vue'
import { TriangleAlert } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Textarea from '@/components/ui/Textarea.vue'
import FormSection from '@/components/admin/FormSection.vue'
import { useKonfirmasiKeluar } from '@/composables/useKonfirmasiKeluar'
import { usePermission } from '@/composables/usePermission'
import { getAdminSettings, updateSettings } from '@/services/adminKontenService'
import { errorField, pesanError } from '@/services/errors'
import { useSettingsStore } from '@/stores/settings'
import { useUiStore } from '@/stores/ui'
import { formatRupiah } from '@/utils/format'

const angka = (teks) => String(teks ?? '').replace(/\D/g, '')

const { can } = usePermission()
const settings = useSettingsStore()
const ui = useUiStore()

const status = ref('loading') // loading | ready | error
const pesan = ref('')
const galatUmum = ref('')
const errors = ref({})
const menyimpan = ref(false)
const konfirmasiRasio = ref(false)
const asli = ref(null)
const awal = ref('')
const formEl = ref(null)
const form = reactive({ address: '', opening_hours: '', maps_url: '', point_ratio_rupiah: '' })

// Rasio poin hanya bisa diubah pemegang izin setting.point_ratio (Super Admin).
const bolehUbahRasio = computed(() => can('setting.point_ratio'))
const kotor = computed(() => status.value === 'ready' && JSON.stringify(form) !== awal.value)
const rasioBaru = computed(() => Number(angka(form.point_ratio_rupiah)))
const rasioBerubah = computed(() => bolehUbahRasio.value && rasioBaru.value !== Number(asli.value?.point_ratio_rupiah))
const keluar = useKonfirmasiKeluar(kotor)

function isi(s) {
  asli.value = s
  Object.assign(form, {
    address: s.address ?? '', opening_hours: s.opening_hours ?? '', maps_url: s.maps_url ?? '',
    point_ratio_rupiah: String(s.point_ratio_rupiah ?? ''),
  })
  awal.value = JSON.stringify(form)
}

async function muat() {
  status.value = 'loading'
  try {
    isi(await getAdminSettings())
    status.value = 'ready'
  } catch (error) {
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}
onMounted(muat)

function validasi() {
  const e = {}
  if (!form.address.trim()) e.address = 'Alamat wajib diisi'
  if (!form.opening_hours.trim()) e.opening_hours = 'Jam buka wajib diisi'
  if (form.maps_url.trim() && !form.maps_url.trim().startsWith('https://')) e.maps_url = 'Link Google Maps harus diawali https://'
  if (rasioBerubah.value && rasioBaru.value < 1000) e.point_ratio_rupiah = 'Rasio poin minimal Rp 1.000 per poin'
  errors.value = e
  return Object.keys(e).length === 0
}

async function fokusGalatPertama() {
  await nextTick()
  const el = formEl.value?.querySelector('[aria-invalid="true"]')
  el?.scrollIntoView({ block: 'center' })
  el?.focus({ preventScroll: true })
}

// Rasio berubah → konfirmasi dulu (berdampak ke semua transaksi poin berikutnya).
function mintaSimpan() {
  galatUmum.value = ''
  if (!validasi()) return fokusGalatPertama()
  if (rasioBerubah.value) konfirmasiRasio.value = true
  else simpan()
}

async function simpan() {
  menyimpan.value = true
  const payload = {
    address: form.address.trim(), opening_hours: form.opening_hours.trim(), maps_url: form.maps_url.trim() || null,
    ...(rasioBerubah.value && { point_ratio_rupiah: rasioBaru.value }),
  }
  try {
    isi(await updateSettings(payload))
    // Header, footer, Info toko, dan semua tombol WhatsApp membaca store ini.
    settings.load({ force: true })
    ui.tampilkanToast({ pesan: 'Pengaturan toko disimpan', jenis: 'success' })
  } catch (error) {
    const statusHttp = error?.response?.status
    if (statusHttp === 422) {
      errors.value = errorField(error)
      galatUmum.value = error.response.data?.message ?? 'Periksa kembali isian yang ditandai.'
      fokusGalatPertama()
    } else if (statusHttp !== 401) {
      galatUmum.value = pesanError(error)
      formEl.value?.scrollIntoView({ block: 'start' })
    }
  } finally {
    menyimpan.value = false
    konfirmasiRasio.value = false
  }
}
</script>

<template>
  <div v-if="status === 'loading'" class="flex max-w-3xl flex-col gap-6" aria-busy="true">
    <span class="sr-only">Memuat pengaturan…</span>
    <Skeleton v-for="n in 2" :key="n" class="h-48 rounded-xl" />
  </div>

  <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

  <form v-else ref="formEl" novalidate class="flex flex-col gap-6" @submit.prevent="mintaSimpan">
    <p v-if="galatUmum" role="alert" class="max-w-3xl rounded-md bg-danger-soft px-4 py-3 text-sm font-semibold text-danger-ink">{{ galatUmum }}</p>

    <div class="flex max-w-3xl flex-col gap-6">
      <FormSection judul="WhatsApp toko" deskripsi="Semua tombol WhatsApp di website mengarah ke nomor ini.">
        <dl class="flex flex-col gap-1">
          <dt class="text-sm font-semibold">Nomor WhatsApp</dt>
          <dd class="font-semibold tabular-nums" data-nomor-wa-server>{{ asli.wa_number || 'Belum dikonfigurasi di server' }}</dd>
        </dl>
        <p class="text-sm text-muted">Dikelola oleh Tim Cyber melalui konfigurasi server. Nomor ini tidak dapat diubah melalui dashboard, termasuk oleh Super Admin.</p>
      </FormSection>

      <FormSection judul="Lokasi dan jam buka" deskripsi="Tampil di header, footer, dan halaman Info toko.">
        <Textarea v-model="form.address" label="Alamat" required :rows="3" :error="errors.address ?? ''" />
        <Input v-model="form.opening_hours" label="Jam buka" required autocomplete="off" placeholder="Senin–Sabtu, 08.00–17.00" :error="errors.opening_hours ?? ''" />
        <Input v-model="form.maps_url" label="Link Google Maps" type="url" autocomplete="off" placeholder="https://maps.app.goo.gl/…" :error="errors.maps_url ?? ''" hint="Opsional. Tombol Buka di Google Maps memakai link ini." />
      </FormSection>

      <FormSection judul="Poin Trans Family" data-bagian-rasio>
        <template v-if="bolehUbahRasio">
          <Input
            v-model="form.point_ratio_rupiah"
            label="Belanja untuk 1 poin"
            prefix="Rp"
            inputmode="numeric"
            autocomplete="off"
            :error="errors.point_ratio_rupiah ?? ''"
            :hint="rasioBaru ? `${formatRupiah(rasioBaru)} belanja = 1 poin` : ''"
          />
          <p class="flex items-start gap-2 rounded-md bg-warning-soft px-3 py-2.5 text-sm text-warning-ink">
            <TriangleAlert class="mt-0.5 size-4 shrink-0" aria-hidden="true" />
            Hanya berlaku untuk transaksi baru. Poin yang sudah tercatat tidak berubah.
          </p>
        </template>
        <p v-else class="text-sm" data-rasio-baca>
          <strong>{{ formatRupiah(asli.point_ratio_rupiah) }}</strong> belanja = 1 poin.
          <span class="text-muted">Hanya Super Admin yang bisa mengubah rasio poin.</span>
        </p>
      </FormSection>
    </div>

    <div class="sticky bottom-0 z-20 -mx-4 flex flex-wrap items-center justify-end gap-2 border-t border-border bg-surface px-4 py-3 sm:-mx-7 sm:px-7">
      <span v-if="kotor" class="mr-auto hidden text-sm text-muted sm:inline">Ada perubahan yang belum disimpan</span>
      <Button v-if="kotor" variant="outline" :disabled="menyimpan" @click="isi(asli)">Batalkan perubahan</Button>
      <Button type="submit" :loading="menyimpan && !konfirmasiRasio">Simpan</Button>
    </div>

    <ConfirmModal
      :open="konfirmasiRasio"
      title="Ubah rasio poin?"
      :message="`Mulai sekarang setiap ${formatRupiah(rasioBaru)} belanja = 1 poin (sebelumnya ${formatRupiah(asli.point_ratio_rupiah)}). Hanya berlaku untuk transaksi baru; poin yang sudah tercatat tidak berubah.`"
      confirm-label="Ubah rasio dan simpan"
      :loading="menyimpan"
      @update:open="(buka) => !buka && (konfirmasiRasio = false)"
      @confirm="simpan"
    />
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

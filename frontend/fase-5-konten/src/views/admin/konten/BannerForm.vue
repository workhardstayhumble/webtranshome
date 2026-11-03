<script setup>
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Switch from '@/components/ui/Switch.vue'
import FormSection from '@/components/admin/FormSection.vue'
import GambarField from '@/components/admin/GambarField.vue'
import StatusBanner from '@/components/admin/konten/StatusBanner.vue'
import VideoField from '@/components/admin/konten/VideoField.vue'
import { useFormBanner } from '@/composables/useFormBanner'
import { useKonfirmasiKeluar } from '@/composables/useKonfirmasiKeluar'
import { tujuanKembali } from '@/composables/useTujuanKembali'
import { useUiStore } from '@/stores/ui'
import { POSISI_BANNER, perkiraanStatus } from '@/utils/banner'
import { formatTanggal } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const ui = useUiStore()
const f = useFormBanner()
const { form, errors } = f
const keluar = useKonfirmasiKeluar(f.kotor, { saatSesiHabis: f.simpanDraf })
const formEl = ref(null)

const OPSI_POSISI = Object.entries(POSISI_BANNER).map(([value, p]) => ({ value, label: p.label }))
const posisi = computed(() => POSISI_BANNER[form.position] ?? POSISI_BANNER.HOME_SLIDER)
const hero = computed(() => form.position === 'HOME_HERO')

// Perkiraan status dari isian, supaya admin tahu sebelum menyimpan apakah banner akan tampil.
const perkiraan = computed(() => perkiraanStatus(form))
const keteranganStatus = computed(() => {
  if (perkiraan.value === 'TERJADWAL') return `Mulai tampil ${formatTanggal(form.start_date)}.`
  if (perkiraan.value === 'BERAKHIR') return 'Periodenya sudah lewat, jadi banner tidak tampil. Ubah tanggal selesai untuk memperpanjang.'
  if (perkiraan.value === 'NONAKTIF') return 'Tidak tampil di Beranda karena "Aktif" dimatikan.'
  return form.end_date ? `Tampil di Beranda sampai ${formatTanggal(form.end_date)}.` : 'Tampil di Beranda tanpa batas waktu.'
})

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
  ui.tampilkanToast({ pesan: 'Banner disimpan', jenis: 'success' })
  f.tandaiTersimpan()
  router.push(tujuanKembali('banner', '/admin/banner'))
}

const batal = () => router.push(tujuanKembali('banner', '/admin/banner'))

onMounted(async () => {
  await f.muat()
  // "Perpanjang" (daftar banner / dashboard): /admin/banner/:id#periode langsung ke tanggal selesai.
  if (route.hash === '#periode' && f.status.value === 'ready') {
    await nextTick()
    document.getElementById('periode')?.scrollIntoView({ block: 'start' })
    document.querySelector('[data-tanggal-selesai]')?.focus({ preventScroll: true })
  }
})
</script>

<template>
  <div v-if="f.status.value === 'loading'" class="grid gap-6 xl:grid-cols-[minmax(0,1fr)_20rem]" aria-busy="true">
    <span class="sr-only">Memuat banner…</span>
    <div class="flex flex-col gap-6">
      <Skeleton class="h-80 rounded-xl" />
      <Skeleton class="h-36 rounded-xl" />
    </div>
    <Skeleton class="h-56 rounded-xl" />
  </div>

  <section v-else-if="f.status.value === 'tidak-ada'" class="mx-auto flex max-w-md flex-col items-center gap-4 py-16 text-center">
    <h2 class="text-2xl font-bold">Banner tidak ditemukan</h2>
    <p class="text-muted">Banner ini tidak ada atau alamatnya salah.</p>
    <Button to="/admin/banner">Kembali ke daftar banner</Button>
  </section>

  <ErrorState v-else-if="f.status.value === 'error'" :message="f.pesan.value" @retry="f.muat" />

  <form v-else ref="formEl" novalidate class="flex flex-col gap-6" @submit.prevent="simpan">
    <p v-if="f.dipulihkan.value" role="status" class="rounded-md bg-info-soft px-4 py-3 text-sm font-medium text-info-ink">
      Isian sebelum sesi habis sudah dipulihkan. Periksa lagi, lalu simpan.
    </p>
    <p v-if="f.galatUmum.value" role="alert" class="rounded-md bg-danger-soft px-4 py-3 text-sm font-semibold text-danger-ink">{{ f.galatUmum.value }}</p>

    <div class="grid items-start gap-6 xl:grid-cols-[minmax(0,1fr)_20rem]">
      <div class="flex min-w-0 flex-col gap-6">
        <FormSection judul="Isi banner">
          <Input
            v-model="form.title"
            label="Judul"
            required
            maxlength="100"
            autocomplete="off"
            :error="errors.title ?? ''"
            :hint="hero ? 'Judul besar hero di bagian paling atas Beranda.' : 'Tampil di atas gambar di Beranda.'"
          />
          <Input v-if="hero" v-model="form.subtitle" label="Subjudul" maxlength="160" autocomplete="off" :error="errors.subtitle ?? ''" hint="Opsional, satu kalimat di bawah judul." />
          <GambarField
            v-model="form.image_url"
            :label="hero ? 'Gambar hero' : 'Gambar banner'"
            required
            :rasio="posisi.rasio"
            lebar="w-64"
            :error="errors.image_url ?? ''"
            :hint="
              hero
                ? `JPG, PNG, atau WebP, maksimal 2 MB. Ukuran disarankan ${posisi.ukuran}. Tampil bila video belum dimuat, dijeda karena “kurangi animasi”, atau tidak ada video.`
                : `JPG, PNG, atau WebP, maksimal 2 MB. Ukuran disarankan ${posisi.ukuran}.`
            "
          />
          <VideoField v-if="hero" v-model="form.video_url" label="Video hero" :error="errors.video_url ?? ''" />
          <template v-if="hero">
            <div class="grid gap-4 sm:grid-cols-2">
              <Input
                v-model="form.cta_label"
                label="Label tombol"
                maxlength="30"
                autocomplete="off"
                placeholder="Jelajahi katalog"
                :error="errors.cta_label ?? ''"
                hint="Opsional. Kosong = “Jelajahi katalog”."
              />
              <Input
                v-model="form.link_url"
                label="Tautan tombol"
                autocomplete="off"
                placeholder="/katalog"
                :error="errors.link_url ?? ''"
                hint="Opsional. Kosong = /katalog. Diawali / atau https://."
              />
              <Input
                v-model="form.cta2_label"
                label="Label tombol kedua"
                maxlength="30"
                autocomplete="off"
                placeholder="Lihat promo"
                :error="errors.cta2_label ?? ''"
                hint="Opsional. Kosong = tidak tampil."
              />
              <Input
                v-model="form.cta2_url"
                label="Tautan tombol kedua"
                autocomplete="off"
                placeholder="/promo"
                :error="errors.cta2_url ?? ''"
                hint="Wajib bila label diisi. Diawali / atau https://."
              />
            </div>
            <p class="text-sm text-muted">Tombol “Lihat selengkapnya” selalu tampil di hero.</p>
          </template>
          <Input
            v-else
            v-model="form.link_url"
            label="Tautan"
            autocomplete="off"
            placeholder="/promo"
            :error="errors.link_url ?? ''"
            hint="Opsional. Halaman Transhome diawali /, situs lain diawali https://."
          />
        </FormSection>

        <FormSection id="periode" judul="Periode tayang" deskripsi="Mulai pukul 00.00 dan selesai pukul 23.59 WIB.">
          <div class="grid gap-4 sm:grid-cols-2">
            <Input v-model="form.start_date" type="date" label="Tanggal mulai" required :error="errors.start_date ?? ''" />
            <Input
              v-model="form.end_date"
              type="date"
              label="Tanggal selesai"
              :min="form.start_date || undefined"
              :error="errors.end_date ?? ''"
              hint="Kosongkan bila tayang terus."
              data-tanggal-selesai
            />
          </div>
        </FormSection>
      </div>

      <div class="flex flex-col gap-6 xl:sticky xl:top-24">
        <FormSection judul="Tampil di Beranda">
          <Select v-model="form.position" label="Posisi" :options="OPSI_POSISI" :error="errors.position ?? ''" />
          <Switch v-model="form.is_active" label="Aktif" />
          <div class="flex flex-col gap-1.5 rounded-md bg-background px-3 py-3 text-sm" aria-live="polite" data-perkiraan-status>
            <span class="flex items-center gap-2">
              <span class="font-semibold">Status setelah disimpan</span>
              <StatusBanner :status="perkiraan" />
            </span>
            <span class="text-muted">{{ keteranganStatus }}</span>
          </div>
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

import { computed, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { createBanner, getAdminBanner, updateBanner } from '@/services/adminKontenService'
import { errorField, pesanError } from '@/services/errors'
import { hariIniWib, tanggalSaja, tautanBannerValid } from '@/utils/banner'

const KUNCI_DRAF = 'transhome.drafBanner'

function formKosong() {
  return {
    title: '', image_url: '', link_url: '', position: 'HOME_SLIDER', start_date: hariIniWib(), end_date: '', is_active: true,
    // Khusus hero Beranda (HOME_HERO).
    subtitle: '', video_url: '', cta_label: '', cta2_label: '', cta2_url: '',
  }
}

/**
 * State dan aksi form banner (baru & ubah): muat, validasi, simpan, dan draf saat sesi habis.
 * Tanggal dikirim YYYY-MM-DD; backend menyimpan mulai 00.00 dan selesai 23.59 WIB.
 */
export function useFormBanner() {
  const route = useRoute()
  const id = route.params.id ? Number(route.params.id) : null

  const form = reactive(formKosong())
  const errors = ref({})
  const status = ref(id ? 'loading' : 'ready') // loading | ready | tidak-ada | error
  const pesan = ref('')
  const galatUmum = ref('')
  const menyimpan = ref(false)
  const dipulihkan = ref(false)
  const awal = ref('')

  const potret = () => JSON.stringify(form)
  const kotor = computed(() => status.value === 'ready' && potret() !== awal.value)

  function isi(b) {
    Object.assign(form, {
      title: b.title, image_url: b.image_url ?? '', link_url: b.link_url ?? '', position: b.position,
      start_date: tanggalSaja(b.start_at), end_date: tanggalSaja(b.end_at), is_active: b.is_active,
      subtitle: b.subtitle ?? '', video_url: b.video_url ?? '', cta_label: b.cta_label ?? '',
      cta2_label: b.cta2_label ?? '', cta2_url: b.cta2_url ?? '',
    })
  }

  function simpanDraf() {
    try {
      sessionStorage.setItem(KUNCI_DRAF, JSON.stringify({ path: route.path, form }))
    } catch {
      // Penyimpanan penuh atau diblokir: draf tidak tersimpan.
    }
  }

  function pulihkanDraf() {
    try {
      const draf = JSON.parse(sessionStorage.getItem(KUNCI_DRAF) ?? 'null')
      if (!draf || draf.path !== route.path) return
      sessionStorage.removeItem(KUNCI_DRAF)
      Object.assign(form, draf.form)
      dipulihkan.value = true
    } catch {
      // Draf rusak: abaikan.
    }
  }

  async function muat() {
    if (!id) return
    status.value = 'loading'
    try {
      isi(await getAdminBanner(id))
      awal.value = potret()
      pulihkanDraf()
      status.value = 'ready'
    } catch (error) {
      if (error?.response?.status === 404) status.value = 'tidak-ada'
      else {
        pesan.value = pesanError(error)
        status.value = 'error'
      }
    }
  }

  function validasi() {
    const e = {}
    if (!form.title.trim()) e.title = 'Judul wajib diisi'
    if (!form.image_url) e.image_url = form.position === 'HOME_HERO' ? 'Unggah gambar hero (juga dipakai sebagai cadangan video)' : 'Unggah gambar banner'
    if (form.position === 'HOME_HERO' && form.subtitle.trim().length > 160) e.subtitle = 'Subjudul maksimal 160 karakter'
    if (form.position === 'HOME_HERO' && form.cta_label.trim().length > 30) e.cta_label = 'Label tombol maksimal 30 karakter'
    if (form.position === 'HOME_HERO') {
      // Tombol kedua opsional, tapi label dan tautannya harus diisi berpasangan.
      const label2 = form.cta2_label.trim()
      const tautan2 = form.cta2_url.trim()
      if (label2.length > 30) e.cta2_label = 'Label tombol kedua maksimal 30 karakter'
      else if (!label2 && tautan2) e.cta2_label = 'Isi label tombol kedua, atau kosongkan tautannya'
      if (label2 && !tautan2) e.cta2_url = 'Isi tautan tombol kedua'
      else if (!tautanBannerValid(tautan2)) e.cta2_url = 'Tautan harus diawali / (halaman Transhome) atau https://'
    }
    if (!form.start_date) e.start_date = 'Isi tanggal mulai'
    else if (form.end_date && form.end_date < form.start_date) e.end_date = 'Tanggal selesai harus sama atau setelah tanggal mulai'
    if (!tautanBannerValid(form.link_url.trim())) e.link_url = 'Tautan harus diawali / (halaman Transhome) atau https://'
    errors.value = e
    return Object.keys(e).length === 0
  }

  const payload = () => ({
    title: form.title.trim(), image_url: form.image_url, link_url: form.link_url.trim() || null, position: form.position,
    start_date: form.start_date, end_date: form.end_date || null, is_active: form.is_active,
    ...(form.position === 'HOME_HERO'
      ? {
          subtitle: form.subtitle.trim() || null, video_url: form.video_url || null, cta_label: form.cta_label.trim() || null,
          cta2_label: form.cta2_label.trim() || null, cta2_url: form.cta2_url.trim() || null,
        }
      : {}),
  })

  /** Kirim ke API. Hasil: banner tersimpan, atau null bila gagal (galat 422 dipetakan ke field). */
  async function kirim() {
    if (menyimpan.value) return null
    galatUmum.value = ''
    if (!validasi()) return null
    menyimpan.value = true
    try {
      return id ? await updateBanner(id, payload()) : await createBanner(payload())
    } catch (error) {
      const statusHttp = error?.response?.status
      if (statusHttp === 422) {
        errors.value = errorField(error)
        galatUmum.value = error.response.data?.message ?? 'Periksa kembali isian yang ditandai.'
      } else if (statusHttp !== 401) {
        galatUmum.value = pesanError(error)
      }
      return null
    } finally {
      menyimpan.value = false
    }
  }

  if (!id) {
    awal.value = potret()
    pulihkanDraf()
  }

  return {
    id, form, errors, status, pesan, galatUmum, menyimpan, dipulihkan, kotor,
    muat, kirim, simpanDraf,
    tandaiTersimpan: () => (awal.value = potret()),
  }
}

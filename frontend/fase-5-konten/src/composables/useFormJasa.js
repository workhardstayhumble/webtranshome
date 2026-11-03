import { computed, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { createService, getAdminService, updateService } from '@/services/adminKontenService'
import { errorField, pesanError } from '@/services/errors'
import { buatSlug, POLA_SLUG } from '@/utils/slug'

const KUNCI_DRAF = 'transhome.drafJasa'

const formKosong = () => ({
  name: '', slug: '', summary: '', description: '', image_url: '', service_area: '', category_ids: [], is_active: true,
})

/**
 * State dan aksi form jasa (baru & ubah): muat, slug otomatis dari nama, validasi, simpan, dan draf saat sesi habis.
 * Harga tidak diisi per jasa: tampilan memakai label tetap "Harga pasti murah" (utils/jasa.js).
 */
export function useFormJasa() {
  const route = useRoute()
  const id = route.params.id ? Number(route.params.id) : null

  const form = reactive(formKosong())
  const errors = ref({})
  const status = ref(id ? 'loading' : 'ready') // loading | ready | tidak-ada | error
  const pesan = ref('')
  const galatUmum = ref('')
  const menyimpan = ref(false)
  const dipulihkan = ref(false)
  const slugManual = ref(Boolean(id))
  const awal = ref(JSON.stringify(form))

  const potret = () => JSON.stringify(form)
  const kotor = computed(() => status.value === 'ready' && potret() !== awal.value)

  function isi(j) {
    Object.assign(form, {
      name: j.name, slug: j.slug, summary: j.summary ?? '', description: j.description ?? '', image_url: j.image_url ?? '',
      service_area: j.service_area ?? '',
      category_ids: [...(j.category_ids ?? [])], is_active: j.is_active,
    })
  }

  function simpanDraf() {
    try {
      sessionStorage.setItem(KUNCI_DRAF, JSON.stringify({ path: route.path, form, slugManual: slugManual.value }))
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
      slugManual.value = draf.slugManual
      dipulihkan.value = true
    } catch {
      // Draf rusak: abaikan.
    }
  }

  async function muat() {
    if (!id) {
      pulihkanDraf()
      return
    }
    status.value = 'loading'
    try {
      isi(await getAdminService(id))
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

  function ubahNama(nilai) {
    form.name = nilai
    if (!slugManual.value) form.slug = buatSlug(nilai)
  }

  function ubahSlug(nilai) {
    slugManual.value = true
    form.slug = nilai
  }

  function ubahKategori(idKategori, pilih) {
    const ada = new Set(form.category_ids)
    if (pilih) ada.add(idKategori)
    else ada.delete(idKategori)
    form.category_ids = [...ada]
  }

  function validasi() {
    const e = {}
    if (!form.name.trim()) e.name = 'Nama jasa wajib diisi'
    if (!form.slug) e.slug = 'Slug wajib diisi'
    else if (!POLA_SLUG.test(form.slug)) e.slug = 'Slug hanya boleh huruf kecil, angka, dan tanda hubung'
    if (!form.summary.trim()) e.summary = 'Ringkasan wajib diisi'
    else if (form.summary.trim().length > 160) e.summary = 'Ringkasan maksimal 160 karakter'
    errors.value = e
    return Object.keys(e).length === 0
  }

  const payload = () => ({
    name: form.name.trim(), slug: form.slug, summary: form.summary.trim(), description: form.description.trim(),
    image_url: form.image_url || null, service_area: form.service_area.trim() || null,
    category_ids: form.category_ids, is_active: form.is_active,
  })

  /** Kirim ke API. Hasil: jasa tersimpan, atau null bila gagal (422/409 dipetakan ke field). */
  async function kirim() {
    if (menyimpan.value) return null
    galatUmum.value = ''
    if (!validasi()) return null
    menyimpan.value = true
    try {
      return id ? await updateService(id, payload()) : await createService(payload())
    } catch (error) {
      const statusHttp = error?.response?.status
      if (statusHttp === 422 || statusHttp === 409) {
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

  const tandaiTersimpan = () => (awal.value = potret())

  return { id, form, errors, status, pesan, galatUmum, menyimpan, dipulihkan, kotor, muat, ubahNama, ubahSlug, ubahKategori, kirim, simpanDraf, tandaiTersimpan }
}

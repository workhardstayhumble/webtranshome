import http from '@/services/http'

// Konten di panel admin (Fase 5): banner, artikel/promo/event, jasa, dan pengaturan toko.
// Pesan gagal ditampilkan di halaman (form/toast), jadi permintaan ubah memakai tanpaToast.
const TANPA_TOAST = { tanpaToast: true }

/** Semua banner (termasuk nonaktif/berakhir) dengan `status`: TAYANG | TERJADWAL | BERAKHIR | NONAKTIF. */
export async function listAdminBanners() {
  const { data } = await http.get('/admin/banners')
  return data.data
}

export async function getAdminBanner(id) {
  const { data } = await http.get(`/admin/banners/${id}`)
  return data.data
}

/** { title, image_url, link_url, position, start_date, end_date, is_active }; tanggal YYYY-MM-DD. */
export async function createBanner(payload) {
  const { data } = await http.post('/admin/banners', payload, TANPA_TOAST)
  return data.data
}

export async function updateBanner(id, payload) {
  const { data } = await http.put(`/admin/banners/${id}`, payload, TANPA_TOAST)
  return data.data
}

/** { is_active } atau { sort_order }. */
export async function patchBanner(id, perubahan) {
  const { data } = await http.patch(`/admin/banners/${id}`, perubahan, TANPA_TOAST)
  return data.data
}

/** Daftar artikel: q, type, status (DRAFT | PUBLISHED), page. Hasil { data, meta }. */
export async function listAdminArticles(params = {}, { signal } = {}) {
  const { data } = await http.get('/admin/articles', { params, signal })
  return data
}

/** Artikel lengkap untuk editor: + content dan preview_token. */
export async function getAdminArticle(id) {
  const { data } = await http.get(`/admin/articles/${id}`)
  return data.data
}

/** { type, title, slug, thumbnail_url, content } → artikel baru berstatus DRAFT. 409 slug dipakai. */
export async function createArticle(payload) {
  const { data } = await http.post('/admin/articles', payload, TANPA_TOAST)
  return data.data
}

export async function updateArticle(id, payload) {
  const { data } = await http.put(`/admin/articles/${id}`, payload, TANPA_TOAST)
  return data.data
}

/** Terbitkan ('PUBLISHED') atau batalkan terbit ('DRAFT'); API: content.article.manage, mock: article.publish. */
export async function setArticleStatus(id, status) {
  const { data } = await http.patch(`/admin/articles/${id}`, { status }, TANPA_TOAST)
  return data.data
}

/** { wa_number, address, opening_hours, maps_url, point_ratio_rupiah }. */
export async function getAdminSettings() {
  const { data } = await http.get('/admin/settings')
  return data.data
}

/** Simpan alamat, jam buka, peta, dan rasio (izin setting.point_ratio). Nomor WhatsApp hanya-baca dari server. */
export async function updateSettings(payload) {
  const { data } = await http.put('/admin/settings', payload, TANPA_TOAST)
  return data.data
}

/** Semua jasa (termasuk nonaktif) + categories [{ id, name, slug, parent_id }]. Izin service.view. */
export async function listAdminServices() {
  const { data } = await http.get('/admin/services')
  return data.data
}

export async function getAdminService(id) {
  const { data } = await http.get(`/admin/services/${id}`)
  return data.data
}

/**
 * { name, slug, summary, description, image_url, service_area, category_ids, is_active }. 409 bila slug dipakai jasa lain; 422 errors per kolom.
 */
export async function createService(payload) {
  const { data } = await http.post('/admin/services', payload, TANPA_TOAST)
  return data.data
}

export async function updateService(id, payload) {
  const { data } = await http.put(`/admin/services/${id}`, payload, TANPA_TOAST)
  return data.data
}

/** { is_active } atau { sort_order }. */
export async function patchService(id, perubahan) {
  const { data } = await http.patch(`/admin/services/${id}`, perubahan, TANPA_TOAST)
  return data.data
}

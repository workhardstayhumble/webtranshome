import http from '@/services/http'

// Konten publik (Fase 5): banner Beranda (termasuk hero), artikel/promo/event, dan jasa.

/**
 * Banner yang sedang tayang di satu posisi: HOME_HERO (hero paling atas; + subtitle, video_url, cta_label),
 * HOME_SLIDER (slider utama), atau HOME_SIDE (samping slider).
 */
export async function listBanners(position) {
  const { data } = await http.get('/banners', { params: { position } })
  return data.data
}

/**
 * Artikel/promo/event yang sudah terbit, terbaru dulu. Parameter: type (mis. 'ARTIKEL' atau 'PROMO,EVENT'),
 * page, per_page. Hasil { data: [{ id, type, title, slug, thumbnail_url, excerpt, published_at }], meta }.
 */
export async function listArticles(params = {}, { signal } = {}) {
  const { data } = await http.get('/articles', { params, signal })
  return data
}

/** Detail artikel. Draf hanya bisa dibuka dengan token pratinjau (?preview=); selain itu 404. */
export async function getArticle(slug, preview = null, { signal } = {}) {
  const { data } = await http.get(`/articles/${encodeURIComponent(slug)}`, { params: preview ? { preview } : {}, signal })
  return data.data
}

/**
 * Jasa aktif, urut sort_order (KEPUTUSAN #20). Parameter: category (slug kategori produk; jasa yang terhubung ke
 * kategori itu atau induknya), limit. Item { id, name, slug, summary, image_url, service_area }.
 */
export async function listServices(params = {}) {
  const { data } = await http.get('/services', { params })
  return data.data
}

/** Detail jasa: item di atas + description dan categories [{ name, slug }]. Jasa nonaktif → 404. */
export async function getService(slug) {
  const { data } = await http.get(`/services/${encodeURIComponent(slug)}`)
  return data.data
}

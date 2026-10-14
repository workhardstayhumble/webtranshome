import http from '@/services/http'
import { bersihkanCachePublik } from './cache'

const MOCK = import.meta.env.VITE_USE_MOCK === 'true'
const BASE = MOCK ? '/admin/products' : '/admin/catalog/products'

function ubahDariAPI(p) {
  if (!p) return p
  const result = { ...p }
  const map = {
    categoryId:'category_id', brandId:'brand_id', priceGeneral:'price_general',
    unitSale:'unit_sale', minOrder:'min_order', stockStatus:'stock_status',
    stockQtyLabel:'stock_qty_label', isActive:'is_active', isFeatured:'is_featured',
    datasheetPdfUrl:'datasheet_pdf_url', metaTitle:'meta_title', metaDescription:'meta_description',
  }
  for (const [key, dest] of Object.entries(map)) if (key in result) result[dest] = result[key]
  if (Array.isArray(result.images)) result.images = result.images.map(i => ({ ...i, image_url: i.imageUrl ?? i.image_url, is_primary: i.isPrimary ?? i.is_primary }))
  if (result.primaryImageUrl) result.primary_image_url = result.primaryImageUrl
  if (typeof result.specifications === 'string') { try { result.specifications = JSON.parse(result.specifications) } catch { result.specifications = {} } }
  return result
}
function ubahKeAPI(p) {
  if (MOCK) return p
  const result = { ...p }
  const map = {
    category_id:'categoryId', brand_id:'brandId', price_general:'priceGeneral',
    unit_sale:'unitSale', min_order:'minOrder', stock_status:'stockStatus',
    stock_qty_label:'stockQtyLabel', is_active:'isActive', is_featured:'isFeatured',
    datasheet_pdf_url:'datasheetPdfUrl', meta_title:'metaTitle', meta_description:'metaDescription',
  }
  for (const [key, dest] of Object.entries(map)) if (key in result) {result[dest] = result[key]; delete result[key]}
  // specifications harus sesuai JSON string yang diharapkan kontrak Go, bukan objek mentah.
  if (result.specifications && typeof result.specifications !== 'string') result.specifications = JSON.stringify(result.specifications)
  return result
}

export async function listAdminProducts(params = {}, { signal } = {}) {
  const { data } = await http.get(BASE, { params, signal })
  if (MOCK) return data
  // Backend saat ini belum menyediakan filter/pagination server-side.
  return { data: (data.products ?? []).map(ubahDariAPI), meta: { page: 1, per_page: (data.products ?? []).length, total: (data.products ?? []).length, total_pages: 1 } }
}
export async function getAdminProduct(id) {
  const { data } = await http.get(`${BASE}/${encodeURIComponent(id)}`)
  return MOCK ? data.data : ubahDariAPI(data)
}
export async function checkSku(sku, exclude = null) {
  if (MOCK) { const { data } = await http.get('/admin/products/check-sku',{params:{sku,exclude:exclude ?? undefined}}); return data.data.available }
  // Backend belum menyediakan check-sku khusus: baca daftar yang ada sebagai bantuan UX.
  // Constraint unik SKU di database tetap menjadi otoritas final saat penyimpanan.
  const { data } = await http.get(BASE)
  return !(data.products ?? []).some((p) => p.sku?.toLowerCase() === sku.toLowerCase() && String(p.id) !== String(exclude ?? ''))
}
export async function createProduct(payload) {
  const { data } = await http.post(BASE, ubahKeAPI(payload))
  bersihkanCachePublik()
  return MOCK ? data.data : ubahDariAPI(data)
}
export async function updateProduct(id, payload) {
  const { data } = await http[MOCK ? 'put' : 'patch'](`${BASE}/${encodeURIComponent(id)}`, ubahKeAPI(payload))
  bersihkanCachePublik()
  return MOCK ? data.data : ubahDariAPI(data)
}
export async function patchProduct(id, perubahan) {
  const { data } = await http.patch(`${BASE}/${encodeURIComponent(id)}`, ubahKeAPI(perubahan),{tanpaToast:true})
  bersihkanCachePublik()
  return MOCK ? data.data : ubahDariAPI(data)
}

/** Produk Favorit: Admin Konten tidak perlu permission edit keseluruhan produk. */
export async function setFeatured(id, aktif) {
  if (MOCK) return patchProduct(id, { is_featured: aktif })
  const { data } = await http.patch(`${BASE}/${encodeURIComponent(id)}/featured`, { isFeatured: Boolean(aktif) })
  bersihkanCachePublik()
  return ubahDariAPI(data)
}
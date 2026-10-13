import http from '@/services/http'
import { bersihkanCachePublik } from './cache'

const MOCK = import.meta.env.VITE_USE_MOCK === 'true'
const JENIS = { categories: 'categories', brands: 'brands', rooms: 'rooms' }
const AKSES = {
  categories: 'catalog.category', brands: 'catalog.brand', rooms: 'catalog.room',
}

function endpoint(jenis, id = null) {
  if (!Object.hasOwn(JENIS, jenis)) throw new Error('Jenis master katalog tidak dikenal.')
  const path = MOCK ? `/admin/${jenis}` : `/admin/catalog/${jenis}`
  return id == null ? path : `${path}/${encodeURIComponent(id)}`
}

function dariAPI(jenis, data) {
  // Backend Go mengirim camelCase; komponen Vue existing memakai snake_case.
  const item = { ...data }
  if ('parentId' in item) item.parent_id = item.parentId || null
  if ('imageUrl' in item) item.image_url = item.imageUrl
  if ('logoUrl' in item) item.logo_url = item.logoUrl
  if ('imageCover' in item) item.image_cover = item.imageCover
  if ('sortOrder' in item) item.sort_order = item.sortOrder
  if ('isActive' in item) item.is_active = item.isActive
  if (item.product_count == null) item.product_count = 0 // belum disediakan backend
  return item
}

function keAPI(payload) {
  if (MOCK) return payload
  const p = { ...payload }
  const keys = { parent_id: 'parentId', image_url: 'imageUrl', logo_url: 'logoUrl', image_cover: 'imageCover', sort_order: 'sortOrder', is_active: 'isActive' }
  for (const [a, b] of Object.entries(keys)) if (a in p) { p[b] = p[a]; delete p[a] }
  // Backend Create Room memerlukan imageCover dan kategori/ruangan sortOrder integer.
  return p
}

function ambil(jenis, response, daftar = false) {
  if (MOCK) return response.data
  if (!daftar) return dariAPI(jenis, response)
  const key = JENIS[jenis]
  return (response[key] ?? []).map((v) => dariAPI(jenis, v))
}

export async function listMaster(jenis) {
  const { data } = await http.get(endpoint(jenis))
  return ambil(jenis, data, true)
}

export async function createMaster(jenis, payload) {
  const { data } = await http.post(endpoint(jenis), keAPI(payload))
  bersihkanCachePublik()
  return ambil(jenis, data)
}

export async function updateMaster(jenis, id, payload) {
  const { data } = await http[MOCK ? 'put' : 'patch'](endpoint(jenis, id), keAPI(payload))
  bersihkanCachePublik()
  return ambil(jenis, data)
}

// Backend ZIP tidak menyediakan DELETE /admin/catalog/brands/:id.
// Fail closed: jangan menganggap berhasil atau menghapus data hanya di UI.
export async function deleteBrand(id) {
  if (!MOCK) throw new Error('Penghapusan brand belum tersedia pada backend Go.')
  await http.delete(endpoint('brands', id), { tanpaToast: true })
  bersihkanCachePublik()
}

export async function patchMaster(jenis, id, perubahan) {
  const { data } = await http.patch(endpoint(jenis, id), keAPI(perubahan), { tanpaToast: true })
  bersihkanCachePublik()
  return ambil(jenis, data)
}
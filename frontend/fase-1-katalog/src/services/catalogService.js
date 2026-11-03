import http from '@/services/http'
import { daftarkanPembersihCache } from '@/services/cache'

let cacheRuangan = null
let ruanganSiap = null
daftarkanPembersihCache(() => {
  cacheRuangan = null
  ruanganSiap = null
})

/** Semua brand aktif: [{ id, name, slug, logo_url, product_count, categories: [{ slug, name }] }] (Beranda, /brand). */
export async function listBrands() {
  const { data } = await http.get('/brands')
  return data.data
}

/** Detail brand: { id, name, slug, logo_url, product_count }. 404 bila tidak ada. */
export async function getBrand(slug) {
  const { data } = await http.get(`/brands/${encodeURIComponent(slug)}`)
  return data.data
}

/**
 * Ruangan aktif berurutan: [{ id, name, slug, image_cover, product_count }].
 * Disimpan di memori selama tab terbuka (jarang berubah); kalau gagal, bisa dicoba lagi.
 */
export function listRooms() {
  if (!cacheRuangan) {
    cacheRuangan = http
      .get('/rooms')
      .then(({ data }) => (ruanganSiap = data.data))
      .catch((error) => {
        cacheRuangan = null
        throw error
      })
  }
  return cacheRuangan
}

/** Daftar ruangan yang sudah pernah dimuat (atau null), untuk render langsung tanpa menunggu. */
export function ruanganTersimpan() {
  return ruanganSiap
}

/** Detail ruangan: { id, name, slug, image_cover, product_count }. 404 bila tidak ada. */
export async function getRoom(slug) {
  const { data } = await http.get(`/rooms/${encodeURIComponent(slug)}`)
  return data.data
}

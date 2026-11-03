// Rute API tiruan untuk mode mock. Bentuk respons mengikuti usulan kontrak di
// docs/FRONTEND_BUILD_PLAN.md (lampiran): sukses { data, meta? }, gagal { message, errors? }.
// TODO: ganti ke API setiap rute di sini begitu endpoint backend tersedia.
import { brands, categories, productImages, productRooms, products, rooms } from './data/catalog'
import { publicSettings } from './data/akun'
import { pulihkanDb } from '@/services/mock/db'
import { ruteAdmin } from '@/services/mock/routesAdmin'
import { ruteAkun } from '@/services/mock/routesAkun'
import { ruteKonten } from '@/services/mock/routesKonten'
import { ruteMember } from '@/services/mock/routesMember'
import { ruteQontak } from '@/services/mock/routesQontak'
import { ruteSistem } from '@/services/mock/routesSistem'

// Perubahan dari panel admin yang tersimpan (lihat db.js) dipakai juga oleh halaman publik.
pulihkanDb()

const STATUS_STOK = ['TERSEDIA', 'SISA_STOK', 'PRE_ORDER', 'HABIS']
// Atribut spesifikasi yang bisa jadi filter. TODO: daftar final dari backend (CLAUDE.md §12).
const ATRIBUT_FILTER = ['ukuran', 'finishing', 'warna', 'bahan', 'daya']

const huruf = (teks) => String(teks ?? '').toLowerCase()
const daftar = (nilai) => String(nilai ?? '').split(',').map((x) => x.trim()).filter(Boolean)
const tidakDitemukan = (pesan) => ({ status: 404, data: { message: pesan } })

// ---------- bentuk data ----------

const merekDari = (id) => brands.find((x) => x.id === id) ?? null
const ringkasMerek = (id) => {
  const b = merekDari(id)
  return b ? { id: b.id, name: b.name, slug: b.slug } : null
}

function ringkasKategori(id, { denganInduk = false } = {}) {
  const c = categories.find((x) => x.id === id)
  if (!c) return null
  const hasil = { id: c.id, name: c.name, slug: c.slug, parent_id: c.parent_id }
  if (denganInduk && c.parent_id) hasil.parent = ringkasKategori(c.parent_id)
  return hasil
}

const fotoUtama = (productId) => productImages.find((i) => i.product_id === productId && i.is_primary)?.image_url || null
const ruanganProduk = (productId) => productRooms.filter((pr) => pr.product_id === productId).map((pr) => rooms.find((r) => r.id === pr.room_id)).filter(Boolean)

// Harga hanya dikirim bila boleh tampil: produk dengan show_price false tidak membocorkan harganya lewat API.
const hargaTerlihat = (p) => p.show_price !== false && p.price_general != null

function itemProduk(p) {
  return {
    id: p.id, sku: p.sku, name: p.name, slug: p.slug,
    price_general: hargaTerlihat(p) ? p.price_general : null, show_price: hargaTerlihat(p), unit_sale: p.unit_sale, min_order: p.min_order,
    stock_status: p.stock_status, stock_qty_label: p.stock_qty_label, is_featured: p.is_featured,
    brand: ringkasMerek(p.brand_id), category: ringkasKategori(p.category_id),
    primary_image_url: fotoUtama(p.id),
  }
}

function idKategoriTurunan(slug) {
  const induk = categories.find((c) => c.slug === slug)
  if (!induk) return []
  return [induk.id, ...categories.filter((c) => c.parent_id === induk.id).map((c) => c.id)]
}

// ---------- katalog ----------

/** Filter dari query: brand, stock, room (dipisah koma), attr[kunci] (dipisah koma). */
function bacaFilter(query) {
  const attr = {}
  for (const [kunci, nilai] of Object.entries(query)) {
    const cocok = kunci.match(/^attr\[(.+)\]$/)
    if (cocok) attr[cocok[1]] = daftar(nilai)
  }
  return { brand: daftar(query.brand), stock: daftar(query.stock), room: daftar(query.room), attr }
}

/**
 * Produk lolos semua filter, kecuali grup `kecuali`. Dipakai untuk hasil (kecuali = null) dan untuk
 * menghitung angka per opsi (faceting disjungtif): angka sebuah opsi mengabaikan filter dari grupnya sendiri.
 */
function lolosFilter(p, f, kecuali = null) {
  if (kecuali !== 'brand' && f.brand.length && !f.brand.includes(merekDari(p.brand_id)?.slug)) return false
  if (kecuali !== 'stock' && f.stock.length && !f.stock.includes(p.stock_status)) return false
  if (kecuali !== 'room' && f.room.length && !ruanganProduk(p.id).some((r) => f.room.includes(r.slug))) return false
  for (const [kunci, nilai] of Object.entries(f.attr)) {
    if (kecuali !== `attr:${kunci}` && nilai.length && !nilai.includes(p.specifications?.[kunci])) return false
  }
  return true
}

// `tanpaHarga`: dasar sebelum filter min/maks, agar rentang harga tidak menyempit oleh filter harga itu sendiri.
function hitungFacet(dasar, f, tanpaHarga) {
  const jumlah = (kecuali, cocok) => dasar.filter((p) => lolosFilter(p, f, kecuali) && cocok(p)).length

  const merek = brands
    .filter((b) => b.is_active)
    .map((b) => ({ slug: b.slug, name: b.name, count: jumlah('brand', (p) => p.brand_id === b.id) }))
    .filter((x) => x.count > 0 || f.brand.includes(x.slug))

  const stok = STATUS_STOK.map((value) => ({ value, count: jumlah('stock', (p) => p.stock_status === value) }))

  const ruangan = rooms
    .filter((r) => r.is_active)
    .map((r) => ({ slug: r.slug, name: r.name, count: jumlah('room', (p) => ruanganProduk(p.id).some((x) => x.id === r.id)) }))
    .filter((x) => x.count > 0 || f.room.includes(x.slug))

  const atribut = ATRIBUT_FILTER.map((kunci) => {
    const terpilih = f.attr[kunci] ?? []
    const nilaiUnik = [...new Set(dasar.map((p) => p.specifications?.[kunci]).filter(Boolean))].sort((a, b) => a.localeCompare(b, 'id', { numeric: true }))
    const values = nilaiUnik
      .map((value) => ({ value, count: jumlah(`attr:${kunci}`, (p) => p.specifications?.[kunci] === value) }))
      .filter((x) => x.count > 0 || terpilih.includes(x.value))
    return { key: kunci, values }
  }).filter((grup) => grup.values.length >= 2 || (f.attr[grup.key] ?? []).length)

  const harga = tanpaHarga.filter((p) => lolosFilter(p, f) && hargaTerlihat(p)).map((p) => Number(p.price_general))
  return {
    brands: merek,
    stock: stok,
    rooms: ruangan,
    attributes: atribut,
    price: harga.length ? { min: Math.min(...harga), max: Math.max(...harga) } : null,
  }
}

function daftarProduk({ query }) {
  // Dasar: aktif + pencarian + kategori + harga + exclude. Filter faset diterapkan di atasnya.
  let dasar = products.filter((p) => p.is_active)
  if (query.q) {
    const q = huruf(query.q)
    dasar = dasar.filter((p) => huruf(p.name).includes(q) || huruf(p.sku).includes(q) || huruf(merekDari(p.brand_id)?.name).includes(q))
  }
  if (query.category) {
    const ids = idKategoriTurunan(query.category)
    dasar = dasar.filter((p) => ids.includes(p.category_id))
  }
  if (query.exclude) dasar = dasar.filter((p) => String(p.id) !== String(query.exclude))
  if (query.featured === '1' || query.featured === 1) dasar = dasar.filter((p) => p.is_featured)
  const tanpaHarga = dasar
  // Filter harga: produk yang harganya disembunyikan tidak ikut (harganya tidak diketahui pelanggan).
  if (query.min) dasar = dasar.filter((p) => hargaTerlihat(p) && Number(p.price_general) >= Number(query.min))
  if (query.max) dasar = dasar.filter((p) => hargaTerlihat(p) && Number(p.price_general) <= Number(query.max))

  const f = bacaFilter(query)
  let hasil = dasar.filter((p) => lolosFilter(p, f))

  // Urut harga: produk tanpa harga yang terlihat selalu di akhir.
  const tanpaHargaAkhir = (a, b) => Number(!hargaTerlihat(a)) - Number(!hargaTerlihat(b))
  const urut = {
    harga_terendah: (a, b) => tanpaHargaAkhir(a, b) || Number(a.price_general) - Number(b.price_general),
    harga_tertinggi: (a, b) => tanpaHargaAkhir(a, b) || Number(b.price_general) - Number(a.price_general),
    nama: (a, b) => a.name.localeCompare(b.name, 'id'),
    terbaru: (a, b) => b.updated_at.localeCompare(a.updated_at) || b.id - a.id,
  }
  hasil = [...hasil].sort(urut[query.sort] ?? urut.terbaru)

  const perPage = Math.min(Number(query.per_page) || 24, 60)
  const total = hasil.length
  const totalPages = Math.max(1, Math.ceil(total / perPage))
  const page = Math.min(Math.max(Number(query.page) || 1, 1), totalPages)
  const data = hasil.slice((page - 1) * perPage, page * perPage).map(itemProduk)
  return { data: { data, meta: { page, per_page: perPage, total, total_pages: totalPages, facets: hitungFacet(dasar, f, tanpaHarga) } } }
}

function detailProduk({ params }) {
  const p = products.find((x) => x.slug === params.slug && x.is_active)
  if (!p) return tidakDitemukan('Produk tidak ditemukan.')
  return {
    data: {
      data: {
        ...itemProduk(p),
        category: ringkasKategori(p.category_id, { denganInduk: true }),
        description: p.description, specifications: p.specifications, datasheet_pdf_url: p.datasheet_pdf_url,
        meta_title: p.meta_title, meta_description: p.meta_description,
        images: productImages.filter((i) => i.product_id === p.id).sort((a, b) => a.sort_order - b.sort_order),
        rooms: ruanganProduk(p.id).map(({ id, name, slug }) => ({ id, name, slug })),
      },
    },
  }
}

function saranProduk({ query }) {
  const q = huruf(query.q).trim()
  if (q.length < 2) return { data: { data: { products: [], categories: [], brands: [] } } }
  const produk = products
    .filter((p) => p.is_active && (huruf(p.name).includes(q) || huruf(p.sku).includes(q)))
    .slice(0, 5)
    .map(({ id, name, slug, sku }) => ({ id, name, slug, sku }))
  const kategori = categories.filter((c) => c.is_active && huruf(c.name).includes(q)).slice(0, 3).map(({ id, name, slug }) => ({ id, name, slug }))
  const merek = brands.filter((b) => b.is_active && huruf(b.name).includes(q)).slice(0, 3).map(({ id, name, slug }) => ({ id, name, slug }))
  return { data: { data: { products: produk, categories: kategori, brands: merek } } }
}

function pohonKategori() {
  const aktif = categories.filter((c) => c.is_active).sort((a, b) => a.sort_order - b.sort_order)
  const pohon = aktif
    .filter((c) => c.parent_id === null)
    .map((induk) => ({
      id: induk.id, name: induk.name, slug: induk.slug, image_url: induk.image_url,
      children: aktif.filter((c) => c.parent_id === induk.id).map(({ id, name, slug, image_url }) => ({ id, name, slug, image_url })),
    }))
  return { data: { data: pohon } }
}

const jumlahProdukRuangan = (roomId) =>
  products.filter((p) => p.is_active && productRooms.some((pr) => pr.product_id === p.id && pr.room_id === roomId)).length

function daftarRuangan() {
  const data = rooms
    .filter((r) => r.is_active)
    .sort((a, b) => a.sort_order - b.sort_order)
    .map(({ id, name, slug, image_cover }) => ({ id, name, slug, image_cover, product_count: jumlahProdukRuangan(id) }))
  return { data: { data } }
}

function detailRuangan({ params }) {
  const r = rooms.find((x) => x.slug === params.slug && x.is_active)
  if (!r) return tidakDitemukan('Ruangan tidak ditemukan.')
  return { data: { data: { id: r.id, name: r.name, slug: r.slug, image_cover: r.image_cover, product_count: jumlahProdukRuangan(r.id) } } }
}

function daftarMerek() {
  const produkAktif = (b) => products.filter((p) => p.is_active && p.brand_id === b.id)
  // Kategori brand = kategori utama dari produk aktifnya (untuk filter di /brand); tanpa kolom baru di tabel brands.
  const kategoriUtama = (b) => {
    const id = new Set(produkAktif(b).map((p) => {
      const c = categories.find((x) => x.id === p.category_id)
      return c?.parent_id ?? c?.id
    }))
    return categories
      .filter((c) => id.has(c.id) && c.is_active)
      .sort((x, y) => x.sort_order - y.sort_order)
      .map(({ slug, name }) => ({ slug, name }))
  }
  const data = brands
    .filter((b) => b.is_active)
    .map((b) => ({ id: b.id, name: b.name, slug: b.slug, logo_url: b.logo_url, product_count: produkAktif(b).length, categories: kategoriUtama(b) }))
  return { data: { data } }
}

function detailMerek({ params }) {
  const b = brands.find((x) => x.slug === params.slug && x.is_active)
  if (!b) return tidakDitemukan('Brand tidak ditemukan.')
  const jumlah = products.filter((p) => p.is_active && p.brand_id === b.id).length
  return { data: { data: { id: b.id, name: b.name, slug: b.slug, logo_url: b.logo_url, product_count: jumlah } } }
}

// ---------- konten & pengaturan ----------

function pengaturanPublik() {
  return { data: { data: publicSettings } }
}

// Banner, artikel, bagian rumah di Beranda, dan pengaturan admin: services/mock/routesKonten.js (Fase 5).

// ---------- tabel rute (yang lebih spesifik ditulis lebih dulu) ----------

const rute = [
  ['get', '/settings/public', pengaturanPublik],
  ['get', '/categories', pohonKategori],
  ['get', '/products/suggest', saranProduk],
  ['get', '/products', daftarProduk],
  ['get', '/products/:slug', detailProduk],
  ['get', '/brands', daftarMerek],
  ['get', '/brands/:slug', detailMerek],
  ['get', '/rooms', daftarRuangan],
  ['get', '/rooms/:slug', detailRuangan],
  ...ruteAkun,
  ...ruteAdmin,
  ...ruteMember,
  ...ruteKonten,
  ...ruteSistem,
  ...ruteQontak,
]

function cocokkan(pola, path) {
  const a = pola.split('/').filter(Boolean)
  const b = path.split('/').filter(Boolean)
  if (a.length !== b.length) return null
  const params = {}
  for (let i = 0; i < a.length; i++) {
    if (a[i].startsWith(':')) params[a[i].slice(1)] = decodeURIComponent(b[i])
    else if (a[i] !== b[i]) return null
  }
  return params
}

export function cariRute(method, path) {
  for (const [m, pola, handler] of rute) {
    if (m !== method) continue
    const params = cocokkan(pola, path)
    if (params) return { handler, params }
  }
  return null
}

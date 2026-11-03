// Rute API tiruan panel admin (Fase 2). Bentuk mengikuti lampiran docs/FRONTEND_BUILD_PLAN.md.
// TODO: ganti ke API begitu endpoint /admin/* tersedia di backend.
import { brands, categories, productImages, productRooms, products, rooms } from '@/services/mock/data/catalog'
import { activityLogs, menus, qontakContacts, syncLogs } from './data/admin'
import { pointTransactions, users } from '@/services/mock/data/akun'
import { banners } from '@/services/mock/data/konten'
import { idBaru, simpanDb } from './db'
import { cekAdmin, izinPengguna, namaPengguna, sekarang } from './sesi'
import { antrekan, prosesKontak } from '@/services/mock/antreanQontak'

const STATUS_STOK = ['TERSEDIA', 'SISA_STOK', 'PRE_ORDER', 'HABIS']
const POLA_SLUG = /^[a-z0-9]+(?:-[a-z0-9]+)*$/
const huruf = (teks) => String(teks ?? '').toLowerCase()
const kosong = (nilai) => nilai == null || String(nilai).trim() === ''
const tidakDitemukan = (pesan) => ({ status: 404, data: { message: pesan } })
const galatIsian = (errors) => ({ status: 422, data: { message: 'Periksa kembali isian yang ditandai.', errors } })
const sama = (a, b) => String(a) === String(b)

/** Catat ke activity_logs. `sebelum` / `sesudah` mengisi before_data / after_data (ditampilkan di log aktivitas, Fase 6). */
export function catat(u, action, module, { sebelum = null, sesudah = null } = {}) {
  activityLogs.unshift({ id: idBaru(activityLogs), admin_id: u.id, action, module, before_data: sebelum, after_data: sesudah, ip_address: '127.0.0.1', created_at: sekarang() })
}

const KOLOM_TEKNIS = ['id', 'created_at', 'updated_at']

/**
 * Isi before_data / after_data: hanya kolom yang berubah, ditambah kolom `identitas` (mis. sku, name)
 * supaya log tetap menunjukkan data mana yang diubah.
 */
export function perubahan(lama, baru, identitas = []) {
  const sebelum = {}
  const sesudah = {}
  for (const k of Object.keys({ ...lama, ...baru })) {
    if (KOLOM_TEKNIS.includes(k)) continue
    const beda = JSON.stringify(lama[k] ?? null) !== JSON.stringify(baru[k] ?? null)
    if (beda || identitas.includes(k)) {
      sebelum[k] = lama[k] ?? null
      sesudah[k] = baru[k] ?? null
    }
  }
  return { sebelum, sesudah }
}

export function halaman(daftar, query, bawaan = 20) {
  const perPage = Math.min(Number(query.per_page) || bawaan, 100)
  const total = daftar.length
  const totalPages = Math.max(1, Math.ceil(total / perPage))
  const page = Math.min(Math.max(Number(query.page) || 1, 1), totalPages)
  return { data: daftar.slice((page - 1) * perPage, page * perPage), meta: { page, per_page: perPage, total, total_pages: totalPages } }
}

// ---------- menu & dashboard ----------

function daftarMenu({ headers }) {
  const { u, gagal } = cekAdmin(headers)
  if (gagal) return gagal
  const izin = izinPengguna(u)
  const anak = menus.filter((m) => m.parent_id !== null && m.is_active && (!m.izin || izin.includes(m.izin)))
  const induk = menus.filter((m) => m.parent_id === null && m.is_active && anak.some((a) => a.parent_id === m.id))
  const data = [...induk, ...anak].map((m) => ({ id: m.id, parent_id: m.parent_id, name: m.name, code: m.code, route: m.route, icon: m.icon, sort_order: m.sort_order }))
  return { data: { data } }
}

const punyaFotoUtama = (productId) => productImages.some((i) => i.product_id === productId && i.is_primary)

function dashboard({ headers }) {
  const { u, gagal } = cekAdmin(headers)
  if (gagal) return gagal
  const izin = new Set(izinPengguna(u))
  const kini = sekarang()
  const aktif = products.filter((p) => p.is_active)
  const habis = aktif.filter((p) => p.stock_status === 'HABIS')
  const tanpaFoto = aktif.filter((p) => !punyaFotoUtama(p.id))
  const gagalSync = qontakContacts.filter((c) => c.sync_status === 'FAILED')
  const batasBanner = new Date(Date.parse(`${kini}Z`) + 7 * 86400 * 1000).toISOString().slice(0, 19)
  const bannerBerakhir = banners.filter((b) => b.is_active && b.end_at && b.end_at >= kini && b.end_at <= batasBanner)

  // Kartu ringkasan sesuai izin user.
  const summary = []
  if (izin.has('product.view')) {
    summary.push({ key: 'produk_aktif', label: 'Produk aktif', value: aktif.length, link: '/admin/produk?active=1' })
    summary.push({ key: 'stok_habis', label: 'Produk stok habis', value: habis.length, link: '/admin/produk?stock=HABIS&active=1' })
  }
  if (izin.has('member.view')) {
    summary.push({ key: 'member', label: 'Member terdaftar', value: users.filter((x) => x.roles.includes('CUSTOMER')).length, link: '/admin/member' })
    const bulan = kini.slice(0, 7)
    const poin = pointTransactions.filter((t) => t.type === 'EARN' && t.created_at.startsWith(bulan)).reduce((n, t) => n + t.points, 0)
    summary.push({ key: 'poin_bulan_ini', label: 'Poin diberikan bulan ini', value: poin, link: null })
  }
  if (izin.has('banner.view')) {
    const tayang = banners.filter((b) => b.is_active && (!b.start_at || b.start_at <= kini) && (!b.end_at || b.end_at >= kini))
    summary.push({ key: 'banner_tayang', label: 'Banner tayang', value: tayang.length, link: '/admin/banner' })
  }
  if (izin.has('qontak.retry')) {
    summary.push({ key: 'sync_gagal', label: 'Sync Qontak gagal', value: gagalSync.length, link: '/admin/sync-qontak?status=FAILED' })
  }

  // "Perlu perhatian": paling banyak 3 per jenis, jumlah lengkapnya di attention_total.
  const attention = []
  const attentionTotal = {}
  const tambah = (jenis, daftar, bentuk) => {
    attentionTotal[jenis] = daftar.length
    attention.push(...daftar.slice(0, 3).map((x) => ({ type: jenis, id: x.id, ...bentuk(x) })))
  }
  if (izin.has('qontak.retry')) {
    tambah('sync_gagal', gagalSync, (c) => {
      const log = syncLogs.filter((l) => l.user_id === c.user_id).at(-1)
      return { title: namaPengguna(c.user_id), sync_status: c.sync_status, retry_count: log?.retry_count ?? 0, status_code: log?.status_code ?? null }
    })
  }
  if (izin.has('product.view')) tambah('stok_habis', habis, (p) => ({ title: p.name, sku: p.sku }))
  if (izin.has('banner.view')) tambah('banner_berakhir', bannerBerakhir, (b) => ({ title: b.title, end_at: b.end_at }))
  if (izin.has('product.view')) tambah('tanpa_foto', tanpaFoto, (p) => ({ title: p.name, sku: p.sku }))

  // Log: semua admin bila punya log.view, selain itu hanya aktivitas sendiri.
  const semua = izin.has('log.view')
  const logs = activityLogs
    .filter((l) => semua || l.admin_id === u.id)
    .slice(0, 6)
    .map((l) => ({ id: l.id, admin_name: namaPengguna(l.admin_id), action: l.action, module: l.module, created_at: l.created_at }))

  return { data: { data: { summary, attention, attention_total: attentionTotal, logs, logs_scope: semua ? 'all' : 'mine' } } }
}

// ---------- produk ----------

function itemProdukAdmin(p) {
  const c = categories.find((x) => x.id === p.category_id)
  const b = brands.find((x) => x.id === p.brand_id)
  const foto = productImages.filter((i) => i.product_id === p.id)
  return {
    id: p.id, sku: p.sku, name: p.name, slug: p.slug, price_general: p.price_general, show_price: p.show_price !== false, unit_sale: p.unit_sale,
    stock_status: p.stock_status, stock_qty_label: p.stock_qty_label, is_featured: p.is_featured, is_active: p.is_active,
    updated_at: p.updated_at, category: c ? { id: c.id, name: c.name } : null, brand: b ? { id: b.id, name: b.name } : null,
    primary_image_url: foto.find((i) => i.is_primary)?.image_url ?? null, image_count: foto.length,
  }
}

// Detail untuk form & pratinjau admin: semua kolom + relasi (brand, kategori + induk, foto, ruangan).
function detailProdukAdmin(p) {
  const b = brands.find((x) => x.id === p.brand_id)
  const c = categories.find((x) => x.id === p.category_id)
  const induk = c?.parent_id ? categories.find((x) => x.id === c.parent_id) : null
  const roomIds = productRooms.filter((pr) => pr.product_id === p.id).map((pr) => pr.room_id)
  const ringkas = (x) => (x ? { id: x.id, name: x.name, slug: x.slug } : null)
  return {
    ...p,
    brand: ringkas(b),
    category: c ? { ...ringkas(c), parent: ringkas(induk) } : null,
    images: productImages.filter((i) => i.product_id === p.id).sort((a, b2) => a.sort_order - b2.sort_order),
    room_ids: roomIds,
    rooms: rooms.filter((r) => roomIds.includes(r.id)).map(ringkas),
  }
}

function daftarProdukAdmin({ headers, query }) {
  const { gagal } = cekAdmin(headers, 'product.view')
  if (gagal) return gagal
  let hasil = [...products]
  const q = huruf(query.q).trim()
  if (q) hasil = hasil.filter((p) => huruf(p.name).includes(q) || huruf(p.sku).includes(q))
  if (query.category) {
    const ids = [Number(query.category), ...categories.filter((c) => sama(c.parent_id, query.category)).map((c) => c.id)]
    hasil = hasil.filter((p) => ids.includes(p.category_id))
  }
  if (query.brand) hasil = hasil.filter((p) => sama(p.brand_id, query.brand))
  if (query.stock) hasil = hasil.filter((p) => p.stock_status === query.stock)
  if (query.active === '1' || query.active === '0') hasil = hasil.filter((p) => p.is_active === (query.active === '1'))
  if (query.foto === '1' || query.foto === '0') hasil = hasil.filter((p) => punyaFotoUtama(p.id) === (query.foto === '1'))
  if (query.room) hasil = hasil.filter((p) => productRooms.some((pr) => pr.product_id === p.id && sama(pr.room_id, query.room)))
  hasil.sort((a, b) => b.updated_at.localeCompare(a.updated_at) || b.id - a.id)
  const { data, meta } = halaman(hasil, query)
  return { data: { data: data.map(itemProdukAdmin), meta } }
}

function produkAdmin({ headers, params }) {
  const { gagal } = cekAdmin(headers, 'product.view')
  if (gagal) return gagal
  const p = products.find((x) => sama(x.id, params.id))
  return p ? { data: { data: detailProdukAdmin(p) } } : tidakDitemukan('Produk tidak ditemukan.')
}

function cekSku({ headers, query }) {
  const { gagal } = cekAdmin(headers, ['product.create', 'product.update'])
  if (gagal) return gagal
  const dipakai = products.some((p) => huruf(p.sku) === huruf(String(query.sku ?? '').trim()) && !sama(p.id, query.exclude))
  return { data: { data: { available: !dipakai } } }
}

function validasiProduk(body, idSendiri = null) {
  const errors = {}
  const wajib = (kunci, pesan) => {
    if (kosong(body[kunci])) errors[kunci] = [pesan]
  }
  wajib('sku', 'SKU wajib diisi')
  wajib('name', 'Nama produk wajib diisi')
  wajib('slug', 'Slug wajib diisi')
  wajib('category_id', 'Pilih kategori')
  wajib('unit_sale', 'Satuan jual wajib diisi')
  if (!errors.slug && !POLA_SLUG.test(body.slug)) errors.slug = ['Slug hanya boleh huruf kecil, angka, dan tanda hubung']
  // Harga wajib hanya bila ditampilkan di website (show_price); bila disembunyikan boleh kosong.
  if (kosong(body.price_general)) {
    if (body.show_price !== false) errors.price_general = ['Harga wajib diisi, atau matikan "Tampilkan harga di website"']
  } else if (!(Number(body.price_general) >= 0)) errors.price_general = ['Harga tidak valid']
  if (!(Number.isInteger(Number(body.min_order)) && Number(body.min_order) >= 1)) errors.min_order = ['Minimal pembelian paling sedikit 1']
  if (!STATUS_STOK.includes(body.stock_status)) errors.stock_status = ['Pilih status stok']
  if (!errors.category_id && !categories.some((c) => sama(c.id, body.category_id))) errors.category_id = ['Kategori tidak ditemukan']
  if (!kosong(body.brand_id) && !brands.some((b) => sama(b.id, body.brand_id))) errors.brand_id = ['Brand tidak ditemukan']
  const foto = Array.isArray(body.images) ? body.images : []
  if (body.is_active && !foto.some((f) => f.is_primary)) errors.images = ['Tambahkan minimal 1 foto utama, atau matikan "Tampil di website".']
  if (Object.keys(errors).length) return galatIsian(errors)

  if (products.some((p) => huruf(p.sku) === huruf(body.sku.trim()) && !sama(p.id, idSendiri))) {
    return { status: 409, data: { message: 'SKU sudah dipakai produk lain.', errors: { sku: ['SKU sudah dipakai produk lain'] } } }
  }
  if (products.some((p) => p.slug === body.slug && !sama(p.id, idSendiri))) {
    return { status: 409, data: { message: 'Slug sudah dipakai produk lain.', errors: { slug: ['Slug sudah dipakai produk lain'] } } }
  }
  return null
}

function simpanRelasiProduk(id, body) {
  const fotoLama = productImages.filter((i) => i.product_id === id)
  for (const f of fotoLama) productImages.splice(productImages.indexOf(f), 1)
  ;(body.images ?? []).forEach((f, i) => {
    productImages.push({ id: idBaru(productImages), product_id: id, image_url: f.image_url, is_primary: Boolean(f.is_primary), sort_order: i, created_at: sekarang() })
  })
  const ruanganLama = productRooms.filter((pr) => pr.product_id === id)
  for (const r of ruanganLama) productRooms.splice(productRooms.indexOf(r), 1)
  for (const roomId of new Set(body.room_ids ?? [])) productRooms.push({ product_id: id, room_id: Number(roomId) })
}

function kolomProduk(body) {
  return {
    sku: body.sku.trim(), name: body.name.trim(), slug: body.slug, category_id: Number(body.category_id),
    brand_id: kosong(body.brand_id) ? null : Number(body.brand_id), description: body.description || null,
    specifications: body.specifications && Object.keys(body.specifications).length ? body.specifications : null,
    datasheet_pdf_url: body.datasheet_pdf_url || null, unit_sale: body.unit_sale.trim(),
    price_general: kosong(body.price_general) ? null : Number(body.price_general).toFixed(2), show_price: body.show_price !== false,
    min_order: Number(body.min_order), stock_status: body.stock_status, stock_qty_label: body.stock_qty_label || null,
    is_featured: Boolean(body.is_featured), is_active: Boolean(body.is_active),
    meta_title: body.meta_title || null, meta_description: body.meta_description || null, updated_at: sekarang(),
  }
}

function buatProduk({ headers, body }) {
  const { u, gagal } = cekAdmin(headers, 'product.create')
  if (gagal) return gagal
  const salah = validasiProduk(body)
  if (salah) return salah
  const p = { id: idBaru(products), ...kolomProduk(body), created_at: sekarang() }
  products.push(p)
  simpanRelasiProduk(p.id, body)
  catat(u, 'product.create', 'Produk', { sesudah: { sku: p.sku, name: p.name, price_general: p.price_general, stock_status: p.stock_status, is_active: p.is_active } })
  simpanDb()
  return { status: 201, data: { data: detailProdukAdmin(p) } }
}

function ubahProduk({ headers, params, body }) {
  const { u, gagal } = cekAdmin(headers, 'product.update')
  if (gagal) return gagal
  const p = products.find((x) => sama(x.id, params.id))
  if (!p) return tidakDitemukan('Produk tidak ditemukan.')
  const salah = validasiProduk(body, p.id)
  if (salah) return salah
  const lama = { ...p }
  Object.assign(p, kolomProduk(body))
  simpanRelasiProduk(p.id, body)
  catat(u, 'product.update', 'Produk', perubahan(lama, p, ['sku', 'name']))
  simpanDb()
  return { data: { data: detailProdukAdmin(p) } }
}

function patchProduk({ headers, params, body }) {
  const { u, gagal } = cekAdmin(headers, 'product.update')
  if (gagal) return gagal
  const p = products.find((x) => sama(x.id, params.id))
  if (!p) return tidakDitemukan('Produk tidak ditemukan.')
  if (body.is_active === true && !punyaFotoUtama(p.id)) {
    return { status: 422, data: { message: 'Tambahkan minimal 1 foto utama dulu.', errors: { is_active: ['Tambahkan minimal 1 foto utama dulu'] } } }
  }
  // Ubah cepat harga dari daftar produk: { price_general, unit_sale, show_price } (sebagian boleh).
  const hargaBaru = 'price_general' in body ? (kosong(body.price_general) ? null : body.price_general) : p.price_general
  const tampilBaru = 'show_price' in body ? body.show_price !== false : p.show_price !== false
  const errors = {}
  if (hargaBaru != null && !(Number(hargaBaru) >= 0)) errors.price_general = ['Harga tidak valid']
  else if (tampilBaru && hargaBaru == null) errors.price_general = ['Isi harga, atau matikan "Tampilkan harga di website"']
  if ('unit_sale' in body && kosong(body.unit_sale)) errors.unit_sale = ['Satuan jual wajib diisi']
  if (Object.keys(errors).length) return galatIsian(errors)
  const lama = { ...p }
  if ('price_general' in body) p.price_general = hargaBaru == null ? null : Number(hargaBaru).toFixed(2)
  if ('show_price' in body) p.show_price = tampilBaru
  if ('unit_sale' in body) p.unit_sale = String(body.unit_sale).trim()
  for (const kunci of ['is_featured', 'is_active', 'stock_status', 'stock_qty_label']) {
    if (kunci in body) p[kunci] = body[kunci]
  }
  p.updated_at = sekarang()
  catat(u, 'product.update', 'Produk', perubahan(lama, p, ['sku', 'name']))
  simpanDb()
  return { data: { data: itemProdukAdmin(p) } }
}

// ---------- kategori, brand, ruangan ----------

const MASTER = {
  categories: { tabel: categories, izin: 'category', modul: 'Kategori', urut: true, gambar: 'image_url', label: 'Kategori' },
  brands: { tabel: brands, izin: 'brand', modul: 'Brand', urut: false, gambar: 'logo_url', label: 'Brand' },
  rooms: { tabel: rooms, izin: 'room', modul: 'Ruangan', urut: true, gambar: 'image_cover', label: 'Ruangan' },
}

function jumlahProduk(jenis, baris) {
  if (jenis === 'brands') return products.filter((p) => p.brand_id === baris.id).length
  if (jenis === 'rooms') return productRooms.filter((pr) => pr.room_id === baris.id).length
  const ids = [baris.id, ...categories.filter((c) => c.parent_id === baris.id).map((c) => c.id)]
  return products.filter((p) => ids.includes(p.category_id)).length
}

function daftarMaster(jenis) {
  return ({ headers }) => {
    const m = MASTER[jenis]
    const { gagal } = cekAdmin(headers, `${m.izin}.view`)
    if (gagal) return gagal
    const urut = m.urut ? (a, b) => a.sort_order - b.sort_order || a.id - b.id : (a, b) => a.name.localeCompare(b.name, 'id')
    const data = [...m.tabel].sort(urut).map((x) => ({ ...x, product_count: jumlahProduk(jenis, x) }))
    return { data: { data } }
  }
}

function validasiMaster(jenis, body, idSendiri = null) {
  const m = MASTER[jenis]
  const errors = {}
  if (kosong(body.name)) errors.name = ['Nama wajib diisi']
  if (kosong(body.slug)) errors.slug = ['Slug wajib diisi']
  else if (!POLA_SLUG.test(body.slug)) errors.slug = ['Slug hanya boleh huruf kecil, angka, dan tanda hubung']
  if (jenis === 'rooms' && kosong(body.image_cover)) errors.image_cover = ['Cover ruangan wajib diunggah']
  if (jenis === 'categories' && !kosong(body.parent_id)) {
    const induk = categories.find((c) => sama(c.id, body.parent_id))
    if (!induk || induk.parent_id !== null) errors.parent_id = ['Induk harus kategori utama']
    else if (sama(induk.id, idSendiri)) errors.parent_id = ['Kategori tidak bisa menjadi induk dirinya sendiri']
    else if (idSendiri && categories.some((c) => sama(c.parent_id, idSendiri))) errors.parent_id = ['Kategori yang punya subkategori tidak bisa dipindah ke bawah kategori lain']
  }
  if (Object.keys(errors).length) return galatIsian(errors)
  if (m.tabel.some((x) => x.slug === body.slug && !sama(x.id, idSendiri))) {
    return { status: 409, data: { message: `Slug sudah dipakai ${m.label.toLowerCase()} lain.`, errors: { slug: [`Slug sudah dipakai ${m.label.toLowerCase()} lain`] } } }
  }
  return null
}

function kolomMaster(jenis, body) {
  const m = MASTER[jenis]
  const kolom = { name: body.name.trim(), slug: body.slug, [m.gambar]: body[m.gambar] || null, is_active: body.is_active !== false, updated_at: sekarang() }
  if (jenis === 'categories') kolom.parent_id = kosong(body.parent_id) ? null : Number(body.parent_id)
  return kolom
}

function buatMaster(jenis) {
  return ({ headers, body }) => {
    const m = MASTER[jenis]
    const { u, gagal } = cekAdmin(headers, `${m.izin}.create`)
    if (gagal) return gagal
    const salah = validasiMaster(jenis, body)
    if (salah) return salah
    const baris = { id: idBaru(m.tabel), ...kolomMaster(jenis, body), created_at: sekarang() }
    if (m.urut) {
      const saudara = jenis === 'categories' ? m.tabel.filter((x) => x.parent_id === baris.parent_id) : m.tabel
      baris.sort_order = saudara.reduce((maks, x) => Math.max(maks, x.sort_order), 0) + 1
    }
    m.tabel.push(baris)
    catat(u, `${m.izin}.create`, m.modul, { sesudah: { name: baris.name, slug: baris.slug, is_active: baris.is_active } })
    simpanDb()
    return { status: 201, data: { data: { ...baris, product_count: 0 } } }
  }
}

function ubahMaster(jenis) {
  return ({ headers, params, body }) => {
    const m = MASTER[jenis]
    const { u, gagal } = cekAdmin(headers, `${m.izin}.update`)
    if (gagal) return gagal
    const baris = m.tabel.find((x) => sama(x.id, params.id))
    if (!baris) return tidakDitemukan(`${m.label} tidak ditemukan.`)
    const salah = validasiMaster(jenis, body, baris.id)
    if (salah) return salah
    const lama = { ...baris }
    Object.assign(baris, kolomMaster(jenis, body))
    catat(u, `${m.izin}.update`, m.modul, perubahan(lama, baris, ['name']))
    simpanDb()
    return { data: { data: { ...baris, product_count: jumlahProduk(jenis, baris) } } }
  }
}

function patchMaster(jenis) {
  return ({ headers, params, body }) => {
    const m = MASTER[jenis]
    const { u, gagal } = cekAdmin(headers, `${m.izin}.update`)
    if (gagal) return gagal
    const baris = m.tabel.find((x) => sama(x.id, params.id))
    if (!baris) return tidakDitemukan(`${m.label} tidak ditemukan.`)
    const lama = { ...baris }
    if ('is_active' in body) baris.is_active = Boolean(body.is_active)
    if (m.urut && 'sort_order' in body) baris.sort_order = Number(body.sort_order)
    baris.updated_at = sekarang()
    catat(u, `${m.izin}.update`, m.modul, perubahan(lama, baris, ['name']))
    simpanDb()
    return { data: { data: { ...baris, product_count: jumlahProduk(jenis, baris) } } }
  }
}

// Hapus permanen hanya brand nonaktif yang tidak dipakai produk mana pun (aktif maupun nonaktif).
// Di database FK products.brand_id ON DELETE SET NULL: tanpa cek ini produk diam-diam kehilangan brand.
function hapusBrand({ headers, params }) {
  const { u, gagal } = cekAdmin(headers, 'brand.delete')
  if (gagal) return gagal
  const b = brands.find((x) => sama(x.id, params.id))
  if (!b) return tidakDitemukan('Brand tidak ditemukan.')
  if (b.is_active) return { status: 422, data: { message: 'Nonaktifkan brand dulu sebelum menghapusnya.' } }
  const dipakai = products.filter((p) => p.brand_id === b.id).length
  if (dipakai) {
    return { status: 409, data: { message: `Brand ini masih dipakai ${dipakai} produk. Pindahkan produknya ke brand lain dulu.`, product_count: dipakai } }
  }
  brands.splice(brands.indexOf(b), 1)
  catat(u, 'brand.delete', 'Brand', { sebelum: { name: b.name, slug: b.slug, is_active: b.is_active } })
  simpanDb()
  return { status: 204, data: null }
}

// ---------- unggah file ----------

const ATURAN_UNGGAH = {
  gambar: { tipe: ['image/jpeg', 'image/png', 'image/webp'], maks: 2 * 1024 * 1024, pesan: 'Foto harus JPG, PNG, atau WebP, maksimal 2 MB' },
  dokumen: { tipe: ['application/pdf'], maks: 10 * 1024 * 1024, pesan: 'Dokumen harus PDF, maksimal 10 MB' },
  // Mode mock: 3 MB (disimpan sebagai data URL di localStorage). Usulan server sungguhan: 20 MB.
  video: { tipe: ['video/mp4', 'video/webm'], maks: 3 * 1024 * 1024, pesan: 'Video harus MP4 atau WebM, maksimal 3 MB (mode contoh)' },
}

async function unggah({ headers, body }) {
  const { gagal } = cekAdmin(headers)
  if (gagal) return gagal
  const file = body?.get?.('file')
  const aturan = ATURAN_UNGGAH[body?.get?.('jenis')] ?? ATURAN_UNGGAH.gambar
  if (!file) return galatIsian({ file: ['Pilih file dulu'] })
  if (!aturan.tipe.includes(file.type) || file.size > aturan.maks) return galatIsian({ file: [aturan.pesan] })
  // Mode mock: file disimpan sebagai data URL (backend asli mengembalikan URL penyimpanan).
  const url = await new Promise((selesai, gagalBaca) => {
    const pembaca = new FileReader()
    pembaca.onload = () => selesai(pembaca.result)
    pembaca.onerror = gagalBaca
    pembaca.readAsDataURL(file)
  })
  return { status: 201, data: { data: { url, name: file.name, size: file.size, type: file.type } } }
}

// ---------- sync Qontak ----------

export function bentukKontak(c) {
  const log = syncLogs.filter((l) => l.user_id === c.user_id).at(-1)
  return { id: c.id, user_id: c.user_id, name: namaPengguna(c.user_id), sync_status: c.sync_status, last_synced_at: c.last_synced_at, retry_count: log?.retry_count ?? 0, status_code: log?.status_code ?? null }
}

function ulangSync({ headers, params }) {
  const { u, gagal } = cekAdmin(headers, 'qontak.retry')
  if (gagal) return gagal
  const c = qontakContacts.find((x) => sama(x.id, params.id))
  if (!c) return tidakDitemukan('Kontak tidak ditemukan.')
  antrekan(c)
  catat(u, 'qontak.retry', 'Sync Qontak', { sesudah: { member: namaPengguna(c.user_id), sync_status: 'PENDING' } })
  simpanDb()
  return { status: 202, data: { data: bentukKontak(c) } }
}

// Proses di background selesai ±1,5 detik setelah retry (services/mock/antreanQontak.js).
// Contoh: Siti berhasil, Andi gagal lagi (Qontak 500).
function detailKontak({ headers, params }) {
  const { gagal } = cekAdmin(headers, 'qontak.retry')
  if (gagal) return gagal
  const c = qontakContacts.find((x) => sama(x.id, params.id))
  if (!c) return tidakDitemukan('Kontak tidak ditemukan.')
  if (prosesKontak(c)) simpanDb()
  return { data: { data: bentukKontak(c) } }
}

// Urutan penting: path tetap (check-sku) sebelum path ber-parameter.
export const ruteAdmin = [
  ['get', '/admin/menus', daftarMenu],
  ['get', '/admin/dashboard', dashboard],
  ['get', '/admin/products/check-sku', cekSku],
  ['get', '/admin/products', daftarProdukAdmin],
  ['post', '/admin/products', buatProduk],
  ['get', '/admin/products/:id', produkAdmin],
  ['put', '/admin/products/:id', ubahProduk],
  ['patch', '/admin/products/:id', patchProduk],
  ...Object.keys(MASTER).flatMap((jenis) => [
    ['get', `/admin/${jenis}`, daftarMaster(jenis)],
    ['post', `/admin/${jenis}`, buatMaster(jenis)],
    ['put', `/admin/${jenis}/:id`, ubahMaster(jenis)],
    ['patch', `/admin/${jenis}/:id`, patchMaster(jenis)],
  ]),
  ['delete', '/admin/brands/:id', hapusBrand],
  ['post', '/admin/uploads', unggah],
  ['post', '/admin/qontak/contacts/:id/retry', ulangSync],
  ['get', '/admin/qontak/contacts/:id', detailKontak],
]

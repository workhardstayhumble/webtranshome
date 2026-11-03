// Rute API tiruan konten (Fase 5): banner (termasuk hero Beranda), artikel/promo/event, jasa, dan pengaturan toko.
// TODO: ganti ke API begitu endpoint /banners, /articles, /admin/banners, /admin/articles,
// dan /admin/settings tersedia di backend.
import { publicSettings } from '@/services/mock/data/akun'
import { categories } from '@/services/mock/data/catalog'
import { articles, banners, services } from '@/services/mock/data/konten'
import { idBaru, simpanDb } from '@/services/mock/db'
import { catat, halaman, perubahan } from '@/services/mock/routesAdmin'
import { cekAdmin, ditolak, izinPengguna, namaPengguna, sekarang } from '@/services/mock/sesi'

// HOME_HERO: hero paling atas Beranda (gambar atau video), pengganti ilustrasi rumah 3D (30 Okt 2026).
const POSISI = ['HOME_HERO', 'HOME_SLIDER', 'HOME_SIDE']
const TIPE = ['ARTIKEL', 'PROMO', 'EVENT', 'HALAMAN']
const POLA_SLUG = /^[a-z0-9]+(?:-[a-z0-9]+)*$/
const POLA_TANGGAL = /^\d{4}-\d{2}-\d{2}$/
const kosong = (nilai) => nilai == null || String(nilai).trim() === ''
const sama = (a, b) => String(a) === String(b)
const tidakDitemukan = (pesan) => ({ status: 404, data: { message: pesan } })
const galatIsian = (errors) => ({ status: 422, data: { message: 'Periksa kembali isian yang ditandai.', errors } })
const tokenAcak = () => [...crypto.getRandomValues(new Uint8Array(8))].map((b) => b.toString(16).padStart(2, '0')).join('')

// ---------- banner ----------

/** Status dihitung dari is_active dan rentang tanggal (waktu WIB server). */
function statusBanner(b, kini = sekarang()) {
  if (!b.is_active) return 'NONAKTIF'
  if (b.start_at && b.start_at > kini) return 'TERJADWAL'
  if (b.end_at && b.end_at < kini) return 'BERAKHIR'
  return 'TAYANG'
}

const urutBanner = (a, b) => POSISI.indexOf(a.position) - POSISI.indexOf(b.position) || a.sort_order - b.sort_order || a.id - b.id

function daftarBannerPublik({ query }) {
  const kini = sekarang()
  const data = banners
    .filter((b) => statusBanner(b, kini) === 'TAYANG' && (!query.position || b.position === query.position))
    .sort(urutBanner)
    .map(({ id, title, image_url, link_url, position, start_at, end_at, subtitle, video_url, cta_label, cta2_label, cta2_url }) => ({
      id, title, image_url, link_url, position, start_at, end_at,
      ...(position === 'HOME_HERO' ? { subtitle, video_url, cta_label, cta2_label, cta2_url } : {}),
    }))
  return { data: { data } }
}

const bentukBannerAdmin = (b, kini) => ({ ...b, status: statusBanner(b, kini) })

function daftarBannerAdmin({ headers }) {
  const { gagal } = cekAdmin(headers, 'banner.view')
  if (gagal) return gagal
  const kini = sekarang()
  return { data: { data: [...banners].sort(urutBanner).map((b) => bentukBannerAdmin(b, kini)) } }
}

function ambilBanner({ headers, params }) {
  const { gagal } = cekAdmin(headers, 'banner.view')
  if (gagal) return gagal
  const b = banners.find((x) => sama(x.id, params.id))
  return b ? { data: { data: bentukBannerAdmin(b) } } : tidakDitemukan('Banner tidak ditemukan.')
}

// Tanggal dikirim YYYY-MM-DD: mulai pukul 00.00, selesai pukul 23.59 WIB. Selesai kosong = tayang terus.
function validasiBanner(body) {
  const errors = {}
  if (kosong(body.title)) errors.title = ['Judul wajib diisi']
  if (kosong(body.image_url)) errors.image_url = ['Unggah gambar banner']
  if (!POSISI.includes(body.position)) errors.position = ['Pilih posisi banner']
  if (kosong(body.start_date) || !POLA_TANGGAL.test(body.start_date)) errors.start_date = ['Isi tanggal mulai']
  if (!kosong(body.end_date) && !POLA_TANGGAL.test(body.end_date)) errors.end_date = ['Tanggal selesai tidak valid']
  else if (!errors.start_date && !kosong(body.end_date) && body.end_date < body.start_date) errors.end_date = ['Tanggal selesai harus sama atau setelah tanggal mulai']
  const link = String(body.link_url ?? '').trim()
  if (link && !link.startsWith('/') && !link.startsWith('https://')) errors.link_url = ['Tautan harus diawali / (halaman Transhome) atau https://']
  if (body.position === 'HOME_HERO') {
    if (String(body.subtitle ?? '').length > 160) errors.subtitle = ['Subjudul maksimal 160 karakter']
    if (String(body.cta_label ?? '').length > 30) errors.cta_label = ['Label tombol maksimal 30 karakter']
    const video = String(body.video_url ?? '').trim()
    if (video && !/^(https:\/\/|\/|data:video\/)/.test(video)) errors.video_url = ['Video tidak valid, unggah ulang']
    // Tombol kedua tampil hanya bila label dan tautannya sama-sama diisi.
    const label2 = String(body.cta2_label ?? '').trim()
    const tautan2 = String(body.cta2_url ?? '').trim()
    if (label2.length > 30) errors.cta2_label = ['Label tombol kedua maksimal 30 karakter']
    else if (!label2 && tautan2) errors.cta2_label = ['Isi label tombol kedua, atau kosongkan tautannya']
    if (label2 && !tautan2) errors.cta2_url = ['Isi tautan tombol kedua']
    else if (tautan2 && !tautan2.startsWith('/') && !tautan2.startsWith('https://')) errors.cta2_url = ['Tautan harus diawali / (halaman Transhome) atau https://']
  }
  return Object.keys(errors).length ? galatIsian(errors) : null
}

const kolomBanner = (body) => ({
  title: body.title.trim(), image_url: body.image_url, link_url: String(body.link_url ?? '').trim() || null, position: body.position,
  // Kolom khusus hero; banner posisi lain selalu kosong.
  subtitle: body.position === 'HOME_HERO' ? String(body.subtitle ?? '').trim() || null : null,
  video_url: body.position === 'HOME_HERO' ? String(body.video_url ?? '').trim() || null : null,
  cta_label: body.position === 'HOME_HERO' ? String(body.cta_label ?? '').trim() || null : null,
  cta2_label: body.position === 'HOME_HERO' ? String(body.cta2_label ?? '').trim() || null : null,
  cta2_url: body.position === 'HOME_HERO' ? String(body.cta2_url ?? '').trim() || null : null,
  start_at: `${body.start_date}T00:00:00`, end_at: kosong(body.end_date) ? null : `${body.end_date}T23:59:59`,
  is_active: body.is_active !== false, updated_at: sekarang(),
})

const urutanTerakhir = (posisi, kecualiId = null) =>
  banners.filter((b) => b.position === posisi && !sama(b.id, kecualiId)).reduce((maks, b) => Math.max(maks, b.sort_order), 0) + 1

function buatBanner({ headers, body }) {
  const { u, gagal } = cekAdmin(headers, 'banner.create')
  if (gagal) return gagal
  const salah = validasiBanner(body)
  if (salah) return salah
  const b = { id: idBaru(banners), ...kolomBanner(body), sort_order: urutanTerakhir(body.position), created_at: sekarang() }
  banners.push(b)
  catat(u, 'banner.create', 'Banner', { sesudah: { title: b.title, start_at: b.start_at, end_at: b.end_at } })
  simpanDb()
  return { status: 201, data: { data: bentukBannerAdmin(b) } }
}

function ubahBanner({ headers, params, body }) {
  const { u, gagal } = cekAdmin(headers, 'banner.update')
  if (gagal) return gagal
  const b = banners.find((x) => sama(x.id, params.id))
  if (!b) return tidakDitemukan('Banner tidak ditemukan.')
  const salah = validasiBanner(body)
  if (salah) return salah
  const sebelum = { title: b.title, start_at: b.start_at, end_at: b.end_at, is_active: b.is_active }
  const pindahPosisi = b.position !== body.position
  Object.assign(b, kolomBanner(body))
  if (pindahPosisi) b.sort_order = urutanTerakhir(b.position, b.id)
  catat(u, 'banner.update', 'Banner', { sebelum, sesudah: { title: b.title, start_at: b.start_at, end_at: b.end_at, is_active: b.is_active } })
  simpanDb()
  return { data: { data: bentukBannerAdmin(b) } }
}

function patchBanner({ headers, params, body }) {
  const { u, gagal } = cekAdmin(headers, 'banner.update')
  if (gagal) return gagal
  const b = banners.find((x) => sama(x.id, params.id))
  if (!b) return tidakDitemukan('Banner tidak ditemukan.')
  if ('is_active' in body) b.is_active = Boolean(body.is_active)
  if ('sort_order' in body) b.sort_order = Number(body.sort_order)
  b.updated_at = sekarang()
  catat(u, 'banner.update', 'Banner')
  simpanDb()
  return { data: { data: bentukBannerAdmin(b) } }
}

// ---------- artikel, promo, event ----------

/** Ringkasan untuk kartu: paragraf pertama (bukan subjudul/daftar), maksimal ±160 karakter. */
function kutipan(isi) {
  const paragraf = String(isi ?? '').split(/\n\s*\n/).map((p) => p.trim()).find((p) => p && !p.startsWith('## ') && !p.startsWith('- ')) ?? ''
  const teks = paragraf.replace(/\s+/g, ' ')
  return teks.length > 160 ? `${teks.slice(0, 157).replace(/\s+\S*$/, '')}…` : teks
}

const ringkasArtikel = (a) => ({
  id: a.id, type: a.type, title: a.title, slug: a.slug, thumbnail_url: a.thumbnail_url, excerpt: kutipan(a.content), published_at: a.published_at,
})

// Publik: hanya yang terbit. Tipe HALAMAN tidak ikut daftar kecuali diminta.
function daftarArtikelPublik({ query }) {
  const tipe = String(query.type ?? '').split(',').filter((t) => TIPE.includes(t))
  const hasil = articles
    .filter((a) => a.status === 'PUBLISHED' && (tipe.length ? tipe.includes(a.type) : a.type !== 'HALAMAN'))
    .sort((a, b) => b.published_at.localeCompare(a.published_at) || b.id - a.id)
  const { data, meta } = halaman(hasil, query, 9)
  return { data: { data: data.map(ringkasArtikel), meta } }
}

// Draf hanya bisa dibuka dengan token pratinjau yang benar; tanpa token → 404 seperti artikel yang tidak ada.
function detailArtikelPublik({ params, query }) {
  const a = articles.find((x) => x.slug === params.slug)
  const boleh = a && (a.status === 'PUBLISHED' || (query.preview && query.preview === a.preview_token))
  if (!boleh) return tidakDitemukan('Artikel tidak ditemukan.')
  return { data: { data: { ...ringkasArtikel(a), content: a.content, status: a.status } } }
}

const bentukArtikelAdmin = (a) => ({
  id: a.id, type: a.type, title: a.title, slug: a.slug, thumbnail_url: a.thumbnail_url, status: a.status,
  published_at: a.published_at, updated_at: a.updated_at, author_name: a.author_id ? namaPengguna(a.author_id) : null,
})

function daftarArtikelAdmin({ headers, query }) {
  const { gagal } = cekAdmin(headers, 'article.view')
  if (gagal) return gagal
  let hasil = [...articles]
  const q = String(query.q ?? '').trim().toLowerCase()
  if (q) hasil = hasil.filter((a) => a.title.toLowerCase().includes(q) || a.slug.includes(q))
  if (TIPE.includes(query.type)) hasil = hasil.filter((a) => a.type === query.type)
  if (['DRAFT', 'PUBLISHED'].includes(query.status)) hasil = hasil.filter((a) => a.status === query.status)
  hasil.sort((a, b) => b.updated_at.localeCompare(a.updated_at) || b.id - a.id)
  const { data, meta } = halaman(hasil, query)
  return { data: { data: data.map(bentukArtikelAdmin), meta } }
}

function ambilArtikelAdmin({ headers, params }) {
  const { gagal } = cekAdmin(headers, 'article.view')
  if (gagal) return gagal
  const a = articles.find((x) => sama(x.id, params.id))
  return a ? { data: { data: { ...bentukArtikelAdmin(a), content: a.content, preview_token: a.preview_token } } } : tidakDitemukan('Artikel tidak ditemukan.')
}

function validasiArtikel(body, idSendiri = null) {
  const errors = {}
  if (!TIPE.includes(body.type)) errors.type = ['Pilih tipe']
  if (kosong(body.title)) errors.title = ['Judul wajib diisi']
  if (kosong(body.slug)) errors.slug = ['Slug wajib diisi']
  else if (!POLA_SLUG.test(body.slug)) errors.slug = ['Slug hanya boleh huruf kecil, angka, dan tanda hubung']
  if (kosong(body.content)) errors.content = ['Isi wajib diisi']
  if (Object.keys(errors).length) return galatIsian(errors)
  if (articles.some((a) => a.slug === body.slug && !sama(a.id, idSendiri))) {
    return { status: 409, data: { message: 'Slug sudah dipakai artikel lain.', errors: { slug: ['Slug sudah dipakai artikel lain'] } } }
  }
  return null
}

const kolomArtikel = (body) => ({
  type: body.type, title: body.title.trim(), slug: body.slug, thumbnail_url: body.thumbnail_url || null, content: body.content.trim(), updated_at: sekarang(),
})

function buatArtikel({ headers, body }) {
  const { u, gagal } = cekAdmin(headers, 'article.create')
  if (gagal) return gagal
  const salah = validasiArtikel(body)
  if (salah) return salah
  const a = { id: idBaru(articles), ...kolomArtikel(body), status: 'DRAFT', published_at: null, author_id: u.id, preview_token: tokenAcak(), created_at: sekarang() }
  articles.push(a)
  catat(u, 'article.create', 'Artikel', { sesudah: { title: a.title, type: a.type } })
  simpanDb()
  return { status: 201, data: ambilArtikelAdmin({ headers, params: { id: a.id } }).data }
}

function ubahArtikel({ headers, params, body }) {
  const { u, gagal } = cekAdmin(headers, 'article.update')
  if (gagal) return gagal
  const a = articles.find((x) => sama(x.id, params.id))
  if (!a) return tidakDitemukan('Artikel tidak ditemukan.')
  const salah = validasiArtikel(body, a.id)
  if (salah) return salah
  Object.assign(a, kolomArtikel(body))
  catat(u, 'article.update', 'Artikel', { sesudah: { title: a.title } })
  simpanDb()
  return ambilArtikelAdmin({ headers, params: { id: a.id } })
}

// Terbitkan / batalkan terbit (izin article.publish). published_at diisi saat terbit.
function terbitkanArtikel({ headers, params, body }) {
  const { u, gagal } = cekAdmin(headers, 'article.publish')
  if (gagal) return gagal
  const a = articles.find((x) => sama(x.id, params.id))
  if (!a) return tidakDitemukan('Artikel tidak ditemukan.')
  if (!['DRAFT', 'PUBLISHED'].includes(body.status)) return galatIsian({ status: ['Status tidak valid'] })
  const sebelum = { status: a.status }
  a.status = body.status
  a.published_at = body.status === 'PUBLISHED' ? sekarang() : null
  a.updated_at = sekarang()
  catat(u, 'article.publish', 'Artikel', { sebelum, sesudah: { status: a.status, title: a.title } })
  simpanDb()
  return ambilArtikelAdmin({ headers, params: { id: a.id } })
}

// ---------- jasa (KEPUTUSAN #20) ----------
// Konsultasi jasa hanya lewat WhatsApp (KEPUTUSAN #12): API cukup menyediakan daftar jasa, tanpa formulir.

const ringkasKategoriJasa = (ids) =>
  ids.map((id) => categories.find((c) => c.id === id)).filter(Boolean).map(({ id, name, slug, parent_id }) => ({ id, name, slug, parent_id }))
const itemJasaPublik = ({ id, name, slug, summary, image_url, service_area }) => ({ id, name, slug, summary, image_url, service_area })
const urutJasa = (a, b) => a.sort_order - b.sort_order || a.id - b.id

// ?category=<slug>: jasa yang terhubung ke kategori itu atau induknya (dipakai halaman detail produk). ?limit=n.
function daftarJasaPublik({ query }) {
  let data = services.filter((s) => s.is_active)
  if (query.category) {
    const c = categories.find((x) => x.slug === query.category)
    if (!c) return { data: { data: [] } }
    const ids = [c.id, c.parent_id].filter(Boolean)
    data = data.filter((s) => s.category_ids.some((id) => ids.includes(id)))
  }
  data = [...data].sort(urutJasa)
  if (Number(query.limit) > 0) data = data.slice(0, Number(query.limit))
  return { data: { data: data.map(itemJasaPublik) } }
}

function detailJasaPublik({ params }) {
  const s = services.find((x) => x.slug === params.slug && x.is_active)
  if (!s) return tidakDitemukan('Jasa tidak ditemukan.')
  const kategoriAktif = ringkasKategoriJasa(s.category_ids).filter((c) => categories.find((x) => x.id === c.id)?.is_active)
  return { data: { data: { ...itemJasaPublik(s), description: s.description, categories: kategoriAktif.map(({ name, slug }) => ({ name, slug })) } } }
}

const bentukJasaAdmin = (s) => ({ ...s, categories: ringkasKategoriJasa(s.category_ids) })

function daftarJasaAdmin({ headers }) {
  const { gagal } = cekAdmin(headers, 'service.view')
  if (gagal) return gagal
  return { data: { data: [...services].sort(urutJasa).map(bentukJasaAdmin) } }
}

function ambilJasaAdmin({ headers, params }) {
  const { gagal } = cekAdmin(headers, 'service.view')
  if (gagal) return gagal
  const s = services.find((x) => sama(x.id, params.id))
  return s ? { data: { data: bentukJasaAdmin(s) } } : tidakDitemukan('Jasa tidak ditemukan.')
}

function validasiJasa(body, kecualiId = null) {
  const errors = {}
  const nama = String(body.name ?? '').trim()
  if (!nama) errors.name = ['Nama jasa wajib diisi']
  else if (nama.length > 100) errors.name = ['Nama jasa maksimal 100 karakter']
  if (kosong(body.slug)) errors.slug = ['Slug wajib diisi']
  else if (!POLA_SLUG.test(body.slug)) errors.slug = ['Slug hanya boleh huruf kecil, angka, dan tanda hubung']
  const ringkas = String(body.summary ?? '').trim()
  if (!ringkas) errors.summary = ['Ringkasan wajib diisi']
  else if (ringkas.length > 160) errors.summary = ['Ringkasan maksimal 160 karakter']
  if (String(body.service_area ?? '').trim().length > 100) errors.service_area = ['Area layanan maksimal 100 karakter']
  if (!Array.isArray(body.category_ids ?? []) || (body.category_ids ?? []).some((id) => !categories.some((c) => c.id === Number(id)))) {
    errors.category_ids = ['Kategori tidak valid']
  }
  if (Object.keys(errors).length) return galatIsian(errors)
  if (services.some((s) => s.slug === body.slug && !sama(s.id, kecualiId))) {
    return { status: 409, data: { message: 'Slug sudah dipakai jasa lain.', errors: { slug: ['Slug sudah dipakai jasa lain'] } } }
  }
  return null
}

const kolomJasa = (body) => ({
  name: String(body.name).trim(), slug: body.slug, summary: String(body.summary).trim(),
  description: String(body.description ?? '').trim(), image_url: body.image_url || null,
  service_area: String(body.service_area ?? '').trim() || null,
  category_ids: [...new Set((body.category_ids ?? []).map(Number))], is_active: body.is_active !== false, updated_at: sekarang(),
})

function buatJasa({ headers, body }) {
  const { u, gagal } = cekAdmin(headers, 'service.create')
  if (gagal) return gagal
  const salah = validasiJasa(body)
  if (salah) return salah
  const s = { id: idBaru(services), ...kolomJasa(body), sort_order: services.reduce((m, x) => Math.max(m, x.sort_order), 0) + 1, created_at: sekarang() }
  services.push(s)
  catat(u, 'service.create', 'Jasa', { sesudah: { name: s.name, category_ids: s.category_ids } })
  simpanDb()
  return { status: 201, data: { data: bentukJasaAdmin(s) } }
}

function ubahJasa({ headers, params, body }) {
  const { u, gagal } = cekAdmin(headers, 'service.update')
  if (gagal) return gagal
  const s = services.find((x) => sama(x.id, params.id))
  if (!s) return tidakDitemukan('Jasa tidak ditemukan.')
  const salah = validasiJasa(body, s.id)
  if (salah) return salah
  const lama = { ...s }
  Object.assign(s, kolomJasa(body))
  catat(u, 'service.update', 'Jasa', perubahan(lama, s, ['name']))
  simpanDb()
  return { data: { data: bentukJasaAdmin(s) } }
}

function patchJasa({ headers, params, body }) {
  const { u, gagal } = cekAdmin(headers, 'service.update')
  if (gagal) return gagal
  const s = services.find((x) => sama(x.id, params.id))
  if (!s) return tidakDitemukan('Jasa tidak ditemukan.')
  const lama = { ...s }
  if ('is_active' in body) s.is_active = Boolean(body.is_active)
  if ('sort_order' in body) s.sort_order = Number(body.sort_order)
  s.updated_at = sekarang()
  catat(u, 'service.update', 'Jasa', perubahan(lama, s, ['name']))
  simpanDb()
  return { data: { data: bentukJasaAdmin(s) } }
}

// ---------- pengaturan toko ----------

function pengaturanAdmin({ headers }) {
  const { gagal } = cekAdmin(headers, 'setting.update')
  if (gagal) return gagal
  return { data: { data: { ...publicSettings } } }
}

// Rasio poin hanya boleh diubah pemegang izin setting.point_ratio (usulan: Super Admin).
function ubahPengaturan({ headers, body }) {
  const { u, gagal } = cekAdmin(headers, 'setting.update')
  if (gagal) return gagal
  if (Object.hasOwn(body, 'wa_number')) return galatIsian({ wa_number: ['Nomor WhatsApp hanya bisa diubah melalui konfigurasi server.'] })
  const errors = {}
  if (kosong(body.address)) errors.address = ['Alamat wajib diisi']
  if (kosong(body.opening_hours)) errors.opening_hours = ['Jam buka wajib diisi']
  const peta = String(body.maps_url ?? '').trim()
  if (peta && !peta.startsWith('https://')) errors.maps_url = ['Link Google Maps harus diawali https://']
  const ubahRasio = body.point_ratio_rupiah != null && Number(body.point_ratio_rupiah) !== Number(publicSettings.point_ratio_rupiah)
  if (ubahRasio && !izinPengguna(u).includes('setting.point_ratio')) return ditolak('Rasio poin hanya bisa diubah Super Admin.')
  if (ubahRasio && !(Number.isInteger(Number(body.point_ratio_rupiah)) && Number(body.point_ratio_rupiah) >= 1000)) {
    errors.point_ratio_rupiah = ['Rasio poin minimal Rp 1.000 per poin']
  }
  if (Object.keys(errors).length) return galatIsian(errors)
  const sebelum = { ...publicSettings }
  Object.assign(publicSettings, {
    address: body.address.trim(), opening_hours: body.opening_hours.trim(), maps_url: peta || null,
    ...(ubahRasio && { point_ratio_rupiah: Number(body.point_ratio_rupiah) }),
  })
  catat(u, 'setting.update', 'Pengaturan toko', { sebelum, sesudah: { ...publicSettings } })
  simpanDb()
  return { data: { data: { ...publicSettings } } }
}

// Urutan penting: path tetap sebelum path ber-parameter.
export const ruteKonten = [
  ['get', '/banners', daftarBannerPublik],
  ['get', '/articles', daftarArtikelPublik],
  ['get', '/articles/:slug', detailArtikelPublik],
  ['get', '/services', daftarJasaPublik],
  ['get', '/services/:slug', detailJasaPublik],
  ['get', '/admin/banners', daftarBannerAdmin],
  ['post', '/admin/banners', buatBanner],
  ['get', '/admin/banners/:id', ambilBanner],
  ['put', '/admin/banners/:id', ubahBanner],
  ['patch', '/admin/banners/:id', patchBanner],
  ['get', '/admin/articles', daftarArtikelAdmin],
  ['post', '/admin/articles', buatArtikel],
  ['get', '/admin/articles/:id', ambilArtikelAdmin],
  ['put', '/admin/articles/:id', ubahArtikel],
  ['patch', '/admin/articles/:id', terbitkanArtikel],
  ['get', '/admin/services', daftarJasaAdmin],
  ['post', '/admin/services', buatJasa],
  ['get', '/admin/services/:id', ambilJasaAdmin],
  ['put', '/admin/services/:id', ubahJasa],
  ['patch', '/admin/services/:id', patchJasa],
  ['get', '/admin/settings', pengaturanAdmin],
  ['put', '/admin/settings', ubahPengaturan],
]

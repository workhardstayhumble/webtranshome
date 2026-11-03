// Tampilan log aktivitas: label aksi, nama kolom, format nilai, dan perbandingan before_data vs after_data.
import { formatAngka, formatRupiah, formatTanggalJam } from '@/utils/format'

const AKSI = {
  'product.create': 'Menambah produk',
  'product.update': 'Mengubah produk',
  'category.create': 'Menambah kategori',
  'category.update': 'Mengubah kategori',
  'brand.create': 'Menambah brand',
  'brand.update': 'Mengubah brand',
  'brand.delete': 'Menghapus brand',
  'room.create': 'Menambah ruangan',
  'room.update': 'Mengubah ruangan',
  'member.update': 'Mengubah status akun member',
  'point.create': 'Mencatat poin',
  'point.adjust': 'Penyesuaian poin',
  'banner.create': 'Membuat banner',
  'banner.update': 'Mengubah banner',
  'article.create': 'Menulis artikel',
  'service.create': 'Menambah jasa',
  'service.update': 'Mengubah jasa',
  'article.update': 'Mengubah artikel',
  'article.publish': 'Menerbitkan / membatalkan terbit artikel',
  'setting.update': 'Mengubah pengaturan toko',
  'user.create': 'Menambah admin',
  'user.update': 'Mengubah admin',
  'user.reset_password': 'Mengirim link reset password',
  'role.update': 'Mengubah izin role',
  'qontak.retry': 'Mencoba ulang sync Qontak',
  'qontak.retry_failed': 'Mencoba ulang semua sync yang gagal',
}

export const labelAksi = (kode) => AKSI[kode] ?? kode

const KOLOM = {
  sku: 'SKU', name: 'Nama', full_name: 'Nama', title: 'Judul', slug: 'Slug', email: 'Email', role: 'Role', roles: 'Role',
  permissions: 'Izin', price_general: 'Harga umum', show_price: 'Tampilkan harga', unit_sale: 'Satuan jual', min_order: 'Minimal pembelian',
  stock_status: 'Status stok', stock_qty_label: 'Label stok', is_active: 'Aktif', is_featured: 'Unggulan',
  category_id: 'ID kategori', brand_id: 'ID brand', parent_id: 'ID induk', description: 'Deskripsi', specifications: 'Spesifikasi',
  datasheet_pdf_url: 'Datasheet', meta_title: 'Judul SEO', meta_description: 'Deskripsi SEO', image_url: 'Gambar', logo_url: 'Logo',
  image_cover: 'Cover', sort_order: 'Urutan', member_code: 'Kode member', member: 'Member', type: 'Jenis', points: 'Poin',
  reference_no: 'Nomor nota', purchase_amount: 'Total belanja', note: 'Catatan', start_at: 'Mulai tayang', end_at: 'Selesai tayang',
  position: 'Posisi', link_url: 'Tautan', subtitle: 'Subjudul', video_url: 'Video', cta_label: 'Label tombol',
  cta2_label: 'Label tombol kedua', cta2_url: 'Tautan tombol kedua', summary: 'Ringkasan',
  service_area: 'Area layanan', category_ids: 'Kategori produk terkait', status: 'Status', published_at: 'Terbit', wa_number: 'Nomor WhatsApp', address: 'Alamat',
  opening_hours: 'Jam buka', maps_url: 'Link Google Maps', point_ratio_rupiah: 'Rasio poin', sync_status: 'Status sync', count: 'Jumlah kontak',
}

export const labelKolom = (kunci) => KOLOM[kunci] ?? kunci

const RUPIAH = ['price_general', 'purchase_amount', 'point_ratio_rupiah']

// Nilai berkode (enum di database) ditampilkan dengan label yang sama seperti di halaman lain.
const LABEL_NILAI = {
  sync_status: { PENDING: 'Menunggu', SYNCED: 'Tersinkron', FAILED: 'Gagal' },
  stock_status: { TERSEDIA: 'Tersedia', SISA_STOK: 'Sisa stok', PRE_ORDER: 'Pre-order', HABIS: 'Habis' },
  status: { DRAFT: 'Draf', PUBLISHED: 'Terbit' },
  type: { EARN: 'Poin belanja', REDEEM: 'Tukar poin', ADJUST: 'Penyesuaian', ARTIKEL: 'Artikel', PROMO: 'Promo', EVENT: 'Event', HALAMAN: 'Halaman' },
  position: { HOME_SLIDER: 'Slider Beranda', HOME_SIDE: 'Samping slider' },
}
const WAKTU = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}/

/** Nilai satu kolom sebagai teks: Rupiah, Ya/Tidak, tanggal, daftar, atau JSON untuk objek. */
export function formatNilai(kunci, nilai) {
  if (nilai == null || nilai === '') return '—'
  if (typeof nilai === 'boolean') return nilai ? 'Ya' : 'Tidak'
  if (LABEL_NILAI[kunci]?.[nilai]) return LABEL_NILAI[kunci][nilai]
  if (RUPIAH.includes(kunci) && !Number.isNaN(Number(nilai))) return formatRupiah(nilai)
  if (typeof nilai === 'number') return formatAngka(nilai)
  if (Array.isArray(nilai)) return nilai.length ? nilai.join(', ') : '—'
  if (typeof nilai === 'object') return JSON.stringify(nilai, null, 2)
  if (typeof nilai === 'string' && WAKTU.test(nilai)) return formatTanggalJam(nilai)
  return String(nilai)
}

const sama = (a, b) => JSON.stringify(a ?? null) === JSON.stringify(b ?? null)

/**
 * Baris perbandingan: [{ kunci, label, sebelum, sesudah, berubah, ditambah, dicabut }].
 * Daftar (mis. izin role) juga diberi `ditambah` / `dicabut` supaya perubahannya mudah dibaca.
 */
export function bandingkanData(sebelum, sesudah) {
  const a = sebelum ?? {}
  const b = sesudah ?? {}
  return [...new Set([...Object.keys(a), ...Object.keys(b)])].map((kunci) => {
    const lama = a[kunci]
    const baru = b[kunci]
    const baris = {
      kunci, label: labelKolom(kunci),
      sebelum: sebelum ? formatNilai(kunci, lama) : null,
      sesudah: sesudah ? formatNilai(kunci, baru) : null,
      berubah: Boolean(sebelum && sesudah) && !sama(lama, baru),
    }
    if (Array.isArray(lama) && Array.isArray(baru)) {
      baris.ditambah = baru.filter((x) => !lama.includes(x))
      baris.dicabut = lama.filter((x) => !baru.includes(x))
    }
    return baris
  })
}

/** Data mana yang diubah, untuk kolom ringkas di tabel log (mis. "KRM-4040-01 · Keramik lantai…"). */
export function subjekLog(log) {
  const d = log.after_data ?? log.before_data ?? {}
  const bagian = [d.sku ?? d.member_code, d.name ?? d.full_name ?? d.title ?? d.role ?? d.member ?? d.email].filter(Boolean)
  if (!bagian.length && d.count != null) return `${d.count} kontak`
  return bagian.join(' · ')
}

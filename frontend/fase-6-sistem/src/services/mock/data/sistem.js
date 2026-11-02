// Data contoh tabel roles dan permissions (docs/transhome_postgres.sql). TODO: ganti ke API (seed menunggu backend).
// Pemetaan role → izin (role_permissions) ada di data/akun.js (`rolePermissions`) dan ikut tersimpan di mockDb.

export const roles = [
  { id: 1, code: 'SUPER_ADMIN', name: 'Super Admin', description: 'Semua akses, termasuk pengguna admin, role, dan rasio poin.', is_active: true },
  { id: 2, code: 'ADMIN_KATALOG', name: 'Admin Katalog', description: 'Produk, kategori, brand, dan ruangan.', is_active: true },
  { id: 3, code: 'ADMIN_MEMBERSHIP', name: 'Admin Membership', description: 'Member, input poin dari nota, dan sync Qontak.', is_active: true },
  { id: 4, code: 'ADMIN_KONTEN', name: 'Admin Konten', description: 'Banner, artikel dan promo, serta pengaturan toko.', is_active: true },
  { id: 5, code: 'CUSTOMER', name: 'Member', description: 'Pelanggan Trans Family. Tidak punya akses admin.', is_active: true },
]

// [menu_id (data/admin.js), kode, aksi, keterangan]. Muka 1 umumnya tidak menghapus data (cukup dinonaktifkan);
// satu-satunya izin `delete` adalah brand.delete (8 Okt 2026). Aksi khusus: publish, adjust, point_ratio, retry.
const DAFTAR = [
  [11, 'product.view', 'view', 'Melihat daftar dan detail produk'],
  [11, 'product.create', 'create', 'Menambah produk'],
  [11, 'product.update', 'update', 'Mengubah produk, foto, harga, dan status'],
  [12, 'category.view', 'view', 'Melihat kategori'],
  [12, 'category.create', 'create', 'Menambah kategori'],
  [12, 'category.update', 'update', 'Mengubah kategori dan urutannya'],
  [13, 'brand.view', 'view', 'Melihat brand'],
  [13, 'brand.create', 'create', 'Menambah brand'],
  [13, 'brand.update', 'update', 'Mengubah brand'],
  [13, 'brand.delete', 'delete', 'Menghapus brand nonaktif yang tidak dipakai produk'],
  [14, 'room.view', 'view', 'Melihat ruangan'],
  [14, 'room.create', 'create', 'Menambah ruangan'],
  [14, 'room.update', 'update', 'Mengubah ruangan dan urutannya'],
  [21, 'member.view', 'view', 'Melihat daftar dan detail member'],
  [21, 'member.update', 'update', 'Menonaktifkan atau mengaktifkan akun member'],
  [22, 'point.create', 'create', 'Mencatat poin belanja dan penukaran poin'],
  [22, 'point.adjust', 'adjust', 'Penyesuaian poin (tambah atau kurangi manual)'],
  [31, 'banner.view', 'view', 'Melihat banner'],
  [31, 'banner.create', 'create', 'Membuat banner'],
  [31, 'banner.update', 'update', 'Mengubah, mengurutkan, dan menonaktifkan banner'],
  [32, 'article.view', 'view', 'Melihat artikel dan promo'],
  [32, 'article.create', 'create', 'Menulis artikel baru'],
  [32, 'article.update', 'update', 'Mengubah artikel'],
  [32, 'article.publish', 'publish', 'Menerbitkan dan membatalkan terbit'],
  [33, 'setting.update', 'update', 'Mengubah alamat, jam buka, dan peta; nomor WhatsApp hanya-baca'],
  [33, 'setting.point_ratio', 'point_ratio', 'Mengubah rasio poin'],
  [41, 'user.view', 'view', 'Melihat pengguna admin'],
  [41, 'user.create', 'create', 'Menambah admin'],
  [41, 'user.update', 'update', 'Mengubah role, reset password, menonaktifkan admin'],
  [42, 'role.view', 'view', 'Melihat role dan izinnya'],
  [42, 'role.update', 'update', 'Mengubah izin role'],
  [43, 'log.view', 'view', 'Melihat log aktivitas semua admin'],
  [44, 'qontak.retry', 'retry', 'Melihat status sync dan mencoba ulang'],
]

export const permissions = DAFTAR.map(([menu_id, code, action, description], i) => ({ id: i + 1, menu_id, code, action, description }))

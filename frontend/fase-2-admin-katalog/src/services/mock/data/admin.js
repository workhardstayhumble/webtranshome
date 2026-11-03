// Data contoh panel admin: menu sidebar, log aktivitas, dan status sync Qontak.
// Kolom mengikuti docs/transhome_postgres.sql. TODO: ganti ke API (seed menus & permissions menunggu backend).
import { MEMBER_TAMBAHAN } from '@/services/mock/data/akun'

// Tabel menus. `izin` = kode permission yang membuat menu terlihat (permissions.menu_id → menus.id);
// tidak ikut dikirim ke frontend. Grup (route null) tampil bila minimal satu anaknya terlihat.
export const menus = [
  { id: 1, parent_id: null, name: 'Utama', code: 'utama', route: null, icon: null, sort_order: 1, is_active: true },
  { id: 2, parent_id: 1, name: 'Dashboard', code: 'dashboard', route: '/admin', icon: 'layout-dashboard', sort_order: 1, is_active: true, izin: null },
  { id: 10, parent_id: null, name: 'Katalog', code: 'katalog', route: null, icon: null, sort_order: 2, is_active: true },
  { id: 11, parent_id: 10, name: 'Produk', code: 'product', route: '/admin/produk', icon: 'package', sort_order: 1, is_active: true, izin: 'product.view' },
  { id: 12, parent_id: 10, name: 'Kategori', code: 'category', route: '/admin/kategori', icon: 'folder-tree', sort_order: 2, is_active: true, izin: 'category.view' },
  { id: 13, parent_id: 10, name: 'Brand', code: 'brand', route: '/admin/brand', icon: 'tag', sort_order: 3, is_active: true, izin: 'brand.view' },
  { id: 14, parent_id: 10, name: 'Ruangan', code: 'room', route: '/admin/ruangan', icon: 'sofa', sort_order: 4, is_active: true, izin: 'room.view' },
  { id: 20, parent_id: null, name: 'Membership', code: 'membership', route: null, icon: null, sort_order: 3, is_active: true },
  { id: 21, parent_id: 20, name: 'Member', code: 'member', route: '/admin/member', icon: 'users', sort_order: 1, is_active: true, izin: 'member.view' },
  { id: 22, parent_id: 20, name: 'Input poin', code: 'point', route: '/admin/poin/input', icon: 'coins', sort_order: 2, is_active: true, izin: 'point.create' },
  { id: 30, parent_id: null, name: 'Konten', code: 'konten', route: null, icon: null, sort_order: 4, is_active: true },
  { id: 31, parent_id: 30, name: 'Banner', code: 'banner', route: '/admin/banner', icon: 'image', sort_order: 1, is_active: true, izin: 'banner.view' },
  { id: 32, parent_id: 30, name: 'Artikel dan promo', code: 'article', route: '/admin/artikel', icon: 'newspaper', sort_order: 2, is_active: true, izin: 'article.view' },
  { id: 34, parent_id: 30, name: 'Jasa', code: 'service', route: '/admin/jasa', icon: 'wrench', sort_order: 3, is_active: true, izin: 'service.view' },
  { id: 33, parent_id: 30, name: 'Pengaturan toko', code: 'setting', route: '/admin/pengaturan', icon: 'settings', sort_order: 4, is_active: true, izin: 'setting.update' },
  { id: 40, parent_id: null, name: 'Sistem', code: 'sistem', route: null, icon: null, sort_order: 5, is_active: true },
  { id: 41, parent_id: 40, name: 'Pengguna admin', code: 'user', route: '/admin/pengguna', icon: 'user-cog', sort_order: 1, is_active: true, izin: 'user.view' },
  { id: 42, parent_id: 40, name: 'Role dan izin', code: 'role', route: '/admin/role', icon: 'shield', sort_order: 2, is_active: true, izin: 'role.view' },
  { id: 43, parent_id: 40, name: 'Log aktivitas', code: 'log', route: '/admin/log-aktivitas', icon: 'scroll-text', sort_order: 3, is_active: true, izin: 'log.view' },
  { id: 44, parent_id: 40, name: 'Sync Qontak', code: 'qontak', route: '/admin/sync-qontak', icon: 'refresh-cw', sort_order: 4, is_active: true, izin: 'qontak.retry' },
]

// Tabel activity_logs (terbaru di depan). admin_id merujuk users di data/akun.js.
export const activityLogs = [
  { id: 6, admin_id: 2, action: 'product.update', module: 'Produk', before_data: { sku: 'KRM-4040-01', name: 'Keramik lantai abu doff 40×40', price_general: '87000.00' }, after_data: { sku: 'KRM-4040-01', name: 'Keramik lantai abu doff 40×40', price_general: '89500.00' }, ip_address: '10.0.0.12', created_at: '2026-10-06T16:40:00' },
  { id: 5, admin_id: 3, action: 'point.create', module: 'Poin', before_data: null, after_data: { member_code: 'TF-000123', type: 'EARN', points: 4, reference_no: 'NT-2026-10-0456', purchase_amount: '4250000.00', note: null }, ip_address: '10.0.0.15', created_at: '2026-10-05T14:20:00' },
  { id: 4, admin_id: 4, action: 'banner.create', module: 'Banner', before_data: null, after_data: { title: 'Paket atap baja ringan', start_at: '2026-09-27T00:00:00', end_at: '2026-11-06T23:59:59' }, ip_address: '10.0.0.18', created_at: '2026-10-04T10:05:00' },
  { id: 3, admin_id: 2, action: 'product.create', module: 'Produk', before_data: null, after_data: { sku: 'SMN-PCC-50', name: 'Semen PCC 50 kg', price_general: '68000.00', stock_status: 'TERSEDIA', is_active: true }, ip_address: '10.0.0.12', created_at: '2026-10-03T09:30:00' },
  { id: 2, admin_id: 1, action: 'role.update', module: 'Role dan izin', before_data: { role: 'Admin Membership', permissions: ['member.view', 'member.update', 'point.create'] }, after_data: { role: 'Admin Membership', permissions: ['member.view', 'member.update', 'point.create', 'qontak.retry'] }, ip_address: '10.0.0.10', created_at: '2026-10-01T15:00:00' },
  { id: 1, admin_id: 3, action: 'point.create', module: 'Poin', before_data: null, after_data: { member_code: 'TF-000124', type: 'EARN', points: 3, reference_no: 'NT-2026-10-0431', purchase_amount: '3600000.00', note: null }, ip_address: '10.0.0.15', created_at: '2026-09-30T11:45:00' },
]

// Tabel qontak_contacts: satu baris per member. Member 102 & 103 gagal tersinkron (untuk tombol "Coba lagi").
export const qontakContacts = [
  { id: 1, user_id: 101, qontak_contact_id: 'QC-0001', sync_status: 'SYNCED', last_synced_at: '2026-08-12T10:01:00', created_at: '2026-08-12T10:00:00', updated_at: '2026-08-12T10:01:00' },
  { id: 2, user_id: 102, qontak_contact_id: null, sync_status: 'FAILED', last_synced_at: null, created_at: '2026-10-05T19:12:00', updated_at: '2026-10-05T19:30:00' },
  { id: 3, user_id: 103, qontak_contact_id: null, sync_status: 'FAILED', last_synced_at: null, created_at: '2026-10-06T08:40:00', updated_at: '2026-10-06T09:10:00' },
]

// Tabel sync_logs: percobaan terakhir per member.
export const syncLogs = [
  { id: 1, user_id: 101, event_type: 'CONTACT_UPSERT', status_code: 200, retry_count: 0, created_at: '2026-08-12T10:01:00' },
  { id: 2, user_id: 102, event_type: 'CONTACT_UPSERT', status_code: 500, retry_count: 3, created_at: '2026-10-05T19:30:00' },
  { id: 3, user_id: 103, event_type: 'CONTACT_UPSERT', status_code: 429, retry_count: 2, created_at: '2026-10-06T09:10:00' },
]

// Kontak Qontak member contoh tambahan (data/akun.js): sebagian besar tersinkron, satu gagal, dua masih menunggu.
MEMBER_TAMBAHAN.forEach(([, , , , dibuat, status], i) => {
  const userId = 104 + i
  const selesai = status === 'PENDING' ? null : dibuat.replace(/:00$/, ':40')
  qontakContacts.push({
    id: qontakContacts.length + 1, user_id: userId, qontak_contact_id: status === 'SYNCED' ? `QC-${String(userId).padStart(4, '0')}` : null,
    sync_status: status, last_synced_at: status === 'SYNCED' ? selesai : null, created_at: dibuat, updated_at: selesai ?? dibuat,
  })
  if (selesai) {
    syncLogs.push({ id: syncLogs.length + 1, user_id: userId, event_type: 'CONTACT_UPSERT', status_code: status === 'SYNCED' ? 200 : 503, retry_count: status === 'SYNCED' ? 0 : 3, created_at: selesai })
  }
})

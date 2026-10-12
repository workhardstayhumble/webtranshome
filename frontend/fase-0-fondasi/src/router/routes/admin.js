import AdminLayout from '@/layouts/AdminLayout.vue'

const ProdukForm = () => import('@/views/admin/katalog/ProdukForm.vue')
const BannerForm = () => import('@/views/admin/konten/BannerForm.vue')
const ArtikelForm = () => import('@/views/admin/konten/ArtikelForm.vue')

// Area admin: role admin, dicek per permission (CLAUDE.md §4). /admin/masuk ada di routes/auth.js.
// `grup` tampil kecil di atas judul header (mis. Katalog › Produk).
export default [
  {
    path: '/admin',
    component: AdminLayout,
    meta: { area: 'admin' },
    // Child '' wajib ada: tanpa itu /admin hanya mencocokkan induk dan RouterView kosong.
    children: [
      { path: '', name: 'admin-dashboard', component: () => import('@/views/admin/Dashboard.vue'), meta: { title: 'Dashboard' } },
      {
        path: 'produk',
        name: 'admin-produk',
        component: () => import('@/views/admin/katalog/ProdukDaftar.vue'),
        meta: { title: 'Produk', grup: 'Katalog', permission: 'catalog.product.read' },
      },
      { path: 'produk/baru', name: 'admin-produk-baru', component: ProdukForm, meta: { title: 'Tambah produk', grup: 'Katalog', permission: 'catalog.product.create' } },
      { path: 'produk/:id(\\d+)', name: 'admin-produk-ubah', component: ProdukForm, meta: { title: 'Ubah produk', grup: 'Katalog', permission: 'catalog.product.update' } },
      {
        path: 'produk/:id(\\d+)/pratinjau',
        name: 'admin-produk-pratinjau',
        component: () => import('@/views/admin/katalog/ProdukPratinjau.vue'),
        meta: { title: 'Pratinjau produk', grup: 'Katalog', permission: 'catalog.product.read' },
      },
      {
        path: 'kategori',
        name: 'admin-kategori',
        component: () => import('@/views/admin/katalog/Kategori.vue'),
        meta: { title: 'Kategori', grup: 'Katalog', permission: 'catalog.category.manage' },
      },
      {
        path: 'brand',
        name: 'admin-brand',
        component: () => import('@/views/admin/katalog/Brand.vue'),
        meta: { title: 'Brand', grup: 'Katalog', permission: 'catalog.brand.manage' },
      },
      {
        path: 'ruangan',
        name: 'admin-ruangan',
        component: () => import('@/views/admin/katalog/Ruangan.vue'),
        meta: { title: 'Ruangan', grup: 'Katalog', permission: 'catalog.room.manage' },
      },
      {
        path: 'member',
        name: 'admin-member',
        component: () => import('@/views/admin/membership/MemberDaftar.vue'),
        meta: { title: 'Member', grup: 'Membership', permission: 'member.view' },
      },
      {
        path: 'member/:id(\\d+)',
        name: 'admin-member-detail',
        component: () => import('@/views/admin/membership/MemberDetail.vue'),
        meta: { title: 'Detail member', grup: 'Membership', permission: 'member.view' },
      },
      {
        path: 'poin/input',
        name: 'admin-poin-input',
        component: () => import('@/views/admin/membership/PoinInput.vue'),
        meta: { title: 'Input poin', grup: 'Membership', permission: 'point.create' },
      },
      {
        path: 'banner',
        name: 'admin-banner',
        component: () => import('@/views/admin/konten/BannerDaftar.vue'),
        meta: { title: 'Banner', grup: 'Konten', permission: 'banner.view' },
      },
      { path: 'banner/baru', name: 'admin-banner-baru', component: BannerForm, meta: { title: 'Buat banner', grup: 'Konten', permission: 'banner.create' } },
      { path: 'banner/:id(\\d+)', name: 'admin-banner-ubah', component: BannerForm, meta: { title: 'Ubah banner', grup: 'Konten', permission: 'banner.update' } },
      {
        path: 'artikel',
        name: 'admin-artikel',
        component: () => import('@/views/admin/konten/ArtikelDaftar.vue'),
        meta: { title: 'Artikel dan promo', grup: 'Konten', permission: 'article.view' },
      },
      { path: 'artikel/baru', name: 'admin-artikel-baru', component: ArtikelForm, meta: { title: 'Tulis artikel', grup: 'Konten', permission: 'article.create' } },
      { path: 'artikel/:id(\\d+)', name: 'admin-artikel-ubah', component: ArtikelForm, meta: { title: 'Ubah artikel', grup: 'Konten', permission: 'article.update' } },
      {
        path: 'pengaturan',
        name: 'admin-pengaturan',
        component: () => import('@/views/admin/konten/Pengaturan.vue'),
        meta: { title: 'Pengaturan toko', grup: 'Konten', permission: 'setting.update' },
      },
      {
        path: 'pengguna',
        name: 'admin-pengguna',
        component: () => import('@/views/admin/akses/PenggunaDaftar.vue'),
        meta: { title: 'Pengguna admin', grup: 'Sistem', permission: 'user.view' },
      },
      {
        path: 'role',
        name: 'admin-role',
        component: () => import('@/views/admin/akses/Role.vue'),
        meta: { title: 'Role dan izin', grup: 'Sistem', permission: 'role.view' },
      },
      {
        path: 'log-aktivitas',
        name: 'admin-log',
        component: () => import('@/views/admin/audit/LogAktivitas.vue'),
        meta: { title: 'Log aktivitas', grup: 'Sistem', permission: 'log.view' },
      },
      {
        path: 'sync-qontak',
        name: 'admin-sync-qontak',
        component: () => import('@/views/admin/integrasi/SyncQontak.vue'),
        meta: { title: 'Sync Qontak', grup: 'Sistem', permission: 'qontak.retry' },
      },
      {
        path: ':bagian(.*)',
        name: 'admin-tidak-ditemukan',
        component: () => import('@/views/admin/TidakDitemukan.vue'),
        meta: { title: 'Halaman tidak ditemukan' },
      },
    ],
  },
]
import PublicLayout from '@/layouts/PublicLayout.vue'

// Galeri komponen hanya ada di mode development (npm run dev), tidak ikut build produksi.
// Kotak email & simulasi Google hanya bila juga memakai mock (email dan OAuth asli dikerjakan backend).
const modeMock = import.meta.env.VITE_USE_MOCK === 'true'
const halamanDev = import.meta.env.DEV
  ? [
      {
        path: 'dev/components',
        name: 'dev-components',
        component: () => import('@/views/dev/Components.vue'),
        meta: { title: 'Galeri komponen' },
      },
      ...(modeMock
        ? [
            { path: 'dev/email', name: 'dev-email', component: () => import('@/views/dev/KotakEmail.vue'), meta: { title: 'Kotak email tiruan' } },
            { path: 'dev/google', name: 'dev-google', component: () => import('@/views/dev/SimulasiGoogle.vue'), meta: { title: 'Simulasi Google' } },
          ]
        : []),
    ]
  : []

// Area publik: semua pengunjung (CLAUDE.md §3).
export default [
  {
    path: '/',
    component: PublicLayout,
    meta: { area: 'public' },
    children: [
      {
        path: '',
        name: 'beranda',
        component: () => import('@/views/public/Beranda.vue'),
        meta: { title: 'Beranda' },
      },
      {
        path: 'artikel',
        name: 'artikel',
        component: () => import('@/views/public/ArtikelDaftar.vue'),
        meta: { title: 'Artikel' },
      },
      {
        path: 'artikel/:slug',
        name: 'artikel-detail',
        component: () => import('@/views/public/ArtikelDetail.vue'),
        meta: { title: 'Artikel' },
      },
      {
        path: 'jasa',
        name: 'jasa',
        component: () => import('@/views/public/JasaDaftar.vue'),
        meta: { title: 'Jasa' },
      },
      {
        path: 'jasa/:slug',
        name: 'jasa-detail',
        component: () => import('@/views/public/JasaDetail.vue'),
        meta: { title: 'Jasa' },
      },
      {
        path: 'promo',
        name: 'promo',
        component: () => import('@/views/public/Promo.vue'),
        meta: { title: 'Promo dan event' },
      },
      {
        path: 'info-toko',
        name: 'info-toko',
        component: () => import('@/views/public/InfoToko.vue'),
        meta: { title: 'Info toko' },
      },
      {
        path: 'katalog/:kategori?',
        name: 'katalog',
        component: () => import('@/views/public/Katalog.vue'),
        meta: { title: 'Katalog' },
      },
      {
        path: 'cari',
        name: 'cari',
        component: () => import('@/views/public/Cari.vue'),
        meta: { title: 'Cari produk' },
      },
      {
        path: 'produk/:slug',
        name: 'produk',
        component: () => import('@/views/public/ProdukDetail.vue'),
        // Di HP halaman ini punya bar harga + WhatsApp sendiri di bawah: tombol WA melayang hanya di desktop,
        // dan PublicLayout memberi ruang di bawah footer agar tidak tertutup bar.
        meta: { title: 'Detail produk', barAksiHp: true },
      },
      {
        path: 'brand',
        name: 'brand-daftar',
        component: () => import('@/views/public/BrandDaftar.vue'),
        meta: { title: 'Semua brand' },
      },
      {
        path: 'brand/:slug',
        name: 'brand',
        component: () => import('@/views/public/Brand.vue'),
        meta: { title: 'Brand' },
      },
      {
        path: 'ruangan',
        name: 'ruangan',
        component: () => import('@/views/public/Ruangan.vue'),
        meta: { title: 'Belanja per ruangan' },
      },
      {
        path: 'ruangan/:slug',
        name: 'ruangan-detail',
        component: () => import('@/views/public/RuanganDetail.vue'),
        meta: { title: 'Ruangan' },
      },
      ...halamanDev,
      {
        path: 'tidak-punya-akses',
        name: 'tidak-punya-akses',
        component: () => import('@/views/public/TidakPunyaAkses.vue'),
        meta: { title: 'Tidak punya akses' },
      },
      {
        path: ':pathMatch(.*)*',
        name: 'tidak-ditemukan',
        component: () => import('@/views/public/TidakDitemukan.vue'),
        meta: { title: 'Halaman tidak ditemukan' },
      },
    ],
  },
]

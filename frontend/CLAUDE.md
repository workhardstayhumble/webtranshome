# Project Rules — Website Transhome Muka 1 (Frontend)

Website katalog digital + membership "Trans Family" untuk CV Trans Home, supermarket bahan bangunan & elektronik, toko tunggal di Sleman, Yogyakarta.

Sumber kebenaran: `docs/KEPUTUSAN.md` (keputusan terbaru, MENGALAHKAN BRD/PRD), lalu BRD v0.4 & PRD v0.4 (`docs/*.pdf`).
Kalau rules ini bertentangan dengan dokumen itu, BERHENTI dan tanya dulu.
Rencana kerja per fase: `docs/FRONTEND_BUILD_PLAN.md`. Bentuk data untuk mock: `docs/transhome_postgres.sql`.
Acuan visual: sketsa "Sketsa Dashboard Transhome" (https://claude.ai/artifact/S8vhfiAUswfBB85neKS1mZ).

## 1. Stack
- Vue 3 (Composition API, `<script setup>`), Vite, Vue Router, Pinia. JavaScript; jangan migrasi ke TypeScript tanpa diminta.
- Tailwind CSS v4 lewat `@tailwindcss/vite`. JANGAN tambah library komponen UI (Nuxt UI, shadcn, Vuetify, PrimeVue, Element, Bootstrap, dll) tanpa izin. Komponen dasar dibuat sendiri di `components/ui/`.
- Ikon: `lucide-vue-next`. Emoji DILARANG dipakai sebagai ikon.
- HTTP: satu instance axios di `services/http.js` (baseURL dari `VITE_API_BASE_URL`, interceptor token, 401, 403). Jangan hardcode URL.
- Backend: Go (DDD) + PostgreSQL, dikerjakan tim lain. Frontend HANYA konsumsi REST API, tanpa logika bisnis penting di browser. JANGAN mengubah folder `backend/`.
- SPA tanpa SSR. `.env` hanya berisi variabel `VITE_*` dan semuanya terbaca di browser, jadi jangan simpan rahasia di sana.

## 2. Struktur Folder
Kode dipisah per fase. Setiap folder fase berisi `src/` dengan susunan yang sama; pengaturan proyek tetap di `frontend/`.
```
frontend/
  index.html, vite.config.js, jsconfig.json, package.json, .env.example, docs/, design-system/, CLAUDE.md, README.md
  fase-0-fondasi/src/            main.js, App.vue, assets/, layouts/, router/, stores/, components/ui/, components/layout/,
                                 services/ (http, errors, mock dasar), composables/, utils/, views/dev/, halaman sistem
  fase-1-katalog/src/            views/public/ (katalog, cari, produk, brand, ruangan), components/katalog/, components/produk/
  fase-2-admin-katalog/src/      views/admin/ (masuk, dashboard, katalog/), components/admin/, services admin, mock admin
  fase-3-akun-member/src/        views/auth/, views/member/, components/akun/, services akun, mock akun
  fase-4-membership-admin/src/   views/admin/membership/, components/admin/member/, components/admin/poin/, mock membership
  fase-5-konten/src/             views/public/ (Beranda, artikel, promo, info toko), views/admin/konten/, components/beranda/, mock konten
  fase-6-sistem/src/             views/admin/akses/, audit/, integrasi/ (pengguna, role, log, sync Qontak), mock sistem
```
- Import antarfile memakai `@/...` (mis. `@/components/ui/Button.vue`). Alias `@` mencari di `src/` semua folder fase secara berurutan (`vite.config.js`, `jsconfig.json`). Satu path di bawah `src/` hanya boleh ada di SATU folder fase.
- Import relatif (`./`, `../`) hanya untuk file di folder fase yang sama.
- File BARU ditaruh di folder fase yang sedang dikerjakan. File lama diubah di tempatnya (folder fase pembuatnya), walau perubahannya untuk fase lain. Daftar per fase: `docs/FILE_PER_FASE.md`.

Susunan di dalam setiap `src/` (gabungan semua fase):
```
src/
  assets/main.css   Tailwind + token desain (@theme)
  layouts/          PublicLayout, AuthLayout, MemberLayout, AdminLayout
  views/public/     Home, Katalog, ProdukDetail, Brand, Ruangan, RuanganDetail, Artikel, ArtikelDetail,
                    Promo, InfoToko, Cari, TidakDitemukan, TidakPunyaAkses
  views/auth/       Masuk, Daftar, Verifikasi, LupaPassword, ResetPassword
  views/member/     Ringkasan, RiwayatPoin, Profil, Keamanan
  views/admin/      Masuk, Dashboard, katalog/, membership/, konten/, akses/, audit/, integrasi/
  views/dev/        Components (hanya mode development)
  components/ui/    komponen dasar (Button, Input, Select, Checkbox, Switch, Badge, StockBadge, PriceTag,
                    CopyButton, Modal, ConfirmModal, Toast, EmptyState, ErrorState, Skeleton, Pagination, Table)
  components/<fitur>/  layout/, katalog/, produk/, akun/, admin/
  composables/      useAuth, usePermission (+ directive v-can), useWhatsApp, useFilter, useSeo
  stores/           auth, menu, ui, settings (Pinia)
  services/         http.js + satu service per modul (productService.js, dst) + mock/
  router/           routes/ per area + guard.js
  utils/            format.js (Rupiah, tanggal, nomor HP), whatsapp.js
```
- Satu komponen = satu tanggung jawab. Pecah jika > ~200 baris.
- `defineProps` / `defineEmits` wajib jelas: props pakai objek `{ type, default }` (JavaScript).
- Data diambil lewat service di `services/`; setiap pemanggilan menangani state loading, kosong, dan error.
- Endpoint yang belum ada di backend: service memakai data dari `services/mock/` (bentuk sesuai `docs/transhome_postgres.sql`) saat `VITE_USE_MOCK=true`, ditandai `// TODO: ganti ke API`.

## 3. Peta Halaman
```
Publik : /  /katalog  /katalog/:kategori  /produk/:slug  /brand  /brand/:slug  /ruangan  /ruangan/:slug
         /artikel  /artikel/:slug  /promo  /info-toko  /cari?q=  /jasa  /jasa/:slug
Auth   : /masuk  /daftar  /verifikasi  /lupa-password  /reset-password
Member : /akun  /akun/poin  /akun/profil  /akun/keamanan
Admin  : /admin/masuk  /admin  /admin/produk  /admin/produk/baru  /admin/produk/:id
         /admin/kategori  /admin/brand  /admin/ruangan  /admin/member  /admin/member/:id
         /admin/poin/input  /admin/banner  /admin/artikel  /admin/jasa  /admin/jasa/baru  /admin/jasa/:id  /admin/pengaturan
         /admin/pengguna  /admin/role  /admin/log-aktivitas  /admin/sync-qontak
Sistem : /tidak-punya-akses  halaman 404   ·   Dev: /dev/components (hanya development)
```
- Route guard membaca `meta.area`: `public`, `guest` (masuk/daftar; yang sudah login diarahkan ke `/akun` atau `/admin`), `member`, `admin` (+ `meta.permission`).
- Belum login: area member → `/masuk?redirect=…`, area admin → `/admin/masuk?redirect=…`.
- Login admin TERPISAH di `/admin/masuk`. Pelanggan tidak pernah melihat halaman login admin.
- Login & daftar OPSIONAL. Seluruh katalog bisa dilihat tanpa login.

## 4. RBAC Admin Panel
- Role: SUPER_ADMIN, ADMIN_KATALOG, ADMIN_MEMBERSHIP, ADMIN_KONTEN, CUSTOMER.
- Sidebar admin dibangun DINAMIS dari API menus yang boleh dilihat user. JANGAN hardcode menu per role.
- Kode permission format `<modul>.<aksi>`, contoh `product.update`. Aksi: view, create, update, delete (+ aksi khusus di rencana build, mis. `article.publish`).
- Cek di UI pakai `usePermission().can('product.update')` / directive `v-can`. Tombol tanpa izin disembunyikan.
- UI hanya menyembunyikan. Backend tetap wajib menolak (403). Navigasi tanpa izin → halaman "Tidak punya akses"; aksi ditolak → Toast.
- Proteksi halaman: `meta: { area: 'admin', permission: 'product.view' }`.

## 5. Fitur Muka 1 (scope FINAL) & Urutan Kerja
**1a:** Master Katalog (FR-A), Smart Filter & Search termasuk filter BRAND (FR-B), Detail Produk (FR-C), Indikator Stok (FR-D), Auth Member (FR-E), Katalog Ruangan (FR-F), Admin Panel RBAC + Audit Log (FR-G).
**1b:** Portal Member: member_code, saldo & riwayat poin (FR-H), Input Poin dari nota (FR-I), CMS Artikel (FR-N), CMS Banner (FR-O), Info Toko (FR-P).
**1c:** Widget WhatsApp Mekari Qontak (FR-R), Context Sync identitas ke widget (FR-S), Sync kontak Qontak (FR-T), Log & retry integrasi (FR-U).

Urutan kerja mengikuti `docs/FRONTEND_BUILD_PLAN.md` (Fase 0–6). Kerjakan SATU fase per permintaan, jalankan checklist QA fase itu, lalu berhenti. Jangan mulai fase berikutnya sebelum diminta.

## 6. DI LUAR SCOPE — JANGAN DIBUAT
Keranjang, checkout, payment gateway, ongkir, kalkulator material, booking konsultasi desain, FORMULIR konsultasi proyek, katalog hadiah, halaman penukaran poin oleh member, wishlist, pengingat stok, ulasan produk, rincian produk per nota, kartu member fisik, harga bertingkat, multi-cabang, upgrade tier.
- Muka 1 = SATU harga umum (`price_general`). Tombol aksi utama produk = "Tanya via WhatsApp", BUKAN "Beli" / "Tambah ke Keranjang".
- Konsultasi proyek & pertanyaan produk = HANYA lewat WhatsApp (Mekari Qontak). Tidak ada form lead.
- Penukaran poin dicatat ADMIN sebagai transaksi REDEEM, bukan dipilih member di website.

## 7. Aturan Bisnis di UI
- **Akun:** EMAIL WAJIB, nomor HP OPSIONAL. Login pakai email + password, atau tombol terpisah "Masuk dengan Google".
- **Daftar:** nama, email, password, HP (opsional), checkbox persetujuan komunikasi (tidak dicentang default). Semua member = Customer Tahap 1.
- **member_code** DITAMPILKAN jelas di `/akun` (besar, bisa disalin) karena disebut ke kasir saat belanja.
- **Poin:** saldo & riwayat dari API. Rasio dari `site_settings.point_ratio_rupiah` (contoh: 1 poin per Rp 1.000.000), JANGAN hardcode.
- **Input poin (Admin Membership):** cari member via member_code / email / HP, input `reference_no` (nomor nota) + total belanja. Nomor nota duplikat ditolak dengan pesan jelas. REDEEM ditolak jika saldo kurang.
- **Status stok** (badge + label teks + ikon, bukan warna saja): TERSEDIA, SISA STOK, PRE-ORDER, HABIS.
- **Harga boleh disembunyikan per produk** (`show_price`, KEPUTUSAN #15): tampil "Harga via WhatsApp" + tombol "Tanya harga"; harga boleh kosong bila disembunyikan. Pakai `hargaDisembunyikan()` / `labelTanyaProduk()` di `utils/whatsapp.js` dan `PriceTag` (prop `show`), jangan cek manual.
- **Tanya via WA:** pesan otomatis berisi SKU + nama produk + link produk; produk HABIS → "Tanya ketersediaan". Nomor WA dari pengaturan toko (API), jangan hardcode.
- **Lupa password:** lewat email. Pesan selalu sama ("Jika email terdaftar, link reset sudah dikirim").
- **Produk:** tidak dihapus permanen, cukup toggle aktif/nonaktif. Minimal 1 foto, 1 foto utama.
- **Brand:** semua brand aktif tampil (Beranda, `/brand`, filter katalog). Hapus permanen hanya untuk brand NONAKTIF yang tidak dipakai produk mana pun, izin `brand.delete` (keputusan pemilik proyek 8 Okt 2026). Brand yang masih dipakai: pindahkan produknya dulu.
- **Detail produk HANYA menampilkan data untuk pelanggan.** DILARANG menampilkan HPP, hnet, diskon (disk1/disk2/proyek), stok per lokasi/gudang, konsinyasi, karantina.
- **Banner:** tampil hanya dalam rentang tanggal & aktif. Admin melihat status "Terjadwal / Tayang / Berakhir". Posisi Hero Beranda (`HOME_HERO`) boleh berisi video; gambar tetap wajib sebagai cadangan (KEPUTUSAN #16).
- **Artikel:** DRAFT → Preview → PUBLISHED.
- **Sync Qontak:** proses background. UI tidak menunggu Qontak. Admin melihat status PENDING / SYNCED / FAILED + tombol retry. Member tanpa HP tetap tersinkron dengan email saja.

## 8. Design System (UI UX Pro Max)
- Token desain diturunkan dari sketsa (lihat atas) di Fase 0 dan dicatat di `design-system/transhome/MASTER.md`.
- Sebelum membuat/mengubah halaman, baca `design-system/transhome/MASTER.md`, lalu `design-system/transhome/pages/<halaman>.md` jika ada (didahulukan).
- Pakai skill ui-ux-pro-max dengan `--stack vue`.
- Token warna/font/radius/shadow didefinisikan SEKALI di `@theme` pada `fase-0-fondasi/src/assets/main.css`. Pakai token semantik: primary, secondary, success, warning, danger, surface, muted, dst.
- DILARANG hardcode hex atau warna arbitrary (`bg-[#1e40af]`) di komponen. Spacing pakai skala Tailwind (grid 8pt).
- **Situs publik:** bersih, terpercaya, foto produk jadi fokus, mobile-first, light mode saja.
- **Admin:** clean, data-dense. Hindari glassmorphism, gradien ungu-pink, animasi berlebihan. Dark mode opsional (class strategy).

## 9. Bahasa & Format
- Semua teks UI Bahasa Indonesia yang ramah dan jelas ("Simpan", "Batal", "Data berhasil disimpan").
- Pakai fungsi di `fase-0-fondasi/src/utils/format.js`:
  - Rupiah: `Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 })` → "Rp 125.000". Selalu tampilkan satuan jual.
  - Tanggal: `id-ID`, zona Asia/Jakarta, contoh "6 Okt 2026".
  - HP: input boleh `08…` / `+62…`, normalisasi ke `62…` sebelum dikirim.

## 10. SEO & Performa
- URL produk/kategori/brand/ruangan/artikel pakai `slug`.
- Set `<title>` & meta description per halaman lewat `useSeo` (pakai `meta_title` / `meta_description` jika ada). `useSeo` juga menerima `canonical`, `robots`, `image`, dan `jsonLd` (schema.org); tag dibersihkan otomatis saat halaman ditutup. Halaman hasil cari/filter = `noindex, follow` + canonical ke halaman bersih (contoh di `/brand`, KEPUTUSAN #17). Alamat situs untuk canonical dari `VITE_SITE_URL`.
- Paginasi daftar publik berupa tautan (`PaginationBar` prop `tautan`) supaya bisa diikuti mesin pencari.
- Karena SPA, pratinjau link di WhatsApp/Facebook memakai meta default di `index.html`. Pratinjau per produk butuh prerender/SSR, di luar stack saat ini.
- Gambar: `loading="lazy"`, `alt` deskriptif, ukuran eksplisit agar tidak layout shift.
- Daftar produk pakai pagination. Filter katalog disimpan di URL query.

## 11. Checklist Sebelum Halaman Selesai
- Kontras teks ≥ 4.5:1. `cursor-pointer`, hover, dan `focus-visible` ring di semua elemen interaktif.
- Responsive 375 / 768 / 1024 / 1440 px.
- Status tidak dibedakan hanya dengan warna (ada label/ikon).
- Setiap tabel/list/form punya state loading (skeleton), kosong (EmptyState), dan error (pesan + coba lagi).
- Form: label jelas, error di bawah field, tombol submit loading + cegah double submit.
- Aksi destruktif pakai modal konfirmasi. Transisi 150–200 ms, hormati `motion-reduce`.

## 12. Masih Terbuka — TANYA DULU
- Daftar 58 kolom Master Barang: mana yang tampil di detail produk (menunggu header kolom).
- Warna token vs logo resmi: lambang resmi (rumah biru + H kuning di latar oranye `#FF2F00`, dari `docs/wireframe-low-fidelity-transhome.xlsx`) sudah dipakai di favicon dan di samping tulisan Transhome (`LogoTranshome.vue`, KEPUTUSAN #18). Token warna situs masih dari sketsa; belum diputuskan apakah ikut warna logo.
- Atribut apa saja di `specifications` JSON (untuk filter dinamis).

Menunggu backend (laporkan, jangan dikerjakan di frontend):
- Migration `users.phone_number` boleh NULL (KEPUTUSAN #11). Database masih `NOT NULL`.
- Seed `menus` dan `permissions` (masih kosong). Daftar kode izin ada di lampiran `docs/FRONTEND_BUILD_PLAN.md`.
- Endpoint REST sesuai lampiran rencana build. Sementara pakai mock.

## 13. Git
- Sejak 3 Nov 2026 (keputusan pemilik project): kerja **langsung di `main`**, tanpa branch fitur + PR. Frontend dan backend sama-sama push ke `main`, jadi **selalu `git pull origin main` sebelum mengubah** dan sebelum push. Commit/push hanya atas permintaan pemilik project. Jangan mengubah folder `backend/` dari sisi frontend.
- Commit kecil dan deskriptif (Conventional Commits: `feat:`, `fix:`, `refactor:`, `style:`).

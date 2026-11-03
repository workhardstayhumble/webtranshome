# Frontend Transhome

Frontend Website Transhome Muka 1 (katalog digital + membership Trans Family).
Stack: **Vue 3 + Vite + Tailwind CSS v4**, dengan Vue Router, Pinia, axios, dan
lucide-vue-next. Aturan lengkap ada di [CLAUDE.md](CLAUDE.md).

## Dokumen

| File | Isi |
|---|---|
| [docs/KEPUTUSAN.md](docs/KEPUTUSAN.md) | Keputusan terbaru, mengalahkan BRD/PRD |
| [docs/FRONTEND_BUILD_PLAN.md](docs/FRONTEND_BUILD_PLAN.md) | Rencana Fase 0–6: tombol, API, checklist QA |
| [docs/MASTER_PROMPT_UI.md](docs/MASTER_PROMPT_UI.md) | Panduan desain untuk Antigravity |
| [docs/transhome_postgres.sql](docs/transhome_postgres.sql) | Bentuk data untuk mock (jangan dijalankan ke database) |
| `docs/*.pdf` | BRD dan PRD v0.4 |

## Menjalankan

Persyaratan: Node.js 20.19+ atau 22.12+, dan npm.

```bash
cd frontend
cp .env.example .env        # isi VITE_API_BASE_URL
npm install
npm run dev                 # http://localhost:5173
```

Perintah lain: `npm run build` (hasil di `dist/`) dan `npm run preview`.

## Data

Frontend hanya memakai REST API backend Go lewat satu instance axios di
`fase-0-fondasi/src/services/http.js`. Tidak ada koneksi database dan tidak ada logika bisnis
penting di sisi browser.

Selama endpoint backend belum ada, `VITE_USE_MOCK=true` membuat axios dijawab
data contoh di `services/mock/` (bentuk sesuai `docs/transhome_postgres.sql`), yang tersebar di folder
fasenya masing-masing (data dasar di Fase 0, data admin di Fase 2, dan seterusnya).
Service tidak perlu diubah saat API siap: cukup set `VITE_USE_MOCK=false`.

Akun contoh mode mock (password semua `rahasia123`): `member@example.com`,
`katalog@example.com`, `membership@example.com`, `konten@example.com`,
`superadmin@example.com`. Panel admin: buka `/admin/masuk`. Cara cepat lain:
`/dev/components`, bagian "Akun contoh dan route guard". Daftar member admin
berisi 25 member contoh (termasuk yang tanpa HP, nonaktif, dan sync gagal).

Di mode mock, perubahan data (produk, kategori, brand, ruangan, akun member
yang baru daftar, transaksi poin) disimpan di localStorage browser, jadi tetap
ada setelah refresh dan sama di semua tab. Tombol "Reset data contoh" di `/dev/components`
mengembalikan data awal; "Simulasikan sesi habis" untuk menguji alur login ulang.

Mode mock tidak mengirim email sungguhan dan tidak memanggil Google:
- `/dev/email`: kotak email tiruan berisi link verifikasi dan reset password.
- "Masuk/Daftar dengan Google" membuka `/dev/google` untuk memilih akun Google
  contoh (akun baru dibuatkan member) atau membatalkan.

## Galeri komponen

`npm run dev`, lalu buka `http://localhost:5173/dev/components`. Semua komponen
dasar tampil dalam setiap state. Halaman ini tidak ikut build produksi.

`.env` hanya berisi variabel `VITE_*`, dan semuanya ikut terbaca di browser.
Jangan simpan rahasia di sana.

## Struktur

Kode dipisah per fase. Setiap folder fase berisi `src/` dengan susunan yang sama, dan aplikasinya tetap satu
(`npm run dev` dari folder `frontend/`). Import ditulis `@/...`; alias `@` mencari di `src/` semua folder fase.
Daftar isi tiap folder ada di [docs/FILE_PER_FASE.md](docs/FILE_PER_FASE.md).

```
fase-0-fondasi/            pengaturan aplikasi, komponen dasar, layout, router, store, HTTP, mock dasar, galeri dev
fase-1-katalog/            katalog, pencarian, detail produk, brand, ruangan
fase-2-admin-katalog/      admin: masuk, dashboard, produk, kategori, brand, ruangan
fase-3-akun-member/        masuk, daftar, verifikasi, lupa password, area /akun
fase-4-membership-admin/   admin: daftar & detail member, input poin
fase-5-konten/             Beranda + hero banner (gambar/video), artikel, promo, info toko; admin banner, artikel, pengaturan
fase-6-sistem/             admin: pengguna admin, role & izin, log aktivitas, sync Qontak
```

Susunan di dalam setiap `src/` (gabungan semua fase):

```
src/
  assets/main.css   Tailwind + token desain dari sketsa (@theme)
  layouts/          PublicLayout, AuthLayout, MemberLayout, AdminLayout
  views/            public/, auth/, member/, admin/<modul>/, dev/
  components/       ui/ (komponen dasar), layout/ (header, footer, WhatsApp), katalog/, produk/, admin/, akun/
  composables/      useAuth, usePermission (+ v-can), useWhatsApp, useFilter, useKatalog, useSeo,
                    useFormProduk, useKonfirmasiKeluar, useMasterData, useTujuanKembali, useJeda,
                    useFormPoin, useMemberTerpilih, useSyncQontak, useMuat, useDaftarBanner, useFormBanner,
                    useFormArtikel, useDaftar, useMatriksIzin
  stores/           auth, menu, ui, settings (Pinia)
  services/         http.js, errors.js, satu service per modul, mock/
  router/           routes per area (public, auth, member, admin) + guard.js
  utils/            format.js (Rupiah, tanggal, nomor HP), whatsapp.js, produk.js (spesifikasi publik), validasi.js,
                    poin.js (jenis transaksi, hitung poin), artikel.js (tipe, blok isi), banner.js (posisi, status),
                    log.js (label aksi & kolom log aktivitas)
design-system/transhome/MASTER.md   aturan desain (baca sebelum membuat halaman)
```

## Status

Fase 0 (fondasi), Fase 1 (katalog, pencarian, detail produk, brand, ruangan),
Fase 2 (admin: login, layout, dashboard, master katalog), Fase 3 (auth dan
akun member), Fase 4 (membership admin: daftar & detail member, input
poin), Fase 5 (konten: Beranda + hero banner, artikel, promo, info toko,
admin banner/artikel/pengaturan), dan Fase 6 (sistem: pengguna admin, role &
izin, log aktivitas, sync Qontak) selesai 7 Okt 2026; hasil QA ada di
[docs/FRONTEND_BUILD_PLAN.md](docs/FRONTEND_BUILD_PLAN.md).

- Publik: `/` (Beranda), `/artikel`, `/artikel/:slug` (`?preview=<token>` untuk draf),
  `/promo` (`?jenis=promo|event`), `/info-toko`, `/katalog`, `/katalog/:kategori`, `/cari?q=`, `/produk/:slug`,
  `/brand` (semua brand A–Z), `/brand/:slug`, `/ruangan`, `/ruangan/:slug`.
- Auth: `/masuk`, `/daftar`, `/verifikasi`, `/lupa-password`, `/reset-password`,
  `/masuk/google` (callback Google).
- Member: `/akun`, `/akun/poin`, `/akun/profil`, `/akun/keamanan`.
- Admin: `/admin/masuk`, `/admin`, `/admin/produk` (+ `/baru`, `/:id`,
  `/:id/pratinjau`), `/admin/kategori`, `/admin/brand`, `/admin/ruangan`,
  `/admin/member`, `/admin/member/:id`, `/admin/poin/input` (`?member=<id>`),
  `/admin/banner` (+ `/baru`, `/:id`, `/:id#periode`), `/admin/artikel` (+ `/baru`, `/:id`),
  `/admin/pengaturan`, `/admin/pengguna`, `/admin/role` (`?role=<kode>`),
  `/admin/log-aktivitas`, `/admin/sync-qontak` (`?status=FAILED`).

Semua fase Muka 1 (0–6) sudah selesai dengan data mock. Langkah berikutnya: sambungkan ke
API backend begitu endpoint di lampiran rencana build tersedia (matikan `VITE_USE_MOCK`).

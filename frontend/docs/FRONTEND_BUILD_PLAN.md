# Rencana Build Frontend — Website Transhome Muka 1

Stack: Vue 3 + Vite + Vue Router + Pinia + Tailwind CSS v4 + axios + lucide-vue-next (JavaScript, SPA). Acuan: `CLAUDE.md`, `docs/KEPUTUSAN.md`, `docs/transhome_postgres.sql`, sketsa "Sketsa Dashboard Transhome".
Cara pakai di Antigravity atau Claude Code: kerjakan SATU fase per percakapan. Selesai satu fase → jalankan checklist QA fase itu → baru lanjut.
Versi ini disesuaikan dari rencana berbasis Nuxt pada 7 Okt 2026 (KEPUTUSAN #7). Path kode (`views/...`, `services/...`, dst.) relatif terhadap `src/` di folder fasenya (`fase-N-.../src/`, lihat `docs/FILE_PER_FASE.md`); path lain relatif terhadap folder `frontend/`.

---

## Aturan Workflow Tombol (berlaku untuk SEMUA tombol)

Setiap tombol yang memanggil API wajib mengikuti pola ini:

1. **Klik** → tombol masuk state `loading` (spinner) dan tidak bisa diklik lagi sampai selesai (cegah double submit).
2. **Validasi client** dulu (field wajib, format). Gagal → tampilkan error di bawah field, JANGAN panggil API.
3. **Panggil API** lewat service di `services/` (memakai axios dari `services/http.js`).
4. **Berhasil** → toast singkat ("Produk disimpan") + pindah halaman / perbarui data di layar.
5. **Gagal** → ditangani terpusat: interceptor di `services/http.js` untuk 401, 403, 429, 5xx/jaringan, dan helper error form untuk 400/422/409:

| Status | Perilaku |
|---|---|
| 400 / 422 | Pesan error dipetakan ke field masing-masing (`errors.sku`, dst) |
| 401 | Sesi habis → logout lokal, redirect ke `/masuk` (area member) atau `/admin/masuk` (area admin) dengan `?redirect=<halaman sekarang>` |
| 403 | Toast "Kamu tidak punya akses untuk aksi ini" |
| 404 | Halaman "Tidak ditemukan" (untuk halaman detail) atau toast (untuk aksi) |
| 409 | Pesan konflik spesifik (SKU sudah dipakai, nota sudah diinput, email terdaftar) |
| 429 | "Terlalu banyak percobaan. Coba lagi beberapa saat." |
| 5xx / jaringan | Toast "Terjadi gangguan. Coba lagi." + tombol **Coba lagi** |

Aturan tambahan:
- **Aksi destruktif** (nonaktifkan, hapus foto, keluar dari form yang belum disimpan) → `ConfirmModal` dulu.
- **Form panjang** → peringatan "Buang perubahan?" saat pindah halaman kalau ada perubahan belum disimpan (`onBeforeRouteLeave`).
- **Tombol berbasis izin** → disembunyikan pakai `can('kode.izin')` atau `v-can`. Backend tetap menolak (403).
- **Toggle cepat** (unggulan, aktif, consent) → optimistic update: UI langsung berubah, dikembalikan + toast error kalau API gagal.
- **Link WhatsApp** → `target="_blank" rel="noopener"`, pesan di-encode `encodeURIComponent`, nomor hanya-baca dari API settings (store `settings`), bersumber dari `WHATSAPP_NUMBER` di server (KEPUTUSAN #21).
- Semua tombol bisa dipakai dengan keyboard (Tab, Enter, Esc untuk modal/menu) dan punya `focus-visible` ring.

---

## FASE 0 — Fondasi (wajib selesai sebelum fase lain)

Tidak ada library komponen UI, jadi semua komponen dasar dibuat sendiri dengan Tailwind.
**Status: SELESAI 7 Okt 2026** — checklist QA di bawah lulus 31/31 (diuji otomatis lewat Chrome DevTools Protocol).

### Yang dibangun
| Bagian | Isi | Status |
|---|---|---|
| Konfigurasi | `vite.config.js` (alias `@`), `.env`: `VITE_API_BASE_URL`, `VITE_USE_MOCK` | Selesai |
| Token desain | ui-ux-pro-max `--design-system` (pola "Trust & Authority + Conversion", gaya "Flat Design"), palet dan font diganti sketsa. Hasil: `design-system/transhome/MASTER.md` + `@theme` di `fase-0-fondasi/src/assets/main.css` | Selesai |
| HTTP & service | `services/http.js` (axios + interceptor 401/403/429/5xx; toast hanya untuk aksi, GET → ErrorState), `services/errors.js` (`errorField`, `pesanError`), service per modul. Mode mock = adapter axios ke `services/mock/`, jadi service tidak perlu tahu sedang mock | Selesai |
| Composable | `useAuth`, `usePermission` + `v-can`, `useWhatsApp` (`productLink`, `bulkLink`, `generalLink`), `useFilter`, `useSeo` | Selesai |
| Utilitas | `utils/format.js`, `utils/whatsapp.js` (pesan tanya produk, "Tanya ketersediaan", jumlah besar, umum) | Selesai |
| Router guard | `meta.area` `public`/`guest`/`member`/`admin` + `meta.permission`, pemulihan sesi setelah refresh, `?redirect=` hanya path internal | Selesai |
| Store | `auth` (+ `restoreSession`), `menu`, `ui`, `settings` | Selesai |
| Layout | `PublicLayout` (skip link, header, footer, WA melayang), `AuthLayout`, `MemberLayout` (tab di HP), `AdminLayout` (sidebar navy) | Selesai |
| Komponen dasar | 18 di `components/ui/`, 7 di `components/layout/` (termasuk `UserMenu`), 2 di `components/katalog/` | Selesai |
| Mock | `services/mock/` (data + rute + adapter), akun contoh untuk uji guard (password `rahasia123`, hanya mode mock) | Selesai |
| Halaman dev | `/dev/components`, tidak ikut build produksi | Selesai |
| Placeholder route | `/masuk`, `/daftar`, `/verifikasi`, `/lupa-password`, `/reset-password`, `/admin/masuk`, `/akun/*`, `/admin/*` menampilkan "sedang dibangun (Fase n)" agar guard dan tautan header bisa diuji | Selesai; diganti halaman asli di fasenya |

### Tombol global (header, footer, melayang)
| Tombol | Aksi | Berhasil | Gagal / catatan |
|---|---|---|---|
| Logo | → `/` | — | — |
| **Kategori** | Buka mega menu (`GET /categories?tree=1`, di-cache di store) | Panel terbuka, fokus ke item pertama | Esc / klik di luar → tutup. Gagal load → teks "Kategori gagal dimuat" + Coba lagi |
| Item kategori di mega menu | → `/katalog/[slug]` | Menu tertutup | — |
| **Kotak cari** (Enter / ikon) | → `/cari?q=<teks>` | Halaman hasil | Teks kosong → tidak terjadi apa-apa |
| Saran pencarian | Ketik ≥ 2 huruf, debounce 300 ms → `GET /products/suggest?q=` | Dropdown Produk / Kategori / Brand; ↑↓ + Enter memilih | Tidak ada hasil → "Tidak ada saran" |
| **Masuk** | → `/masuk?redirect=<halaman sekarang>` | — | Hanya tampil saat belum login |
| **Daftar member** | → `/daftar` | — | Hanya tampil saat belum login |
| Avatar (sudah login) | Buka menu: Akun saya, Riwayat poin, Keluar | — | — |
| **Keluar** | `POST /auth/logout` | State user dikosongkan, toast "Kamu sudah keluar", → `/` | Gagal jaringan → tetap kosongkan state lokal |
| **WhatsApp melayang** | Buka `wa.me/<nomor>?text=` pesan umum | Tab WhatsApp terbuka | Kalau login: pesan diawali nama + member_code (FR-S) |
| Link footer | Navigasi biasa | — | — |

### QA Fase 0 (hasil 7 Okt 2026)
- [x] `/dev/components` tampil benar di 375 / 768 / 1024 / 1440 px: konten ter-render, tanpa scroll horizontal, tidak ada teks/tombol yang keluar dari kartu produk.
- [x] Mega menu bisa dibuka/tutup pakai mouse dan keyboard (Enter membuka dan memindahkan fokus ke panel, Esc menutup dan mengembalikan fokus, klik di luar menutup).
- [x] Mematikan API (mock off + server mati) → muncul ErrorState + Coba lagi, header "Info toko belum bisa dimuat.", bukan layar putih.
- [x] `VITE_USE_MOCK=true` → semua data dari `services/mock/`.
- [x] Buka `/akun` tanpa login → `/masuk?redirect=/akun`; buka `/admin` → `/admin/masuk?redirect=/admin`.
- [x] `npm run build` berhasil tanpa error dan tanpa warna hex di komponen; galeri dev tidak ikut build produksi.
- [x] Tambahan: saran pencarian (↓ + Enter), member ditolak di `/admin`, admin ditolak di `/akun`, sesi pulih setelah reload, modal mengunci dan mengembalikan fokus, toast diumumkan lewat live region, konsol tanpa error.

---

## FASE 1 — Katalog, Pencarian, Detail Produk, Brand, Ruangan (FR-A–FR-D, FR-F)

**Status: SELESAI 7 Okt 2026** — checklist QA di bawah lulus 85/85 (Chrome DevTools Protocol, mode mock + mode API mati), regresi Fase 0 tetap 31/31.

### Yang dibangun
| Bagian | Isi |
|---|---|
| Halaman | `views/public/`: `Katalog` (`/katalog/:kategori?`), `Cari`, `Brand`, `Ruangan`, `RuanganDetail`, `ProdukDetail` |
| Komponen katalog | `components/katalog/`: `KatalogDaftar` (dipakai ulang di katalog, cari, brand, detail ruangan), `FilterPanel`, `FilterOpsi`, `RentangHarga`, `FilterSheet` (bottom sheet HP), `ActiveFilterChips`, `KatalogToolbar`, `ChipNav`, `RoomCard`, `RoomCarousel` |
| Komponen produk | `components/produk/`: `ProductGallery`, `Lightbox`, `ProductSummary`, `ProductTabs`, `RelatedProducts` |
| Komponen dasar baru | `ui/Breadcrumb` (di HP diringkas jadi "‹ induk"), `ui/NotFoundState` (dipakai semua 404), `ui/CopyButton`; `Modal` mendapat mode `sheet` dan ukuran `xl` |
| Data | `composables/useKatalog` (filter dari URL, cache 60 detik untuk tombol Back, chip filter aktif), `services/catalogService` (brand, ruangan), `utils/produk.js` (`spesifikasiPublik` membuang HPP, hnet, diskon, stok gudang/lokasi, konsinyasi, karantina) |
| Mock | 38 produk (1 nonaktif, 1 HABIS, 1 berisi `hpp` + `stok_gudang` untuk menguji penyaringan), faset disjungtif, rentang harga, ruangan dan brand dengan `product_count` |

Catatan implementasi:
- Mengubah filter tidak menggulir halaman (router hanya menggulir ke atas saat path berubah). Paginasi menggulir ke atas daftar dan memindahkan fokus ke judul daftar.
- Judul tab dari `useSeo` tidak ditimpa saat hanya query yang berubah.
- Di HP, detail produk memakai bar bawah (harga + WhatsApp); tombol WhatsApp melayang disembunyikan di halaman itu (`meta.barAksiHp`), dan layout memberi ruang di bawah footer.
- "Chat tim proyek" memakai gaya tombol WhatsApp (hijau), sesuai aturan semua tombol WhatsApp di `MASTER.md`.
- Kotak cari di header menampilkan kata kunci dari `/cari?q=` tanpa membuka saran.

### `/katalog/[kategori]`, `/cari?q=`, `/brand/[slug]`
Data: `GET /products?category=&brand=&stock=&room=&attr[ukuran]=&min=&max=&sort=&page=&per_page=24` (hanya produk `is_active = true`).
**Semua filter disimpan di URL query** (`useFilter`) → bisa dibagikan, tombol Back browser mengembalikan filter.

| Tombol | Aksi | Berhasil | Gagal / catatan |
|---|---|---|---|
| Checkbox filter (brand, status stok, ruangan, atribut spesifikasi) | Ubah query (`router.replace`), `page` kembali ke 1, data dimuat ulang | Hasil + jumlah produk + angka per opsi diperbarui | Skeleton saat memuat |
| Harga min / maks | Diterapkan saat Enter atau keluar dari field | Sama seperti di atas | Min > maks → error "Harga minimum melebihi maksimum" |
| **Reset** (panel filter) | Hapus semua filter kecuali kategori | Semua produk kategori tampil | — |
| Chip filter aktif (×) | Hapus satu nilai filter | Chip hilang, data dimuat ulang | — |
| **Hapus semua** | Sama dengan Reset | — | — |
| **Urutkan** (select) | Ubah `sort` di query | Urutan berubah | — |
| **Grid / Daftar** | Ganti tampilan (disimpan di `localStorage` dengan try/catch, tanpa memuat ulang data) | Tampilan berganti | — |
| Paginasi (Sebelumnya, angka, Berikutnya) | Ubah `page` | Data halaman baru, scroll ke atas daftar | Tombol di ujung disembunyikan, bukan dinonaktifkan |
| **Filter (n)** (HP) | Buka bottom sheet filter | Sheet naik, fokus terkunci di dalam | — |
| **Terapkan** (HP) | Terapkan filter draf | Sheet tertutup, data dimuat ulang | — |
| **Batal** (HP) | Buang filter draf | Sheet tertutup, filter lama tetap | — |
| Klik foto / nama produk | → `/produk/[slug]` | — | — |
| **Tanya via WhatsApp** (kartu) | `useWhatsApp().productLink(produk)` | WA terbuka dengan SKU + nama + URL produk | — |
| EmptyState: **Reset filter** / **Tanya via WhatsApp** | Seperti di atas | — | Muncul saat 0 hasil |
| ErrorState: **Coba lagi** | Muat ulang data | Data tampil | — |

### `/produk/[slug]`
Data: `GET /products/:slug` (berisi images, rooms, brand, category). 404 atau produk nonaktif → halaman "Produk tidak ditemukan" + tombol **Kembali ke katalog**.

| Tombol | Aksi | Berhasil | Gagal / catatan |
|---|---|---|---|
| Thumbnail foto | Ganti foto utama | Thumbnail aktif bergaris warna primary | — |
| Thumbnail video | Putar video di lightbox | — | Tampil hanya jika ada video. **Belum dibuat:** `products` dan `product_images` tidak punya kolom video (lihat lampiran, kebutuhan backend) |
| **Perbesar** / klik foto utama | Buka lightbox | ← → ganti foto, Esc tutup, fokus kembali ke tombol | — |
| **Salin SKU** | `navigator.clipboard.writeText(sku)` | Label jadi "SKU disalin" 2 detik | Clipboard ditolak → pilih teks SKU otomatis |
| **Tanya via WhatsApp** | `productLink(produk)` | WA terbuka | Jika `HABIS` → teks tombol "Tanya ketersediaan", pesan menyesuaikan |
| **Bagikan** | `navigator.share()` jika tersedia, selain itu salin URL | Toast "Link disalin" | — |
| **Unduh datasheet** | Buka `datasheet_pdf_url` di tab baru | — | Disembunyikan jika kosong |
| **Chat tim proyek** | `bulkLink(produk)` (pesan kebutuhan proyek) | WA terbuka | — |
| Tab Deskripsi / Spesifikasi / Dokumen | Ganti isi tab (← → di keyboard) | — | Tab Dokumen disembunyikan jika tidak ada dokumen |
| Chip ruangan | → `/ruangan/[slug]` | — | — |
| Nama brand | → `/brand/[slug]` | — | — |
| Kartu produk terkait | → `/produk/[slug]` | — | `GET /products?category=<slug kategori>&exclude=<id produk>&per_page=4`. Tidak ada produk serupa → bagian disembunyikan |

### `/ruangan`, `/ruangan/[slug]`
| Tombol | Aksi | Catatan |
|---|---|---|
| Kartu ruangan | → `/ruangan/[slug]` | — |
| ‹ › carousel | Geser kartu | Hormati `prefers-reduced-motion` |
| Halaman detail ruangan | Daftar produk = katalog dengan filter `room` terkunci | Pakai ulang komponen katalog |

### QA Fase 1 (hasil 7 Okt 2026)
- [x] Pilih 2 brand + status stok → refresh browser → filter tetap sama (URL, centang, chip, jumlah).
- [x] Back browser setelah buka detail produk → kembali ke katalog dengan filter & halaman yang sama, posisi gulir pulih.
- [x] Kombinasi filter tanpa hasil → EmptyState tampil dengan Reset filter + Tanya via WhatsApp.
- [x] Pesan WA berisi SKU, nama, dan URL yang benar (3 halaman detail + 3 kartu katalog).
- [x] Produk `HABIS` → tombol berubah jadi "Tanya ketersediaan" (detail, bar HP, kartu katalog), pesan menanyakan ketersediaan.
- [x] Slug ngawur → halaman 404 yang rapi (produk, produk nonaktif, kategori, brand, ruangan).
- [x] Tidak ada HPP / diskon / stok gudang yang muncul di mana pun (termasuk isi tab yang tersembunyi).
- [x] `useSeo` mengisi title dan meta description di halaman detail (`meta_title`/`meta_description`, atau dibuat otomatis).
- [x] Tambahan: faset disjungtif, chip × dan Hapus semua, urutkan, harga min/maks (Enter, keluar field, error min > maks), paginasi (gulir + fokus, tombol ujung disembunyikan), Grid/Daftar tersimpan, filter tidak menggulir halaman, bottom sheet HP (draf, pratinjau jumlah, Terapkan/Batal/Esc, fokus kembali), pencarian (reset mempertahankan kata kunci), lightbox (← → Esc, fokus kembali), tab (← → Home End), Salin SKU + cadangan bila clipboard ditolak, Bagikan (share/salin), carousel ruangan (‹ › nonaktif di ujung, reduced motion), tanpa scroll horizontal di 375/768/1024/1440 px untuk 8 halaman, API mati → ErrorState, konsol bersih.

---

## FASE 2 — Admin: Login, Layout, Dashboard, Master Katalog (FR-G)

**Status: SELESAI 7 Okt 2026.** Checklist QA di bawah lulus 79/79 (Chrome DevTools Protocol, mode mock + API mati). Uji rute mock admin lulus 38/38, dan regresi Fase 0 (31/31) serta Fase 1 (85/85) tetap lulus.

### Yang dibangun
| Bagian | Isi |
|---|---|
| Halaman | `views/admin/`: `Masuk`, `Dashboard`, `TidakDitemukan`; `views/admin/katalog/`: `ProdukDaftar`, `ProdukForm` (baru & ubah), `ProdukPratinjau`, `Kategori`, `Brand`, `Ruangan` |
| Layout | `AdminLayout` sesuai sketsa: sidebar navy dari `GET /admin/menus` (`components/admin/AdminSidebar`), header grup › judul + badge role, ☰ drawer di HP (fokus terkunci, Esc menutup) |
| Komponen admin | `PerhatianList` (termasuk Coba lagi sync Qontak), `AktivitasTerbaru`, `FormSection`, `GambarField`, `MasterFormModal`, `TombolUrut`, `BrandCepatModal`; `admin/produk/`: `ProdukFilterBar`, `ProdukTabel`, `FotoProduk`, `DatasheetField`, `SpesifikasiEditor`, `RuanganPicker` |
| Komponen dasar | Baru: `PasswordInput`, `Textarea`. Diperluas: `Input` (prefix + slot ikon kanan), `Select` (opsi berkelompok), `Table` (tanpa bingkai, lebar minimum) |
| Composable | `useFormProduk` (state, slug otomatis, cek SKU, validasi, draf saat sesi habis), `useKonfirmasiKeluar` ("Buang perubahan?" + peringatan tutup tab), `useMasterData`, `useTujuanKembali` (Batal/Simpan kembali ke daftar dengan filter yang sama) |
| Service | `adminService` (menu, dashboard), `adminProductService`, `adminMasterService`, `uploadService` (aturan file dicek juga di browser), `qontakService`, `cache.js` (cache publik dibersihkan setelah admin mengubah data) |
| Router | Rute admin + `meta.permission`. Halaman Fase 4–6 sudah dipasangi izin, jadi role yang tidak berhak langsung ke "Tidak punya akses" |
| Mock | `routesAdmin.js` (menu per izin, dashboard per role, CRUD produk/master, upload, retry Qontak), `sesi.js`, `db.js`: perubahan admin disimpan di localStorage agar sama di semua tab. Di `/dev/components` ada tombol "Simulasikan sesi habis" dan "Reset data contoh" |

Catatan implementasi:
- Login admin memakai pesan di form (bukan toast). Setelah sesi habis, form produk menyimpan draf di `sessionStorage`; setelah masuk lagi isian dipulihkan dengan pemberitahuan.
- Produk aktif wajib punya foto utama (dicek di browser dan di backend, 422). Produk tanpa foto tetap boleh disimpan sebagai nonaktif.
- Kunci spesifikasi disimpan snake_case (`Isi per dus` → `isi_per_dus`), sama dengan label di halaman publik.
- Pratinjau produk aktif membuka halaman publik; produk nonaktif membuka `/admin/produk/:id/pratinjau` (label "Tidak tampil di website").
- Daftar produk punya filter tambahan "Foto" (`foto=0/1`) untuk tautan "Lihat semua produk tanpa foto" dari dashboard.

### `/admin/masuk`
| Tombol | Aksi | Berhasil | Gagal |
|---|---|---|---|
| Ikon mata | Tampil / sembunyikan password | — | — |
| **Masuk** | `POST /admin/auth/login {email, password}` → `GET /me` | Punya role admin → ke `?redirect` atau `/admin` | 401 → "Email atau password salah". Bukan admin → logout + "Akun ini tidak punya akses admin". 429 → pesan batas percobaan |
| Lupa password | → `/lupa-password` | — | — |

### Layout admin
| Tombol | Aksi | Catatan |
|---|---|---|
| Item sidebar | Navigasi | Daftar dari `GET /admin/menus` (sudah difilter backend sesuai izin), disimpan di store `menu`. Item aktif disorot |
| ☰ (HP) | Buka / tutup sidebar | — |
| **Keluar** | `POST /auth/logout` → `/admin/masuk` | — |

### `/admin` (Dashboard)
Data: `GET /admin/dashboard` (ringkasan, perlu perhatian, log terbaru).

| Tombol | Izin | Aksi |
|---|---|---|
| **+ Tambah produk** | `product.create` | → `/admin/produk/baru` |
| **Input poin member** | `point.create` | → `/admin/poin/input` |
| **Buat banner** | `banner.create` | → `/admin/banner/baru` |
| **Tulis artikel** | `article.create` | → `/admin/artikel/baru` |
| Kartu ringkasan | sesuai modul | → daftar terkait dengan filter (mis. Sync gagal → `/admin/sync-qontak?status=FAILED`) |
| **Ubah stok** (stok habis) | `product.update` | → `/admin/produk/:id#stok` |
| **Perpanjang** (banner) | `banner.update` | → `/admin/banner/:id` |
| **Coba lagi** (sync gagal) | `qontak.retry` | `POST /admin/qontak/contacts/:id/retry` → baris berubah PENDING lalu SYNCED/FAILED |
| **Lengkapi** (produk tanpa foto) | `product.update` | → `/admin/produk/:id#foto` |
| Log lengkap | `log.view` | → `/admin/log-aktivitas` |

### `/admin/produk` (daftar)
Data: `GET /admin/products?q=&category=&brand=&stock=&active=&page=`

| Tombol | Izin | Aksi | Berhasil | Gagal / catatan |
|---|---|---|---|---|
| Cari (debounce 400 ms) & filter select | `product.view` | Ubah query | Tabel diperbarui | — |
| **+ Tambah produk** | `product.create` | → `/admin/produk/baru` | — | — |
| ★ Unggulan | `product.update` | `PATCH /admin/products/:id {is_featured}` (optimistic) | Bintang berubah | Gagal → kembali + toast |
| Switch **Aktif** → mati | `product.update` | ConfirmModal "Produk tidak akan tampil di website" → `PATCH {is_active:false}` | Switch mati | — |
| Switch **Aktif** → nyala | `product.update` | `PATCH {is_active:true}` | Switch nyala | 422 "Tambahkan minimal 1 foto utama dulu" |
| **Ubah** | `product.update` | → `/admin/produk/:id` | — | — |
| **Pratinjau** | `product.view` | Buka `/produk/[slug]` di tab baru | — | Produk nonaktif → pratinjau dengan label "Tidak tampil di website" |
| Paginasi | — | Ubah `page` | — | — |

### `/admin/produk/baru` dan `/admin/produk/:id` (form)
| Tombol / field | Aksi | Berhasil | Gagal / catatan |
|---|---|---|---|
| SKU (keluar field) | `GET /admin/products/check-sku?sku=&exclude=<id>` | Centang hijau | "SKU sudah dipakai produk lain" |
| Nama → slug otomatis | Slug dibuat dari nama (bisa diedit) | — | Slug duplikat → 409 saat simpan |
| **+ Tambah brand** (di select brand) | Modal cepat → `POST /admin/brands` | Brand baru langsung terpilih | — |
| Spesifikasi: **+ Tambah baris** / × hapus baris / geser urutan | Repeater key-value → `specifications` JSON | — | Key kosong → error baris |
| **Unggah foto** (multi) | `POST /admin/uploads` (jpg/png/webp ≤ 2 MB) | Pratinjau muncul | Format/ukuran salah → error per file |
| **Jadikan utama** | Tandai `is_primary` | Label "Utama" pindah | — |
| × hapus foto | ConfirmModal → hapus dari daftar | — | Foto utama dihapus → foto pertama jadi utama |
| Geser foto | Ubah `sort_order` | — | — |
| **Unggah datasheet** | `POST /admin/uploads` (PDF ≤ 10 MB) | Nama file tampil | — |
| × hapus datasheet | Kosongkan `datasheet_pdf_url` | — | — |
| Pilih ruangan (multi-select) | → `product_rooms` | — | — |
| **Simpan** | Validasi → `POST /admin/products` atau `PUT /admin/products/:id` (1 transaksi di backend) | Toast "Produk disimpan" → kembali ke daftar | 422 → error per field. 409 → SKU/slug duplikat |
| **Simpan & tambah lagi** (form baru) | Sama, lalu form dikosongkan | Fokus ke field SKU | — |
| **Batal** | Ada perubahan → "Buang perubahan?" | → daftar | — |

### `/admin/kategori`, `/admin/brand`, `/admin/ruangan`
| Tombol | Izin | Aksi | Gagal / catatan |
|---|---|---|---|
| **+ Tambah** | `<modul>.create` | Modal form → `POST` | Slug duplikat → 409 |
| **Ubah** | `<modul>.update` | Modal form → `PUT` | — |
| Switch **Aktif** | `<modul>.update` | ConfirmModal saat menonaktifkan → `PATCH` | Kategori yang masih punya produk tetap boleh dinonaktifkan; tampilkan peringatan jumlah produk |
| ↑ ↓ urutan | `<modul>.update` | `PATCH {sort_order}` | — |
| Unggah gambar / logo / cover | `<modul>.update` | `POST /admin/uploads` | Cover ruangan wajib |
| **Lihat produk** (ruangan) | `product.view` | → `/admin/produk?room=<id>` | — |

### QA Fase 2 (hasil 7 Okt 2026)
- [x] Login dengan akun Admin Katalog → sidebar hanya menampilkan menu katalog (Utama + Katalog).
- [x] Buka `/admin/poin/input` langsung dengan akun Admin Katalog → halaman "Tidak punya akses" (juga `/admin/pengguna`).
- [x] Tambah produk lengkap (3 foto, datasheet, 2 ruangan) → muncul di katalog publik, termasuk brand baru, spesifikasi, dan filter atribut.
- [x] SKU duplikat → ditolak dengan pesan jelas ("SKU sudah dipakai produk lain" di bawah field, fokus ke SKU).
- [x] Nonaktifkan produk → hilang dari katalog publik (404 + tidak muncul di pencarian), tetap ada di admin.
- [x] Ubah form lalu klik menu lain → muncul "Buang perubahan?" ("Tetap di sini" mempertahankan isian).
- [x] Sesi habis di tengah mengisi form → diarahkan ke login, lalu kembali ke halaman semula dengan isian dipulihkan.
- [x] Tambahan: validasi & pesan login (401, akun member, API mati), ikon mata; dashboard per role (aksi cepat, kartu, perhatian, log), Coba lagi sync PENDING → SYNCED/FAILED; daftar produk (cari debounce, filter, ★ optimistic, switch Aktif + konfirmasi, 422 tanpa foto, paginasi, pratinjau); form (slug otomatis, cek SKU, unggah 4 foto + tolak GIF, jadikan utama, geser, hapus + konfirmasi, datasheet, ruangan, brand cepat, ubah produk, `#foto`, Simpan & tambah lagi, Batal); master (tambah subkategori, slug duplikat, peringatan jumlah produk, ↑↓ urutan dipakai halaman publik, cover ruangan wajib, Lihat produk); HP 375 px (drawer, tanpa scroll horizontal); konsol bersih.

---

## FASE 3 — Auth Member & Dashboard Member (FR-E, FR-H)

Ketergantungan backend: `users.phone_number` harus sudah boleh NULL (KEPUTUSAN #11) agar daftar tanpa HP berhasil di API asli.

**Status: SELESAI 7 Okt 2026.** Checklist QA di bawah lulus 52/52 (Chrome DevTools Protocol, mode mock + API mati). Uji rute mock akun lulus 33/33; regresi Fase 0 (31/31), Fase 1 (85/85), dan Fase 2 (79/79) tetap lulus.

### Yang dibangun
| Bagian | Isi |
|---|---|
| Halaman auth | `views/auth/`: `Masuk`, `Daftar`, `Verifikasi` (setelah daftar & dari link email), `LupaPassword`, `ResetPassword`, `MasukGoogle` (callback Google) |
| Area member | `views/member/`: `Ringkasan` (kartu member navy + Salin, saldo, transaksi terakhir, bantuan WhatsApp), `RiwayatPoin` (filter di URL, Muat lebih banyak), `Profil` (lihat/ubah, + Tambah nomor HP), `Keamanan` (persetujuan promo, ubah/buat password) |
| Komponen | `components/akun/`: `KartuAuth`, `TombolGoogle`, `KartuMember`, `DaftarPoin` (tabel di desktop, daftar di HP), `PetunjukMock`; `CopyButton` mendapat mode `gelap` |
| Service & util | `authService` (daftar, verifikasi, lupa/reset password, URL Google), `memberService` (profil, persetujuan, password, poin), `utils/validasi.js` (email, HP Indonesia), `formatNomorHp`, composable `useJeda` (tombol kirim ulang 60 detik) |
| Mock | `routesAkun.js` (login dengan cek verifikasi, daftar, token verifikasi & reset sekali pakai, Google, `/me/*`), password akun baru disimpan sebagai hash. Halaman dev (mode mock): `/dev/email` (kotak email tiruan) dan `/dev/google` (simulasi pilih akun Google) |

Catatan implementasi:
- `/verifikasi`, `/reset-password`, dan `/masuk/google` memakai area `public` (bukan `guest`), supaya link dari email tetap bisa dibuka walau sedang masuk.
- Akun admin yang masuk lewat `/masuk` diarahkan ke `/admin`; `?redirect=` hanya dipakai bila tujuannya bukan halaman admin.
- Token Google diterima lewat fragment URL (`#token=`), lalu langsung dihapus dari riwayat browser.
- Member tanpa password (dibuat lewat Google) melihat "Buat password" di `/akun/keamanan`.
- Rasio poin di kartu saldo diambil dari `site_settings.point_ratio_rupiah`.

### `/daftar`
| Tombol | Aksi | Berhasil | Gagal |
|---|---|---|---|
| Ikon mata | Tampil / sembunyikan password | — | — |
| **Daftar** | Validasi (nama, email, password ≥ 8, HP opsional & format valid) → `POST /auth/register` | → `/verifikasi?email=` | 409 email → "Email sudah terdaftar" + link **Masuk**. 409 HP → "Nomor HP sudah dipakai akun lain" |
| **Daftar dengan Google** | Pindah halaman penuh ke `<VITE_API_BASE_URL>/auth/google/redirect` (OAuth di backend) | Callback → `/akun` (member baru otomatis dibuat) | Dibatalkan → kembali ke `/daftar` + toast |
| Checkbox persetujuan komunikasi | Tidak dicentang default → `communication_consent` | — | Opsional |

### `/masuk`
| Tombol | Aksi | Berhasil | Gagal |
|---|---|---|---|
| **Masuk** | `POST /auth/login {email, password}` | → `?redirect` atau `/akun` | 401 → "Email atau password salah". Email belum diverifikasi → banner + tombol **Kirim ulang verifikasi** |
| **Masuk dengan Google** | Seperti di atas | → `/akun` | — |
| Lupa password | → `/lupa-password` | — | — |
| Belum punya akun? Daftar | → `/daftar` | — | — |

### `/lupa-password`, `/reset-password`, `/verifikasi`
| Tombol | Aksi | Berhasil | Gagal |
|---|---|---|---|
| **Kirim link reset** | `POST /auth/forgot-password {email}` | Pesan SELALU sama: "Jika email terdaftar, link reset sudah dikirim." Tombol jeda 60 detik | — |
| **Simpan password baru** | `POST /auth/reset-password {token, password}` | → `/masuk` + toast "Password diperbarui" | Token tidak valid/kedaluwarsa → pesan + tombol **Minta link baru** |
| **Kirim ulang email verifikasi** | `POST /auth/email/resend` | Toast + jeda 60 detik | — |

### `/akun` (Ringkasan)
| Tombol | Aksi | Catatan |
|---|---|---|
| **Salin** (member_code) | Salin ke clipboard | Label "Disalin" 2 detik |
| Lihat riwayat poin | → `/akun/poin` | — |
| **Chat WhatsApp CS** | `generalLink()` dengan nama + member_code | — |
| Menu samping (Ringkasan, Riwayat poin, Profil, Keamanan) | Navigasi | HP: jadi tab horizontal |
| **Keluar** | Seperti tombol Keluar global | — |

### `/akun/poin`
| Tombol | Aksi | Catatan |
|---|---|---|
| Semua / Didapat / Ditukar | `GET /me/points?type=&page=` | Simpan di query |
| **Muat lebih banyak** | Halaman berikutnya ditambahkan ke daftar | Hilang jika data habis |

### `/akun/profil` & `/akun/keamanan`
| Tombol | Aksi | Berhasil | Gagal |
|---|---|---|---|
| **Ubah** | Mode edit (nama, HP) | — | Email hanya-baca |
| **+ Tambah nomor HP** | Buka mode edit fokus ke HP | — | — |
| **Simpan** | `PATCH /me {full_name, phone_number}` | Toast "Profil diperbarui" | 409 HP dipakai |
| **Batal** | Kembali ke mode lihat | — | — |
| Checkbox **Terima info promo dan poin** | `PATCH /me/consent` (optimistic) | — | Gagal → kembali + toast |
| **Ubah password** | `POST /me/password {old, new}` | Toast + semua sesi lain keluar | 422 password lama salah |
| **Buat password** (akun Google tanpa password) | `POST /me/password {new}` | Toast | — |

### QA Fase 3 (hasil 7 Okt 2026)
- [x] Daftar tanpa HP → berhasil. Daftar dengan HP yang sudah dipakai → ditolak ("Nomor HP sudah dipakai akun lain").
- [x] Lupa password dengan email tidak terdaftar → pesan sama persis dengan email terdaftar (dan hanya email terdaftar yang menerima link).
- [x] Link reset dipakai 2x → kedua kalinya ditolak, dengan tombol "Minta link baru".
- [x] Buka `/akun` tanpa login → `/masuk?redirect=/akun` → setelah login kembali ke `/akun`.
- [x] member_code tampil dan tombol Salin bekerja (label "Disalin").
- [x] Login Google pertama kali → akun member baru + member_code dibuat (simulasi mock); dibatalkan → kembali ke halaman asal + toast.
- [x] Tambahan: validasi daftar & email terdaftar (+ link Masuk), masuk sebelum verifikasi (banner + kirim ulang berjeda 60 detik), verifikasi lewat link (dipakai ulang → tidak berlaku), Ringkasan (saldo, rasio, transaksi terakhir, WhatsApp dengan nama + member_code), riwayat poin (Muat lebih banyak, filter tersimpan setelah refresh), profil (HP bentrok 409, Tambah nomor HP), persetujuan tersimpan, ubah password (lama salah / berhasil / masuk dengan password baru), buat password akun Google, member di `/masuk` → `/akun`, admin di `/masuk` → `/admin`, HP 375 px tanpa scroll horizontal, API mati → pesan di form, konsol bersih.

---

## FASE 4 — Membership Admin (FR-I)

**Status: SELESAI 7 Okt 2026.** Checklist QA di bawah lulus 95/95 (Chrome DevTools Protocol, mode mock). Uji rute mock membership lulus 52/52; regresi Fase 0 (31/31), Fase 1 (85/85), Fase 2 (79/79), dan Fase 3 (52/52) tetap lulus.

### Yang dibangun
| Bagian | Isi |
|---|---|
| Halaman | `views/admin/membership/`: `MemberDaftar` (cari + filter di URL, paginasi), `MemberDetail` (profil + Salin kode, data akun, saldo, sync Qontak + Coba sync ulang, riwayat poin dengan filter jenis, Nonaktifkan/Aktifkan akun), `PoinInput` (alur sketsa "AdminPoin": 1. cari member → 2. catat transaksi, riwayat member di samping) |
| Komponen | `components/admin/member/`: `MemberFilterBar`, `MemberTabel`, `KartuProfilMember`, `DataMember`, `SyncMember`, `RiwayatMember`, `RiwayatPoinAdmin`; `components/admin/poin/`: `CariMember`, `MemberTerpilih`, `PilihanJenis`, `IsianBelanja`, `IsianTukarPoin`, `RingkasanPoin`, `RiwayatRingkas`; `components/admin/StatusSync` (status sync berlabel + ikon) |
| Composable & util | `useFormPoin` (validasi, cek nota, pratinjau poin, ringkasan konfirmasi, simpan, draf saat sesi habis), `useMemberTerpilih` (member terpilih di URL, saldo), `useSyncQontak` (retry + cek status, dipakai juga dashboard), `utils/poin.js`, `formatAngka`, `formatTanggalJam`, `parseRupiah`, `inisialNama` |
| Service | `adminMemberService` (daftar, cari, detail, ubah member; riwayat poin, cek nota, simpan transaksi) |
| Mock | `routesMember.js`; 25 member contoh (tambahan: tanpa HP, nonaktif, belum verifikasi, sync menunggu/gagal) dengan riwayat poin. Transaksi poin dan login terakhir ikut tersimpan di localStorage |

Catatan implementasi:
- Member dikenali dengan `users.id` di semua endpoint admin; `point_transactions.member_id` tetap `member_profiles.id` (dipetakan backend).
- Poin belanja dihitung backend: `floor(total ÷ point_ratio_rupiah)`. Angka di layar hanya pratinjau dari `/settings/public`; bila rasio belum termuat, poin tetap dihitung saat disimpan.
- Nomor nota dirapikan (spasi di ujung dibuang, huruf besar) lalu dicek saat keluar field dan sebelum ringkasan konfirmasi. Server tetap menolak nota ganda (409) dan saldo kurang (422); setelah ditolak, saldo member dimuat ulang karena mungkin baru diubah admin lain.
- REDEEM dikirim sebagai angka positif dan disimpan negatif. Penyesuaian (ADJUST) hanya tampil untuk izin `point.adjust` (Super Admin). Saldo tidak boleh minus.
- Member terpilih disimpan di URL (`?member=<id>`), jadi tetap terpilih setelah refresh. Setelah simpan, form kosong dan member tetap terpilih untuk nota berikutnya.
- Akun nonaktif tidak bisa masuk (sesi yang berjalan ikut berakhir) dan poinnya tidak bisa dicatat; saldo dan riwayat tetap. Mengaktifkan lagi tanpa konfirmasi.
- Menu "Sync Qontak" ikut tampil untuk Admin Membership karena memakai izin `qontak.retry` (halamannya Fase 6).
- Perbaikan aksesibilitas fase sebelumnya: tombol tabel "Ubah", "Pratinjau", "Lihat produk", dan "Lihat semua" tidak lagi terbaca menempel dengan nama item oleh pembaca layar ("UbahGranit…").

### `/admin/member`
| Tombol | Izin | Aksi |
|---|---|---|
| Cari (member_code / email / HP / nama) & filter (persetujuan promo, punya HP, status sync, status akun) | `member.view` | `GET /admin/members?...` (disimpan di URL) |
| **Input poin** | `point.create` | → `/admin/poin/input` |
| **Lihat** | `member.view` | → `/admin/member/:id` |

### `/admin/member/:id`
| Tombol | Izin | Aksi | Catatan |
|---|---|---|---|
| **Input poin untuk member ini** | `point.create` | → `/admin/poin/input?member=<id>` | Member langsung terpilih. Tidak tampil untuk akun nonaktif |
| **Salin** (member_code) | — | Salin ke clipboard | Label "Disalin" |
| **Coba sync ulang** | `qontak.retry` | `POST /admin/qontak/contacts/:id/retry` (`id` = `qontak.id` dari detail member) | Tampil hanya jika FAILED; status dicek sampai SYNCED/FAILED |
| **Nonaktifkan akun** | `member.update` | ConfirmModal → `PATCH /admin/members/:id {is_active:false}` | Member tidak bisa login |
| **Aktifkan akun** | `member.update` | `PATCH /admin/members/:id {is_active:true}` | — |
| Jenis transaksi & paginasi riwayat | `member.view` | `GET /admin/points?user_id=&type=&page=` | Disimpan di URL (`?jenis=`) |

### `/admin/poin/input`
| Tombol / field | Aksi | Berhasil | Gagal / catatan |
|---|---|---|---|
| **Cari** (Enter) | `GET /admin/members/search?q=` | 1 hasil → langsung terpilih. >1 → daftar pilihan | 0 hasil → "Member tidak ditemukan" |
| **Ganti member** | Kosongkan pilihan | Fokus ke kotak cari | — |
| Jenis: **Belanja (EARN)** | Tampilkan field nota + total belanja | Poin dihitung otomatis: `floor(total ÷ point_ratio_rupiah)` | — |
| Jenis: **Tukar poin (REDEEM)** | Tampilkan field jumlah poin + catatan (wajib) | — | Poin > saldo → error sebelum kirim |
| Jenis: **Penyesuaian (ADJUST)** | Field poin (+/−) + alasan (wajib) | — | Izin `point.adjust` (usulan: Super Admin saja) |
| Nomor nota (keluar field) | `GET /admin/points/check-reference?ref=` | "Nomor nota belum pernah diinput" | "Nomor nota ini sudah pernah diinput pada <tanggal> (<kode member>)." |
| **Simpan transaksi** | ConfirmModal ringkasan ("Tambah +4 poin ke TF-000123?") → `POST /admin/points` | Toast, form kosong tapi member tetap terpilih, riwayat & saldo diperbarui | 409 nota duplikat. 422 saldo kurang. Pesan di field, saldo dimuat ulang |
| **Batal** | Kosongkan form | — | — |
| Pindah halaman dengan isian belum disimpan | ConfirmModal "Buang isian transaksi?" | — | Sesi habis → isian disimpan dan dipulihkan setelah masuk lagi |

### QA Fase 4 (hasil 7 Okt 2026)
- [x] Input nota yang sama 2x → kedua kalinya ditolak di client (saat keluar field dan sebelum simpan, beda huruf besar/kecil tetap terdeteksi) DAN di server (409, juga bila nota keburu diinput admin lain setelah dicek).
- [x] Total Rp 4.250.000 dengan rasio Rp 1.000.000 → +4 poin (pratinjau di layar dan hasil server).
- [x] REDEEM melebihi saldo → ditolak di client, dan di server (422) bila saldo keburu berkurang.
- [x] Saldo di `/akun` member berubah setelah admin input poin; riwayat member ikut memuat transaksinya.
- [x] Setiap transaksi tercatat di log aktivitas (`point.create` / `point.adjust`, modul Poin, isi transaksi di `after_data`) dan tampil di "Aktivitas terbaru" dashboard.
- [x] Tambahan: akses per role (Admin Katalog → "Tidak punya akses"); daftar member (cari dengan jeda & Enter, HP format 08…, filter di URL, reset, paginasi, kosong, kembali dengan filter sama); detail (data akun, saldo + rasio, riwayat + filter jenis + paginasi, Coba sync ulang PENDING → SYNCED, nonaktifkan dengan konfirmasi → tidak bisa masuk → aktifkan lagi, id salah → tidak ditemukan); input poin (cari 0/1/banyak hasil, Ganti member, member dari URL, ringkasan konfirmasi + "Periksa lagi", form kosong setelah simpan, Batal, "Buang isian transaksi?", sesi habis → isian dipulihkan, ADJUST khusus Super Admin); gangguan API → ErrorState/pesan + Coba lagi; HP 375 px tanpa scroll horizontal; konsol bersih.

---

## FASE 5 — Konten (FR-N, FR-O, FR-P) + Beranda & Landing 3D

**Status: SELESAI 7 Okt 2026.** Checklist QA di bawah lulus 217/217 (Chrome DevTools Protocol, mode mock + API mati). Uji rute mock konten lulus 40/40; regresi Fase 0 (31/31), Fase 1 (85/85), Fase 2 (79/79), Fase 3 (52/52), dan Fase 4 (95/95; satu cek waktu tunggu sempat gagal saat dijalankan bersama, lulus 2× saat diulang sendiri) tetap lulus. Mock Fase 1–4: 17, 38, 33, 52. Skrip `run-qa-fase-5.mjs` (uji rute mock tanpa browser, dibuat di sesi lain) juga lulus 26/26.

### Yang dibangun
| Bagian | Isi |
|---|---|
| Halaman publik | `views/public/`: `Beranda` (hero rumah 3D + seksi sketsa "Main"; hero diganti hero banner 30 Okt 2026, KEPUTUSAN #16), `ArtikelDaftar` (`/artikel`, Muat lebih banyak), `Promo` (`/promo`, chip Semua/Promo/Event di `?jenis=`), `ArtikelDetail` (`/artikel/:slug`, pratinjau draf `?preview=`, 404, WhatsApp, artikel lainnya), `InfoToko` (`/info-toko`). `Home.vue` sementara dihapus |
| Halaman admin | `views/admin/konten/`: `BannerDaftar` (per posisi, status, ↑ ↓, switch Aktif, Perpanjang), `BannerForm` (baru/ubah, `#periode`), `ArtikelDaftar` (filter di URL, paginasi), `ArtikelForm` (simpan draf, pratinjau, terbitkan, batalkan terbit), `Pengaturan` (WA, alamat, jam buka, peta, rasio poin) |
| Komponen | `components/beranda/`: `HeroBanner` (sejak 30 Okt 2026; menggantikan `HeroRumah`, `RumahTigaDimensi`, `geometriRumah.js`, `bagianRumah.js` yang dihapus), `BannerPromo`, `BannerSlider`, `KategoriTahap`, `ProdukUnggulan`, `PromoEvent`, `InspirasiRuangan`, `BrandPilihan`, `Keunggulan`, `AjakanMember`, `JudulSeksi`, `PesanGagal`; `components/konten/`: `KartuArtikel`, `DaftarArtikel`, `IsiArtikel`, `ArtikelLainnya`; `components/admin/konten/`: `StatusBanner`, `StatusArtikel`, `ArtikelFilterBar`, `VideoField` (30 Okt 2026) |
| Composable & util | `useMuat`, `useDaftarBanner`, `useFormBanner`, `useFormArtikel`; `utils/artikel.js`, `utils/banner.js` |
| Service | `kontenService` (banner, artikel; bagian rumah dihapus 30 Okt 2026), `adminKontenService` (banner, artikel, pengaturan) |
| Mock | `routesKonten.js` (`data/beranda.js` dihapus 30 Okt 2026); banner (tayang, terjadwal, berakhir, nonaktif) dan artikel (promo, event, artikel, satu draf) di `data/konten.js`, ikut tersimpan di localStorage |

Catatan implementasi:
- Status banner (`TAYANG` / `TERJADWAL` / `BERAKHIR` / `NONAKTIF`) dihitung server dari `is_active` + periode WIB. Form memakai tanggal saja: mulai 00.00, selesai 23.59 WIB, jadi tanggal selesai boleh sama dengan tanggal mulai (banner satu hari). Form menampilkan perkiraan status sebelum disimpan.
- Urutan banner berlaku per posisi (`HOME_HERO`, `HOME_SLIDER`, `HOME_SIDE`); ↑ ↓ memberi nomor ulang 1..n dan hanya mengirim yang berubah.
- "Simpan draf" tidak mengubah status. Artikel yang sudah terbit disimpan dengan "Simpan perubahan" dan tetap terbit (tidak ikut jadi draf); menarik artikel lewat "Batalkan terbit". "Terbitkan" menyimpan isi dulu bila ada perubahan atau artikel masih baru.
- Isi artikel ditulis teks biasa: baris kosong = paragraf baru, `## ` = subjudul, `- ` = daftar. Dirender per blok tanpa `v-html` (aman dari skrip sisipan).
- Rasio poin hanya bisa diubah pemegang izin baru `setting.point_ratio` (usulan: Super Admin), dengan konfirmasi. Admin lain melihatnya sebagai teks.
- Setelah pengaturan disimpan, store `settings` dimuat ulang, jadi header, footer, Info toko, dan semua tombol WhatsApp langsung memakai data baru.
- (Sampai 30 Okt 2026) Hero 3D: semua warna dari token `--color-rumah-*` (`@theme static` di main.css), animasi 0,8 detik (pengecualian di MASTER.md) dan mati dengan "kurangi animasi"; slider banner tidak berputar otomatis bila "kurangi animasi" aktif.
- Setiap seksi Beranda memuat datanya sendiri; bila API gagal, seksi itu menampilkan pesan + Coba lagi dan hero tetap tampil dengan teks bawaan + "Jelajahi katalog".

### `/admin/banner`
| Tombol | Izin | Aksi | Gagal / catatan |
|---|---|---|---|
| **+ Buat banner** | `banner.create` | → form | — |
| Form: **Unggah gambar** | `banner.create` | `POST /admin/uploads` | Wajib |
| Form posisi **Hero Beranda**: Subjudul, **Unggah video**, Label & Tautan tombol, Label & Tautan tombol kedua | `banner.create` | `POST /admin/uploads` (`jenis=video`) | 30 Okt 2026. Gambar tetap wajib (cadangan video). Video MP4/WebM ≤ 20 MB (mock 3 MB); daftar menandai "Dengan video" |
| Form: tanggal mulai & selesai | — | — | Selesai sebelum mulai → error (tanggal saja: selesai 23.59, jadi sama dengan mulai = banner satu hari) |
| Form: **Simpan** | `banner.create/update` | `POST/PUT /admin/banners` | — |
| **Ubah** / **Perpanjang** | `banner.update` | → form | — |
| Switch **Aktif**, ↑ ↓ urutan | `banner.update` | `PATCH` | — |
| Label status | — | Dihitung: Terjadwal / Tayang / Berakhir / Nonaktif | — |

### `/admin/artikel`
| Tombol | Izin | Aksi | Catatan |
|---|---|---|---|
| Filter tipe (Artikel, Promo, Event, Halaman) & status | `article.view` | Query | — |
| **+ Tulis artikel** | `article.create` | → editor | — |
| **Simpan draf** | `article.create/update` | `POST` (baru = `DRAFT`) / `PUT` (status tetap) | Artikel terbit: tombol jadi **Simpan perubahan**, tetap terbit |
| **Pratinjau** | `article.view` | Buka `/artikel/[slug]?preview=<token>` di tab baru | — |
| **Terbitkan** | `article.publish` | ConfirmModal → `PATCH {status:'PUBLISHED'}` (`published_at` = sekarang) | — |
| **Batalkan terbit** | `article.publish` | ConfirmModal → `PATCH {status:'DRAFT'}` | — |

### `/admin/pengaturan`
| Field / tombol | Izin | Catatan |
|---|---|---|
| Alamat, jam buka, link Google Maps | `setting.update` | Nomor WhatsApp hanya ditampilkan; perubahan hanya lewat konfigurasi server (KEPUTUSAN #21) |
| Rasio poin (`point_ratio_rupiah`) | Super Admin | Peringatan: "Hanya berlaku untuk transaksi baru" |
| **Simpan** | `setting.update` | `PUT /admin/settings` |

### Publik: Beranda `/` (hero banner + section Beranda), `/artikel`, `/promo`, `/info-toko`
| Tombol | Aksi | Catatan |
|---|---|---|
| Hero: tombol utama (bawaan **Jelajahi katalog**) | → tautan banner, bawaan `/katalog` | Label & tautan dari banner `HOME_HERO` |
| Hero: tombol kedua (mis. **Lihat promo**) | → tautan kedua banner | Tampil bila admin mengisi label + tautan |
| Hero: **Lihat selengkapnya** | Gulir ke seksi berikutnya + fokus pindah | Hero setinggi layar |
| Hero: **Jeda video / Putar video** | Jeda/putar video latar | Hanya bila banner punya video; "kurangi animasi" → tidak berputar sendiri |
| ~~Hero 3D: bagian rumah, Rakit/Pisahkan, Putar rumah, Atur ulang, Lihat produk~~ | Dihapus 30 Okt 2026, KEPUTUSAN #16 | — |
| Banner: ‹ › titik | Ganti slide | Autoplay berhenti saat hover/fokus |
| Kartu kategori / produk unggulan / promo / ruangan / brand | Navigasi ke halaman terkait | Pakai komponen dari Fase 1 |
| **Daftar sekarang** / **Masuk dengan Google** | → `/daftar` / OAuth | Disembunyikan jika sudah login |
| Info toko: **Buka di Google Maps** | Tab baru | — |
| Info toko: **Salin alamat** | Clipboard + toast | — |
| Info toko: **Chat WhatsApp** | `generalLink()` | — |

### QA Fase 5 (hasil 7 Okt 2026)
- [x] Banner dengan `end_at` kemarin → tidak tampil di beranda, status "Berakhir" di admin (begitu juga banner terjadwal dan nonaktif). Setelah diperpanjang → Tayang dan tampil lagi.
- [x] Artikel DRAFT tidak bisa dibuka publik tanpa token pratinjau (token salah juga 404); dengan token benar tampil bertanda "Pratinjau draf".
- [x] Ganti nomor WA di pengaturan → semua tombol WhatsApp memakai nomor baru (Beranda 11 tautan, Info toko, detail artikel, katalog 27 tautan, detail produk), juga setelah halaman dimuat ulang.
- [x] (Hero 3D, dihapus 30 Okt 2026) Hero 3D lancar di HP kelas menengah: 375 px dengan CPU diperlambat 4× → rata-rata 56–60 fps saat animasi rakit/pisah (Chrome headless, render software). Dengan "kurangi animasi" aktif → bagian rumah langsung ke posisi akhir, slider tidak berganti otomatis.
- [x] Tambahan: hero (pilih bagian + kartu + "Lihat produk" ke kategori dari API, rakit/pisah, putar, atur ulang); slider (‹ › titik, aria-current, jeda, autoplay 6 detik berhenti saat fokus); banner samping + kartu Trans Family terlihat utuh; urutan seksi sesuai MASTER.md; semua gambar termuat; ajakan member hilang setelah masuk; `/artikel` + Muat lebih banyak; `/promo` + chip di URL; detail (breadcrumb, isi per blok, pesan WhatsApp berisi judul + tautan, artikel lainnya, isi berisi HTML tetap jadi teks); Info toko (peta di tab baru, Salin alamat → clipboard + toast); akses per role; admin banner (status, nonaktifkan dengan konfirmasi, ↑ ↓ tersimpan + fokus, Perpanjang `#periode`, validasi, unggah gambar, Buang perubahan?); admin artikel (filter di URL, slug otomatis, slug kembar 409, simpan draf → pratinjau, terbitkan/batalkan terbit dengan konfirmasi, tanpa izin `article.publish`); pengaturan (validasi WA & peta, rasio khusus Super Admin + konfirmasi); log aktivitas; API mati → pesan + Coba lagi; 375/768/1024 px tanpa scroll horizontal; konsol bersih.

---

## FASE 6 — Sistem Admin & Integrasi Qontak (FR-G lanjutan, FR-R–FR-U)

**Status: SELESAI 7 Okt 2026.** Checklist QA di bawah lulus 120/120 (Chrome DevTools Protocol, mode mock). Uji rute mock sistem lulus 54/54; regresi browser Fase 0 (31/31), Fase 1 (85/85), Fase 2 (79/79), Fase 3 (52/52), Fase 4 (95/95; satu cek disesuaikan dengan label aksi baru di dashboard), dan Fase 5 (217/217) tetap lulus, begitu juga uji mock Fase 1–5.

### Yang dibangun
| Bagian | Isi |
|---|---|
| Halaman | `views/admin/akses/`: `PenggunaDaftar` (cari + filter role/status di URL, Tambah admin, Ubah role, Reset password, switch Aktif), `Role` (daftar role + matriks izin, role terpilih di `?role=`); `views/admin/audit/LogAktivitas` (filter admin/modul/tanggal di URL, drawer perbandingan data); `views/admin/integrasi/SyncQontak` (chip status + jumlah, cari, Coba lagi per baris & semua yang gagal, drawer log sync) |
| Komponen | `components/admin/akses/`: `PenggunaFilterBar`, `PenggunaTabel`, `TambahAdminModal`, `UbahRoleModal`, `DaftarRole`, `MatriksIzin`; `components/admin/audit/`: `LogFilterBar`, `LogTabel`, `DetailLog`, `PerbandinganData`; `components/admin/integrasi/`: `KontakFilterBar`, `KontakTabel`, `LogSyncDrawer`; `components/admin/SkeletonTabel` |
| Composable & util | `useDaftar` (daftar ber-paginasi + muat ulang diam-diam), `useMatriksIzin` (centang, ketergantungan, pilih semua baris/kolom, perubahan); `utils/log.js` (label aksi & kolom, format nilai, perbandingan sebelum/sesudah) |
| Service | `adminAksesService` (admin, role, izin), `adminLogService`; `qontakService` (Fase 2) ditambah `listContacts`, `retryFailed`, `getContactLogs` |
| Mock | `routesSistem.js`, `routesQontak.js`, `antreanQontak.js` (proses sync di background, bisa "dimatikan"), `data/sistem.js` (tabel `roles` & `permissions`). Izin role ikut tersimpan di localStorage |

Catatan implementasi:
- Matriks izin: baris = menu (dikelompokkan seperti sidebar), kolom Lihat / Tambah / Ubah / Lainnya + Semua. Muka 1 tidak menghapus data, jadi tidak ada kolom Hapus. "Lainnya" = aksi khusus (Terbitkan, Penyesuaian poin, Rasio poin, Lihat & coba ulang sync).
- Setiap izin di satu menu butuh izin yang membuat menunya terlihat (`menus.izin`): mencentang Tambah/Ubah ikut mencentang Lihat, mencabut Lihat mencabut seluruh baris. Server juga menolak (422) kombinasi yang tidak lengkap dan body tanpa daftar izin.
- Super Admin terkunci (semua izin) supaya selalu ada yang bisa mengatur akses. Perubahan izin berlaku saat admin terkait memuat ulang halaman; sebelum itu server sudah menolak aksinya (403).
- Pengguna admin: satu role per admin di form (API menerima daftar `roles`). Akun sendiri tidak bisa dinonaktifkan atau diganti role-nya; Super Admin aktif terakhir juga dilindungi. Menonaktifkan pakai konfirmasi, mengaktifkan tidak. "Reset password" mengirim link sekali pakai (1 jam) ke email admin, sama seperti lupa password.
- Log aktivitas memakai `before_data` / `after_data`: hanya kolom yang berubah, ditambah penanda data (mis. SKU dan nama produk) karena `activity_logs` tidak punya kolom id data. Log produk, kategori, brand, dan ruangan (Fase 2) kini ikut menyimpan data ini. Nilai ditampilkan dengan format yang sama seperti di halaman lain (Rupiah, Ya/Tidak, label status).
- Sync Qontak: urutan kontak tetap (terdaftar terbaru) supaya baris tidak melompat saat status berubah. Selama ada kontak PENDING daftar dimuat ulang diam-diam tiap 5 detik (berhenti saat tab tidak terlihat). "Coba lagi" per baris dipantau tiap 1 detik seperti di detail member.
- Mode mock: tombol "Matikan Qontak" di `/dev/components` membuat pengiriman berikutnya gagal (503). Pendaftaran member tetap berhasil.
- Modal dasar mendapat varian `side` (drawer kanan setinggi layar, selebar layar di HP) untuk detail log dan log sync. Kolom Aksi di "Aktivitas terbaru" dashboard memakai label yang sama dengan log aktivitas.

### `/admin/pengguna` (Super Admin)
| Tombol | Aksi | Catatan |
|---|---|---|
| **+ Tambah admin** | Modal (nama, email, role, password sementara) → `POST /admin/users` | Email duplikat → 409 |
| **Ubah role** | `PUT /admin/users/:id/roles` | Tidak bisa mencabut role Super Admin dari diri sendiri |
| **Reset password** | ConfirmModal → `POST /admin/users/:id/reset-password` | — |
| Switch **Aktif** | ConfirmModal → `PATCH` | — |

### `/admin/role`
| Tombol | Aksi | Catatan |
|---|---|---|
| **Atur izin** | Matriks menu × aksi (view/create/update/delete) berupa checkbox | Role SUPER_ADMIN terkunci |
| **Pilih semua** per baris / kolom | Centang massal | — |
| **Simpan izin** | `PUT /admin/roles/:id/permissions` | Berlaku saat admin terkait memuat ulang halaman |

### `/admin/log-aktivitas`
| Tombol | Aksi |
|---|---|
| Filter admin / modul / rentang tanggal | Query `GET /admin/activity-logs` |
| **Lihat detail** | Drawer: perbandingan `before_data` vs `after_data` |

### `/admin/sync-qontak`
| Tombol | Aksi | Catatan |
|---|---|---|
| Filter status (PENDING, SYNCED, FAILED) | Query | — |
| **Coba lagi** (per baris) | `POST /admin/qontak/contacts/:id/retry` | Status berubah otomatis (polling 5 detik selama PENDING) |
| **Coba lagi semua yang gagal** | ConfirmModal → `POST /admin/qontak/retry-failed` | — |
| **Lihat log** | Drawer `sync_logs` (request, response, status code, retry_count) | — |

### QA Fase 6 (hasil 7 Okt 2026)
- [x] Hapus izin `product.update` dari Admin Katalog → selagi halaman belum dimuat ulang, simpan produk ditolak API (403, toast) dan harga tidak berubah; setelah reload tombol Ubah hilang, switch "Tampil di website" nonaktif, `/admin/produk/1` → "Tidak punya akses", dan API menolak 403. Izin lain tetap ada.
- [x] Ubah harga produk → log aktivitas menampilkan harga lama dan baru (Rp 89.500 → Rp 92.500, bertanda "Berubah"), hanya kolom yang berubah + SKU dan nama.
- [x] Qontak dimatikan → registrasi member tetap berhasil, kontak berstatus FAILED (503, dikirim dengan email saja karena tanpa HP), log sync menampilkan isi permintaan & respons; setelah Qontak hidup, "Coba lagi" → SYNCED.
- [x] Tambahan: akses per role (menu Sistem & "Tidak punya akses"); pengguna admin (filter di URL, validasi per field, email kembar 409, password acak + Salin, ubah role, nonaktifkan dengan konfirmasi → tidak bisa masuk → aktifkan lagi, reset password → email, akun sendiri terkunci, admin baru masuk dengan menu sesuai role); role (Super Admin terkunci, ketergantungan izin, indeterminate, Buang perubahan? saat ganti role/menu, Batal, simpan + muat ulang, tercatat di log); log (drawer kanan, Esc mengembalikan fokus, filter modul/admin/tanggal + validasi rentang, EmptyState); sync Qontak (chip + jumlah, Coba lagi semua dengan konfirmasi, polling 5 detik mengubah status sendiri, fokus setelah retry, tercatat di log); API gangguan → ErrorState + Coba lagi; 375/768/1024 px tanpa scroll horizontal, drawer selebar layar di HP, area sentuh matriks ≥ 44 px; konsol bersih.

---

## Setelah Fase 6 — Penyempurnaan (8 Okt 2026)

Diminta pemilik proyek setelah semua fase selesai. Semua tanpa library baru.

| Bagian | Isi |
|---|---|
| Header publik menempel | `AppHeader` `sticky`. Mulai 1024 px seluruh header menempel; di bawahnya baris alamat ikut tergulir (yang menempel: logo, Masuk/Daftar, Kategori, cari). Tinggi bagian yang menempel → `scroll-padding-top` supaya lompatan gulir/anchor tidak tertutup |
| Animasi situs publik | View Transitions (pindah halaman), `v-reveal` (scroll reveal + stagger), tombol ditekan `scale(0.97)`, modal/sheet dari arah logis. Panel admin tidak dianimasikan; semua mati saat "kurangi animasi". Rincian di `design-system/transhome/MASTER.md` |
| Brand untuk jumlah besar | Halaman `/brand` (A–Z, cari, lompat huruf); "Brand pilihan" di Beranda berisi semua brand aktif, logo berjalan (jeda, berhenti saat disorot/difokus; diganti susunan lingkaran 2 Nov 2026, KEPUTUSAN #19); filter Brand katalog punya kotak cari bila > 10 opsi; admin: cari brand, pilihan brand berupa combobox, hapus permanen brand nonaktif yang tidak dipakai produk (`brand.delete`) |

| Harga boleh disembunyikan (15 Okt 2026) | Per produk switch "Tampilkan harga di website" di form produk dan di modal **Ubah harga** (tombol baru di daftar produk admin: harga, satuan, tampil/sembunyi, tanpa membuka form lengkap). Harga disembunyikan → "Harga via WhatsApp" + tombol "Tanya harga" di kartu, baris katalog, detail produk, dan bar bawah HP; pesan WhatsApp "saya ingin menanyakan harga produk ini". Tabel admin menandai "Disembunyikan" / "Belum ada harga". KEPUTUSAN #15 |
| Hero banner + halaman Brand (30 Okt 2026) | Hero rumah 3D dihapus, diganti `HeroBanner` sesuai wireframe XLSX: gambar/video dari banner posisi baru **Hero Beranda** (`HOME_HERO`) + judul, subjudul; hero setinggi layar dengan tombol utama, tombol kedua opsional (admin), dan Lihat selengkapnya (tombol WhatsApp di hero dihapus 2 Nov 2026); video tanpa suara berulang, tombol Jeda/Putar, tidak berputar sendiri saat "kurangi animasi", tidak dimuat saat hemat data. Admin: kolom Subjudul, Video hero (`VideoField`), Label & Tautan tombol. Seksi Beranda lain tidak berubah. `/brand` mengikuti sheet "Brand": "N Brand Partner", Cari (tombol), filter kategori, A–Z, kartu 6/3/2 kolom, 24 per halaman, filter di URL. KEPUTUSAN #16 |
| SEO halaman Brand (30 Okt 2026) | `useSeo` diperluas (canonical, robots, Open Graph, `image`, JSON-LD) tanpa library. `/brand` dan `/brand/:slug`: judul/deskripsi dinamis, BreadcrumbList + CollectionPage/ItemList/Brand, hasil cari/filter `noindex, follow` + canonical ke halaman bersih, brand tanpa produk `noindex`. Paginasi `/brand` dan katalog berupa `<a href>` (klik tetap tanpa muat ulang). Breadcrumb detail brand jadi Beranda › Brand › nama. `public/robots.txt`, og:site_name/og:locale di `index.html`, `VITE_SITE_URL`. KEPUTUSAN #17 |
| Brand pilihan berbentuk lingkaran (2 Nov 2026) | Desain dari contoh pemilik: judul tengah "Siap untuk menjelajah lebih banyak Brand?" + "Telusuri N brand", 18 lingkaran logo brand (HP 10), brand berlogo dulu. Logo diatur di admin › Brand (petunjuk baru di kolom Logo brand). Logo berjalan + tombol Jeda dihapus. Tombol WhatsApp di hero juga dihapus. KEPUTUSAN #19 |
| Media sosial di footer (2 Nov 2026) | Kolom pertama footer: "Ikuti kami" + tombol Instagram (instagram.com/transhome) dan TikTok (tiktok.com/@transhome_id), buka di tab baru. Alamat ditulis tetap di `AppFooter.vue`; bila nanti perlu diubah dari admin, tambahkan kunci `instagram_url` / `tiktok_url` di `site_settings` |
| Artikel panduan isi rumah (2 Nov 2026) | Materi pemilik: artikel utama "5 Langkah Mudah Mengisi Rumah: Dari Lantai Granit, Kamar Mandi, hingga Sofa" + 5 artikel pilar (lantai, dinding, kamar mandi, lampu & elektronik, sofa). Menggantikan 4 artikel contoh lama; pilar sofa berstatus draf (contoh alur pratinjau → terbitkan). Isinya baru ada di data contoh (`data/konten.js`): setelah backend siap, Admin Konten perlu memasukkannya lewat admin › Artikel, atau backend men-seed tabel `articles` dari teks yang sama |
| Fitur Jasa (2 Nov 2026) | Halaman publik `/jasa` (cara pesan 3 langkah + kartu jasa) dan `/jasa/:slug` (penjelasan, area, bahan terkait, kotak konsultasi menempel). Harga tidak ditampilkan: semua jasa bertuliskan "Harga pasti murah". Daftar contoh: pasang keramik & granit, pengecatan, perlengkapan kamar mandi, AC & elektronik, service & perbaikan, pintu & jendela, antar & rakit perabot. Konsultasi hanya lewat WhatsApp (`jasaLink`, pesan berisi nama jasa + tautan + produk asal + isian lokasi/kebutuhan). Menu "Jasa" di header, menu HP, footer; seksi Beranda "Butuh tukang? Pakai jasa Transhome" (3 teratas, setelah Inspirasi ruangan); kotak "Butuh pemasangan?" di ringkasan detail produk. Admin Konten + Super Admin: menu Konten › Jasa (daftar, urutan ↑ ↓, aktif dengan konfirmasi, tambah/ubah, kategori produk terkait). SEO: canonical, JSON-LD `Service` (tanpa `Offer`). Parser isi (`blokIsi`) kini memisah subjudul `##` yang langsung diikuti baris lain. KEPUTUSAN #20; kontrak di "Detail kontrak jasa" |
| Perbaikan akses login admin (2 Nov 2026) | Bug: browser yang masih masuk sebagai member tidak bisa membuka login admin (`/admin` → "Tidak punya akses", `/admin/masuk` → `/akun`). Sekarang: bukan admin → diarahkan ke `/admin/masuk?redirect=…` dengan pemberitahuan akun member akan diganti; admin tanpa izin halaman → "Tidak punya akses" dengan tombol **Masuk dengan akun lain** (keluar lalu ke login yang sesuai, `?dari=`) dan **Kembali ke dashboard**. Akun admin tetap hanya dibuat Super Admin (tidak ada pendaftaran admin sendiri) |

Hasil uji: header 33/33, animasi 22/22, brand 41/41, harga 34/34, hero banner + `/brand` + SEO 64/64, uji asap SEO semua halaman 23/23, video hero 2,87 MB tersimpan & tampil setelah muat ulang (Chrome DevTools Protocol); regresi browser Fase 0–6 tetap lulus (harga: lihat catatan uji di bawah).

Katalog produk sudah memakai paginasi server (24 per halaman), jadi jumlah produk dibatasi backend, bukan browser. Untuk ribuan produk backend perlu index di kolom pencarian & filter (`sku`, `name`, `category_id`, `brand_id`, `is_active`).

## Cara Mulai Satu Fase
Buka folder `frontend/` di Antigravity atau Claude Code, lalu kirim:
```
Baca CLAUDE.md, docs/KEPUTUSAN.md, dan docs/FRONTEND_BUILD_PLAN.md.
Kerjakan FASE 0 saja sesuai dokumen. Pakai data mock. Berhenti setelah
selesai dan laporkan hasil checklist QA Fase 0.
```
Untuk fase berikutnya, ganti nomor fasenya.

---

## Lampiran — Usulan Kontrak API untuk Backend Go

| Method | Endpoint | Dipakai di |
|---|---|---|
| GET | `/categories?tree=1` | Mega menu, filter |
| GET | `/products` · `/products/:slug` · `/products/suggest` | Fase 1 |
| GET | `/brands` · `/brands/:slug` · `/rooms` · `/rooms/:slug` | Fase 1, `/brands` juga Beranda & `/brand` |
| GET | `/banners?position=HOME_SLIDER` · `/articles` · `/articles/:slug` · `/settings/public` | Fase 0 (settings), Fase 5 |
| POST | `/auth/register` · `/auth/login` · `/auth/logout` · `/auth/forgot-password` · `/auth/reset-password` · `/auth/email/resend` · `/auth/email/verify` | Fase 3 |
| GET | `/auth/google/redirect` · `/auth/google/callback` | Fase 3 |
| GET / PATCH | `/me` · `/me/consent` · POST `/me/password` · GET `/me/points` | Fase 3 |
| POST | `/admin/auth/login` | Fase 2 |
| GET | `/admin/menus` · `/admin/dashboard` | Fase 2 |
| CRUD | `/admin/products` (+ `check-sku`), `/admin/categories`, `/admin/brands`, `/admin/rooms` | Fase 2 |
| POST | `/admin/uploads` | Fase 2, 5 |
| GET / PATCH | `/admin/members` · `/admin/members/:id` · `/admin/members/search` | Fase 4 |
| GET / POST | `/admin/points` · `/admin/points/check-reference` | Fase 4 |
| CRUD | `/admin/banners` · `/admin/articles` · PUT `/admin/settings` | Fase 5 |
| GET | `/services` · `/services/:slug` | Jasa (2 Nov 2026) |
| CRUD | `/admin/services` (tanpa hapus; nonaktifkan lewat PATCH) | Jasa (2 Nov 2026) |
| GET / POST / PUT / PATCH | `/admin/users` · `/admin/users/:id/roles` · `/admin/users/:id/reset-password` · `/admin/roles` · `/admin/permissions` · `/admin/roles/:id/permissions` · `/admin/activity-logs` (+ `/options`) | Fase 6 |
| GET / POST | `/admin/qontak/contacts` · `/:id` · `/:id/retry` · `/:id/logs` · `/admin/qontak/retry-failed` | Fase 2, 6 |

### Detail kontrak katalog (dipakai mock Fase 1)
`GET /products` — hanya produk aktif.

| Parameter | Isi |
|---|---|
| `q` | Cari di nama, SKU, dan nama brand |
| `category` | Slug kategori; kategori induk ikut menyertakan produk anak-anaknya |
| `brand`, `stock`, `room` | Banyak nilai dipisah koma: `brand=merek-a,merek-b`, `stock=TERSEDIA,SISA_STOK` |
| `attr[<kunci>]` | Atribut `specifications`, banyak nilai dipisah koma: `attr[ukuran]=40 × 40 cm` |
| `min`, `max` | Rentang `price_general` |
| `exclude` | ID produk yang dikecualikan (produk serupa) |
| `featured=1` | Hanya produk unggulan |
| `sort` | `terbaru` (default) · `harga_terendah` · `harga_tertinggi` · `nama` |
| `page`, `per_page` | Default 24 |

Respons:
```json
{
  "data": [{ "id": 1, "sku": "…", "name": "…", "slug": "…", "price_general": "89500.00", "unit_sale": "dus", "min_order": 1,
             "stock_status": "TERSEDIA", "stock_qty_label": null, "is_featured": true,
             "brand": { "id": 1, "name": "…", "slug": "…" }, "category": { "id": 21, "name": "…", "slug": "…", "parent_id": 2 },
             "primary_image_url": "…" }],
  "meta": {
    "page": 1, "per_page": 24, "total": 37, "total_pages": 2,
    "facets": {
      "brands": [{ "slug": "merek-a", "name": "Merek A", "count": 6 }],
      "stock": [{ "value": "TERSEDIA", "count": 21 }],
      "rooms": [{ "slug": "dapur", "name": "Dapur", "count": 11 }],
      "attributes": [{ "key": "ukuran", "values": [{ "value": "40 × 40 cm", "count": 4 }] }],
      "price": { "min": 9500, "max": 3250000 }
    }
  }
}
```
- Angka faset **disjungtif**: angka tiap opsi menghitung semua filter kecuali filter dari grupnya sendiri (memilih Merek A tidak membuat angka Merek B jadi 0). Opsi yang sedang dipilih tetap dikirim walau angkanya 0.
- `facets.price` dihitung tanpa filter `min`/`max`, supaya rentangnya tidak menyempit oleh filter harga itu sendiri.
- `GET /products/:slug`: seperti item di atas + `description`, `specifications`, `datasheet_pdf_url`, `meta_title`, `meta_description`, `images` (urut `sort_order`), `rooms` (`id`, `name`, `slug`), dan `category.parent`. Produk nonaktif → 404.
- `GET /brands/:slug` dan `GET /rooms/:slug` → objek + `product_count`; `GET /rooms` → daftar ruangan aktif berurutan + `product_count`.
- Frontend tetap membuang field internal Master Barang di `specifications`, tapi sebaiknya backend tidak mengirimkannya sama sekali.

### Detail kontrak admin (dipakai mock Fase 2)
Semua `/admin/*`: tanpa sesi → 401, bukan admin atau tidak punya izin → 403.

| Endpoint | Isi |
|---|---|
| `GET /admin/menus` | Daftar datar `{ id, parent_id, name, code, route, icon, sort_order }` yang boleh dilihat user. Grup = baris tanpa `route`; hanya dikirim bila ada anak yang terlihat |
| `GET /admin/dashboard` | `{ summary: [{ key, label, value, link }], attention: [{ type, id, title, … }], attention_total: { <type>: n }, logs: [{ id, admin_name, action, module, created_at }], logs_scope: 'all' \| 'mine' }`. Isi sesuai izin; `attention` maks 3 per jenis (`stok_habis`, `tanpa_foto` + `sku`; `banner_berakhir` + `end_at`; `sync_gagal` + `sync_status`, `retry_count`, `status_code`). Tanpa `log.view` → hanya log milik sendiri |
| `GET /admin/products` | Param `q, category (id, induk ikut anak), brand (id), stock, active (1/0), foto (1/0: punya foto utama), room (id), page, per_page (20)`. Item: kolom ringkas + `category {id,name}`, `brand`, `primary_image_url`, `image_count` |
| `GET /admin/products/:id` | Semua kolom + `images`, `room_ids`, `rooms`, `brand`, `category` (+ `parent`), termasuk produk nonaktif |
| `GET /admin/products/check-sku?sku=&exclude=` | `{ available }`, tidak peka huruf besar |
| `POST` / `PUT /admin/products(/:id)` | Kolom produk + `specifications` (objek) + `images: [{ image_url, is_primary, sort_order }]` + `room_ids`. 422 per field (`images` bila aktif tanpa foto utama), 409 `sku` / `slug` |
| `PATCH /admin/products/:id` | `{ is_featured }`, `{ is_active }` (422 tanpa foto utama), `{ stock_status, stock_qty_label }` |
| `GET` / `POST` / `PUT` / `PATCH /admin/categories`, `/admin/brands`, `/admin/rooms` | Daftar berisi yang nonaktif + `product_count`. Simpan `{ name, slug, is_active, parent_id \| image_url \| logo_url \| image_cover }`; 409 slug; kategori maks 2 tingkat (422 `parent_id`); cover ruangan wajib (422). `PATCH { is_active }` / `{ sort_order }` |
| `POST /admin/uploads` | multipart `file` + `jenis` (`gambar`: JPG/PNG/WebP ≤ 2 MB, `dokumen`: PDF ≤ 10 MB, `video`: MP4/WebM ≤ 20 MB sejak 30 Okt 2026) → `{ url, name, size, type }`; 422 `errors.file` |
| `POST /admin/qontak/contacts/:id/retry` | 202 `{ sync_status: 'PENDING', … }`; status akhir dibaca lewat `GET /admin/qontak/contacts/:id` (frontend mengecek tiap 1 detik, maks 15 kali) |

Setiap perubahan admin mencatat `activity_logs` (`action` = kode izin, mis. `product.update`).

### Detail kontrak akun (dipakai mock Fase 3)
| Endpoint | Isi |
|---|---|
| `POST /auth/login` | `{ email, password }` → `{ token, user, permissions }`. 401 salah; 403 `code: EMAIL_BELUM_VERIFIKASI` bila email belum diverifikasi |
| `POST /auth/register` | `{ full_name, email, password (≥ 8), phone_number (62…, boleh null), communication_consent }` → 201 + email verifikasi. 409 `errors.email` / `errors.phone_number` |
| `POST /auth/email/resend` | `{ email }` → selalu 200 dengan pesan yang sama |
| `POST /auth/email/verify` | `{ token }` dari link `/verifikasi?token=` → `{ email }`. Token sekali pakai; salah/terpakai → 422 `code: TOKEN_TIDAK_VALID` |
| `POST /auth/forgot-password` | `{ email }` → selalu 200, pesan "Jika email terdaftar, link reset sudah dikirim."; link email ke `/reset-password?token=` (berlaku 1 jam) |
| `POST /auth/reset-password` | `{ token, password }` → 200. Token sekali pakai; salah/terpakai/kedaluwarsa → 422 `code: TOKEN_TIDAK_VALID` |
| `GET /auth/google/redirect?return_to=<url>` | Halaman penuh ke Google. Setelah OAuth backend kembali ke `return_to` + `#token=<token>` (tambah `&new=1` bila member baru dibuat), atau `return_to` + `?error=access_denied` bila dibatalkan |
| `GET /me` | `user` berisi juga `member_since`, `email_verified`, `has_password` (false untuk akun Google tanpa password), `google_connected` |
| `PATCH /me` | `{ full_name, phone_number }` → `{ user }`. 409 HP dipakai akun lain |
| `PATCH /me/consent` | `{ communication_consent }` → `{ user }` (backend mengisi `consent_at`) |
| `POST /me/password` | `{ old_password, new_password }` (akun Google tanpa password: `{ new_password }`). 422 `errors.old_password` bila salah; backend mengeluarkan sesi lain |
| `GET /me/points?type=&page=&per_page=` | `type` = `EARN` / `REDEEM` / kosong. `meta: { balance, page, per_page, total, total_pages }`, urut terbaru |

### Detail kontrak membership admin (dipakai mock Fase 4)
Member dikenali dengan `id` = `users.id` (bukan `member_profiles.id`).

| Endpoint | Isi |
|---|---|
| `GET /admin/members` | Param `q` (kode member, email, nama, nomor HP `08…`/`62…` atau potongan angka), `consent` (1/0), `hp` (1/0: punya nomor HP), `sync` (`PENDING`/`SYNCED`/`FAILED`), `active` (1/0), `page`, `per_page` (20). Item `{ id, member_code, full_name, email, phone_number, tier, is_active, email_verified, communication_consent, balance, sync_status, created_at }`, urut terdaftar terbaru. Izin `member.view` |
| `GET /admin/members/search?q=` | Minimal 3 karakter (422). Kode member, email, atau HP yang persis sama → hanya member itu; selain itu cocok sebagian, maks 10 (`meta.total` = jumlah semua). Izin `member.view` atau `point.create` |
| `GET /admin/members/:id` | Item di atas + `consent_at`, `google_connected`, `last_login_at`, `qontak: { id, sync_status, last_synced_at, retry_count, status_code } \| null`. Bukan member → 404. Izin `member.view` atau `point.create` |
| `PATCH /admin/members/:id` | `{ is_active }` → detail. Nonaktif: login ditolak 403 "Akun ini dinonaktifkan.", sesi yang berjalan → 401. Izin `member.update` |
| `GET /admin/points?user_id=&type=&page=&per_page=` | Riwayat satu member, terbaru dulu: `{ id, user_id, member_code, type, points, reference_no, purchase_amount, note, created_by_name, created_at }`, `meta.balance`. Izin `member.view` atau `point.create` |
| `GET /admin/points/check-reference?ref=` | Nota EARN, tidak peka huruf besar/kecil dan spasi di ujung → `{ available: true }` atau `{ available: false, used_at, member_code }`. Izin `point.create` |
| `POST /admin/points` | EARN `{ user_id, type, reference_no, purchase_amount, note? }`: poin dihitung backend dari `point_ratio_rupiah`. REDEEM `{ user_id, type, points (positif), note }`. ADJUST `{ user_id, type, points (+/−), note }`, izin `point.adjust`. → 201 `{ transaction, balance }`. 409 `errors.reference_no` (nota sudah dipakai, pesan berisi tanggal + kode member); 422 `errors.points` (saldo kurang / jadi minus), `errors.purchase_amount` (di bawah rasio), `errors.note`, `errors.user_id` (akun nonaktif) |

Dicatat di `activity_logs`: `point.create` / `point.adjust` (modul Poin, isi transaksi di `after_data`) dan `member.update` (`before_data` / `after_data` berisi `is_active`).

### Detail kontrak konten (dipakai mock Fase 5)
Waktu `start_at`, `end_at`, `published_at` = TIMESTAMP WIB tanpa zona, sama dengan kolom lain.

| Endpoint | Isi |
|---|---|
| `GET /banners?position=` | Hanya banner tayang (aktif dan dalam periode), urut `sort_order`: `{ id, title, image_url, link_url, position, start_at, end_at }`. `position` = `HOME_HERO` / `HOME_SLIDER` / `HOME_SIDE`. Khusus `HOME_HERO` + `subtitle`, `video_url`, `cta_label`, `cta2_label`, `cta2_url` (boleh `null`); Beranda memakai yang pertama |
| `GET /articles?type=&page=&per_page=` | Hanya `PUBLISHED`, terbaru dulu, `per_page` bawaan 9. `type` boleh beberapa dipisah koma (`PROMO,EVENT`); tanpa `type` → semua kecuali `HALAMAN`. Item `{ id, type, title, slug, thumbnail_url, excerpt, published_at }` (`excerpt` = paragraf pertama, ±160 karakter) |
| `GET /articles/:slug?preview=` | Item di atas + `content`, `status`. Draf hanya dengan token pratinjau yang benar; selain itu 404 seperti artikel yang tidak ada |
| ~~`GET /home/house-parts`~~ (tidak dipakai lagi sejak 30 Okt 2026) | `[{ key, category: { slug, name } \| null }]` untuk bagian `pondasi`, `lantai`, `dinding`, `bukaan`, `sanitasi`, `atap` (pemetaan ke kategori ditentukan backend; kategori nonaktif → `null`) |
| `GET /admin/banners` | Semua banner + `status` (`TAYANG` / `TERJADWAL` / `BERAKHIR` / `NONAKTIF`, dihitung server). Izin `banner.view` |
| `POST /admin/banners`, `PUT /admin/banners/:id` | `{ title, image_url, link_url?, position, start_date, end_date?, is_active }`, tanggal `YYYY-MM-DD` → disimpan `T00:00:00` / `T23:59:59`. 422: `title`, `image_url`, `position`, `start_date`, `end_date` (sebelum mulai), `link_url` (harus diawali `/` atau `https://`). Untuk `HOME_HERO` juga `subtitle?` (≤ 160), `video_url?` (URL hasil unggah), `cta_label?` (≤ 30), `cta2_label?` (≤ 30) + `cta2_url?` (berpasangan: label tanpa tautan → 422 `cta2_url`, tautan tanpa label → 422 `cta2_label`; diawali `/` atau `https://`); posisi lain mengabaikannya. Banner baru masuk urutan terakhir di posisinya. Izin `banner.create` / `banner.update` |
| `PATCH /admin/banners/:id` | `{ is_active }` atau `{ sort_order }`. Izin `banner.update` |
| `GET /admin/articles?q=&type=&status=&page=` | `{ id, type, title, slug, thumbnail_url, status, published_at, updated_at, author_name }`, urut terakhir diubah. Izin `article.view` |
| `GET /admin/articles/:id` | + `content`, `preview_token` |
| `POST /admin/articles`, `PUT /admin/articles/:id` | `{ type, title, slug, thumbnail_url?, content }`; artikel baru selalu `DRAFT`, PUT tidak mengubah status. 409 `errors.slug` (dipakai artikel lain), 422 `type`, `title`, `slug`, `content`. Izin `article.create` / `article.update` |
| `PATCH /admin/articles/:id` | `{ status: 'PUBLISHED' \| 'DRAFT' }`; terbit → `published_at` = sekarang, draf → `null`. Izin `article.publish` |
| `GET /admin/settings`, `PUT /admin/settings` | GET: `{ wa_number, address, opening_hours, maps_url, point_ratio_rupiah }`. PUT hanya menerima `{ address, opening_hours, maps_url, point_ratio_rupiah? }`; `wa_number` ditolak (422), termasuk untuk Super Admin. Nomor hanya-baca dari `WHATSAPP_NUMBER` server. `maps_url` harus `https://`. Mengubah `point_ratio_rupiah` butuh izin `setting.point_ratio` (403), minimal 1000. Izin `setting.update` |

Dicatat di `activity_logs`: `banner.create`, `banner.update`, `article.create`, `article.update`, `article.publish`, `setting.update` (nilai sebelum/sesudah di `before_data` / `after_data`).

### Detail kontrak jasa (dipakai mock sejak 2 Nov 2026, KEPUTUSAN #20)
Fitur Jasa di frontend sudah jadi dengan data contoh (`fase-5-konten/src/services/mock/routesKonten.js`). Supaya backend bisa menyesuaikan:

| Endpoint | Isi |
|---|---|
| `GET /services?category=&limit=` | Publik. Hanya jasa aktif, urut `sort_order` lalu `id`: `[{ id, name, slug, summary, image_url, service_area }]`. `category` = slug kategori produk; kembalikan jasa yang terhubung ke kategori itu **atau induknya** (dipakai kotak "Butuh pemasangan?" di detail produk dengan `limit=2`; Beranda memakai `limit=3`). Slug tidak dikenal → `[]` |
| `GET /services/:slug` | Publik. Item di atas + `description` (teks berformat sama dengan artikel: baris kosong = paragraf, `## ` subjudul, `- ` daftar) + `categories: [{ name, slug }]` (hanya kategori aktif). Jasa nonaktif / tidak ada → 404 |
| `GET /admin/services` | Semua jasa (termasuk nonaktif), urut `sort_order`: kolom tabel + `category_ids` + `categories: [{ id, name, slug, parent_id }]`. Izin `service.view` |
| `GET /admin/services/:id` | Satu jasa, bentuk sama. 404 bila tidak ada. Izin `service.view` |
| `POST /admin/services`, `PUT /admin/services/:id` | `{ name, slug, summary, description?, image_url?, service_area?, category_ids: number[], is_active }`. 422 per kolom: `name` (wajib, ≤ 100), `slug` (wajib, pola `^[a-z0-9]+(?:-[a-z0-9]+)*$`), `summary` (wajib, ≤ 160), `service_area` (≤ 100), `category_ids` (id kategori yang ada). 409 `{ errors: { slug } }` bila slug dipakai jasa lain. Jasa baru masuk urutan terakhir. Izin `service.create` / `service.update` |
| `PATCH /admin/services/:id` | `{ is_active }` atau `{ sort_order }`. Izin `service.update`. Tidak ada DELETE: jasa dinonaktifkan |
| `POST /admin/uploads` (`jenis=gambar`) | Dipakai untuk `image_url` (sama seperti gambar banner/artikel) |

Harga jasa tidak disimpan maupun ditampilkan (keputusan pemilik 2 Nov 2026): frontend selalu menulis "Harga pasti murah" (`utils/jasa.js`), biaya dibahas lewat WhatsApp.

Dicatat di `activity_logs` dengan `module = 'Jasa'`: `service.create` (after_data: name, category_ids) dan `service.update` (kolom yang berubah + `name` sebagai penanda), termasuk PATCH aktif/urutan.

Izin baru (menu Konten › Jasa, `menu_id` usulan 34, `route` `/admin/jasa`, `sort_order` 3 di grup Konten; "Pengaturan toko" bergeser ke 4): `service.view`, `service.create`, `service.update`. Role: **Admin Konten** dan **Super Admin**. Total izin jadi 36.

Usulan tabel (juga ditulis di `docs/transhome_postgres.sql`, bagian "USULAN: jasa"):
```sql
CREATE TABLE services (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    summary VARCHAR(160) NOT NULL,
    description TEXT NULL,
    image_url VARCHAR(255) NULL,
    service_area VARCHAR(100) NULL,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE service_categories (
    service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (service_id, category_id)
);
```
Konsultasi jasa tidak membuat data di backend (hanya tautan WhatsApp ke nomor `site_settings.wa_number`, masuk Mekari Qontak seperti chat lain). Bila nanti perlu laporan "berapa orang klik konsultasi jasa", itu kebutuhan baru (mis. endpoint pencatat klik), belum dibuat.

### Detail kontrak sistem (dipakai mock Fase 6)
| Endpoint | Isi |
|---|---|
| `GET /admin/users?q=&role=&active=&page=` | Hanya akun ber-role admin. Item `{ id, full_name, email, is_active, roles: [{ code, name }], last_login_at, created_at }`, urut aktif lalu nama. `meta.roles` = pilihan role admin `[{ code, name, description }]`. Izin `user.view` |
| `POST /admin/users` | `{ full_name, email, role, password }` (password sementara ≥ 8). Email langsung terverifikasi. 409 `errors.email` (dipakai akun mana pun), 422 per field. Izin `user.create` |
| `PUT /admin/users/:id/roles` | `{ roles: ['ADMIN_KATALOG'] }` (min. 1 role admin). 422 `errors.roles`: mencabut Super Admin dari diri sendiri, atau Super Admin aktif terakhir. Bukan admin → 404. Izin `user.update` |
| `PATCH /admin/users/:id` | `{ is_active }`. Nonaktif: login 403 "Akun ini dinonaktifkan.", sesi berjalan → 401. 422 untuk akun sendiri / Super Admin aktif terakhir. Izin `user.update` |
| `POST /admin/users/:id/reset-password` | Kirim link reset sekali pakai (1 jam) ke email admin → `{ message }`. Akun nonaktif → 422. Izin `user.update` |
| `GET /admin/roles` | Role admin (tanpa CUSTOMER): `{ id, code, name, description, is_locked, user_count (admin aktif), permissions: [kode] }`. Izin `role.view` |
| `GET /admin/permissions` | Izin per menu, urut sidebar: `[{ menu_id, menu, group, menu_permission, permissions: [{ id, code, action, description }] }]`. `menu_permission` = izin yang membuat menu terlihat. Izin `role.view` |
| `PUT /admin/roles/:id/permissions` | `{ permissions: [kode] }` daftar lengkap (bukan selisih) → role. SUPER_ADMIN → 403. 422: body tanpa daftar, kode tidak dikenal, atau izin menu tanpa `menu_permission`-nya. Izin `role.update` |
| `GET /admin/activity-logs?admin_id=&module=&from=&to=&page=` | `from` / `to` = tanggal WIB `YYYY-MM-DD` (inklusif), terbaru dulu, 20 per halaman. Item `{ id, admin_id, admin_name, action, module, before_data, after_data, ip_address, created_at }`. Izin `log.view` |
| `GET /admin/activity-logs/options` | `{ admins: [{ id, name }], modules: [nama modul] }` untuk filter. Izin `log.view` |
| `GET /admin/qontak/contacts?status=&q=&page=` | Item `{ id, user_id, name, email, phone_number, member_code, sync_status, retry_count, status_code, last_synced_at, created_at, updated_at }`, urut terdaftar terbaru. `meta.counts = { ALL, PENDING, SYNCED, FAILED }` untuk pencarian yang sama. Izin `qontak.retry` |
| `POST /admin/qontak/retry-failed` | Semua FAILED masuk antrean lagi → 202 `{ count }` (200 `{ count: 0 }` bila tidak ada). Izin `qontak.retry` |
| `GET /admin/qontak/contacts/:id/logs` | `sync_logs` kontak itu, terbaru dulu: `{ id, event_type, status_code, retry_count, request_payload, response_payload, created_at }`. Izin `qontak.retry` |

Dicatat di `activity_logs`: `user.create`, `user.update` (role / status), `user.reset_password`, `role.update` (daftar izin sebelum/sesudah), `qontak.retry`, `qontak.retry_failed`.

Format error yang disarankan agar helper error form bisa memetakan ke field:
```json
{ "message": "Validasi gagal", "errors": { "sku": ["SKU sudah dipakai produk lain"] } }
```

### Kebutuhan data dari backend (per 7 Okt 2026)
- Tabel `menus` dan `permissions` masih kosong di database; yang sudah di-seed baru 5 role. Kode izin yang dipakai rencana ini: `product.view/create/update`, `category.view/create/update`, `brand.view/create/update/delete`, `room.view/create/update`, `member.view/update`, `point.create`, `point.adjust`, `banner.view/create/update`, `article.view/create/update`, `article.publish`, `setting.update`, `setting.point_ratio`, `user.view/create/update`, `role.view/update`, `log.view`, `qontak.retry` (32 izin, daftar lengkap dengan `menu_id` di `fase-6-sistem/src/services/mock/data/sistem.js`).
- `users.phone_number` perlu boleh NULL sebelum Fase 3 (KEPUTUSAN #11).
- `site_settings` perlu kunci `point_ratio_rupiah`, alamat, jam buka, dan link Google Maps. Nomor WhatsApp hanya dari `WHATSAPP_NUMBER` di server; kunci `wa_number` lama tidak lagi digunakan (KEPUTUSAN #21).
- Katalog (Fase 1): endpoint di "Detail kontrak katalog" termasuk `meta.facets`. Daftar atribut `specifications` yang boleh jadi filter masih ditunggu (CLAUDE.md §12); mock memakai `ukuran`, `finishing`, `warna`, `bahan`, `daya`.
- Video produk: belum ada kolom di `products` maupun `product_images`. Thumbnail video di galeri dibuat setelah kolomnya ada (mis. `products.video_url`).
- Admin (Fase 2): endpoint di "Detail kontrak admin". Seed `menus` yang dipakai mock (grup Utama, Katalog, Membership, Konten, Sistem; route `/admin/...`) ada di `fase-2-admin-katalog/src/services/mock/data/admin.js`, termasuk kode izin per menu. Menu "Sync Qontak" memakai izin `qontak.retry` karena belum ada `qontak.view`.
- `POST /admin/uploads` perlu penyimpanan file (lokal/objek) yang mengembalikan URL publik; mock memakai data URL.
- Akun (Fase 3): endpoint di "Detail kontrak akun", termasuk `POST /auth/email/verify` (belum ada di daftar awal), pengiriman email verifikasi & reset password, dan alur Google yang kembali ke `return_to#token=…`.
- Konten (Fase 5): endpoint di "Detail kontrak konten". Tabel `articles` belum punya kolom untuk token pratinjau draf: tambah kolom `preview_token` atau buat token bertanda tangan (mis. HMAC id artikel + kedaluwarsa) di `GET /admin/articles/:id`. Status banner dihitung server dengan waktu WIB. Izin baru `setting.point_ratio` untuk rasio poin.
- Harga (15 Okt 2026): migration `products.show_price BOOLEAN NOT NULL DEFAULT TRUE` dan `products.price_general` boleh NULL. Validasi server: harga wajib bila `show_price` true; bila false harga boleh kosong. API publik (`/products`, `/products/:slug`) **tidak mengirim harga** produk dengan `show_price` false (`price_general: null`, `show_price: false`) supaya harga tidak bocor; produk itu tidak ikut filter `min`/`max` dan `facets.price`, dan diletakkan di akhir `sort=harga_terendah|harga_tertinggi`. `PATCH /admin/products/:id` menerima `{ price_general, unit_sale, show_price }` (sebagian boleh) untuk ubah cepat harga, izin `product.update`, dicatat `product.update` dengan harga lama/baru.
- Hero banner (30 Okt 2026): kolom baru `banners.subtitle VARCHAR(160) NULL`, `banners.video_url VARCHAR(255) NULL`, `banners.cta_label VARCHAR(30) NULL`, `banners.cta2_label VARCHAR(30) NULL`, `banners.cta2_url VARCHAR(255) NULL`; `position` (VARCHAR) menerima `HOME_HERO`. `POST /admin/uploads` menerima `jenis=video` (MP4/WebM, usulan ≤ 20 MB) dan perlu penyimpanan file yang bisa melayani video (HTTP range request). `GET /home/house-parts` tidak dipakai lagi.
- SEO (30 Okt 2026, KEPUTUSAN #17): `GET /sitemap.xml` dari backend (Beranda, `/katalog/*`, `/produk/*`, `/brand`, `/brand/*` yang punya produk, `/ruangan/*`, `/artikel/*`, `/promo`, `/info-toko`; `lastmod` dari `updated_at`), lalu aktifkan baris `Sitemap:` di `public/robots.txt` dengan domain produksi dan isi `VITE_SITE_URL`. Server sebaiknya mengembalikan HTTP 404 untuk slug yang tidak ada. Pratinjau tautan per halaman di WhatsApp/Facebook butuh prerender/SSR (mis. backend menyisipkan meta per URL), belum dikerjakan.
- Jasa (2 Nov 2026, KEPUTUSAN #20): tabel `services` + `service_categories`, endpoint publik & admin, 3 izin baru `service.*` dan menu Konten › Jasa, log `service.create/update`. Rincian lengkap di "Detail kontrak jasa" di atas. Sitemap juga perlu memuat `/jasa` dan `/jasa/:slug` yang aktif.
- Brand (8 Okt 2026): `GET /brands` mengirim semua brand aktif + `product_count` (produk aktif) + `categories: [{ slug, name }]` (kategori utama dari produk aktifnya, untuk filter di `/brand`, sejak 30 Okt 2026), tanpa paginasi (ratusan brand masih kecil). `DELETE /admin/brands/:id` (izin baru `brand.delete`, usulan: Super Admin + Admin Katalog): 422 bila brand masih aktif, 409 `{ message, product_count }` bila masih dipakai produk mana pun (aktif atau nonaktif) — penting karena FK `products.brand_id` memakai `ON DELETE SET NULL`, jadi tanpa cek ini produk diam-diam kehilangan brand. 204 bila berhasil; dicatat `brand.delete` dengan `before_data`.
- Sistem (Fase 6): endpoint di "Detail kontrak sistem". Seed tabel `permissions` (32 izin dengan `menu_id`) dan `role_permissions` yang dipakai mock ada di `fase-6-sistem/src/services/mock/data/sistem.js` dan `data/akun.js` (`rolePermissions`). Tabel `users` belum punya kolom nama untuk admin (nama member ada di `member_profiles`): tambah `users.full_name` atau tabel profil admin. `activity_logs` belum punya kolom id data yang diubah, jadi `before_data` / `after_data` perlu memuat penanda (SKU, nama, kode member); usulan: tambah `entity_type` + `entity_id`. Pengiriman ke Qontak perlu worker/antrean yang mengisi `sync_logs` (termasuk `request_payload` / `response_payload`) dan `qontak_contacts.sync_status`. Admin belum punya halaman ganti password sendiri; sementara lewat "Reset password" (link email).
- Membership (Fase 4): endpoint di "Detail kontrak membership admin". `users.last_login_at` perlu diisi saat login. Nomor nota sebaiknya disimpan dalam huruf besar, karena constraint `UNIQUE (type, reference_no)` di PostgreSQL membedakan huruf besar/kecil. Izin `point.adjust` dipakai untuk Penyesuaian (usulan: Super Admin saja).

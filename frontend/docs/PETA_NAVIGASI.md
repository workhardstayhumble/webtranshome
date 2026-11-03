# Peta Navigasi Web Transhome (E-Catalog + Membership + Admin)

Sumber: `frontend/fase-0-fondasi/src/router/routes/*.js`, layout, header/footer, dan tombol di setiap view.
Notasi: `[Tombol/Link]` → tujuan. `:slug` / `:id` = parameter dinamis.

## 1. Gambaran besar: 4 area (tab) website

| # | Area | Layout | Siapa yang boleh | Prefix URL |
|---|---|---|---|---|
| 1 | **Publik** (landing + katalog + konten) | `PublicLayout` (header + footer + WA melayang) | Semua pengunjung | `/` |
| 2 | **Auth** (masuk/daftar/reset) | `AuthLayout` (logo → `/`) | Tamu (sudah login dialihkan) | `/masuk`, `/daftar`, ... |
| 3 | **Member** (akun pelanggan) | `MemberLayout` (menu samping) | Login + role `CUSTOMER` | `/akun` |
| 4 | **Admin panel** | `AdminLayout` (sidebar dinamis) | Login admin + permission per halaman | `/admin` |

Aturan guard (`router/guard.js`):
- Belum login buka `/akun/*` → `/masuk?redirect=<halaman>`; buka `/admin/*` → `/admin/masuk?redirect=<halaman>`.
- Sudah login buka halaman tamu (`/masuk`, `/daftar`, `/lupa-password`) → admin ke `/admin`, member ke `/akun` (atau `redirect`). Pengecualian (2 Nov 2026): member yang sedang masuk **tetap boleh membuka `/admin/masuk`** (ganti akun; ada pemberitahuan akun member akan diganti).
- Member/pelanggan buka `/admin/*` → `/admin/masuk?redirect=<halaman>` (bukan jalan buntu). Admin tanpa permission halaman itu → `/tidak-punya-akses?dari=<halaman>`. Admin buka `/akun/*` → `/tidak-punya-akses?dari=<halaman>`.
- Akun admin hanya dibuat Super Admin di `/admin/pengguna` (tidak ada pendaftaran admin sendiri). Email tidak terdaftar di `/admin/masuk` → "Email atau password salah." (sama dengan password salah).
- URL tak dikenal → `/…` 404 publik, `/admin/…` 404 admin, `/akun/…` 404 member.

---

## 2. Navigasi global (muncul di banyak halaman)

### Header publik
| Elemen | Tujuan |
|---|---|
| Logo Transhome | `/` |
| Baris atas: `Info toko` | `/info-toko` |
| Baris atas: `WhatsApp CS` | WhatsApp (tab baru) |
| `Kategori` (mega menu) → judul induk / sub kategori | `/katalog/:slug-kategori` |
| Mega menu (HP/tablet) chip `Ruangan`, `Jasa`, `Promo`, `Artikel`, `Info toko` | `/ruangan`, `/jasa`, `/promo`, `/artikel`, `/info-toko` |
| Kotak pencarian: Enter / tombol cari | `/cari?q=...` |
| Kotak pencarian: klik saran | halaman saran (produk/kategori/brand) |
| Menu desktop (xl): `Ruangan`, `Brand`, `Jasa`, `Promo`, `Artikel` | `/ruangan`, `/brand`, `/jasa`, `/promo`, `/artikel` |
| Belum login: `Masuk` | `/masuk?redirect=<halaman sekarang>` |
| Belum login: `Daftar member` | `/daftar` |
| Sudah login (member): menu avatar → `Akun saya` / `Riwayat poin` | `/akun` / `/akun/poin` |
| Sudah login (admin): menu avatar → `Admin panel` | `/admin` |
| Menu avatar → `Keluar` | logout → `/` |

### Footer publik
`Katalog` `/katalog` · `Ruangan` `/ruangan` · `Brand` `/brand` · `Jasa` `/jasa` · `Promo` `/promo` · `Artikel` `/artikel` · `Info toko` `/info-toko` · `Buka di Google Maps` (eksternal) · `Chat WhatsApp CS` (WA) · "Ikuti kami": `Instagram` (instagram.com/transhome, tab baru) · `TikTok` (tiktok.com/@transhome_id, tab baru).

### Tombol WhatsApp melayang
Desktop saja → WhatsApp CS. Di halaman produk HP digantikan bar harga + WA di bawah layar.

---

## 3. TAB 1: Area Publik

### 3.1 Landing page `/` (Beranda)
Urutan section dari atas ke bawah dan tujuan tiap tombol:

| Section | Tombol / link | Tujuan |
|---|---|---|
| Hero banner (setinggi layar, gambar/video dari Admin > Banner posisi Hero Beranda; menggantikan hero rumah 3D sejak 30 Okt 2026) | tombol utama (bawaan `Jelajahi katalog`) | tautan banner, bawaan `/katalog` |
| Hero banner | tombol kedua (opsional, diisi admin) | tautan kedua banner |
| Hero banner | `Jeda video` / `Putar video` (bila ada video) | di halaman |
| Hero banner | `Lihat selengkapnya` | gulir ke section berikutnya |
| Banner slider/promo | klik banner | `link_url` banner (diatur dari Admin > Banner) |
| Banner promo member | tombol ajakan | belum login `/daftar`, sudah login `/akun` |
| Kategori per tahap bangun | kartu kategori | `/katalog/:slug` |
| Produk unggulan | kartu produk | `/produk/:slug` |
| Promo & event | `Lihat semua` | `/promo` |
| Promo & event | kartu artikel | `/artikel/:slug` |
| Inspirasi ruangan | `Lihat semua` / kartu ruangan | `/ruangan` / `/ruangan/:slug` |
| Jasa ("Butuh tukang? Pakai jasa Transhome", 3 teratas) | `Lihat semua jasa` / kartu jasa / `Konsultasi via WhatsApp` | `/jasa` / `/jasa/:slug` / WhatsApp |
| Brand pilihan (lingkaran logo + "Siap untuk menjelajah lebih banyak Brand?") | `Telusuri N brand` / lingkaran logo | `/brand` / `/brand/:slug` |
| Keunggulan | (informasi, tanpa tombol) | – |
| Ajakan member | `Daftar sekarang` | `/daftar` |

### 3.2 Katalog `/katalog` dan `/katalog/:kategori`
- Breadcrumb: `Beranda` `/` → `Katalog` `/katalog` → induk `/katalog/:induk` → sub kategori.
- Chip navigasi kategori (semua induk, atau `Semua <induk>` + sub): pindah `/katalog/:slug`.
- Toolbar: `Urutkan` (terbaru, harga terendah/tertinggi, nama A-Z) dan toggle `Grid/Daftar` (hanya ubah query/tampilan); `Filter` (HP) buka panel filter.
- Filter (brand, harga, stok, dll.) + chip filter aktif: hanya ubah query URL, tidak pindah halaman. Paginasi: query `page`.
- Kartu produk (gambar/nama) → `/produk/:slug`.
- Kosong/tidak ketemu: `Lihat semua katalog` `/katalog`, atau `Tanya via WhatsApp`.

### 3.3 Pencarian `/cari?q=`
Breadcrumb `Beranda` → `Pencarian`. Hasil memakai komponen katalog yang sama (filter, urutan, kartu → `/produk/:slug`). Kosong → tanya WhatsApp.

### 3.4 Detail produk `/produk/:slug`
| Elemen | Tujuan |
|---|---|
| Breadcrumb | `Beranda` `/` → `Katalog` `/katalog` → induk `/katalog/:induk` → kategori `/katalog/:kategori` |
| Link brand | `/brand/:slug` |
| Chip ruangan (produk cocok untuk ruangan) | `/ruangan/:slug` |
| Galeri foto | buka Lightbox (tetap di halaman) |
| `Salin SKU`, `Bagikan` (salin link) | aksi di halaman |
| `Tanya via WhatsApp` / `Tanya ketersediaan` (stok habis) | WhatsApp dengan pesan terisi SKU + link produk |
| `Chat tim proyek` (jumlah besar) | WhatsApp |
| "Butuh pemasangan?" (jasa dari kategori produk ini, bila ada): nama jasa / `Konsultasi jasa` | `/jasa/:slug` / WhatsApp (pesan berisi jasa + SKU produk) |
| `Unduh datasheet` (jika ada) | PDF eksternal |
| Tab info: `Deskripsi` · `Spesifikasi` · `Dokumen` (jika ada datasheet) | ganti tab di halaman |
| Produk terkait: kartu | `/produk/:slug` lain |
| Produk terkait: `Lihat kategori` | `/katalog/:kategori` |
| Produk tidak ada | `Lihat katalog` `/katalog` |

### 3.5 Brand
- `/brand` (daftar brand): breadcrumb `Beranda`; kartu brand → `/brand/:slug`.
- `/brand/:slug`: breadcrumb `Beranda` → `Katalog` → nama brand; berisi daftar produk brand (filter, urutan, kartu → `/produk/:slug`); brand tak ada → `/katalog`.

### 3.6 Ruangan
- `/ruangan` ("Belanja per ruangan"): carousel/kartu ruangan → `/ruangan/:slug`.
- `/ruangan/:slug`: breadcrumb `Beranda` → `Ruangan` `/ruangan`; chip pindah antar ruangan `/ruangan/:slug`; kartu produk → `/produk/:slug`; ruangan tak ada → `/ruangan`.

### 3.7 Konten
- `/artikel`: daftar artikel (kartu → `/artikel/:slug`).
- `/artikel/:slug`: breadcrumb `Beranda` → asal (`Artikel` atau `Promo`) → judul; `Artikel lainnya` → artikel lain; tanya WhatsApp tentang artikel; tidak ada → `/artikel`.
- `/promo` ("Promo dan event"): filter jenis (query `jenis`), kartu → `/artikel/:slug`.
- `/info-toko`: alamat, jam buka, peta (Google Maps eksternal), WhatsApp.

### 3.7a Jasa (2 Nov 2026, KEPUTUSAN #20)
- `/jasa`: breadcrumb `Beranda` → `Jasa`; "Cara pesan jasa" 3 langkah; kartu jasa (judul → `/jasa/:slug`, `Konsultasi via WhatsApp` → WhatsApp); `Tanya via WhatsApp` untuk jasa yang belum ada.
- `/jasa/:slug`: breadcrumb `Beranda` → `Jasa` → nama; kotak konsultasi ("Harga pasti murah", area layanan, `Konsultasi via WhatsApp`); chip bahan → `/katalog/:kategori`; `Lihat jasa lainnya` → `/jasa`; tidak ada → `Lihat semua jasa` `/jasa`.
- Tidak ada formulir: konsultasi hanya lewat WhatsApp (KEPUTUSAN #12).

### 3.8 Halaman sistem publik
- `/tidak-punya-akses` → `Masuk dengan akun lain` (keluar, lalu `/admin/masuk` atau `/masuk` sesuai `?dari=`, kembali ke halaman itu setelah masuk) · `Kembali ke dashboard` `/admin` (admin) atau `Kembali ke beranda` `/` (lainnya).
- `/…` (tidak ditemukan) → tombol kembali ke `/`.
- Dev saja (`npm run dev`): `/dev/components`, `/dev/email`, `/dev/google`. Tidak ikut produksi.

---

## 4. TAB 2: Area Auth

Semua halaman memakai logo → `/`.

| Halaman | Tombol / link | Tujuan |
|---|---|---|
| `/masuk` | `Masuk` (berhasil) | member → `redirect` atau `/akun`; admin → `/admin` |
| | `Masuk dengan Google` | `/masuk/google` → balik ke `redirect`/`/akun` |
| | `Lupa password?` | `/lupa-password` |
| | `Daftar member` | `/daftar` |
| `/daftar` | `Daftar` (berhasil) | `/verifikasi?email=` |
| | email sudah terdaftar → link | `/masuk?email=` |
| | `Masuk` | `/masuk` |
| `/verifikasi` (dari link email) | `Masuk sekarang` | `/masuk?email=` |
| | `Masuk` / `Daftar lagi` | `/masuk` / `/daftar` |
| `/lupa-password` | `Kembali ke masuk` | `/masuk` |
| `/reset-password` (dari link email) | berhasil → otomatis | `/masuk` |
| | link kedaluwarsa → `Minta link baru` | `/lupa-password` |
| `/masuk/google` | callback Google | `/akun` atau `/admin` |
| `/admin/masuk` | `Masuk` (berhasil) | `redirect` atau `/admin` (sesi member yang sedang aktif diganti) |
| | `Lupa password?` | `/lupa-password` |

---

## 5. TAB 3: Area Member `/akun`

Menu samping (HP: tab geser horizontal): `Ringkasan` `/akun` · `Riwayat poin` `/akun/poin` · `Profil` `/akun/profil` · `Keamanan` `/akun/keamanan`.
Header: logo → `/`, link kembali ke beranda `/`, `Keluar` → logout → `/`.

| Halaman | Isi | Tombol / link → tujuan |
|---|---|---|
| `/akun` Ringkasan | kartu member (kode, poin), transaksi terakhir | `Lihat riwayat poin` / `Lihat semua` → `/akun/poin` |
| `/akun/poin` | riwayat poin + filter jenis (query `jenis`) | – |
| `/akun/profil` | ubah data profil | simpan (tetap di halaman) |
| `/akun/keamanan` | ganti password | simpan (tetap di halaman) |
| `/akun/…` lain | 404 member | `Kembali ke ringkasan` `/akun` |

---

## 6. TAB 4: Admin panel `/admin`

### 6.1 Struktur sidebar (dari API `GET /admin/menus`, difilter permission)
Header halaman menampilkan `grup › judul` dan label role. HP: sidebar jadi drawer. Logo sidebar → `/admin`. Tombol `Keluar` → `/admin/masuk`.

| Grup | Menu | Route | Permission |
|---|---|---|---|
| Utama | Dashboard | `/admin` | – |
| Katalog | Produk | `/admin/produk` | `product.view` |
| | Kategori | `/admin/kategori` | `category.view` |
| | Brand | `/admin/brand` | `brand.view` |
| | Ruangan | `/admin/ruangan` | `room.view` |
| Membership | Member | `/admin/member` | `member.view` |
| | Input poin | `/admin/poin/input` | `point.create` |
| Konten | Banner | `/admin/banner` | `banner.view` |
| | Artikel dan promo | `/admin/artikel` | `article.view` |
| | Jasa | `/admin/jasa` | `service.view` |
| | Pengaturan toko | `/admin/pengaturan` | `setting.update` |
| Sistem | Pengguna admin | `/admin/pengguna` | `user.view` |
| | Role dan izin | `/admin/role` | `role.view` |
| | Log aktivitas | `/admin/log-aktivitas` | `log.view` |
| | Sync Qontak | `/admin/sync-qontak` | `qontak.retry` |

Role contoh: **Super Admin** (semua) · **Admin Katalog** (grup Katalog) · **Admin Membership** (grup Membership) · **Admin Konten** (grup Konten).

### 6.2 Jalur per halaman admin

**Dashboard `/admin`**
- Aksi cepat (muncul sesuai izin): `Tambah produk` → `/admin/produk/baru` · `Input poin member` → `/admin/poin/input` · `Buat banner` → `/admin/banner/baru` · `Tulis artikel` → `/admin/artikel/baru`.
- Kartu ringkasan: Produk aktif → `/admin/produk?active=1` · Stok habis → `/admin/produk?stock=HABIS&active=1` · Member → `/admin/member` · Banner tayang → `/admin/banner` · Sync gagal → `/admin/sync-qontak?status=FAILED` · Poin bulan ini (tanpa link).
- "Perlu perhatian": `Ubah stok` → `/admin/produk/:id#stok` · `Lengkapi` (tanpa foto) → `/admin/produk/:id#foto` · `Perpanjang` banner → `/admin/banner/:id#periode` · `Lihat semua …` → daftar terfilter.
- "Aktivitas terbaru": `Lihat semua` → `/admin/log-aktivitas`.

**Katalog**
- `/admin/produk` (daftar + filter query): `Tambah produk` → `/admin/produk/baru`; baris `Ubah` → `/admin/produk/:id`; baris/aksi lain → `/admin/produk/:id/pratinjau`.
- `/admin/produk/baru` & `/admin/produk/:id`: form produk (foto, spesifikasi, ruangan, datasheet, tambah brand cepat via modal). `Simpan`/`Batal` → kembali ke daftar sebelumnya (`/admin/produk` + filter terakhir). Produk tidak ada → `Kembali ke daftar produk`.
- `/admin/produk/:id/pratinjau`: `Ubah produk` → `/admin/produk/:id`; `Kembali ke daftar produk` → `/admin/produk`.
- `/admin/kategori`, `/admin/brand`, `/admin/ruangan`: tabel + modal tambah/ubah (tanpa halaman terpisah). Di Ruangan baris: `Lihat produk` → `/admin/produk?room=:id`.

**Membership**
- `/admin/member`: tabel + filter; `Input poin` → `/admin/poin/input`; baris `Lihat/Detail` → `/admin/member/:id`.
- `/admin/member/:id`: profil, riwayat poin, status sync Qontak; `Input poin` → `/admin/poin/input?member=:id`; `Kembali` → `/admin/member`.
- `/admin/poin/input`: pilih member (cari), jenis (belanja / tukar poin), simpan; kartu member terpilih → `/admin/member/:id`; riwayat ringkas → `/admin/member/:id`.

**Konten**
- `/admin/banner`: `Buat banner` → `/admin/banner/baru`; `Ubah` → `/admin/banner/:id`; banner berakhir `Perpanjang` → `/admin/banner/:id#periode`.
- `/admin/banner/baru` & `/admin/banner/:id`: simpan/batal → `/admin/banner`.
- `/admin/artikel`: `Tulis artikel` → `/admin/artikel/baru`; `Ubah` → `/admin/artikel/:id`.
- `/admin/artikel/baru` & `/admin/artikel/:id`: simpan pertama kali → pindah ke `/admin/artikel/:id`; `Kembali` → `/admin/artikel`.
- `/admin/jasa`: `Tambah jasa` → `/admin/jasa/baru`; `Ubah` → `/admin/jasa/:id`; `Lihat` → `/jasa/:slug` (tab baru); urutan ↑ ↓ dan sakelar aktif (konfirmasi saat menonaktifkan) di tabel.
- `/admin/jasa/baru` & `/admin/jasa/:id`: simpan/batal → `/admin/jasa` (batal dengan perubahan → "Buang perubahan?").
- `/admin/pengaturan`: alamat, jam buka, nomor WA, rasio poin (tetap di halaman).

**Sistem**
- `/admin/pengguna`: tabel, `Tambah admin` (modal), `Ubah role` (modal).
- `/admin/role?role=KODE`: daftar role + matriks izin (ganti role = ubah query).
- `/admin/log-aktivitas`: tabel + filter, klik baris → panel detail (perbandingan data sebelum/sesudah).
- `/admin/sync-qontak`: tabel kontak, panel log sync, `Coba ulang` sync; baris kontak → `/admin/member/:user_id`.

**Lainnya**: `/admin/…` tidak dikenal → `Kembali ke dashboard` `/admin`.

---

## 7. Alur utama (journey)

1. **Pengunjung → produk → WhatsApp**: `/` → `Lihat katalog` `/katalog` → pilih kategori `/katalog/:slug` → kartu `/produk/:slug` → `Tanya via WhatsApp`.
2. **Pengunjung → cari**: header cari → `/cari?q=` → `/produk/:slug`.
3. **Daftar member**: `/` → `Daftar member` `/daftar` → `/verifikasi` → (link email) → `/masuk` → `/akun`.
4. **Member lihat poin**: `/akun` → `Riwayat poin` `/akun/poin`.
5. **Admin masuk**: `/admin/masuk` → `/admin` (atau halaman `redirect`).
6. **Admin tambah produk**: `/admin` → `Tambah produk` → `/admin/produk/baru` → simpan → `/admin/produk`.
7. **Admin input poin**: `/admin/poin/input` → pilih member → simpan; atau `/admin/member/:id` → `Input poin`.
8. **Admin kelola landing page**: `/admin/banner` + `/admin/artikel` + `/admin/jasa` + logo di `/admin/brand` → tampil di `/`, `/promo`, `/jasa`.
9. **Pengunjung → jasa → WhatsApp**: `/` atau menu `Jasa` → `/jasa` → `/jasa/:slug` → `Konsultasi via WhatsApp`; atau `/produk/:slug` → "Butuh pemasangan?" → `Konsultasi jasa`.
10. **Member membuka admin**: `/admin` → `/admin/masuk?redirect=/admin` (pemberitahuan ganti akun) → masuk dengan akun admin → `/admin`.

---

## 8. Temuan UX yang layak dibahas

1. **Admin lupa password**: `/admin/masuk` mengarah ke `/lupa-password` (halaman pelanggan, lalu `Kembali ke masuk` ke `/masuk`, bukan `/admin/masuk`). Perlu diputuskan: alur reset khusus admin atau minimal kembali ke `/admin/masuk`.
2. **`Katalog` tidak ada di menu header**: hanya lewat tombol `Kategori` (mega menu), Hero, dan footer. Pertimbangkan link `Katalog` di baris menu utama.
3. **Menu desktop vs HP berbeda**: `Brand` hanya di menu desktop (xl); chip HP di mega menu tidak punya `Brand`.
4. ~~**Halaman `/tidak-punya-akses`** hanya menawarkan `Kembali ke beranda`~~ Selesai 2 Nov 2026: ada `Kembali ke dashboard` (admin) dan `Masuk dengan akun lain`.
5. **Kategori/Brand/Ruangan admin** memakai modal (tanpa URL per item), sehingga tidak bisa dibagikan atau di-deep-link; Produk, Banner, Artikel punya URL sendiri. Konsistensi pola perlu dipertimbangkan.
6. **Pengguna admin** tidak punya halaman detail; hanya tabel + modal.
7. **Tidak ada e-commerce** (keranjang/checkout) di Muka 1: semua CTA produk berujung ke WhatsApp. Tombol utama di halaman produk adalah satu-satunya konversi.

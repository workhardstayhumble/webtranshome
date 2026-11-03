# Log Keputusan — Website Transhome Muka 1

Keputusan di file ini MENGALAHKAN BRD v0.4 / PRD v0.4 bila bertentangan.
Diputuskan: 6 Oktober 2026. Diperbarui: 7 Oktober 2026 (#7 stack frontend, #13 lokasi kerja), 8 Oktober 2026 (#14 brand), 15 Oktober 2026 (#15 harga), 2 November 2026 (#19 brand lingkaran, #20 Jasa, #21 nomor WhatsApp).

| # | Topik | Keputusan |
|---|---|---|
| 1 | Login Google | Tetap ada, sebagai tombol sendiri. Login dan daftar opsional; katalog bisa dilihat tanpa login. |
| 2 | member_code | Ditampilkan ke member di `/akun`; disebut ke kasir sebagai pengganti kartu fisik. |
| 3 | Nota, hadiah, penukaran, wishlist, pengingat stok, ulasan, konsultasi | Out of scope Muka 1. Sistem hanya mencatat poin per member (`point_transactions`). |
| 4 | Ulasan & verifikasi pembeli | Out of scope. `reference_no` hanya memuat total belanja; rincian produk per nota nanti dari sistem POS terpisah, tidak masuk website. |
| 5 | 58 kolom Master Barang | BELUM FINAL — menunggu header kolom. Aturan sementara: data internal (HPP, hnet, diskon, stok per lokasi, konsinyasi, karantina) tidak pernah tampil. |
| 6 | Rasio poin | Disimpan di `site_settings.point_ratio_rupiah`, bisa diubah Super Admin. Contoh: 1 poin per Rp 1.000.000. Angka final dari Transhome. |
| 7 | Stack frontend | Vue 3 + Vite (SPA) + Tailwind CSS v4, dengan Vue Router, Pinia, axios, dan lucide-vue-next, ditulis dalam JavaScript. Tanpa library komponen UI; Nuxt dan Nuxt UI tidak dipakai. Dikonfirmasi pemilik proyek 7 Okt 2026, menggantikan catatan "tetap Nuxt 4" sebelumnya. |
| 8 | Peta halaman | Lihat `CLAUDE.md` §3. Login admin terpisah di `/admin/masuk`. |
| 9 | Katalog Hadiah | Tidak dibuat. Penukaran poin dicatat admin sebagai REDEEM. |
| 10 | Brand | Masuk scope: filter brand di katalog, halaman `/brand/[slug]`, menu admin brand. |
| 11 | Email & HP | Email WAJIB, HP OPSIONAL. Reset password lewat email. |
| 12 | Konsultasi proyek | Hanya via WhatsApp Mekari Qontak. Tidak ada formulir konsultasi di website. |
| 13 | Lokasi kerja | Semua kode dan dokumen frontend ada di worktree `rafaelplatinum-ideal-memory`, folder `frontend/`: aturan di `frontend/CLAUDE.md`, dokumen di `frontend/docs/`. Salinan dokumen di clone `C:\Users\Victus\webtranshome` tidak dipakai lagi. |
| 14 | Brand: tampil & hapus | Semua brand aktif tampil di Beranda ("Brand pilihan"), halaman `/brand`, dan filter katalog; tidak ada pilihan "tampilkan di Beranda" per brand. Admin boleh MENGHAPUS PERMANEN brand yang sudah nonaktif dan tidak dipakai produk mana pun (izin `brand.delete`); brand yang masih dipakai harus dipindahkan produknya dulu. Diputuskan pemilik proyek 8 Okt 2026. |
| 15 | Harga boleh disembunyikan | Per produk ada opsi "Tampilkan harga di website" (`products.show_price`, bawaan aktif). Bila dimatikan, harga boleh kosong; pelanggan melihat "Harga via WhatsApp" dan tombol "Tanya harga" (stok HABIS tetap "Tanya ketersediaan"). Produk tanpa harga tidak ikut filter dan rentang harga, dan berada di akhir urutan harga. Admin Katalog dan Super Admin bisa mengubah harga cepat dari daftar produk. Diputuskan pemilik proyek 15 Okt 2026. |
| 16 | Hero Beranda & halaman Brand | Hero rumah 3D ("Pilih bagian rumah") DIHAPUS, diganti hero banner sesuai `docs/wireframe-low-fidelity-transhome.xlsx` sheet "Beranda": gambar atau video selebar layar, judul + subjudul, satu tombol. Isinya diatur Admin Konten lewat menu Banner, posisi baru **Hero Beranda** (`HOME_HERO`): judul, subjudul, gambar (wajib, juga cadangan video), video MP4/WebM opsional, label & tautan tombol. Bila beberapa tayang, urutan pertama yang tampil; tanpa banner tayang → teks bawaan. Urutan seksi Beranda lain dan seksi "Brand pilihan" TIDAK berubah. Halaman `/brand` mengikuti sheet "Brand" (jumlah brand, Cari, filter kategori, A–Z, kartu 6/3/2 kolom, 24 per halaman). Diputuskan pemilik proyek 30 Okt 2026. Tambahan hari yang sama: hero **setinggi layar** (landing page; header + hero = satu layar) dengan tombol utama, **tombol kedua** opsional yang diatur admin (label + tautan), dan **Lihat selengkapnya** (gulir ke seksi berikutnya). Tombol "Tanya via WhatsApp" di hero dihapus 2 Nov 2026 atas permintaan pemilik; WhatsApp tetap di header dan tombol melayang. |
| 17 | SEO halaman Brand | Tanpa library: `/brand` dan `/brand/:slug` punya judul & deskripsi dinamis, canonical, Open Graph, robots, dan data terstruktur schema.org (BreadcrumbList, CollectionPage + ItemList brand, Brand). Hasil cari/filter/urutan = `noindex, follow` + canonical ke halaman bersih; `?page=N` boleh diindeks (canonical ke dirinya); brand tanpa produk dan brand tidak ditemukan = `noindex`. Paginasi berupa tautan `<a href>` agar bisa diikuti Google. `robots.txt` memblokir `/admin`, `/akun`, `/cari`. Sitemap dan pratinjau tautan per halaman (prerender/SSR) menunggu backend. Diputuskan pemilik proyek 30 Okt 2026. |
| 18 | Logo & kartu member | Lambang resmi Transhome (sheet "Logo Transhome" di wireframe XLSX, bagian rumah + H tanpa tulisan) dipakai sebagai ikon tab (`public/favicon.ico`, `apple-touch-icon.png`) dan di samping tulisan "Transhome" di header, footer, halaman masuk/daftar, area member, dan sidebar admin. Ajakan "Jadi member Trans Family" (dengan gambaran kartu) dan kartu Trans Family di Beranda disembunyikan/diubah hanya untuk member (role CUSTOMER); admin yang sedang masuk tetap melihatnya seperti pengunjung. Diputuskan pemilik proyek 30 Okt 2026. |
| 19 | "Brand pilihan" Beranda berbentuk lingkaran | Mengganti logo berjalan (KEPUTUSAN #14): judul di tengah "Siap untuk menjelajah lebih banyak Brand?" + tombol "Telusuri N brand" → `/brand`, diapit lingkaran berwarna berisi logo brand (desktop 9 kiri + 9 kanan, baris tengah bergeser, lingkaran terluar terpotong tepi panel; HP 5 di atas + 5 di bawah). Logo diatur di admin › Katalog › Brand (Super Admin & Admin Katalog). Urutan: brand berlogo dulu, lalu produk terbanyak, lalu nama; sisanya tetap lengkap di `/brand`. Brand tanpa logo tampil sebagai nama; lingkaran yang tidak terisi brand jadi hiasan dua warna. Diputuskan pemilik proyek 2 Nov 2026. |
| 20 | Fitur Jasa | Website punya halaman **Jasa** (`/jasa`, `/jasa/:slug`) berisi layanan Transhome untuk kebutuhan rumah (contoh awal: pasang keramik & granit, pengecatan, pasang perlengkapan kamar mandi, pasang AC & elektronik, service & perbaikan, pasang pintu & jendela, antar & rakit perabot seperti sofa dan lemari). Konsultasi jasa **hanya lewat WhatsApp** (tetap sesuai KEPUTUSAN #12, tanpa formulir): tombol "Konsultasi via WhatsApp" dengan pesan otomatis berisi nama jasa, tautan, produk asal (bila dari halaman produk), dan isian lokasi & kebutuhan. Daftar jasa dikelola **Admin Konten + Super Admin** di menu Konten › Jasa (izin `content.service.read` / `content.service.manage`; alias mock `service.view/create/update`): nama, ringkasan, penjelasan, gambar, area layanan, kategori produk terkait, aktif, urutan. **Harga tidak ditampilkan**: semua jasa bertuliskan tetap "Harga pasti murah"; biaya dibahas lewat WhatsApp. Tampil di menu header & footer, seksi Beranda (3 teratas), dan kotak "Butuh pemasangan?" di detail produk dari kategori terkait. Daftar awal berupa contoh. Diputuskan pemilik proyek 2 Nov 2026. |
| 21 | Nomor WhatsApp toko | Nomor hanya diatur melalui environment backend/server (`WHATSAPP_NUMBER`) oleh Tim Cyber/operator server, lalu backend direstart. Tidak dapat diubah melalui dashboard atau API oleh role mana pun, termasuk Super Admin. Frontend menerima nomor hanya-baca dari API settings; `wa_number` tidak diterima pada PUT settings. Nilai lama di database tidak lagi menjadi sumber nomor. Alamat, jam buka, peta, dan pembatasan rasio poin tetap seperti sebelumnya. |

## Dampak ke Database (`docs/transhome_postgres.sql`)
- `users.phone_number` → boleh NULL (tetap UNIQUE). Baru diubah di `docs/transhome_postgres.sql`. Per 7 Okt 2026, `backend/migrations` dan database `CRM_TH` masih `NOT NULL`, jadi daftar tanpa HP baru bisa jalan setelah tim backend membuat migration baru.
- Seed `site_settings`: `point_ratio_rupiah` = `1000000` (contoh, ganti sesuai keputusan Transhome). Sama seperti di atas, baru ada di file SQL dokumen.
- Tabel `brands` tetap dipakai.
- Tidak perlu tabel: receipts, rewards, redemptions, wishlists, stock_reminders, product_reviews, project_consultations.
- Artifact ERD & Workflow perlu disesuaikan: registrasi (HP opsional), alur Google tanpa wajib isi HP.

## Revisi yang Perlu Dilakukan di BRD v0.4
- §1 Ringkasan Eksekutif, pilar ke-3 "Inbound Lead Generator via Mekari Qontak": hapus "formulir konsultasi proyek", sisakan widget chat WhatsApp.
- §1 pilar ke-2 dan §4 "Identifikasi Digital Tunggal": ubah "Nomor HP (utama) atau Email" → "Email (wajib), Nomor HP (opsional)".
- §2 tabel scope: baris "Formulir Konsultasi Proyek & Chat Widget WhatsApp" → "Chat Widget WhatsApp (Mekari Qontak)". Tambah ke out-of-scope: formulir konsultasi, katalog hadiah, wishlist, pengingat stok, ulasan.
- §3 peran Admin Konten: hapus "moderasi ulasan produk", "galeri inspirasi", "landing page kampanye" jika tidak dibuat. Peran Member Customer: hapus wishlist, pengingat stok, ulasan, submit lead konsultasi.
- §4 "Tunneling Lead Mekari Qontak": ganti "formulir konsultasi proyek" → "chat WhatsApp".

## Revisi yang Perlu Dilakukan di PRD v0.4
- Hapus FR-J (Katalog Hadiah), FR-K (Wishlist), FR-L (Pengingat Stok), FR-M (Ulasan), FR-Q (Formulir Konsultasi).
- FR-B: tambahkan filter brand. FR-E: email wajib + HP opsional + login Google.
- FR-H: tampilkan member_code. FR-I: tambahkan rasio poin dari pengaturan.
- FR-T: hapus "pembuatan tiket lead otomatis" dari formulir; sisakan sinkronisasi kontak.
- §2 tabel database: tambahkan `brands`, `product_images`, `point_transactions`, `banners`, `articles`, `site_settings`, `qontak_contacts`, `sync_logs` (total 21 tabel, sesuai SQL).
- §3 roadmap: hapus tabel receipts, rewards, redemptions, wishlists, stock_reminders, product_reviews, project_consultations; ganti `integration_logs` → `sync_logs`.

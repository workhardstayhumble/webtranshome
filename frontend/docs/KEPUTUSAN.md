# Log Keputusan — Website Transhome Muka 1

Keputusan di file ini MENGALAHKAN BRD v0.4 / PRD v0.4 bila bertentangan.
Diputuskan: 6 Oktober 2026. Diperbarui: 7 Oktober 2026 (#7 stack frontend, #13 lokasi kerja), 8 Oktober 2026 (#14 brand), 2 November 2026 (#15 nomor WhatsApp).

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
| 15 | Nomor WhatsApp toko | Nomor hanya diatur melalui environment backend/server (`WHATSAPP_NUMBER`) oleh Tim Cyber/operator server, lalu backend direstart. Tidak dapat diubah melalui dashboard atau API oleh role mana pun, termasuk Super Admin. Frontend menerima nomor hanya-baca dari API settings; `wa_number` tidak diterima pada PUT settings. Nilai lama di database tidak lagi menjadi sumber nomor. Alamat, jam buka, peta, dan pembatasan rasio poin tetap seperti sebelumnya. |

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

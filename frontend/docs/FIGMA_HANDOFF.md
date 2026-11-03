# Figma Handoff — Transhome E-Catalog UI/UX

File Figma: https://www.figma.com/design/3xWRuQOXPlxHNXdVJQnjGm (draft, paket Starter)

## Status

Pekerjaan berhenti karena kuota MCP paket Starter (20 panggilan/bulan) habis. Lanjutan butuh Pro/Full seat atau bulan berikutnya.

| Bagian | Status |
|---|---|
| Variabel (Primitives, Color, Spacing, Radius), 15 text style Inter | Selesai |
| Design System: Button, WhatsApp CTA, Input, Search, Badge, Tabs, Checkbox, Chip, Select, Breadcrumb, Pagination, Spec table, kartu Produk/Kategori/Brand/Ruangan/Artikel, Header, Mega menu, Footer, State (empty/error/401/403/5xx/skeleton), Toast, Dialog, Drawer filter, Admin sidebar item/Toggle/Data table | Selesai, diperiksa visual (kecuali Header/Footer mobile detail) |
| Homepage desktop + mobile (Section 02) | Dibuat, **belum diperiksa visual** |
| Catalog, Search, Product detail, Inspirasi, Brand, Promo/Artikel, Info toko | Belum |
| Section 03 (Auth + Member) | Kosong |
| Halaman C: Admin Portal (04) dan Prototype/Flows (05) | Kosong |
| Section 00 Overview (sitemap, journey, matriks audit) | Kosong (isi ada di dokumen ini dan `PETA_NAVIGASI.md`) |
| Prototype reactions | Belum |

## Penyimpangan dari prompt

- Starter hanya 3 halaman. Enam halaman diganti 3 halaman dengan Section bernomor 00–05: A (00, 01), B (02, 03), C (04, 05).
- Foto arsitektur belum ada: dipakai placeholder netral berlabel (generate_image memakai kredit AI, belum disetujui).
- Logo resmi belum ada: wordmark "Transhome" placeholder.

## Audit backend (sumber kebenaran)

| Fitur | Route frontend | Backend (`backend/api/transhome.api`) | Catatan |
|---|---|---|---|
| Login | `/masuk`, `/admin/masuk` | `POST /auth/login` | Tidak ada endpoint register, refresh, logout, lupa password di .api |
| Akses/menu | sidebar admin | `GET /auth/me/access` | Frontend mock memakai `/admin/menus`; selaraskan |
| Kategori, Brand, Ruangan (admin) | `/admin/kategori` dst | GET list, GET id, POST, PATCH | Tanpa delete (nonaktifkan via PATCH) |
| Produk (admin) | `/admin/produk` | GET list/id, POST, PATCH, `PUT /products/:id/rooms` | Field: sku, name, slug, categoryId, brandId, description, specifications (string), datasheetPdfUrl, priceGeneral, unitSale, minOrder, stockStatus, stockQtyLabel, isFeatured, isActive, metaTitle, metaDescription |
| RBAC | `/admin/pengguna`, `/admin/role` | access-control catalog, users, roles, user-roles, role-permissions | |
| Katalog publik, pencarian, filter, urut, paginasi | `/katalog`, `/cari`, `/produk/:slug` | **Tidak ada** | List admin tanpa parameter query dan tanpa metadata halaman; desain filter/paginasi mengikuti kontrak frontend (mock) |
| Foto produk, varian, produk terkait | detail produk | **Tidak ada di API** | Command foto ada di kode, route belum terdaftar. Jangan desain varian |
| Member, poin, artikel, banner, info toko, audit log, Qontak | `/akun`, `/admin/member` dst | **Tidak ada** | Hanya mock frontend; tandai dependensi |

## Konflik yang ditemukan

1. Prompt menyebut **3** tipe admin; database punya **4** role admin (SUPER_ADMIN, ADMIN_KATALOG, ADMIN_MEMBERSHIP, ADMIN_KONTEN) plus CUSTOMER. Desain admin memakai 4.
2. Permission backend yang di-seed: `rbac.manage` (SUPER_ADMIN), `catalog.category.manage`, `catalog.brand.manage`, `catalog.room.manage` (SUPER_ADMIN dan ADMIN_KATALOG). Frontend memakai kode granular (`product.view`, `product.update`, dll.). Kode izin produk tidak ditemukan di migrasi 000001–000004; verifikasi dengan tim backend.
3. Brand Red `#FF3300` tidak lolos AA untuk teks (3.67:1 dengan putih, 4.26:1 dengan Ink). Dipakai hanya sebagai aksen; teks merah memakai `#B82500` (6.36:1). Border `#E3E6EC` hanya dekoratif (1.25:1); border input memakai Slate `#647084`.

## Sisa pekerjaan (urut prioritas)

1. Periksa visual Homepage (desktop `10:27`, mobile `10:367`).
2. Catalog desktop/mobile + drawer filter + state kosong; Product detail desktop/mobile (harga tampil vs "konfirmasi via WhatsApp", stok habis, bar WhatsApp sticky di HP).
3. Search, Ruangan (list/detail), Brand (list/detail), Promo/Artikel, Info toko.
4. Section 03: Masuk, Daftar, Verifikasi, Lupa password, `/akun`, `/akun/poin`.
5. Halaman C: Admin login, Dashboard, daftar + form (Produk, Brand), Role matrix, 401/403/5xx per role.
6. Section 00 (sitemap, journey, matriks ini) dan prototype reactions untuk tiga alur utama.

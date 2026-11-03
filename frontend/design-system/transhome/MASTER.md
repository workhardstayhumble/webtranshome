# Design System Master File — Transhome

> **LOGIC:** When building a specific page, first check `design-system/transhome/pages/[page-name].md`.
> If that file exists, its rules **override** this Master file.
> If not, strictly follow the rules below.

---

**Project:** Transhome (Website Muka 1: katalog bahan bangunan & elektronik + membership Trans Family)
**Generated:** 2026-10-07 dengan ui-ux-pro-max (`"hardware store ecommerce product catalog" --design-system`), lalu diselaraskan dengan sketsa pemilik.
**Visual source of truth:** sketsa "Sketsa Dashboard Transhome" — https://claude.ai/artifact/S8vhfiAUswfBB85neKS1mZ
**Implementasi token:** `fase-0-fondasi/src/assets/main.css` (`@theme`). Komponen hanya memakai nama token, tanpa hex.

Dari hasil generator yang DIPAKAI: pola halaman "Trust & Authority + Conversion", gaya "Flat Design", skala spacing, anti-pattern, dan checklist.
Yang DIGANTI sketsa: palet (generator memberi hijau farmasi) dan font (generator memberi Rubik + Nunito Sans).

---

## Global Rules

### Color Palette

| Token Tailwind | Hex | Dipakai untuk |
|---|---|---|
| `foreground` | `#1C2B3A` | Teks utama (navy) |
| `secondary` | `#1C2B3A` | Tombol isi sekunder, top bar, footer, sidebar admin |
| `secondary-hover` | `#24364A` | Hover tombol sekunder, permukaan di atas navy |
| `secondary-foreground` | `#FFFFFF` | Teks di atas `secondary` |
| `primary` | `#C2410C` | CTA isi dengan teks putih (5.2:1), link, garis aktif |
| `primary-hover` | `#9A3412` | Hover CTA dan link |
| `primary-soft` | `#FFEDD5` | Latar chip/menu aktif, kartu ajakan |
| `primary-foreground` | `#FFFFFF` | Teks di atas `primary` |
| `primary-on-dark` | `#FB923C` | Aksen dan ring fokus di atas navy (logo "home", ikon) |
| `accent` | `#FACC15` | Label PROMO, Trans Family — HANYA dengan teks `accent-foreground` |
| `accent-foreground` | `#1C2B3A` | Teks di atas `accent` |
| `whatsapp` | `#15803D` | Semua tombol WhatsApp, teks putih |
| `whatsapp-hover` | `#166534` | Hover tombol WhatsApp |
| `background` | `#FAFAF7` | Latar situs publik |
| `background-admin` | `#F5F5F2` | Latar admin panel |
| `surface` | `#FFFFFF` | Kartu, header, input |
| `subtle` | `#F3F1EC` | Placeholder gambar, chip netral, hover baris |
| `border` | `#E7E5E0` | Garis pemisah, tepi kartu |
| `border-strong` | `#D6D3CD` | Tepi input dan tombol outline |
| `muted` | `#57534E` | Teks sekunder (7.4:1 di latar publik) |
| `faint` | `#78716C` | Caption kecil, hanya di atas `surface` atau `background` (≥ 4.5:1) |
| `on-dark` | `#E7E5E0` | Teks di atas navy |
| `on-dark-muted` | `#A8A29E` | Label kecil di atas navy (sidebar, footer) |
| `ring` | `#C2410C` | Ring fokus di latar terang |

Warna status (selalu bersama label teks + ikon, tidak boleh warna saja):

| Status | Latar (`*-soft`) | Teks (`*-ink`) | Warna isi |
|---|---|---|---|
| `success` (TERSEDIA, berhasil) | `#DCFCE7` | `#14532D` | `#15803D` |
| `warning` (SISA STOK, perlu perhatian) | `#FFEDD5` | `#9A3412` | `#B45309` |
| `info` (PRE-ORDER, info) | `#DBEAFE` | `#1E3A8A` | `#1D4ED8` |
| `danger` (HABIS, error, destruktif) | `#FEE2E2` | `#7F1D1D` | `#B91C1C` |

**Color Notes:** navy + oranye bata dari sketsa; kuning hanya untuk promo dan Trans Family; hijau khusus WhatsApp. Warna logo resmi Transhome masih ditanyakan (CLAUDE.md §12).

### Typography

- **Font utama:** Plus Jakarta Sans (400, 500, 600, 700, 800) — judul, UI, dan isi.
- **Font kode:** IBM Plex Mono (500, 600) — SKU, `member_code`, nomor nota.
- **Angka harga:** `tabular-nums`.
- **Skala:** 12 / 13 / 14 / 16 (isi) / 18 / 20 / 24 / 30 / 40. Isi 16px (sketsa 15px; dinaikkan agar input di HP tidak memicu zoom).
- **Judul:** tebal 700–800, `letter-spacing: -0.02em`, line-height 1.15. Isi line-height 1.5.
- **Google Fonts:** dimuat di `index.html` dengan `display=swap`.

### Spacing Variables

Grid 8pt, memakai skala bawaan Tailwind.

| Token | Value | Usage |
|-------|-------|-------|
| `--space-xs` | `4px` / `0.25rem` | Jarak ikon dan teks |
| `--space-sm` | `8px` / `0.5rem` | Jarak inline, antar chip |
| `--space-md` | `16px` / `1rem` | Padding standar, gutter HP |
| `--space-lg` | `24px` / `1.5rem` | Padding section, gutter desktop |
| `--space-xl` | `32px` / `2rem` | Jarak besar |
| `--space-2xl` | `48px` / `3rem` | Antar section |
| `--space-3xl` | `64px` / `4rem` | Hero |

### Layout
- Container maksimal 1280px, gutter 16px (HP) / 24px (desktop).
- Breakpoint: 375 (dasar) · 640 · 768 · 1024 · 1280 · 1440.

### Radius & Shadow
| Token | Nilai | Dipakai untuk |
|---|---|---|
| `rounded-md` | 6px | Tombol, input, select |
| `rounded-lg` | 8px | Thumbnail, item menu |
| `rounded-xl` | 12px | Kartu, panel |
| `rounded-2xl` | 16px | Kartu member, modal, bottom sheet |
| `rounded-full` | — | Chip filter, avatar |

Flat: kartu memakai border, bukan bayangan. Bayangan hanya untuk lapisan yang melayang (dropdown, modal, toast, tombol WhatsApp melayang).

---

## Component Specs

### Buttons
| Varian | Latar | Teks | Dipakai untuk |
|---|---|---|---|
| `primary` | `primary` | `primary-foreground` | Aksi utama (Daftar member, Simpan, Jelajahi katalog) |
| `secondary` | `secondary` | `secondary-foreground` | Aksi tegas netral (Cari di admin, Unduh) |
| `outline` | `surface` + `border-strong` | `foreground` | Aksi kedua (Masuk, Batal, Bagikan) |
| `ghost` | transparan | `foreground` / `primary` | Aksi ringan (Reset, Lihat semua) |
| `whatsapp` | `whatsapp` | putih | SEMUA tombol WhatsApp |
| `danger` | `danger` | putih | Aksi destruktif di modal konfirmasi |

- Tinggi: 44px (default), 40px (`sm`, tabel admin), 56px (`lg`, CTA detail produk). Area sentuh minimal 44×44px di HP.
- Radius 6px, font 600 14–16px, ikon 18px dengan jarak 8px.
- Loading: spinner + `aria-busy`, tombol nonaktif sampai selesai (cegah double submit).
- Hover: perubahan warna 150ms; tanpa `translate` (flat). Fokus: `focus-visible:ring-2 ring-ring ring-offset-2`.
- Ditekan (situs publik saja): mengecil `scale(0.97)`, 150ms. Otomatis lewat atribut `data-tombol` di `Button`.

### Cards
```
bg-surface border border-border rounded-xl
hover (kartu yang bisa diklik): border-border-strong, transisi 150ms
```

### Inputs
```
h-11 rounded-md border border-border-strong bg-surface px-3 text-base
focus-visible: ring-2 ring-ring, border-primary
error: border-danger, pesan di bawah field (text-danger-ink), aria-invalid + aria-describedby
```
Label selalu terlihat di atas field. Placeholder bukan pengganti label.

### Badges & Chips
- Badge: 12px semibold, `px-2 py-0.5`, radius 4px, warna `*-soft` + `*-ink`.
- Chip filter aktif: `bg-primary-soft text-primary-hover border-primary`, rounded-full, tinggi ≥ 34px, tombol hapus dengan `aria-label`.

### Modals
- Lapisan latar `bg-foreground/50` (tanpa blur), modal `bg-surface rounded-2xl`, lebar maks 500px.
- `role="dialog" aria-modal="true"`, fokus pindah ke modal dan kembali ke pemicu saat ditutup, Esc menutup.

### Toast
- Pojok kanan bawah (HP: bawah, di atas tombol WhatsApp melayang), `role="status"` dengan pesan utuh ("Produk disimpan"), hilang otomatis 4 detik, ada tombol tutup.

---

## Style Guidelines

**Style:** Flat Design

**Keywords:** 2D, minimalist, clean lines, simple shapes, typography-focused, modern

**Key Effects:** Tanpa gradien; hover hanya perubahan warna/opacity; transisi 150–200ms ease-out; hormati `prefers-reduced-motion`. (Hero rumah 3D di Beranda yang dulu jadi pengecualian sudah dihapus 30 Okt 2026, KEPUTUSAN #16.)

### Hero Beranda (30 Okt 2026)
Wireframe 01 (`docs/wireframe-low-fidelity-transhome.xlsx`): media selebar layar, judul + subjudul di kiri, tombol di kanan bawah ("CTA kanan bawah"). Komponen `fase-5-konten/src/components/beranda/HeroBanner.vue`, isi dari banner `HOME_HERO`.

- Hero = landing page: **setinggi layar**. Tinggi `max(26rem, 100svh − tinggi header)` (`--tinggi-header-penuh` dari `AppHeader`), jadi header + hero pas satu layar saat Beranda dibuka; judul di tengah tinggi sisa. Satu `h1` (judul banner), maksimal 100 karakter; subjudul satu kalimat (≤ 160).
- Tombol (kanan bawah, di HP bertumpuk selebar layar, semua `lg` 56px): **tombol utama** (`primary`, label & tautan dari banner, bawaan "Jelajahi katalog" → `/katalog`), lalu **tombol kedua** opsional (`outline`, tampil bila admin mengisi label + tautan). Di tengah bawah: **Lihat selengkapnya** + ikon panah, menggulir ke seksi berikutnya dan memindahkan fokus ke sana. Tidak ada tombol WhatsApp di hero (2 Nov 2026); WhatsApp ada di header dan tombol melayang.
- Media: gambar 1920 × 800 px (wajib) dan video MP4/WebM opsional. Selalu ditutup lapisan `secondary` 70% supaya teks putih tetap kontras di atas foto/video apa pun.
- Video: `autoplay muted loop playsinline`, `preload="metadata"`, gambar jadi `poster`. Dekoratif (`aria-hidden`, tidak bisa difokus). Tombol **Jeda video / Putar video** (≥ 44px) selalu ada bila video tampil.
- "Kurangi animasi" aktif → video tidak berputar sendiri (pengguna masih bisa menekan Putar). Mode hemat data (`navigator.connection.saveData`) atau video gagal dimuat → video tidak dipakai, gambar saja.
- Tanpa banner tayang: latar `secondary` polos + teks bawaan + "Jelajahi katalog".

### Animasi situs publik (8 Okt 2026)
Tanpa library: CSS + View Transitions API + IntersectionObserver (`fase-0-fondasi/src/composables/useAnimasi.js`, gaya di `main.css`). Hanya `transform` dan `opacity`. Panel admin TIDAK dianimasikan (`<html data-area="admin">`). Semua mati saat `prefers-reduced-motion: reduce`.

| Efek | Durasi & easing | Dipakai di |
|---|---|---|
| Pindah halaman (crossfade, halaman baru naik 8px) | keluar 150ms ease-in, masuk 300ms `cubic-bezier(0.22, 1, 0.36, 1)` | Semua pindah halaman publik; ganti filter/query di halaman yang sama tidak. Header & tombol WhatsApp melayang tetap diam (`view-transition-name`) |
| Scroll reveal `v-reveal` (fade + naik 20px) | 500ms ease-out-quint, jeda bergantian 60ms (maks 5 langkah) | Seksi Beranda di bawah hero, kartu katalog (`v-reveal="i % 4"`), kartu artikel/promo (`i % 3`). Elemen yang sudah terlihat saat halaman dibuka langsung tampil |
| Tombol ditekan | `scale(0.97)`, 150ms | `Button` di situs publik |
| Modal muncul | 220ms ease-out-quint, membesar dari 0.97 + naik 8px | Modal publik |
| Sheet (HP) | 300ms ease-out-quint, naik dari bawah | Modal `sheet` (filter katalog) di bawah 640px |

Tidak dipakai: hover lift (`translate` saat hover), smooth scroll buatan (Lenis dkk.) karena mengambil alih scroll bawaan dan bentrok dengan header yang menempel.

### Page Pattern

**Pattern Name:** Trust & Authority + Conversion (disesuaikan Muka 1)

- **Kepercayaan:** toko fisik di Sleman, brand resmi, garansi/SNI di detail produk, poin Trans Family.
- **Konversi:** satu CTA utama per produk, yaitu "Tanya via WhatsApp". Tidak ada "Beli", "Keranjang", atau form lead.
- **Brand pilihan (2 Nov 2026, KEPUTUSAN #19):** panel putih bersudut 16px, judul `h2` di tengah + tombol `primary` `lg`, diapit lingkaran (diameter 80/96/112/128px per breakpoint) berlatar token lembut (`primary-soft`, `info-soft`, `accent`, `success-soft`, `danger-soft`, `on-dark`), tetangga selalu beda warna. Logo `object-contain` dengan `mix-blend-multiply` supaya latar putih logo melebur. Hover/fokus: cincin `primary`/`ring`, tanpa gerak. Hiasan pengisi: dua warna datar (garis lurus, bukan gradien halus). Komponen `BrandPilihan.vue` + `LingkaranBrand.vue`.
- **Urutan Beranda (sketsa):** hero banner → banner promo → kategori per tahap proyek → produk unggulan → promo & event → inspirasi ruangan → brand pilihan → keunggulan → footer.
- **Admin:** padat data, sidebar navy, latar `background-admin`, tabel dengan header kecil huruf kapital.

---

## Anti-Patterns (Do NOT Use)

- ❌ Glassmorphism, blur latar, gradien ungu-pink
- ❌ Emoji sebagai ikon (pakai `lucide-vue-next`)
- ❌ Hex atau warna arbitrary di komponen (`bg-[#…]`)
- ❌ Status hanya dibedakan warna
- ❌ Tombol "Beli", "Tambah ke Keranjang", "Wishlist", "Ingatkan saya", "Tulis ulasan"
- ❌ Teks `accent` (kuning) dengan warna putih
- ❌ Interaksi yang hanya bisa lewat hover

### Additional Forbidden Patterns

- ❌ **Missing cursor:pointer** — All clickable elements must have cursor:pointer
- ❌ **Layout-shifting hovers** — Avoid scale transforms that shift layout
- ❌ **Low contrast text** — Maintain 4.5:1 minimum contrast ratio
- ❌ **Instant state changes** — Always use transitions (150-300ms)
- ❌ **Invisible focus states** — Focus states must be visible for a11y

---

## Pre-Delivery Checklist

Before delivering any UI code, verify:

- [ ] No emojis used as icons (use SVG instead)
- [ ] All icons from `lucide-vue-next`
- [ ] `cursor-pointer` on all clickable elements
- [ ] Hover states with smooth transitions (150-300ms)
- [ ] Light mode: text contrast 4.5:1 minimum
- [ ] Focus states visible for keyboard navigation (on navy: `ring-primary-on-dark`)
- [ ] `prefers-reduced-motion` respected
- [ ] Responsive: 375px, 768px, 1024px, 1440px
- [ ] No content hidden behind fixed navbars or the floating WhatsApp button
- [ ] No horizontal scroll on mobile
- [ ] Loading (skeleton), empty, and error states present

# MASTER PROMPT — UI/UX Website Transhome Muka 1 (Antigravity)

> Cara pakai: paste bagian "MASTER PROMPT" ke Antigravity sekali di awal, dari folder `frontend/`.
> Lalu kerjakan per fase memakai prompt di `docs/FRONTEND_BUILD_PLAN.md` (urutan Fase 0–6 hanya ada di sana).
> Jangan minta semua fase sekaligus — hasilnya akan dangkal.
> Disesuaikan ke Vue 3 + Vite dan token sketsa pada 7 Okt 2026 (KEPUTUSAN #7).

---

## MASTER PROMPT

```
[Context — READ FIRST]
Before doing anything, read:
- CLAUDE.md (project rules — these OVERRIDE anything in this prompt)
- docs/KEPUTUSAN.md (latest decisions — OVERRIDE BRD/PRD)
- docs/FRONTEND_BUILD_PLAN.md (phases, buttons, API calls, QA checklists)
- docs/transhome_postgres.sql (data shapes; UI fields and mocks must map to these columns)
- docs/BRD_Website_Transhome_Muka1_v0.4_DRAFT.pdf and docs/PRD_Website_Transhome_Muka1_v0.4_DRAFT.pdf
Use the ui-ux-pro-max skill with --stack vue for every UI decision.

[Role & Goal]
Act as a Principal UI/UX Designer and Lead Design System Architect for building-material retail.
Design and build the UI for "Transhome" — a single-store building materials & home-goods retailer
in Sleman, Yogyakarta, Indonesia. This is MUKA 1: a digital catalog + "Trans Family" membership website.
There is NO online transaction in Muka 1. The ONLY conversion goal is a WhatsApp chat via Mekari Qontak:
"Tanya via WhatsApp" with the product SKU prefilled. Project consultation also happens in WhatsApp (no lead form).
Stack: Vue 3 (<script setup>, JavaScript) + Vite + Vue Router + Pinia + Tailwind CSS v4 + lucide-vue-next.
No UI component library: build base components yourself in src/components/ui/.
All UI copy in Bahasa Indonesia.

[Design References]
- Visual source of truth: the owner's sketch "Sketsa Dashboard Transhome"
  (https://claude.ai/artifact/S8vhfiAUswfBB85neKS1mZ) — 5 customer pages and 3 admin pages.
- Material Bank → speed of material discovery, dense but calm filtering, instant spec access.
- Kohler / Porcelanosa → refined product photography, generous whitespace on PDP, premium restraint.
- Leroy Merlin → deep category hierarchy, clear mega menu, practical catalog grid.
Take the PRINCIPLES, not the branding. No copying logos, layouts 1:1, or brand assets.

[Scope Guard — Muka 1 adaptation of the original brief]
| Original idea                     | Muka 1 decision                                                        |
|-----------------------------------|------------------------------------------------------------------------|
| Fast checkout / cart              | REMOVED. Primary CTA = "Tanya via WhatsApp"                            |
| Material calculator (m² → boxes)  | REMOVED (eliminated in BRD v0.3). Show static spec "isi per dus / coverage" only |
| Volume pricing matrix / tiers     | REMOVED (single price_general). Replace with "Butuh jumlah besar? Chat tim proyek" card → WhatsApp with bulk-inquiry message |
| Floating RFQ bar                  | → Floating WhatsApp widget only (FR-R, FR-S). NO consultation form     |
| Wishlist, stock reminder, reviews | REMOVED (out of scope Muka 1)                                          |
| Reward catalog / member redeem UI | REMOVED. Admin records REDEEM; member only sees balance and history    |
| 3D product visualizer             | NOT in Muka 1. Photo gallery + zoom + optional video. (The exploded-house hero was removed on 30 Oct 2026; the home page now uses an image/video hero banner, KEPUTUSAN #16.) |
| CAD/BIM downloads                 | Only datasheet PDF exists (datasheet_pdf_url). Design the drawer to hold more doc types later |
| "Pesan Sampel Fisik"              | Not in BRD. Do NOT build. Listed in Open Questions                     |
| Filter by stock location          | REMOVED (single store)                                                 |
| Tier/member pricing               | REMOVED (all members = Tahap 1)                                        |
Never add buttons like "Beli", "Tambah ke Keranjang", "Checkout", "Wishlist", "Ingatkan saya", "Tulis ulasan", or price tiers.

[Personas]
1. Homeowner & Tukang/Kontraktor — mostly on mobile, low patience, wants: find product fast,
   see price + satuan + stock status, one tap to WhatsApp.
2. Arsitek / Kontraktor proyek / B2B — mostly desktop, wants: spec table, datasheet download,
   compare by technical attributes, chat the project team on WhatsApp for bulk needs.
3. Member Trans Family — wants: see member_code (to tell the cashier), points balance, riwayat poin.
4. Admin (4 roles) — wants: dense, fast data tables.

[Key Layouts]
1. Header & Navigation
   - Top utility bar: jam buka, alamat singkat, link Info Toko, separate "Masuk" and "Daftar" buttons
     (optional — never gate the catalog), or member avatar + saldo poin when logged in.
   - Mega menu "Kategori" organized by project phase, mapped to categories (parent → child):
     Struktur · Dinding & Lantai · Sanitasi · Atap · Perkakas & Hardware · (others from API).
     Each column: subcategories + 1 featured image tile. Second entry "Ruangan" (Dapur, Kamar Mandi, Fasad, ...).
   - Omnibox search: matches SKU, product name, brand, and spec attributes (incl. SNI if present in
     specifications JSON). Suggestion dropdown grouped: Produk / Kategori / Brand. Keyboard navigable.
   - Mobile: hamburger → full-height drawer with accordion categories; search is a full-screen overlay.
2. Home
   - Exploded-house hero (Fase 5), hero banner slider (banners with start_at/end_at), category tiles by
     project phase, "Inspirasi Ruangan" carousel (FR-F), produk unggulan (is_featured), artikel terbaru,
     trust strip (toko fisik, konsultasi via WhatsApp, poin Trans Family).
3. Catalog & Filtering (FR-A, FR-B, FR-D)
   - Desktop: sticky left sidebar filter (280px). Mobile: bottom-sheet/drawer filter with "Terapkan (n)".
   - Filters: kategori (tree), brand, status stok, ruangan, rentang harga, and dynamic spec attributes
     from specifications JSON (ketebalan, finishing, area aplikasi, ketahanan cuaca, ukuran).
   - Active filter chips above results + "Hapus semua". Result count. Sort (terbaru, harga, nama).
   - Grid/List toggle. Grid card: foto, brand, nama (2 lines max), harga + "/satuan", stock badge,
     garansi/SNI badge if present. List row adds key specs + datasheet icon.
   - Filter state synced to URL query (shareable, back-button safe).
4. Product Detail Page (FR-C)
   - Left: gallery (main image + thumbs), click-to-zoom lightbox, optional video tab.
   - Right action block: brand, nama, SKU (copyable), price_general + unit_sale, min order,
     stock badge (TERSEDIA / SISA STOK + stock_qty_label / PRE-ORDER / HABIS).
     Primary: "Tanya via WhatsApp" (prefill: SKU + nama + URL). Secondary (ghost): "Salin SKU", "Bagikan".
     If HABIS → show "Tanya ketersediaan" WhatsApp message variant instead.
   - "Butuh jumlah besar?" card → WhatsApp with bulk-inquiry message (SKU + "kebutuhan proyek").
   - Tabs/sections: Deskripsi · Spesifikasi (2-col key/value table from JSON) · Dokumen (drawer/list:
     Datasheet PDF now; Sertifikat & Panduan Pasang slots hidden when empty) · Cocok untuk Ruangan.
   - NEVER show internal fields: HPP, hnet, diskon, stok per lokasi, konsinyasi, karantina.
   - Related products from same category. Mobile: sticky bottom action bar (WhatsApp full width).
5. WhatsApp Widget (FR-R, FR-S — sub-fase 1c, design now, wire later)
   - Floating WhatsApp button bottom-right (mobile: above sticky bar, never overlapping it).
   - When logged in, pass nama + member_code as context to the chat (FR-S).
   - No consultation form anywhere.
6. Auth & Member (FR-E, FR-H)
   - /masuk: email + password, plus a separate "Masuk dengan Google" button. /daftar: nama, EMAIL (wajib),
     password, nomor HP (opsional), consent checkbox. /lupa-password via email.
   - /akun: member_code shown large + copy button, saldo poin, riwayat poin (EARN/REDEEM with reference_no).
   - Admin login lives separately at /admin/masuk.

[Design Tokens — from the sketch; a token not listed here still follows the sketch]
- Spacing: 8pt grid (Tailwind 2/4/6/8/12/16 = 8/16/24/32/48/64px). 4px only for icon gaps.
- Layout: 12-column grid, container max 1280px, gutter 24px desktop / 16px mobile.
- Breakpoints: 375 (base) · 640 sm · 768 md · 1024 lg · 1280 xl · 1440 2xl.
- Color (define once in Tailwind @theme as semantic tokens; NO raw hex in components):
  ink          #1C2B3A  (navy — text, dark top bar/footer, admin sidebar, filled secondary buttons)
  background   #FAFAF7  (public page bg) · admin bg #F5F5F2 · card #FFFFFF · subtle fill #F3F1EC
  border       #E7E5E0  (dividers) · #D6D3CD (inputs, outline buttons) · muted text #57534E
  primary      #C2410C  (filled CTA with white text, 5.2:1) · primary-hover #9A3412 · primary-soft #FFEDD5
  accent       #FACC15  (promo label, Trans Family; ONLY with ink text, never white text)
  whatsapp     #15803D  (every WhatsApp button, white text)
  stock badges, always with label text + icon:
    tersedia bg #DCFCE7 / text #14532D · sisa-stok #FFEDD5 / #9A3412 · pre-order #DBEAFE / #1E3A8A · habis #FEE2E2 / #7F1D1D
  Official Transhome logo colors are still an open question; these come from the sketch.
- Typography: Plus Jakarta Sans for all UI and body text; IBM Plex Mono for SKU, member_code and nomor nota;
  tabular-nums for prices. Scale: 12 / 14 / 16 (body) / 18 / 20 / 24 / 30 / 36 / 48.
  Line-height 1.5 body, 1.2 headings.
- Radius: 6px inputs/buttons, 12px cards, 16px sheets/modals. Shadow: sm (rest), md (hover), lg (drawer/modal).
- Motion: 150ms hover, 200–250ms drawers (ease-out), card hover = translateY(-2px) + shadow-md.
  Respect prefers-reduced-motion (disable translate, keep color change).

[Micro-interactions]
- Card hover: elevation + image scale 1.03 inside overflow-hidden frame.
- Filter apply: result count animates, skeleton grid while loading, focus moves to results heading.
- Mobile filter: bottom sheet slides up, backdrop fades, focus trapped, Esc/drag-down closes.
- WA click: brief "Membuka WhatsApp…" toast. Copy SKU: "SKU disalin" toast.
- Form submit: button spinner + disabled, inline field errors, success state replaces form.

[Deliverables — per phase, saved in repo]
For each page/component in the current phase:
1. Spec file: design-system/transhome/pages/<page>.md containing
   - purpose & persona, data mapping to SQL columns,
   - ASCII wireframe for desktop (1440) AND mobile (375),
   - component list with props, states (default/hover/focus/active/disabled/loading/empty/error),
   - breakpoint behavior table,
   - Mermaid stateDiagram for every interactive component (filter drawer, gallery, search, auth forms).
2. Vue implementation in src/ using mock data in src/services/mock/ shaped exactly like docs/transhome_postgres.sql.
3. Self-review against the checklist in CLAUDE.md §11 and the QA checklist of the phase; report pass/fail per item.
STOP after each phase and wait for my approval before continuing.

[Open Questions — ask, never assume]
- Brand colors / logo Transhome resmi (tokens above come from the sketch).
- Which of the 58 Master Barang columns appear on the PDP (customer-facing only).
- Which spec attributes exist in specifications JSON (needed for dynamic filters).
- Sertifikat / panduan pasang: will columns be added beyond datasheet_pdf_url?
```

---

## PROMPT PER FASE

Urutan dan isi fase ada di `docs/FRONTEND_BUILD_PLAN.md` (Fase 0 Fondasi s.d. Fase 6 Sistem). Prompt untuk tiap fase:

```
Baca CLAUDE.md, docs/KEPUTUSAN.md, dan docs/FRONTEND_BUILD_PLAN.md.
Kerjakan FASE <n> saja sesuai dokumen. Pakai data mock. Berhenti setelah
selesai dan laporkan hasil checklist QA Fase <n>.
```

Admin panel mengikuti gaya data-dense di `CLAUDE.md` §8 dan tiga artboard sisi admin di sketsa.

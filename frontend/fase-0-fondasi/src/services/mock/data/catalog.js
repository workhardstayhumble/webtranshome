// Data contoh katalog. Kolom mengikuti docs/transhome_postgres.sql. TODO: ganti ke API.
// Nama merek sengaja generik. Foto berupa gambar SVG contoh (bukan foto produk asli).

const dibuat = '2026-09-20T09:00:00'

/** Gambar contoh polos berlabel, dipakai sebagai foto produk/cover ruangan di mode mock. */
function gambarContoh(judul, keterangan, latar = '#E7E5E0') {
  const svg =
    `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="800" viewBox="0 0 800 800">` +
    `<rect width="800" height="800" fill="${latar}"/>` +
    `<rect x="40" y="40" width="720" height="720" fill="none" stroke="#FFFFFF" stroke-opacity="0.6" stroke-width="4"/>` +
    `<text x="400" y="380" text-anchor="middle" font-family="sans-serif" font-size="44" font-weight="700" fill="#1C2B3A">${judul}</text>` +
    `<text x="400" y="440" text-anchor="middle" font-family="sans-serif" font-size="28" fill="#57534E">${keterangan}</text>` +
    `<text x="400" y="700" text-anchor="middle" font-family="sans-serif" font-size="22" fill="#57534E">CONTOH FOTO</text>` +
    `</svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

const slugify = (teks) =>
  teks.toLowerCase().replace(/×/g, 'x').replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')

const kategoriInduk = [
  [1, 'Struktur', 'struktur'],
  [2, 'Dinding & Lantai', 'dinding-lantai'],
  [3, 'Sanitasi', 'sanitasi'],
  [4, 'Atap', 'atap'],
  [5, 'Cat', 'cat'],
  [6, 'Listrik', 'listrik'],
  [7, 'Perkakas & Hardware', 'perkakas-hardware'],
  [8, 'Elektronik', 'elektronik'],
]

const kategoriAnak = [
  [11, 1, 'Semen', 'semen'],
  [12, 1, 'Besi beton', 'besi-beton'],
  [13, 1, 'Bata ringan', 'bata-ringan'],
  [21, 2, 'Keramik lantai', 'keramik-lantai'],
  [22, 2, 'Keramik dinding', 'keramik-dinding'],
  [23, 2, 'Granit', 'granit'],
  [24, 2, 'Perekat keramik', 'perekat-keramik'],
  [31, 3, 'Kloset', 'kloset'],
  [32, 3, 'Kran', 'kran'],
  [33, 3, 'Wastafel', 'wastafel'],
  [41, 4, 'Genteng', 'genteng'],
  [42, 4, 'Rangka baja ringan', 'rangka-baja-ringan'],
  [43, 4, 'Talang', 'talang'],
  [51, 5, 'Cat tembok', 'cat-tembok'],
  [52, 5, 'Cat kayu & besi', 'cat-kayu-besi'],
  [61, 6, 'Lampu', 'lampu'],
  [62, 6, 'Kabel', 'kabel'],
  [63, 6, 'Saklar & stop kontak', 'saklar-stop-kontak'],
  [71, 7, 'Engsel & kunci', 'engsel-kunci'],
  [72, 7, 'Alat tukang', 'alat-tukang'],
  [81, 8, 'Pemanas air', 'pemanas-air'],
  [82, 8, 'Kipas angin', 'kipas-angin'],
]

export const categories = [
  ...kategoriInduk.map(([id, name, slug], i) => ({
    id, name, slug, parent_id: null, image_url: null, sort_order: i + 1, is_active: true, created_at: dibuat, updated_at: dibuat,
  })),
  ...kategoriAnak.map(([id, parent_id, name, slug], i) => ({
    id, name, slug, parent_id, image_url: null, sort_order: i + 1, is_active: true, created_at: dibuat, updated_at: dibuat,
  })),
]

// Merek A–F dipakai produk contoh. Sisanya brand fiktif tanpa produk (nama huruf Yunani, bukan merek sungguhan)
// supaya tampilan banyak brand (Beranda berjalan, /brand A–Z, cari di admin) bisa dicoba. "Merek Lama" nonaktif.
const BRAND_TAMBAHAN = [
  'Alfa Keramik', 'Beta Sanitair', 'Gama Cat', 'Delta Pipa', 'Epsilon Baja', 'Zeta Lampu', 'Eta Kabel', 'Teta Atap',
  'Iota Kunci', 'Kapa Kaca', 'Lambda Granit', 'Mu Engsel', 'Nu Semen', 'Ksi Plafon', 'Omikron Pintu', 'Pi Perkakas',
  'Rho Elektrik', 'Sigma Kran', 'Tau Bata', 'Upsilon Lantai', 'Fi Pompa', 'Khi Lem', 'Psi Panel', 'Omega Tangga',
]
// Logo contoh berupa SVG persegi (tanda + nama), BUKAN logo merek asli. 10 brand pertama saja, supaya tampilan
// berlogo dan tanpa logo (nama brand) sama-sama bisa dicoba di "Brand pilihan" Beranda dan /brand.
const WARNA_LOGO = ['#C2410C', '#1D4ED8', '#15803D', '#7C3AED', '#B91C1C', '#0F766E', '#B45309', '#1C2B3A', '#BE185D', '#4D7C0F']
function logoContoh(nama, i) {
  const w = WARNA_LOGO[i % WARNA_LOGO.length]
  const huruf = nama.split(' ').map((k) => k[0]).join('').slice(0, 2).toUpperCase()
  const svg =
    `<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200" viewBox="0 0 200 200">` +
    `<rect x="58" y="28" width="84" height="84" rx="${i % 2 ? 42 : 18}" fill="${w}"/>` +
    `<text x="100" y="84" text-anchor="middle" font-family="sans-serif" font-size="38" font-weight="800" fill="#fff">${huruf}</text>` +
    `<text x="100" y="152" text-anchor="middle" font-family="sans-serif" font-size="24" font-weight="800" fill="${w}">${nama.toUpperCase()}</text>` +
    `</svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

export const brands = [
  ...['Merek A', 'Merek B', 'Merek C', 'Merek D', 'Merek E', 'Merek F'],
  ...BRAND_TAMBAHAN,
  'Merek Lama',
].map((name, i) => ({
  id: i + 1, name, slug: slugify(name), logo_url: i < 10 ? logoContoh(name, i) : null, is_active: name !== 'Merek Lama', created_at: dibuat, updated_at: dibuat,
}))

const warnaRuangan = { Dapur: '#E7D7C1', 'Kamar mandi': '#D4E2EA', Fasad: '#D6D3CD', 'Ruang keluarga': '#E9DFCE', Teras: '#DCE3D2' }

export const rooms = ['Dapur', 'Kamar mandi', 'Fasad', 'Ruang keluarga', 'Teras'].map((name, i) => ({
  id: i + 1, name, slug: slugify(name), image_cover: gambarContoh(name, 'Inspirasi ruangan', warnaRuangan[name]),
  sort_order: i + 1, is_active: true, created_at: dibuat, updated_at: dibuat,
}))

// [sku, nama, kategori, merek, harga, satuan, min order, status stok, label stok, unggulan, aktif, ruangan, spesifikasi]
const produkUtama = [
  ['KRM-4040-01', 'Keramik lantai abu doff 40×40', 21, 1, '89500.00', 'dus', 1, 'TERSEDIA', null, true, true, [4, 1],
    // hpp & stok_gudang SENGAJA ada: menguji frontend menyembunyikan field internal (CLAUDE.md §7).
    { ukuran: '40 × 40 cm', finishing: 'Doff', warna: 'Abu', isi_per_dus: '6 pcs (0,96 m²)', area_aplikasi: 'Lantai dalam ruangan', garansi: '5 tahun', sni: 'Ada', hpp: '61000', stok_gudang: '120 dus' }],
  ['KRM-4040-07', 'Keramik lantai motif kayu 40×40', 21, 1, '96000.00', 'dus', 1, 'SISA_STOK', '3 dus', false, true, [4, 5],
    { ukuran: '40 × 40 cm', finishing: 'Doff', warna: 'Cokelat kayu', isi_per_dus: '6 pcs (0,96 m²)', area_aplikasi: 'Lantai dalam ruangan' }],
  ['KRM-5050-02', 'Keramik lantai putih glossy 50×50', 21, 2, '112000.00', 'dus', 1, 'TERSEDIA', null, false, true, [4],
    { ukuran: '50 × 50 cm', finishing: 'Glossy', warna: 'Putih', isi_per_dus: '4 pcs (1 m²)', area_aplikasi: 'Lantai dalam ruangan' }],
  ['KRM-6060-03', 'Keramik lantai teraso 60×60', 21, 1, '168000.00', 'dus', 1, 'PRE_ORDER', null, false, true, [4, 5],
    { ukuran: '60 × 60 cm', finishing: 'Doff', warna: 'Teraso', isi_per_dus: '4 pcs (1,44 m²)', area_aplikasi: 'Lantai dalam dan luar ruangan' }],
  ['SAN-KLS-02', 'Kloset duduk monoblok putih', 31, 3, '3250000.00', 'unit', 1, 'SISA_STOK', '3 unit', true, true, [2],
    { tipe: 'Monoblok', warna: 'Putih', sistem_bilas: 'Dual flush', garansi: '2 tahun' }],
  ['STR-SMN-50', 'Semen PCC 50 kg', 11, 4, '68000.00', 'sak', 1, 'TERSEDIA', null, true, true, [3],
    { berat: '50 kg', jenis: 'PCC', sni: 'Ada' }],
  ['SAN-KRN-11', 'Kran wastafel stainless', 32, 3, '245000.00', 'pcs', 1, 'PRE_ORDER', null, true, true, [1, 2],
    { bahan: 'Stainless steel', ukuran_drat: '1/2 inch' }],
  ['CAT-INT-05', 'Cat tembok interior 5 kg', 51, 5, '185000.00', 'galon', 1, 'TERSEDIA', null, true, true, [4],
    { berat: '5 kg', finishing: 'Matt', warna: 'Putih', area_aplikasi: 'Dinding dalam ruangan' }],
  ['HRD-ENG-04', 'Engsel pintu 4 inch', 71, 6, '32000.00', 'pasang', 1, 'HABIS', null, true, true, [],
    { ukuran: '4 inch', bahan: 'Stainless steel' }],
  ['ATP-GMP-01', 'Genteng metal pasir', 41, 2, '52000.00', 'lembar', 1, 'TERSEDIA', null, false, false, [3],
    { bahan: 'Baja lapis pasir', warna: 'Merah bata' }],
  ['ELK-LED-12', 'Lampu LED 12 watt', 61, 5, '45000.00', 'pcs', 1, 'SISA_STOK', '8 pcs', true, true, [1, 4],
    { daya: '12 watt', warna: 'Putih (6500K)' }],
  ['STR-BTR-10', 'Bata ringan 10 cm', 13, 4, '9500.00', 'pcs', 50, 'TERSEDIA', null, false, true, [3],
    { ukuran: '60 × 20 × 10 cm' }],
  ['ELK-PMA-15', 'Pemanas air listrik 15 liter', 81, 6, '1450000.00', 'unit', 1, 'TERSEDIA', null, false, true, [2],
    { kapasitas: '15 liter', daya: '350 watt', garansi: '1 tahun' }],
  ['KRM-3060-05', 'Keramik dinding putih 30×60', 22, 2, '104000.00', 'dus', 1, 'TERSEDIA', null, false, true, [1, 2],
    { ukuran: '30 × 60 cm', finishing: 'Glossy', warna: 'Putih', isi_per_dus: '8 pcs (1,44 m²)', area_aplikasi: 'Dinding' }],
]

// Varian tambahan agar katalog punya lebih dari satu halaman (24 per halaman) dan filter berisi.
const POLA_STOK = ['TERSEDIA', 'TERSEDIA', 'SISA_STOK', 'TERSEDIA', 'PRE_ORDER', 'TERSEDIA', 'HABIS']
const varian = [
  ...[['40 × 40 cm', 'Krem', 'Glossy', '78000'], ['40 × 40 cm', 'Putih', 'Doff', '82000'], ['50 × 50 cm', 'Abu tua', 'Doff', '118000'],
    ['50 × 50 cm', 'Batu alam', 'Doff', '124500'], ['60 × 60 cm', 'Abu', 'Glossy', '159000'], ['60 × 60 cm', 'Krem', 'Doff', '162000']]
    .map(([ukuran, warna, finishing, harga], i) => [`KRM-VR${i + 1}`, `Keramik lantai ${warna.toLowerCase()} ${finishing.toLowerCase()} ${ukuran.replace(' cm', '').replace(/ /g, '')}`, 21, (i % 3) + 1, `${harga}.00`, 'dus', 1, [4, 5],
      { ukuran, finishing, warna, area_aplikasi: 'Lantai dalam ruangan' }]),
  ...[['Abu muda', '245000'], ['Hitam', '265000'], ['Putih', '255000']]
    .map(([warna, harga], i) => [`GRN-6060-${i + 1}`, `Granit lantai ${warna.toLowerCase()} 60×60`, 23, (i % 2) + 2, `${harga}.00`, 'dus', 1, [4],
      { ukuran: '60 × 60 cm', finishing: 'Polished', warna, area_aplikasi: 'Lantai dalam ruangan' }]),
  ...[['Krem', '178000'], ['Abu muda', '182000'], ['Hijau sage', '189000'], ['Putih tulang', '176000']]
    .map(([warna, harga], i) => [`CAT-INT-V${i + 1}`, `Cat tembok interior ${warna.toLowerCase()} 5 kg`, 51, 5, `${harga}.00`, 'galon', 1, [4],
      { berat: '5 kg', finishing: 'Matt', warna, area_aplikasi: 'Dinding dalam ruangan' }]),
  ...[['5 watt', '22000'], ['9 watt', '32000'], ['15 watt', '52000'], ['18 watt', '61000']]
    .map(([daya, harga], i) => [`ELK-LED-V${i + 1}`, `Lampu LED ${daya}`, 61, (i % 2) + 5, `${harga}.00`, 'pcs', 1, [1, 4],
      { daya, warna: 'Putih (6500K)' }]),
  ...[['Kran dinding kamar mandi', '185000'], ['Kran dapur leher angsa', '320000']]
    .map(([nama, harga], i) => [`SAN-KRN-V${i + 1}`, nama, 32, 3, `${harga}.00`, 'pcs', 1, [i === 0 ? 2 : 1],
      { bahan: 'Stainless steel', ukuran_drat: '1/2 inch' }]),
  ['SAN-KLS-V1', 'Kloset jongkok putih', 31, 3, '285000.00', 'unit', 1, [2], { tipe: 'Jongkok', warna: 'Putih' }],
  ...[['Perekat keramik 25 kg', '98000'], ['Perekat granit 25 kg', '125000']]
    .map(([nama, harga], i) => [`PRK-25-${i + 1}`, nama, 24, 4, `${harga}.00`, 'sak', 1, [4, 2], { berat: '25 kg' }]),
  ...[['Keramik dinding motif batu 25×40', '96000', 'Batu alam'], ['Keramik dinding hijau 20×40', '88000', 'Hijau']]
    .map(([nama, harga, warna], i) => [`KRM-DND-${i + 1}`, nama, 22, (i % 2) + 1, `${harga}.00`, 'dus', 1, [2, 1],
      { ukuran: i === 0 ? '25 × 40 cm' : '20 × 40 cm', finishing: 'Glossy', warna, area_aplikasi: 'Dinding' }]),
].map(([sku, nama, kategori, merek, harga, satuan, min, ruangan, spesifikasi], i) => {
  const status = POLA_STOK[i % POLA_STOK.length]
  return [sku, nama, kategori, merek, harga, satuan, min, status, status === 'SISA_STOK' ? `${(i % 5) + 2} ${satuan}` : null, false, true, ruangan, spesifikasi]
})

const semua = [...produkUtama, ...varian]

export const products = semua.map(
  ([sku, name, category_id, brand_id, price_general, unit_sale, min_order, stock_status, stock_qty_label, is_featured, is_active, , specifications], i) => ({
    id: i + 1, sku, name, slug: slugify(name), category_id, brand_id,
    description: `${name}. Deskripsi lengkap diisi dari Master Barang.\nCocok untuk proyek rumah tinggal maupun renovasi.`,
    specifications,
    datasheet_pdf_url: i === 0 ? 'https://example.com/datasheet/keramik-lantai-abu-doff.pdf' : null,
    price_general, show_price: true, unit_sale, min_order, stock_status, stock_qty_label, is_featured, is_active,
    meta_title: i === 0 ? 'Keramik lantai abu doff 40×40 — harga & spesifikasi' : null,
    meta_description: i === 0 ? 'Keramik lantai abu doff 40×40 cm, isi 6 pcs per dus. Cek harga, stok, dan tanya langsung lewat WhatsApp.' : null,
    created_at: dibuat,
    updated_at: `2026-10-${String((i % 6) + 1).padStart(2, '0')}T${String(8 + (i % 10)).padStart(2, '0')}:00:00`,
  }),
)

// Harga tidak ditampilkan (products.show_price = false, 15 Okt 2026): pemanas air harganya tetap disimpan untuk
// catatan admin; engsel pintu (stok habis) sama sekali belum diberi harga. Pelanggan melihat "Harga via WhatsApp".
for (const [sku, harga] of [['ELK-PMA-15', '1450000.00'], ['HRD-ENG-04', null]]) {
  Object.assign(products.find((p) => p.sku === sku), { show_price: false, price_general: harga })
}

// Produk 1–4: galeri 4 foto; produk 5–9: 1 foto; sisanya tanpa foto (menguji placeholder).
const LABEL_FOTO = ['Foto utama', 'Detail tekstur', 'Terpasang di ruangan', 'Kemasan']
const LATAR_FOTO = ['#D6D3CD', '#E7DFD3', '#DCE3D2', '#E7E5E0']
export const productImages = products.flatMap((p) => {
  const jumlah = p.id <= 4 ? 4 : p.id <= 9 ? 1 : 0
  return Array.from({ length: jumlah }, (_, i) => ({
    id: p.id * 10 + i, product_id: p.id, image_url: gambarContoh(LABEL_FOTO[i], p.sku, LATAR_FOTO[i]),
    is_primary: i === 0, sort_order: i, created_at: dibuat,
  }))
})

export const productRooms = semua.flatMap((baris, i) => baris[11].map((room_id) => ({ product_id: i + 1, room_id })))

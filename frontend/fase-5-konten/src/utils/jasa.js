// Jasa Transhome (KEPUTUSAN #20): tautan halaman dan label harga. Dipakai halaman publik, Beranda, detail produk, dan admin.

export const tautanJasa = (slug) => `/jasa/${slug}`

// Harga jasa tidak ditampilkan per jasa (permintaan pemilik, 2 Nov 2026): semua memakai label tetap ini.
export const LABEL_HARGA_JASA = 'Harga pasti murah'

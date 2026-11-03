import { normalisasiNomorHp } from './format'

/** Link wa.me dengan pesan yang sudah di-encode. Nomor diambil dari pengaturan toko, jangan hardcode. */
export function buatLinkWhatsApp(nomorWa, pesan) {
  return `https://wa.me/${normalisasiNomorHp(nomorWa)}?text=${encodeURIComponent(pesan)}`
}

/**
 * Harga produk tidak ditampilkan: admin mematikan "Tampilkan harga di website" (`show_price` false)
 * atau produk memang belum diberi harga. Pelanggan melihat "Harga via WhatsApp".
 */
export const hargaDisembunyikan = (produk) => produk?.show_price === false || produk?.price_general == null || produk?.price_general === ''

/** Label tombol WhatsApp produk: stok HABIS paling penting, lalu harga yang disembunyikan. */
export function labelTanyaProduk(produk) {
  if (produk?.stock_status === 'HABIS') return 'Tanya ketersediaan'
  if (hargaDisembunyikan(produk)) return 'Tanya harga'
  return 'Tanya via WhatsApp'
}

/** "Tanya via WhatsApp" di produk; HABIS → tanya ketersediaan; harga disembunyikan → tanya harga. */
export function pesanTanyaProduk({ sku, nama, url, habis = false, tanyaHarga = false }) {
  const pembuka = habis
    ? 'Halo Transhome, apakah produk ini akan tersedia lagi?'
    : tanyaHarga
      ? 'Halo Transhome, saya ingin menanyakan harga produk ini:'
      : 'Halo Transhome, saya ingin bertanya tentang produk ini:'
  return [pembuka, nama, `SKU: ${sku}`, url].join('\n')
}

/** "Chat tim proyek" untuk kebutuhan jumlah besar. */
export function pesanJumlahBesar({ sku, nama, url }) {
  return [
    'Halo tim proyek Transhome, saya butuh produk ini dalam jumlah besar:',
    nama,
    `SKU: ${sku}`,
    url,
    'Perkiraan kebutuhan: ',
  ].join('\n')
}

/** Dari katalog/pencarian yang kosong: barang yang dicari belum ketemu di website. */
export function pesanCariProduk(topik) {
  return `Halo Transhome, saya mencari ${topik} tapi belum ketemu di website. Apakah ada?`
}

/** Dari halaman promo, event, atau artikel: sebut judul dan tautannya. */
export function pesanTanyaKonten({ judul, url }) {
  return ['Halo Transhome, saya ingin bertanya tentang ini:', judul, url].join('\n')
}

/**
 * Konsultasi jasa (KEPUTUSAN #20): nama jasa + tautan halamannya, produk asal bila dari halaman produk,
 * lalu baris isian yang tinggal dilengkapi pelanggan. Bila sudah masuk, diawali nama dan member_code.
 */
export function pesanKonsultasiJasa({ jasa, url, produk = null, nama, memberCode } = {}) {
  const sapaan = nama ? `Halo Transhome, saya ${nama}${memberCode ? ` (member ${memberCode})` : ''}.` : 'Halo Transhome,'
  return [
    `${sapaan} Saya ingin konsultasi ${jasa}.`,
    url,
    ...(produk ? [`Untuk produk: ${produk.sku} - ${produk.nama}`] : []),
    'Lokasi pekerjaan: ',
    'Kebutuhan singkat: ',
  ].join('\n')
}

/** Pesan umum ke CS; bila sudah masuk, diawali nama dan member_code (FR-S). */
export function pesanUmum({ nama, memberCode } = {}) {
  if (nama && memberCode) return `Halo Transhome, saya ${nama} (member ${memberCode}).`
  if (nama) return `Halo Transhome, saya ${nama}.`
  return 'Halo Transhome, saya ingin bertanya.'
}

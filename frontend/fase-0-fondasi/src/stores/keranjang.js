import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

const KUNCI = 'transhome_keranjang'
const JUMLAH_MAKS = 9999

const jumlahMinimal = (item) => Math.max(1, Number(item.min_order) || 1)
const batasi = (item, jumlah) => Math.min(JUMLAH_MAKS, Math.max(jumlahMinimal(item), Math.trunc(Number(jumlah)) || jumlahMinimal(item)))

// Isi localStorage bisa rusak atau diubah pengguna: baris yang tidak lengkap dibuang.
function baca() {
  try {
    const isi = JSON.parse(localStorage.getItem(KUNCI))
    if (!Array.isArray(isi)) return []
    return isi
      .filter((item) => item?.id && item.slug && item.sku && item.name)
      .map((item) => ({ ...item, jumlah: batasi(item, item.jumlah) }))
  } catch {
    return []
  }
}

/**
 * Daftar barang yang ingin ditanyakan lewat WhatsApp (KEPUTUSAN #22). Bukan pesanan: tanpa harga, stok, dan checkout.
 * Disimpan di browser pengguna, jadi tidak perlu login dan tidak ikut ke perangkat lain.
 */
export const useKeranjangStore = defineStore('keranjang', () => {
  const items = ref(baca())
  const jumlahBarang = computed(() => items.value.length)

  function simpan() {
    try {
      localStorage.setItem(KUNCI, JSON.stringify(items.value))
    } catch {
      // Penyimpanan penuh atau diblokir: keranjang tetap berjalan sampai halaman ditutup.
    }
  }

  /** Menambah 1 bila sudah ada, atau sebanyak minimal pembelian bila baru. Hasil: jumlah barang itu sekarang. */
  function tambah(produk) {
    const ada = items.value.find((item) => item.id === produk.id)
    if (ada) {
      ada.jumlah = batasi(ada, ada.jumlah + 1)
      simpan()
      return ada.jumlah
    }
    const baru = {
      id: produk.id,
      slug: produk.slug,
      sku: produk.sku,
      name: produk.name,
      image_url: produk.primary_image_url || produk.images?.[0]?.image_url || null,
      unit_sale: produk.unit_sale,
      min_order: produk.min_order,
    }
    baru.jumlah = jumlahMinimal(baru)
    items.value.push(baru)
    simpan()
    return baru.jumlah
  }

  function ubahJumlah(id, jumlah) {
    const item = items.value.find((x) => x.id === id)
    if (!item) return
    item.jumlah = batasi(item, jumlah)
    simpan()
  }

  function hapus(id) {
    items.value = items.value.filter((item) => item.id !== id)
    simpan()
  }

  // Produk sering dibuka di beberapa tab: tanpa ini, tab yang lama menimpa barang dari tab lain.
  window.addEventListener('storage', (event) => {
    if (event.key === KUNCI) items.value = baca()
  })

  return { items, jumlahBarang, tambah, ubahJumlah, hapus, jumlahMinimal }
})

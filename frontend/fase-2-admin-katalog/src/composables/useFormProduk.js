import { computed, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { checkSku, createProduct, getAdminProduct, updateProduct } from '@/services/adminProductService'
import { errorField, pesanError } from '@/services/errors'
import { labelSpesifikasi } from '@/utils/produk'
import { buatSlug, POLA_SLUG } from '@/utils/slug'

const KUNCI_DRAF = 'transhome.drafProduk'
const angka = (teks) => String(teks ?? '').replace(/\D/g, '')
let nomorBaris = 0
const barisSpek = (label = '', nilai = '') => ({ id: `spek-${++nomorBaris}`, label, nilai })
// "Isi per dus" → "isi_per_dus" (kebalikan labelSpesifikasi, jadi kunci lama tetap sama).
const kunciSpek = (label) => label.trim().toLowerCase().replace(/\s+/g, '_')

function formKosong() {
  return {
    sku: '', name: '', slug: '', category_id: '', brand_id: '', description: '',
    price_general: '', show_price: true, unit_sale: '', min_order: 1, stock_status: 'TERSEDIA', stock_qty_label: '',
    is_active: true, is_featured: false, images: [], room_ids: [], spesifikasi: [],
    datasheet_pdf_url: '', meta_title: '', meta_description: '',
  }
}

/**
 * State dan aksi form produk admin (baru & ubah): muat, slug otomatis, cek SKU, validasi,
 * simpan, dan draf saat sesi habis. Tampilan ada di views/admin/katalog/ProdukForm.vue.
 */
export function useFormProduk() {
  const route = useRoute()
  const id = route.params.id ? Number(route.params.id) : null

  const form = reactive(formKosong())
  const errors = ref({})
  // Harga disembunyikan dari website: pesan "harga wajib" tidak berlaku lagi.
  watch(
    () => form.show_price,
    (tampil) => {
      if (tampil || !errors.value.price_general) return
      const { price_general: _, ...sisa } = errors.value
      errors.value = sisa
    },
  )
  const status = ref(id ? 'loading' : 'ready') // loading | ready | tidak-ada | error
  const pesan = ref('')
  const galatUmum = ref('')
  const statusSku = ref('') // '' | memeriksa | tersedia | dipakai
  const menyimpan = ref('') // '' | simpan | tambah-lagi
  const dipulihkan = ref(false)
  const slugManual = ref(Boolean(id))
  const awal = ref('')
  let skuAwal = ''

  const potret = () => JSON.stringify(form)
  const kotor = computed(() => status.value === 'ready' && potret() !== awal.value)

  function isi(p) {
    Object.assign(form, formKosong(), {
      sku: p.sku, name: p.name, slug: p.slug, category_id: p.category_id ?? '', brand_id: p.brand_id ?? '',
      description: p.description ?? '', unit_sale: p.unit_sale, show_price: p.show_price !== false,
      price_general: p.price_general == null ? '' : angka(Number(p.price_general).toFixed(0)),
      min_order: Number(p.min_order ?? 1), stock_status: p.stock_status, stock_qty_label: p.stock_qty_label ?? '',
      is_active: p.is_active, is_featured: p.is_featured, datasheet_pdf_url: p.datasheet_pdf_url ?? '',
      meta_title: p.meta_title ?? '', meta_description: p.meta_description ?? '',
      images: (p.images ?? []).map(({ id: idFoto, image_url, is_primary }) => ({ id: idFoto, image_url, is_primary })),
      room_ids: [...(p.room_ids ?? [])],
      spesifikasi: Object.entries(p.specifications ?? {}).map(([k, v]) => barisSpek(labelSpesifikasi(k), String(v ?? ''))),
    })
  }

  // Draf hanya untuk halaman yang sama (mis. /admin/produk/baru) dan dipakai sekali.
  function simpanDraf() {
    try {
      sessionStorage.setItem(KUNCI_DRAF, JSON.stringify({ path: route.path, form, slugManual: slugManual.value }))
    } catch {
      // Kuota penuh (foto besar): draf tidak tersimpan.
    }
  }

  function pulihkanDraf() {
    try {
      const draf = JSON.parse(sessionStorage.getItem(KUNCI_DRAF) ?? 'null')
      if (!draf || draf.path !== route.path) return
      sessionStorage.removeItem(KUNCI_DRAF)
      Object.assign(form, draf.form)
      slugManual.value = draf.slugManual
      dipulihkan.value = true
    } catch {
      // Draf rusak: abaikan.
    }
  }

  async function muat() {
    if (!id) return
    status.value = 'loading'
    try {
      const p = await getAdminProduct(id)
      isi(p)
      skuAwal = p.sku
      awal.value = potret()
      pulihkanDraf()
      status.value = 'ready'
    } catch (error) {
      if (error?.response?.status === 404) status.value = 'tidak-ada'
      else {
        pesan.value = pesanError(error)
        status.value = 'error'
      }
    }
  }

  function ubahNama(nilai) {
    form.name = nilai
    if (!slugManual.value) form.slug = buatSlug(nilai)
  }

  function ubahSlug(nilai) {
    slugManual.value = true
    form.slug = nilai
  }

  async function periksaSku() {
    const sku = form.sku.trim()
    if (!sku || sku.toLowerCase() === skuAwal.toLowerCase()) {
      statusSku.value = ''
      return
    }
    statusSku.value = 'memeriksa'
    try {
      const tersedia = await checkSku(sku, id)
      if (form.sku.trim() !== sku) return
      statusSku.value = tersedia ? 'tersedia' : 'dipakai'
      errors.value = { ...errors.value, sku: tersedia ? '' : 'SKU sudah dipakai produk lain' }
    } catch {
      statusSku.value = ''
    }
  }

  function validasi() {
    const e = {}
    if (!form.sku.trim()) e.sku = 'SKU wajib diisi'
    else if (statusSku.value === 'dipakai') e.sku = 'SKU sudah dipakai produk lain'
    if (!form.name.trim()) e.name = 'Nama produk wajib diisi'
    if (!form.slug) e.slug = 'Slug wajib diisi'
    else if (!POLA_SLUG.test(form.slug)) e.slug = 'Slug hanya boleh huruf kecil, angka, dan tanda hubung'
    if (!form.category_id) e.category_id = 'Pilih kategori'
    if (form.show_price && !angka(form.price_general)) e.price_general = 'Harga wajib diisi, atau matikan "Tampilkan harga di website"'
    if (!form.unit_sale.trim()) e.unit_sale = 'Satuan jual wajib diisi, mis. dus atau sak'
    if (!(Number(form.min_order) >= 1 && Number.isInteger(Number(form.min_order)))) e.min_order = 'Minimal pembelian paling sedikit 1'
    if (form.is_active && !form.images.some((f) => f.is_primary)) e.images = 'Tambahkan minimal 1 foto utama, atau matikan "Tampil di website".'
    const dipakai = new Set()
    const spek = form.spesifikasi.map((b) => {
      const label = b.label.trim()
      const nilai = b.nilai.trim()
      if (!label && !nilai) return null
      if (!label) return { label: 'Nama atribut wajib diisi' }
      if (!nilai) return { nilai: 'Isi atribut wajib diisi' }
      if (dipakai.has(kunciSpek(label))) return { label: 'Nama atribut sudah dipakai di baris lain' }
      dipakai.add(kunciSpek(label))
      return null
    })
    if (spek.some(Boolean)) e.spesifikasi = spek
    errors.value = e
    return Object.keys(e).length === 0
  }

  function payload() {
    const spesifikasi = {}
    for (const b of form.spesifikasi) if (b.label.trim() && b.nilai.trim()) spesifikasi[kunciSpek(b.label)] = b.nilai.trim()
    return {
      sku: form.sku.trim(), name: form.name.trim(), slug: form.slug,
      category_id: Number(form.category_id), brand_id: form.brand_id ? Number(form.brand_id) : null,
      description: form.description.trim() || null, specifications: spesifikasi,
      datasheet_pdf_url: form.datasheet_pdf_url || null, show_price: form.show_price,
      price_general: angka(form.price_general) ? Number(angka(form.price_general)) : null,
      unit_sale: form.unit_sale.trim(), min_order: Number(form.min_order), stock_status: form.stock_status,
      stock_qty_label: form.stock_qty_label.trim() || null, is_active: form.is_active, is_featured: form.is_featured,
      meta_title: form.meta_title.trim() || null, meta_description: form.meta_description.trim() || null,
      images: form.images.map((f, i) => ({ image_url: f.image_url, is_primary: f.is_primary, sort_order: i })),
      room_ids: [...form.room_ids],
    }
  }

  /** Kirim ke API. Hasil: true bila tersimpan. Galat 422/409 dipetakan ke field. */
  async function kirim(mode = 'simpan') {
    if (menyimpan.value) return false
    galatUmum.value = ''
    if (!validasi()) return false
    menyimpan.value = mode
    try {
      if (id) await updateProduct(id, payload())
      else await createProduct(payload())
      return true
    } catch (error) {
      const statusHttp = error?.response?.status
      if (statusHttp === 422 || statusHttp === 409) {
        const galat = errorField(error)
        errors.value = galat
        if (galat.sku) statusSku.value = 'dipakai'
        galatUmum.value = error.response.data?.message ?? 'Periksa kembali isian yang ditandai.'
      } else if (statusHttp !== 401) {
        galatUmum.value = pesanError(error)
      }
      return false
    } finally {
      menyimpan.value = ''
    }
  }

  // Form baru langsung siap: titik awal "belum ada perubahan" dicatat sekarang, lalu draf (bila ada) dipulihkan.
  if (!id) {
    awal.value = potret()
    pulihkanDraf()
  }

  /** Setelah "Simpan & tambah lagi": form kosong untuk produk berikutnya. */
  function kosongkan() {
    Object.assign(form, formKosong())
    errors.value = {}
    statusSku.value = ''
    slugManual.value = false
    dipulihkan.value = false
    awal.value = potret()
  }

  return {
    id, form, errors, status, pesan, galatUmum, statusSku, menyimpan, dipulihkan, kotor,
    muat, ubahNama, ubahSlug, periksaSku, kirim, kosongkan, simpanDraf,
    tandaiTersimpan: () => (awal.value = potret()),
    tambahBarisSpek: barisSpek,
  }
}

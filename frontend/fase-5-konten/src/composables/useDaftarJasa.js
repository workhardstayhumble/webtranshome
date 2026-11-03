import { computed, nextTick, onMounted, ref } from 'vue'
import { listAdminServices, patchService } from '@/services/adminKontenService'
import { pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'

/**
 * Halaman /admin/jasa: muat semua jasa, aktif/nonaktif (konfirmasi saat menonaktifkan), dan geser urutan.
 * Urutan menentukan jasa yang tampil di Beranda (3 teratas) dan di /jasa.
 */
export function useDaftarJasa() {
  const ui = useUiStore()
  const status = ref('loading') // loading | ready | error
  const daftar = ref([])
  const pesan = ref('')
  const konfirmasi = ref(null) // jasa yang akan dinonaktifkan
  const sedangNonaktif = ref(false)

  async function muat() {
    status.value = 'loading'
    try {
      daftar.value = await listAdminServices()
      status.value = 'ready'
    } catch (error) {
      pesan.value = pesanError(error)
      status.value = 'error'
    }
  }
  onMounted(muat)

  const urut = computed(() => [...daftar.value].sort((a, b) => a.sort_order - b.sort_order || a.id - b.id))
  const ganti = (baris, hasil) => Object.assign(baris, hasil)

  async function aktifkan(baris) {
    baris.is_active = true
    baris.sibuk = true
    try {
      ganti(baris, await patchService(baris.id, { is_active: true }))
      ui.tampilkanToast({ pesan: `“${baris.name}” tampil lagi di website.`, jenis: 'success' })
    } catch (error) {
      baris.is_active = false
      ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    } finally {
      baris.sibuk = false
    }
  }

  function ubahAktif(baris, nilai) {
    if (nilai) aktifkan(baris)
    else konfirmasi.value = baris
  }

  async function nonaktifkan() {
    const baris = konfirmasi.value
    sedangNonaktif.value = true
    try {
      ganti(baris, await patchService(baris.id, { is_active: false }))
      konfirmasi.value = null
      ui.tampilkanToast({ pesan: `“${baris.name}” tidak tampil lagi di website.`, jenis: 'success' })
    } catch (error) {
      ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    } finally {
      sedangNonaktif.value = false
    }
  }

  /** Tukar dengan tetangga lalu nomori ulang 1..n; hanya yang berubah yang dikirim. Fokus ikut baris yang dipindah. */
  async function geser(baris, arah) {
    const urutan = [...urut.value]
    const i = urutan.indexOf(baris)
    if (!urutan[i + arah]) return
    ;[urutan[i], urutan[i + arah]] = [urutan[i + arah], urutan[i]]
    const berubah = urutan.map((b, n) => ({ b, lama: b.sort_order, baru: n + 1 })).filter((x) => x.lama !== x.baru)
    for (const x of berubah) x.b.sort_order = x.baru
    await nextTick()
    const tombol = (a) => document.querySelector(`[data-urut="${baris.id}:${a}"]:not(:disabled)`)
    ;(tombol(arah) ?? tombol(-arah))?.focus()
    try {
      await Promise.all(berubah.map((x) => patchService(x.b.id, { sort_order: x.baru })))
    } catch (error) {
      for (const x of berubah) x.b.sort_order = x.lama
      ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    }
  }

  return { status, daftar, pesan, muat, urut, konfirmasi, sedangNonaktif, ubahAktif, nonaktifkan, geser }
}

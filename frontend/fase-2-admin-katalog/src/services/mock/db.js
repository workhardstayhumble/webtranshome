// Mode mock: perubahan dari panel admin disimpan di localStorage, jadi sama di semua tab dan tetap ada
// setelah refresh, seperti database sungguhan. Kalau data contoh di kode berubah, simpanan lama dibuang
// otomatis (sidik berbeda). Tombol "Reset data contoh" di /dev/components menghapusnya.
import { brands, categories, productImages, productRooms, products, rooms } from '@/services/mock/data/catalog'
import { activityLogs, qontakContacts, syncLogs } from './data/admin'
import { emailKeluar, pointTransactions, publicSettings, rolePermissions, tokenReset, tokenVerifikasi, users } from '@/services/mock/data/akun'
import { articles, banners, services } from '@/services/mock/data/konten'

const KUNCI = 'transhome.mockDb'
const TABEL = {
  products, productImages, productRooms, categories, brands, rooms, activityLogs, qontakContacts, syncLogs,
  users, emailKeluar, tokenVerifikasi, tokenReset, pointTransactions, banners, articles, services,
  // Objek (bukan daftar): site_settings yang bisa diubah di /admin/pengaturan, izin role dari /admin/role.
  publicSettings, rolePermissions,
}
const sidik = String(JSON.stringify(TABEL).length)

export function pulihkanDb() {
  try {
    const simpanan = JSON.parse(localStorage.getItem(KUNCI) ?? 'null')
    if (!simpanan || simpanan.sidik !== sidik) return
    for (const [nama, baris] of Object.entries(simpanan.tabel ?? {})) {
      if (Array.isArray(TABEL[nama]) && Array.isArray(baris)) TABEL[nama].splice(0, TABEL[nama].length, ...baris)
      else if (TABEL[nama] && baris && typeof baris === 'object') Object.assign(TABEL[nama], baris)
    }
  } catch {
    // Simpanan rusak: tetap memakai data contoh bawaan.
  }
}

/** Dipanggil setelah setiap perubahan. Kuota penuh / storage diblokir: perubahan hanya bertahan di tab ini. */
export function simpanDb() {
  try {
    localStorage.setItem(KUNCI, JSON.stringify({ sidik, tabel: TABEL }))
  } catch {
    // diabaikan
  }
}

export function resetDb() {
  try {
    localStorage.removeItem(KUNCI)
  } catch {
    // diabaikan
  }
}

export const idBaru = (tabel) => tabel.reduce((maks, baris) => Math.max(maks, baris.id ?? 0), 0) + 1

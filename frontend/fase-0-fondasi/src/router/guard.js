import { useAuthStore } from '@/stores/auth'
import { boleh } from '@/composables/usePermission'

const DESKRIPSI_DEFAULT = document.querySelector('meta[name="description"]')?.getAttribute('content') ?? ''

/** Hanya terima redirect ke path internal, supaya ?redirect= tidak bisa membawa keluar situs. */
function redirectAman(nilai) {
  return typeof nilai === 'string' && nilai.startsWith('/') && !nilai.startsWith('//') ? nilai : null
}

/**
 * Guard per area (CLAUDE.md §3–4). Baca `meta.area` dan `meta.permission` dari route.
 * Ini hanya menjaga tampilan; backend tetap wajib menolak dengan 401/403.
 */
export function pasangGuard(router) {
  router.beforeEach(async (to) => {
    const area = to.meta.area ?? 'public'
    const auth = useAuthStore()
    const perluPulihkan = auth.isLoggedIn && !auth.user

    // Halaman publik tidak menunggu profil dimuat; header menyusul setelah selesai.
    if (area === 'public') {
      if (perluPulihkan) auth.restoreSession()
      return true
    }

    // Token masih ada setelah refresh: muat ulang profil dan izin sebelum memutuskan akses.
    if (perluPulihkan) await auth.restoreSession()

    if (area === 'guest') {
      if (!auth.isLoggedIn) return true
      if (auth.isAdmin) return { path: '/admin' }
      return { path: redirectAman(to.query.redirect) ?? '/akun' }
    }

    if (!auth.isLoggedIn) {
      // Login admin terpisah dari login pelanggan (KEPUTUSAN #8).
      const halamanMasuk = area === 'admin' ? '/admin/masuk' : '/masuk'
      return { path: halamanMasuk, query: { redirect: to.fullPath } }
    }

    if (area === 'member' && !auth.hasRole('CUSTOMER')) {
      return { name: 'tidak-punya-akses' }
    }

    if (area === 'admin') {
      const izin = to.meta.permission
      if (!auth.isAdmin || (izin && !boleh(izin))) {
        return { name: 'tidak-punya-akses' }
      }
    }

    return true
  })

  router.afterEach((to, from) => {
    // Perubahan query saja (filter katalog) tidak mengganti judul yang sudah diatur useSeo di halaman.
    if (from.matched.length && to.path === from.path) return
    document.title = to.meta.title ? `${to.meta.title} | Transhome` : 'Transhome'
    document.querySelector('meta[name="description"]')?.setAttribute('content', to.meta.description ?? DESKRIPSI_DEFAULT)
  })
}
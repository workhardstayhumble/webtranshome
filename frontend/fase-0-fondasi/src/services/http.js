import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

const PAKAI_MOCK = import.meta.env.VITE_USE_MOCK === 'true'
const METODE_AKSI = ['post', 'put', 'patch', 'delete']

/**
 * Satu-satunya instance axios (CLAUDE.md §1). Service per modul memakai instance ini.
 * Mode mock (VITE_USE_MOCK=true): permintaan dijawab services/mock/ lewat adapter,
 * jadi service dan interceptor tidak perlu tahu sedang memakai mock atau API asli.
 */
const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 15000,
  headers: { Accept: 'application/json' },
  ...(PAKAI_MOCK && {
    // TODO: ganti ke API — matikan VITE_USE_MOCK begitu endpoint backend tersedia.
    adapter: async (config) => (await import('./mock/adapter')).mockAdapter(config),
  }),
})

http.interceptors.request.use((config) => {
  const { token } = useAuthStore()
  // Hormati token eksplisit saat login mengambil profil sebelum store diperbarui.
  if (token && !config.headers.Authorization) config.headers.Authorization = `Bearer ${token}`
  return config
})

/**
 * Penanganan error terpusat (docs/FRONTEND_BUILD_PLAN.md, Aturan Workflow Tombol):
 * - 401 di mana pun: sesi dihapus, area member/admin diarahkan ke halaman masuk.
 * - Aksi (POST/PUT/PATCH/DELETE): 403, 429, 5xx, dan gangguan jaringan memunculkan toast.
 * - Permintaan data (GET) tidak memunculkan toast; halaman menampilkan ErrorState.
 * - 400/422/409 ditangani form lewat errorField() di services/errors.js.
 * Kirim `{ tanpaToast: true }` di config untuk mematikan toast pada satu permintaan.
 */
http.interceptors.response.use(
  (response) => response,
  async (error) => {
    const status = error.response?.status
    const aksi = METODE_AKSI.includes((error.config?.method ?? 'get').toLowerCase())

    if (status === 401) {
      useAuthStore().logout()
      // Import di sini untuk menghindari import melingkar router -> store -> http.
      const { default: router } = await import('@/router')
      const sekarang = router.currentRoute.value
      const area = sekarang.meta.area
      if (area === 'member' || area === 'admin') {
        const halamanMasuk = area === 'admin' ? '/admin/masuk' : '/masuk'
        await router.push({ path: halamanMasuk, query: { redirect: sekarang.fullPath } })
      }
    } else if (aksi && !error.config?.tanpaToast) {
      const ui = useUiStore()
      if (status === 403) {
        ui.tampilkanToast({ pesan: 'Kamu tidak punya akses untuk aksi ini.', jenis: 'danger' })
      } else if (status === 429) {
        ui.tampilkanToast({ pesan: 'Terlalu banyak percobaan. Coba lagi beberapa saat.', jenis: 'warning' })
      } else if (!error.response || status >= 500) {
        ui.tampilkanToast({ pesan: 'Terjadi gangguan. Coba lagi.', jenis: 'danger' })
      }
    }

    return Promise.reject(error)
  },
)

export default http
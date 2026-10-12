import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { getMe } from '@/services/authService'

const KUNCI_TOKEN = 'transhome.token'
const ROLE_ADMIN = ['SUPER_ADMIN', 'ADMIN_KATALOG', 'ADMIN_MEMBERSHIP', 'ADMIN_KONTEN']

function bacaToken() {
  try {
    return localStorage.getItem(KUNCI_TOKEN)
  } catch {
    return null
  }
}

function simpanToken(token) {
  try {
    if (token) localStorage.setItem(KUNCI_TOKEN, token)
    else localStorage.removeItem(KUNCI_TOKEN)
  } catch {
    // Storage diblokir browser: sesi hanya bertahan sampai halaman ditutup.
  }
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref(bacaToken())
  // Bentuk dari API: { id, email, phone_number, full_name, member_code, roles: ['CUSTOMER'] }
  const user = ref(null)
  // Permission backend, misalnya catalog.product.read atau catalog.category.manage.
  const permissions = ref([])

  const isLoggedIn = computed(() => Boolean(token.value))
  const roles = computed(() => user.value?.roles ?? [])
  const isAdmin = computed(() => roles.value.some((role) => ROLE_ADMIN.includes(role)))

  function hasRole(kode) {
    return roles.value.includes(kode)
  }

  function can(kode) {
    return permissions.value.includes(kode)
  }

  function setSession(sesi) {
    if (!sesi?.token || !sesi?.user?.id) throw new Error('Data sesi login tidak lengkap.')
    token.value = sesi.token
    user.value = sesi.user
    permissions.value = sesi.permissions ?? []
    simpanToken(sesi.token)
  }

  function setUser(data) {
    if (!data?.user?.id) throw new Error('Data profil user tidak lengkap.')
    user.value = data.user
    permissions.value = data.permissions ?? []
  }

  function logout() {
    token.value = null
    user.value = null
    permissions.value = []
    simpanToken(null)
  }

  /** Setelah profil diubah (PATCH /me): perbarui data user tanpa memuat ulang izin. */
  function perbaruiUser(data) {
    user.value = { ...user.value, ...data }
  }

  /** Setelah refresh token masih ada tapi profil belum dimuat: ambil ulang dari API. */
  let pemulihan = null
  function restoreSession() {
    if (!token.value) return Promise.resolve(false)
    if (user.value) return Promise.resolve(true)
    pemulihan ??= getMe()
      .then((data) => {
        setUser(data)
        return true
      })
      .catch(() => {
        logout()
        return false
      })
      .finally(() => {
        pemulihan = null
      })
    return pemulihan
  }

  /** Callback "Masuk dengan Google": simpan token dari backend, lalu muat profil. Hasil true bila berhasil. */
  function masukDenganToken(tokenBaru) {
    token.value = tokenBaru
    user.value = null
    permissions.value = []
    simpanToken(tokenBaru)
    return restoreSession()
  }

  return { token, user, permissions, isLoggedIn, roles, isAdmin, hasRole, can, setSession, setUser, perbaruiUser, logout, restoreSession, masukDenganToken }
})
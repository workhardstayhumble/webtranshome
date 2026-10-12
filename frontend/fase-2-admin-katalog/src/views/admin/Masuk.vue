<script setup>
import { nextTick, reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { ShieldCheck } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import PasswordInput from '@/components/ui/PasswordInput.vue'
import { useAuth } from '@/composables/useAuth'
import { pesanError } from '@/services/errors'
import { useAuthStore } from '@/stores/auth'

// Login admin terpisah dari login pelanggan (KEPUTUSAN #8). Pesan gagal tampil di form, bukan toast.
const POLA_EMAIL = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const modeMock = import.meta.env.DEV && import.meta.env.VITE_USE_MOCK === 'true'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const { login, logout } = useAuth()

const form = reactive({ email: '', password: '' })
const errors = reactive({ email: '', password: '' })
const galat = ref('')
const memuat = ref(false)
const formEl = ref(null)

// Hanya kembali ke halaman admin; redirect lain diabaikan.
function tujuan() {
  const r = route.query.redirect
  return typeof r === 'string' && r.startsWith('/admin') && !r.startsWith('/admin/masuk') && !r.startsWith('//') ? r : '/admin'
}

function validasi() {
  const email = form.email.trim()
  errors.email = !email ? 'Email wajib diisi' : POLA_EMAIL.test(email) ? '' : 'Format email belum benar, contoh nama@toko.com'
  errors.password = form.password ? '' : 'Password wajib diisi'
  return !errors.email && !errors.password
}

async function masuk() {
  if (memuat.value) return
  galat.value = ''
  if (!validasi()) {
    await nextTick()
    formEl.value?.querySelector('[aria-invalid="true"]')?.focus()
    return
  }
  memuat.value = true
  try {
    await login({ email: form.email.trim(), password: form.password }, { admin: true })
    // Pertahanan kedua: backend seharusnya sudah menolak akun non-admin dengan 403.
    if (!auth.isAdmin) {
      await logout()
      galat.value = 'Akun ini tidak punya akses admin.'
      return
    }
    await router.replace(tujuan())
  } catch (error) {
    const status = error?.response?.status
    if (status === 401) galat.value = 'Email atau password salah.'
    else if (status === 403) galat.value = error.response.data?.message || 'Akun ini tidak punya akses admin.'
    else galat.value = pesanError(error)
  } finally {
    memuat.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-6 rounded-xl border border-border bg-surface p-6 shadow-[0_20px_50px_rgba(28,36,48,0.07)] sm:p-9">
    <div class="flex flex-col gap-2">
      <span class="inline-flex items-center gap-1.5 self-start rounded-full bg-primary-soft px-3 py-1.5 text-xs font-bold text-primary">
        <ShieldCheck class="size-3.5" aria-hidden="true" />
        Panel admin
      </span>
      <h1 class="text-3xl font-bold tracking-tight">Masuk admin</h1>
      <p class="text-sm text-muted">Khusus pengelola Transhome. Akun admin dibuat oleh Super Admin.</p>
    </div>

    <form ref="formEl" class="flex flex-col gap-4" novalidate @submit.prevent="masuk">
      <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ galat }}</p>
      <p v-else-if="tujuan() !== '/admin'" class="rounded-md bg-info-soft px-3 py-2.5 text-sm font-medium text-info-ink">
        Masuk lagi untuk melanjutkan. Setelah masuk, kamu kembali ke halaman tadi.
      </p>
      <Input v-model="form.email" label="Email" type="email" autocomplete="username" inputmode="email" required :error="errors.email" />
      <PasswordInput v-model="form.password" autocomplete="current-password" required :error="errors.password" />
      <RouterLink
        to="/lupa-password"
        class="-mt-1 inline-flex min-h-11 items-center self-end rounded-sm text-sm font-semibold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-8"
      >
        Lupa password?
      </RouterLink>
      <Button type="submit" block :loading="memuat">Masuk</Button>
    </form>

    <p v-if="modeMock" class="rounded-md bg-subtle px-3 py-2.5 text-xs text-muted">
      Mode mock: <span class="font-mono">katalog@example.com</span> atau <span class="font-mono">superadmin@example.com</span>, password
      <span class="font-mono">rahasia123</span>.
    </p>
  </div>
</template>
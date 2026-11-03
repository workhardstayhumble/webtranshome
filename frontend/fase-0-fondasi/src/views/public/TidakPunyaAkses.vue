<script setup>
import { useRoute, useRouter } from 'vue-router'
import { ShieldX } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { useAuth } from '@/composables/useAuth'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const { logout } = useAuth()

// Ganti akun: keluar dulu, lalu ke login yang sesuai halaman yang tadi dibuka (?dari= dari guard), dan kembali ke sana.
async function masukAkunLain() {
  const dari = typeof route.query.dari === 'string' && route.query.dari.startsWith('/') && !route.query.dari.startsWith('//') ? route.query.dari : null
  const keAdmin = dari ? dari.startsWith('/admin') : auth.isAdmin
  await logout()
  router.push({ path: keAdmin ? '/admin/masuk' : '/masuk', query: dari ? { redirect: dari } : {} })
}
</script>

<template>
  <section class="mx-auto flex max-w-md flex-col items-center gap-4 px-4 py-20 text-center">
    <ShieldX class="size-10 text-danger" aria-hidden="true" />
    <h1 class="text-2xl font-bold">Tidak punya akses</h1>
    <p class="text-muted">Akun kamu tidak punya izin untuk membuka halaman ini. Hubungi Super Admin kalau kamu membutuhkan akses.</p>
    <div class="flex flex-wrap justify-center gap-2">
      <Button v-if="auth.isLoggedIn" variant="outline" data-ganti-akun @click="masukAkunLain">Masuk dengan akun lain</Button>
      <Button :to="auth.isAdmin ? '/admin' : '/'">{{ auth.isAdmin ? 'Kembali ke dashboard' : 'Kembali ke beranda' }}</Button>
    </div>
  </section>
</template>

<script setup>
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { History, LayoutDashboard, LogOut, ShieldCheck, UserRound } from 'lucide-vue-next'
import UserMenu from '@/components/layout/UserMenu.vue'
import LogoTranshome from '@/components/layout/LogoTranshome.vue'
import { useAuth } from '@/composables/useAuth'

// Menu akun dari sketsa "Dashboard member". Di HP berubah jadi tab horizontal.
const MENU = [
  { label: 'Ringkasan', to: '/akun', ikon: LayoutDashboard, exact: true },
  { label: 'Riwayat poin', to: '/akun/poin', ikon: History },
  { label: 'Profil', to: '/akun/profil', ikon: UserRound },
  { label: 'Keamanan', to: '/akun/keamanan', ikon: ShieldCheck },
]

const router = useRouter()
const { logout } = useAuth()

async function keluar() {
  await logout()
  router.push('/')
}
</script>

<template>
  <div class="flex min-h-dvh flex-col">
    <a
      href="#konten"
      class="sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-60 focus:rounded-md focus:bg-surface focus:px-4 focus:py-2 focus:font-semibold focus:shadow-lg"
    >
      Langsung ke konten
    </a>
    <header class="border-b border-border bg-surface">
      <div class="mx-auto flex h-16 max-w-7xl items-center gap-4 px-4 sm:px-6">
        <RouterLink to="/" class="rounded-md text-2xl font-extrabold tracking-tight focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" aria-label="Transhome, ke beranda">
          <LogoTranshome />
        </RouterLink>
        <RouterLink to="/" class="hidden rounded-sm text-sm font-semibold transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:inline">
          Kembali belanja
        </RouterLink>
        <div class="ml-auto"><UserMenu /></div>
      </div>
    </header>

    <div class="mx-auto flex w-full max-w-7xl flex-1 flex-col gap-6 px-4 py-6 sm:px-6 lg:flex-row lg:items-start lg:py-8">
      <nav aria-label="Menu akun" class="-mx-4 overflow-x-auto px-4 scrollbar-none lg:mx-0 lg:w-64 lg:shrink-0 lg:overflow-visible lg:px-0">
        <ul class="flex gap-1 lg:flex-col lg:rounded-xl lg:border lg:border-border lg:bg-surface lg:p-2">
          <li v-for="item in MENU" :key="item.to">
            <RouterLink
              :to="item.to"
              class="flex h-11 items-center gap-3 rounded-lg px-3.5 text-sm font-semibold whitespace-nowrap transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              :exact-active-class="item.exact ? 'bg-primary-soft text-primary-hover' : ''"
              :active-class="item.exact ? '' : 'bg-primary-soft text-primary-hover'"
            >
              <component :is="item.ikon" class="size-4.5" aria-hidden="true" />
              {{ item.label }}
            </RouterLink>
          </li>
          <li class="hidden lg:block"><hr class="my-1 border-border" /></li>
          <li class="hidden lg:block">
            <button
              type="button"
              class="flex h-11 w-full items-center gap-3 rounded-lg px-3.5 text-sm font-semibold text-danger transition-colors duration-150 hover:bg-danger-soft focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              @click="keluar"
            >
              <LogOut class="size-4.5" aria-hidden="true" />
              Keluar
            </button>
          </li>
        </ul>
      </nav>
      <main id="konten" tabindex="-1" class="min-w-0 flex-1 focus:outline-none">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import { Menu, X } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import AdminSidebar from '@/components/admin/AdminSidebar.vue'
import { useAuth } from '@/composables/useAuth'
import { useAuthStore } from '@/stores/auth'
import { useMenuStore } from '@/stores/menu'

// Sidebar dibangun dari API menus sesuai izin user, tanpa hardcode per role (CLAUDE.md §4).
const LABEL_ROLE = {
  SUPER_ADMIN: 'Super Admin',
  ADMIN_KATALOG: 'Admin Katalog',
  ADMIN_MEMBERSHIP: 'Admin Membership',
  ADMIN_KONTEN: 'Admin Konten',
}
const SELEKTOR_FOKUS = 'a[href], button:not([disabled])'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const menu = useMenuStore()
const { logout } = useAuth()

const labelRole = computed(() => auth.roles.map((r) => LABEL_ROLE[r]).filter(Boolean).join(', '))
const idDrawer = useId()
const drawerBuka = ref(false)
const drawer = ref(null)
const tombolMenu = ref(null)

// Menu dimuat setelah profil admin siap (guard memulihkan sesi lebih dulu).
watch(
  () => auth.user?.id,
  (id) => {
    if (id) menu.muat(id)
  },
  { immediate: true },
)

async function bukaDrawer() {
  drawerBuka.value = true
  document.body.style.overflow = 'hidden'
  await nextTick()
  drawer.value?.querySelector(SELEKTOR_FOKUS)?.focus()
}

function tutupDrawer({ kembalikanFokus = true } = {}) {
  drawerBuka.value = false
  document.body.style.overflow = ''
  if (kembalikanFokus) tombolMenu.value?.focus()
}

// Esc menutup; Tab berputar di dalam drawer selama terbuka.
function onKeydownDrawer(event) {
  if (event.key === 'Escape') {
    event.preventDefault()
    tutupDrawer()
    return
  }
  if (event.key !== 'Tab') return
  const daftar = [...(drawer.value?.querySelectorAll(SELEKTOR_FOKUS) ?? [])]
  if (!daftar.length) return
  const pertama = daftar[0]
  const terakhir = daftar[daftar.length - 1]
  if (event.shiftKey && document.activeElement === pertama) {
    event.preventDefault()
    terakhir.focus()
  } else if (!event.shiftKey && document.activeElement === terakhir) {
    event.preventDefault()
    pertama.focus()
  }
}

watch(
  () => route.fullPath,
  () => {
    if (drawerBuka.value) tutupDrawer({ kembalikanFokus: false })
  },
)
onBeforeUnmount(() => {
  document.body.style.overflow = ''
})

async function keluar() {
  await logout()
  router.push('/admin/masuk')
}
</script>

<template>
  <div class="min-h-dvh bg-background-admin lg:flex">
    <a
      href="#konten"
      class="sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-60 focus:rounded-md focus:bg-surface focus:px-4 focus:py-2 focus:font-semibold focus:shadow-lg"
    >
      Langsung ke konten
    </a>

    <aside class="hidden w-65 shrink-0 bg-secondary lg:block">
      <div class="sticky top-0 max-h-dvh overflow-y-auto">
        <AdminSidebar />
      </div>
    </aside>

    <!-- HP & tablet: sidebar sebagai drawer dari kiri. -->
    <div v-if="drawerBuka" class="fixed inset-0 z-50 flex lg:hidden" @keydown="onKeydownDrawer">
      <div class="absolute inset-0 bg-foreground/50" aria-hidden="true" @click="tutupDrawer()" />
      <div :id="idDrawer" ref="drawer" role="dialog" aria-modal="true" aria-label="Menu admin" class="relative flex w-72 max-w-[85vw] flex-col overflow-y-auto bg-secondary">
        <button
          type="button"
          class="absolute top-3 right-3 inline-flex size-11 items-center justify-center rounded-md text-on-dark transition-colors duration-150 hover:bg-secondary-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark"
          aria-label="Tutup menu"
          @click="tutupDrawer()"
        >
          <X class="size-5" aria-hidden="true" />
        </button>
        <AdminSidebar />
      </div>
    </div>

    <div class="flex min-w-0 flex-1 flex-col">
      <header class="sticky top-0 z-30 flex min-h-19 items-center gap-3 border-b border-border bg-surface px-4 py-3 sm:px-7">
        <button
          ref="tombolMenu"
          type="button"
          class="-ml-1 inline-flex size-11 shrink-0 items-center justify-center rounded-md transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:hidden"
          aria-label="Buka menu"
          :aria-expanded="drawerBuka ? 'true' : 'false'"
          :aria-controls="idDrawer"
          @click="bukaDrawer"
        >
          <Menu class="size-5" aria-hidden="true" />
        </button>
        <div class="flex min-w-0 flex-col">
          <span v-if="route.meta.grup" class="text-xs text-muted">{{ route.meta.grup }}</span>
          <h1 class="truncate text-xl font-bold tracking-tight">{{ route.meta.title ?? 'Admin panel' }}</h1>
        </div>
        <div class="ml-auto flex shrink-0 items-center gap-3">
          <span v-if="labelRole" class="hidden rounded-full border border-primary/20 bg-primary-soft px-3 py-1.5 text-xs font-bold text-primary sm:inline">{{ labelRole }}</span>
          <span class="hidden text-sm font-semibold md:inline">{{ auth.user?.full_name ?? auth.user?.email }}</span>
          <Button variant="outline" size="sm" @click="keluar">Keluar</Button>
        </div>
      </header>
      <main id="konten" tabindex="-1" class="mx-auto w-full max-w-[1440px] min-w-0 flex-1 px-4 py-7 focus:outline-none sm:px-7 lg:py-9">
        <!-- Key per path: pindah dari /admin/produk/1 ke /2 membuat form baru (bukan sisa data lama). -->
        <RouterView v-slot="{ Component, route: rute }">
          <component :is="Component" :key="rute.path" />
        </RouterView>
      </main>
    </div>
  </div>
</template>
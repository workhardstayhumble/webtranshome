<script setup>
import { RouterLink, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useMenuStore } from '@/stores/menu'

// Isi sidebar dari store menu (GET /admin/menus, sudah difilter backend sesuai izin). Dipakai di desktop dan drawer HP.
const route = useRoute()
const auth = useAuthStore()
const menu = useMenuStore()

// Dashboard aktif hanya di /admin persis; menu lain juga aktif di sub-halamannya (mis. /admin/produk/12).
function aktif(item) {
  if (item.route === '/admin') return route.path === '/admin'
  return route.path === item.route || route.path.startsWith(`${item.route}/`)
}

const kelasFokus = 'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-on-dark'
</script>

<template>
  <div class="flex flex-col gap-8 px-3 py-6 text-sm text-on-dark">
    <RouterLink to="/admin" :class="['self-start rounded-md px-2.5 text-xl font-extrabold tracking-tight text-secondary-foreground', kelasFokus]">
      Trans<span class="text-primary-on-dark">home</span>
      <span class="ml-2 rounded border border-white/15 px-2 py-0.5 align-middle text-[10px] font-semibold tracking-widest text-on-dark-muted uppercase">Admin</span>
    </RouterLink>

    <nav aria-label="Menu admin" class="flex flex-col gap-5">
      <div v-if="menu.status === 'error'" class="flex flex-col items-start gap-2 px-2.5">
        <p>Menu gagal dimuat.</p>
        <button
          type="button"
          :class="['inline-flex h-10 items-center rounded-md border border-on-dark-muted px-3 font-semibold transition-colors duration-150 hover:bg-secondary-hover', kelasFokus]"
          @click="menu.muat(auth.user?.id, { paksa: true })"
        >
          Coba lagi
        </button>
      </div>

      <div v-else-if="menu.status !== 'ready'" class="flex flex-col gap-3 px-2.5" aria-busy="true">
        <span class="sr-only">Memuat menu…</span>
        <div v-for="n in 6" :key="n" class="h-4 animate-pulse rounded bg-secondary-hover motion-reduce:animate-none" :class="n % 3 === 1 ? 'w-16' : 'w-36'" />
      </div>

      <template v-else>
        <div v-for="g in menu.grup" :key="g.id" class="flex flex-col gap-0.5">
          <p class="px-2.5 pb-1 text-xs font-bold tracking-wider text-on-dark-muted uppercase">{{ g.name }}</p>
          <ul class="flex flex-col gap-0.5">
            <li v-for="item in g.anak" :key="item.id">
              <RouterLink
                :to="item.route"
                :aria-current="aktif(item) ? 'page' : undefined"
                :class="[
                  'flex h-11 items-center rounded-md px-3 font-semibold transition-colors duration-150',
                  aktif(item) ? 'bg-primary text-primary-foreground' : 'text-on-dark hover:bg-secondary-hover',
                  kelasFokus,
                ]"
              >
                {{ item.name }}
              </RouterLink>
            </li>
          </ul>
        </div>
      </template>
    </nav>
  </div>
</template>
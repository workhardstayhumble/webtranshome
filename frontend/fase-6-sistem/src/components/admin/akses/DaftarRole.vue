<script setup>
import { Lock } from 'lucide-vue-next'
import { formatAngka } from '@/utils/format'

defineProps({
  // GET /admin/roles
  roles: { type: Array, required: true },
  aktif: { type: String, default: '' },
})

const emit = defineEmits(['pilih'])
</script>

<template>
  <nav aria-label="Pilih role">
    <ul class="grid gap-2 sm:grid-cols-2 lg:grid-cols-1">
      <li v-for="r in roles" :key="r.code">
        <button
          type="button"
          :aria-current="r.code === aktif ? 'true' : undefined"
          :data-role="r.code"
          class="flex w-full cursor-pointer flex-col items-start gap-1 rounded-xl border p-4 text-left transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          :class="r.code === aktif ? 'border-secondary bg-surface shadow-sm ring-1 ring-secondary' : 'border-border bg-surface hover:border-border-strong'"
          @click="emit('pilih', r.code)"
        >
          <span class="flex items-center gap-1.5 font-bold">
            {{ r.name }}
            <Lock v-if="r.is_locked" class="size-4 text-muted" aria-hidden="true" />
            <span v-if="r.is_locked" class="sr-only">(terkunci)</span>
          </span>
          <!-- Di bawah 1024 px deskripsi cukup tampil di atas matriks, supaya daftar role tetap ringkas. -->
          <span class="hidden text-sm text-muted lg:block">{{ r.description }}</span>
          <span class="text-xs font-semibold text-muted">
            <template v-if="r.user_count != null">{{ formatAngka(r.user_count) }} admin aktif · </template>{{ formatAngka(r.permissions.length) }} izin
          </span>
        </button>
      </li>
    </ul>
  </nav>
</template>
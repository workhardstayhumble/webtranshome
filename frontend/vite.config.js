import { existsSync, statSync } from 'node:fs'
import { resolve } from 'node:path'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig, normalizePath } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// Kode dipisah per fase: setiap <folder fase>/src memakai susunan yang sama (components/, views/, services/, ...).
// Import `@/...` dicari di src semua folder fase, berurutan. Satu path (mis. components/ui/Button.vue) hanya boleh ada di satu fase.
const FOLDER_FASE = [
  'fase-0-fondasi',
  'fase-1-katalog',
  'fase-2-admin-katalog',
  'fase-3-akun-member',
  'fase-4-membership-admin',
  'fase-5-konten',
  'fase-6-sistem',
]
const DASAR = fileURLToPath(new URL('.', import.meta.url))
const AKHIRAN = ['', '.js', '.vue', '.json', '/index.js']

function aliasFase() {
  return {
    name: 'transhome-alias-fase',
    enforce: 'pre',
    resolveId(sumber) {
      if (!sumber.startsWith('@/')) return null
      const [jalur, query] = sumber.slice(2).split('?')
      for (const folder of FOLDER_FASE) {
        for (const akhiran of AKHIRAN) {
          const calon = resolve(DASAR, folder, 'src', jalur + akhiran)
          // Path dinormalisasi seperti hasil resolve Vite, supaya satu file tidak termuat sebagai dua modul.
          if (existsSync(calon) && statSync(calon).isFile()) return normalizePath(calon) + (query ? `?${query}` : '')
        }
      }
      return null
    },
  }
}

export default defineConfig({
  plugins: [aliasFase(), vue(), tailwindcss()],
  // API menyimpan foto unggahan sebagai alamat relatif, sedangkan saat development API ada di port lain.
  server: { proxy: { '/api/v1/uploads': 'http://localhost:8888' } },
})

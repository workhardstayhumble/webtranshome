<script setup>
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ImageOff, Minus, Plus, ShoppingCart, Trash2 } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import Modal from '@/components/ui/Modal.vue'
import WhatsAppButton from './WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { useAuthStore } from '@/stores/auth'
import { useKeranjangStore } from '@/stores/keranjang'

const route = useRoute()
const auth = useAuthStore()
const keranjang = useKeranjangStore()
const { keranjangLink } = useWhatsApp()

const terbuka = ref(false)
const linkWhatsApp = computed(() => keranjangLink(keranjang.items))
const linkMasuk = computed(() => ({ path: '/masuk', query: { redirect: route.fullPath } }))
const labelTombol = computed(() => (keranjang.jumlahBarang ? `Keranjang, ${keranjang.jumlahBarang} barang` : 'Keranjang, kosong'))

// Tautan di dalam panel berpindah halaman; panel tidak boleh tertinggal menutupi halaman baru.
watch(
  () => route.fullPath,
  () => (terbuka.value = false),
)

// Isian bisa kosong atau di luar batas; kolomnya disamakan lagi dengan nilai yang benar-benar tersimpan.
function ubahDariIsian(item, event) {
  keranjang.ubahJumlah(item.id, event.target.value)
  event.target.value = item.jumlah
}

const kelasStepper =
  'inline-flex size-11 items-center justify-center text-foreground transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent md:size-9'
</script>

<template>
  <button
    type="button"
    class="relative inline-flex size-11 shrink-0 items-center justify-center rounded-md text-foreground transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
    :aria-label="labelTombol"
    aria-haspopup="dialog"
    data-tombol-keranjang
    @click="terbuka = true"
  >
    <ShoppingCart class="size-6" aria-hidden="true" />
    <span
      v-if="keranjang.jumlahBarang"
      class="absolute top-0.5 right-0.5 flex h-5 min-w-5 items-center justify-center rounded-full bg-primary px-1 text-xs font-bold text-white"
      aria-hidden="true"
      data-jumlah-keranjang
    >
      {{ keranjang.jumlahBarang > 99 ? '99+' : keranjang.jumlahBarang }}
    </span>
  </button>

  <Modal v-model:open="terbuka" side title="Keranjang" :description="keranjang.jumlahBarang ? `${keranjang.jumlahBarang} barang` : ''">
    <EmptyState v-if="!keranjang.jumlahBarang" title="Belum ada barang di keranjang" description="Buka halaman produk, lalu pilih Tambah ke keranjang." data-keranjang-kosong>
      <template #icon><ShoppingCart class="size-7" aria-hidden="true" /></template>
      <template #actions>
        <Button to="/katalog" @click="terbuka = false">Lihat katalog</Button>
      </template>
    </EmptyState>

    <ul v-else class="flex flex-col divide-y divide-border" data-daftar-keranjang>
      <li v-for="item in keranjang.items" :key="item.id" class="flex gap-3 py-4 first:pt-0">
        <img
          v-if="item.image_url"
          :src="item.image_url"
          :alt="item.name"
          width="72"
          height="72"
          loading="lazy"
          class="size-18 shrink-0 rounded-md border border-border object-cover"
        />
        <div v-else class="flex size-18 shrink-0 items-center justify-center rounded-md border border-border bg-subtle text-muted">
          <ImageOff class="size-6" aria-hidden="true" />
        </div>

        <div class="flex min-w-0 flex-1 flex-col gap-2">
          <div class="flex items-start justify-between gap-2">
            <div class="flex min-w-0 flex-col gap-0.5">
              <RouterLink :to="`/produk/${item.slug}`" class="rounded-sm font-semibold hover:text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
                {{ item.name }}
              </RouterLink>
              <span class="font-mono text-xs text-muted">SKU {{ item.sku }}</span>
            </div>
            <button
              type="button"
              class="-mt-2 -mr-2 inline-flex size-11 shrink-0 items-center justify-center rounded-md text-muted transition-colors duration-150 hover:bg-subtle hover:text-danger focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              :aria-label="`Hapus ${item.name} dari keranjang`"
              @click="keranjang.hapus(item.id)"
            >
              <Trash2 class="size-5" aria-hidden="true" />
            </button>
          </div>

          <div class="flex flex-wrap items-center gap-x-3 gap-y-1">
            <div class="inline-flex items-center overflow-hidden rounded-md border border-border-strong">
              <button
                type="button"
                :class="kelasStepper"
                :disabled="item.jumlah <= keranjang.jumlahMinimal(item)"
                :aria-label="`Kurangi jumlah ${item.name}`"
                @click="keranjang.ubahJumlah(item.id, item.jumlah - 1)"
              >
                <Minus class="size-4" aria-hidden="true" />
              </button>
              <input
                :value="item.jumlah"
                type="text"
                inputmode="numeric"
                autocomplete="off"
                :aria-label="`Jumlah ${item.name}`"
                class="h-11 w-14 border-x border-border-strong bg-surface text-center text-sm font-semibold focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring md:h-9"
                @change="ubahDariIsian(item, $event)"
              />
              <button type="button" :class="kelasStepper" :aria-label="`Tambah jumlah ${item.name}`" @click="keranjang.ubahJumlah(item.id, item.jumlah + 1)">
                <Plus class="size-4" aria-hidden="true" />
              </button>
            </div>
            <span class="text-sm text-muted">{{ item.unit_sale }}</span>
            <span v-if="keranjang.jumlahMinimal(item) > 1" class="w-full text-xs text-muted">Minimal {{ keranjang.jumlahMinimal(item) }} {{ item.unit_sale }}</span>
          </div>
        </div>
      </li>
    </ul>

    <template v-if="keranjang.jumlahBarang" #footer>
      <div class="flex w-full flex-col gap-3 pb-[env(safe-area-inset-bottom)]">
        <p v-if="!auth.isLoggedIn" class="text-sm text-muted">
          <RouterLink :to="linkMasuk" class="rounded-sm font-semibold text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Masuk</RouterLink>
          agar nama dan nomor telepon terisi otomatis, atau lengkapi sendiri di chat.
        </p>
        <WhatsAppButton size="lg" block :href="linkWhatsApp" label="Tanya ke WhatsApp" />
      </div>
    </template>
  </Modal>
</template>

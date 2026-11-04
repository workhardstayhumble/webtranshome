<script setup>
import { ShoppingCart } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { useKeranjangStore } from '@/stores/keranjang'
import { useUiStore } from '@/stores/ui'

const props = defineProps({
  produk: { type: Object, required: true },
  // true: hanya ikon, untuk bar bawah halaman produk di HP yang sempit.
  ikonSaja: { type: Boolean, default: false },
})

const keranjang = useKeranjangStore()
const ui = useUiStore()

function tambah() {
  const jumlah = keranjang.tambah(props.produk)
  ui.tampilkanToast({ pesan: `Masuk keranjang. Jumlah ${jumlah} ${props.produk.unit_sale}.`, jenis: 'success' })
}
</script>

<template>
  <Button v-if="ikonSaja" variant="outline" class="size-11 px-0!" aria-label="Tambah ke keranjang" data-tambah-keranjang @click="tambah">
    <ShoppingCart class="size-5" aria-hidden="true" />
  </Button>
  <Button v-else variant="outline" size="lg" block data-tambah-keranjang @click="tambah">
    <template #icon><ShoppingCart class="size-4.5" aria-hidden="true" /></template>
    Tambah ke keranjang
  </Button>
</template>

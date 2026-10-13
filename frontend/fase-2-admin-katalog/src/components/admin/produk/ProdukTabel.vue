<script setup>
import { ImageOff, Star } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import StockBadge from '@/components/ui/StockBadge.vue'
import Switch from '@/components/ui/Switch.vue'
import Table from '@/components/ui/Table.vue'
import { usePermission } from '@/composables/usePermission'
import { formatRupiah } from '@/utils/format'

defineProps({
  // Item GET /admin/products (+ `sibuk` saat toggle sedang dikirim)
  produk: { type: Array, required: true },
})

const emit = defineEmits(['unggulan', 'aktif'])
const { can } = usePermission()

const KOLOM = [
  { key: 'produk', label: 'Produk', class: 'min-w-64' },
  { key: 'kategori', label: 'Kategori' },
  { key: 'brand', label: 'Brand', class: 'whitespace-nowrap' },
  { key: 'harga', label: 'Harga umum', align: 'right' },
  { key: 'stok', label: 'Stok' },
  { key: 'unggulan', label: 'Unggulan', align: 'center' },
  { key: 'aktif', label: 'Aktif', align: 'center' },
  { key: 'aksi', label: 'Aksi', align: 'right' },
]

// Produk aktif dipratinjau di halaman publiknya; yang nonaktif lewat pratinjau admin (tidak tampil di website).
const tautanPratinjau = (p) => (p.is_active ? `/produk/${p.slug}` : `/admin/produk/${p.id}/pratinjau`)
</script>

<template>
  <Table :columns="KOLOM" :rows="produk" caption="Daftar produk" bare min-width="min-w-240">
    <template #cell-produk="{ row }">
      <div class="flex items-center gap-3">
        <span class="flex size-11 shrink-0 items-center justify-center overflow-hidden rounded-md bg-subtle text-faint">
          <img v-if="row.primary_image_url" :src="row.primary_image_url" alt="" width="44" height="44" loading="lazy" class="size-full object-cover" />
          <ImageOff v-else class="size-4.5" aria-hidden="true" />
        </span>
        <span class="flex min-w-0 flex-col">
          <strong class="font-semibold">{{ row.name }}</strong>
          <span class="font-mono text-xs text-faint">{{ row.sku }}</span>
        </span>
      </div>
    </template>
    <template #cell-kategori="{ row }">{{ row.category?.name ?? '–' }}</template>
    <template #cell-brand="{ row }">{{ row.brand?.name ?? '–' }}</template>
    <template #cell-harga="{ row }">
      <span class="whitespace-nowrap"><strong class="tabular-nums">{{ formatRupiah(row.price_general) }}</strong> <span class="text-muted">/{{ row.unit_sale }}</span></span>
    </template>
    <template #cell-stok="{ row }">
      <StockBadge :status="row.stock_status" :keterangan="row.stock_qty_label ?? ''" />
    </template>
    <template #cell-unggulan="{ row }">
      <button
        type="button"
        class="inline-flex size-11 items-center justify-center rounded-md border transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-60 md:size-9"
        :class="row.is_featured ? 'border-primary text-primary' : 'border-border-strong text-faint hover:border-foreground hover:text-foreground'"
        :aria-pressed="row.is_featured ? 'true' : 'false'"
        :aria-label="`Produk unggulan: ${row.name}`"
        :disabled="!can('product.featured') || row.sibuk"
        @click="emit('unggulan', row)"
      >
        <Star class="size-4.5" :fill="row.is_featured ? 'currentColor' : 'none'" aria-hidden="true" />
      </button>
    </template>
    <template #cell-aktif="{ row }">
      <Switch
        :model-value="row.is_active"
        :label="`Tampil di website: ${row.name}`"
        hide-label
        :disabled="!can('product.update') || row.sibuk"
        @update:model-value="(nilai) => emit('aktif', row, nilai)"
      />
    </template>
    <template #cell-aksi="{ row }">
      <div class="flex items-center justify-end gap-1.5 whitespace-nowrap">
        <Button v-if="can('product.update')" variant="outline" size="sm" :to="`/admin/produk/${row.id}`">
          Ubah <span class="sr-only">{{ row.name }}</span>
        </Button>
        <a
          :href="tautanPratinjau(row)"
          target="_blank"
          rel="noopener"
          class="inline-flex h-10 items-center rounded-md px-2.5 text-sm font-semibold text-muted transition-colors duration-150 hover:bg-subtle hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          Pratinjau <span class="sr-only">{{ row.name }} (tab baru)</span>
        </a>
      </div>
    </template>
  </Table>
</template>
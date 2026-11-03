<script setup>
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { BadgeCheck, FileDown, Share2, ShieldCheck } from 'lucide-vue-next'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import CopyButton from '@/components/ui/CopyButton.vue'
import PriceTag from '@/components/ui/PriceTag.vue'
import StockBadge from '@/components/ui/StockBadge.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import JasaTerkait from '@/components/jasa/JasaTerkait.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { hargaDisembunyikan, labelTanyaProduk } from '@/utils/whatsapp'
import { useUiStore } from '@/stores/ui'
import { salinTeks } from '@/utils/clipboard'

const props = defineProps({
  // GET /products/:slug
  produk: { type: Object, required: true },
  // 'h2' di pratinjau admin (judul halaman admin sudah h1).
  tagJudul: { type: String, default: 'h1' },
})

const ui = useUiStore()
const { productLink, bulkLink } = useWhatsApp()
const teksSku = ref(null)

const habis = computed(() => props.produk.stock_status === 'HABIS')
const garansi = computed(() => props.produk.specifications?.garansi || null)
const sni = computed(() => {
  const nilai = String(props.produk.specifications?.sni ?? '').trim()
  return nilai !== '' && !/^(tidak|tidak ada|-)$/i.test(nilai)
})

// Clipboard ditolak browser: SKU dipilih otomatis supaya tinggal disalin manual.
function skuGagalDisalin() {
  if (teksSku.value) window.getSelection()?.selectAllChildren(teksSku.value)
  ui.tampilkanToast({ pesan: 'SKU sudah dipilih. Salin secara manual.', jenis: 'info' })
}

async function bagikan() {
  const url = new URL(`/produk/${props.produk.slug}`, window.location.origin).href
  if (navigator.share) {
    try {
      await navigator.share({ title: props.produk.name, url })
      return
    } catch (error) {
      if (error?.name === 'AbortError') return
    }
  }
  const tersalin = await salinTeks(url)
  ui.tampilkanToast(
    tersalin
      ? { pesan: 'Link disalin', jenis: 'success' }
      : { pesan: 'Link gagal disalin. Salin dari kolom alamat browser.', jenis: 'warning' },
  )
}
</script>

<template>
  <div class="flex min-w-0 flex-col gap-5">
    <div class="flex flex-col gap-2">
      <RouterLink
        v-if="produk.brand"
        :to="`/brand/${produk.brand.slug}`"
        class="self-start rounded-sm text-sm font-bold tracking-wide text-muted uppercase transition-colors duration-150 hover:text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        {{ produk.brand.name }}
      </RouterLink>
      <component :is="tagJudul" class="text-2xl leading-tight font-extrabold sm:text-3xl">{{ produk.name }}</component>
      <div class="flex flex-wrap items-center gap-x-3 gap-y-2">
        <span class="font-mono text-sm text-muted">SKU <span ref="teksSku" class="font-semibold text-foreground">{{ produk.sku }}</span></span>
        <CopyButton :text="produk.sku" label="Salin SKU" copied-label="SKU disalin" @failed="skuGagalDisalin" />
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <PriceTag size="lg" :price="produk.price_general" :show="produk.show_price !== false" :unit="produk.unit_sale" />
      <p v-if="produk.min_order > 1" class="text-sm text-muted">Minimal pembelian {{ produk.min_order }} {{ produk.unit_sale }}</p>
    </div>

    <div class="flex flex-wrap gap-2">
      <StockBadge :status="produk.stock_status" :keterangan="produk.stock_qty_label ?? ''" />
      <Badge v-if="garansi" variant="info"><ShieldCheck class="size-3.5" aria-hidden="true" />Garansi {{ garansi }}</Badge>
      <Badge v-if="sni" variant="success"><BadgeCheck class="size-3.5" aria-hidden="true" />SNI</Badge>
    </div>

    <div class="flex flex-col gap-2">
      <WhatsAppButton size="lg" block wrap :href="productLink(produk)" :label="labelTanyaProduk(produk)" />
      <p class="text-sm text-muted">
        {{
          habis
            ? 'Stok sedang habis. Tanyakan kapan tersedia lagi atau minta rekomendasi pengganti.'
            : hargaDisembunyikan(produk)
              ? 'Harga produk ini diberikan lewat WhatsApp. Pesan otomatis berisi nama produk, SKU, dan link halaman ini.'
              : 'Pesan otomatis berisi nama produk, SKU, dan link halaman ini, jadi tim kami langsung tahu barangnya.'
        }}
      </p>
    </div>

    <div class="flex flex-wrap gap-2">
      <Button variant="outline" @click="bagikan">
        <template #icon><Share2 class="size-4.5" aria-hidden="true" /></template>
        Bagikan
      </Button>
      <Button v-if="produk.datasheet_pdf_url" variant="outline" :href="produk.datasheet_pdf_url">
        <template #icon><FileDown class="size-4.5" aria-hidden="true" /></template>
        Unduh datasheet
      </Button>
    </div>

    <div class="flex flex-col gap-3 rounded-xl bg-primary-soft p-4 sm:flex-row sm:items-center">
      <div class="flex min-w-0 flex-1 flex-col gap-0.5">
        <p class="font-bold">Butuh jumlah besar?</p>
        <p class="text-sm text-muted">Untuk proyek atau kontraktor, chat langsung dengan tim proyek kami.</p>
      </div>
      <WhatsAppButton :href="bulkLink(produk)" label="Chat tim proyek" />
    </div>

    <!-- Jasa pemasangan dari kategori produk ini (KEPUTUSAN #20); tidak tampil bila tidak ada. -->
    <JasaTerkait :produk="produk" />

    <div v-if="produk.rooms?.length" class="flex flex-wrap items-center gap-2">
      <span class="text-sm font-semibold">Cocok untuk:</span>
      <RouterLink
        v-for="r in produk.rooms"
        :key="r.id"
        :to="`/ruangan/${r.slug}`"
        class="inline-flex h-11 items-center rounded-full border border-border-strong bg-surface px-4 text-sm font-semibold transition-colors duration-150 hover:border-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:h-9"
      >
        {{ r.name }}
      </RouterLink>
    </div>
  </div>
</template>

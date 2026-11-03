<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ImageOff } from 'lucide-vue-next'
import PriceTag from '@/components/ui/PriceTag.vue'
import StockBadge from '@/components/ui/StockBadge.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { labelTanyaProduk } from '@/utils/whatsapp'

const props = defineProps({
  product: { type: Object, required: true },
})

const { productLink } = useWhatsApp()
const tautan = computed(() => `/produk/${props.product.slug}`)
const linkWhatsApp = computed(() => productLink(props.product))
</script>

<template>
  <article class="flex flex-wrap items-center gap-4 rounded-xl border border-border bg-surface p-3 transition-colors duration-150 hover:border-border-strong">
    <RouterLink :to="tautan" class="size-28 shrink-0 overflow-hidden rounded-lg bg-subtle" tabindex="-1" aria-hidden="true">
      <img v-if="product.primary_image_url" :src="product.primary_image_url" :alt="product.name" loading="lazy" width="112" height="112" class="size-full object-cover" />
      <span v-else class="flex size-full items-center justify-center text-faint"><ImageOff class="size-6" aria-hidden="true" /></span>
    </RouterLink>

    <div class="flex min-w-0 flex-[1_1_16rem] flex-col gap-1">
      <span v-if="product.brand" class="text-xs font-bold tracking-wide text-muted uppercase">{{ product.brand.name }}</span>
      <RouterLink :to="tautan" class="rounded-sm text-base font-bold transition-colors duration-150 hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
        {{ product.name }}
      </RouterLink>
      <span class="font-mono text-xs text-faint">SKU {{ product.sku }}</span>
    </div>

    <div class="flex flex-col items-start gap-2 sm:items-end">
      <PriceTag :price="product.price_general" :show="product.show_price !== false" :unit="product.unit_sale" />
      <StockBadge class="sm:self-end" :status="product.stock_status" :keterangan="product.stock_qty_label ?? ''" />
      <WhatsAppButton size="sm" :href="linkWhatsApp" :label="labelTanyaProduk(product)" />
    </div>
  </article>
</template>

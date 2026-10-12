<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowUpRight, ImageOff } from 'lucide-vue-next'
import PriceTag from '@/components/ui/PriceTag.vue'
import StockBadge from '@/components/ui/StockBadge.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
const props = defineProps({ product: { type: Object, required: true } })
const { productLink } = useWhatsApp()
const tautan = computed(() => `/produk/${props.product.slug}`)
const linkWhatsApp = computed(() => productLink(props.product))
const habis = computed(() => props.product.stock_status === 'HABIS')
</script>
<template>
  <article class="group th-surface-card th-hover-lift flex min-w-0 flex-col overflow-hidden">
    <RouterLink :to="tautan" class="relative block aspect-[4/3] overflow-hidden bg-[linear-gradient(145deg,#f1eee6,#e0d9c8)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-inset" tabindex="-1" aria-hidden="true">
      <img v-if="product.primary_image_url" :src="product.primary_image_url" :alt="product.name" loading="lazy" width="400" height="300" class="size-full object-cover transition-transform duration-300 group-hover:scale-[1.035]" />
      <span v-else class="flex size-full flex-col items-center justify-center gap-2 text-muted"><ImageOff class="size-7" aria-hidden="true" /><span class="text-xs">Foto produk belum tersedia</span></span>
      <span class="absolute top-3 right-3 flex size-8 items-center justify-center rounded-full bg-white/90 text-primary shadow-sm"><ArrowUpRight class="size-4" aria-hidden="true" /></span>
    </RouterLink>
    <div class="flex flex-1 flex-col gap-2 p-3.5 sm:p-4">
      <span v-if="product.brand" class="text-xs font-bold tracking-wide text-muted uppercase">{{ product.brand.name }}</span>
      <RouterLink :to="tautan" class="line-clamp-2 min-h-[2.5em] rounded-sm font-semibold leading-snug text-foreground hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">{{ product.name }}</RouterLink>
      <span class="text-xs text-muted">SKU {{ product.sku }}</span>
      <PriceTag class="mt-0.5" size="card" :price="product.price_general" :unit="product.unit_sale" />
      <StockBadge :status="product.stock_status" :keterangan="product.stock_qty_label ?? ''" />
      <WhatsAppButton class="mt-auto pt-2" :href="linkWhatsApp" :label="habis ? 'Tanya ketersediaan' : 'Tanya via WhatsApp'" block wrap />
    </div>
  </article>
</template>
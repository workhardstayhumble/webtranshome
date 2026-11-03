<script setup>
import { RouterLink } from 'vue-router'
import { BadgeCheck, MapPin, Wrench } from 'lucide-vue-next'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { LABEL_HARGA_JASA, tautanJasa } from '@/utils/jasa'

defineProps({
  // Item GET /services: { name, slug, summary, image_url, service_area }
  jasa: { type: Object, required: true },
  // Tingkat judul kartu: h3 di bawah judul seksi (Beranda), h2 bila kartu langsung di bawah h1 (/jasa).
  tingkat: { type: String, default: 'h3' },
})

const { jasaLink } = useWhatsApp()
</script>

<template>
  <article class="relative flex h-full flex-col overflow-hidden rounded-xl border border-border bg-surface transition-colors duration-150 hover:border-border-strong" data-kartu-jasa>
    <span class="flex aspect-video items-center justify-center bg-subtle text-faint">
      <img v-if="jasa.image_url" :src="jasa.image_url" alt="" width="800" height="450" loading="lazy" class="size-full object-cover" />
      <Wrench v-else class="size-8" aria-hidden="true" />
    </span>
    <div class="flex flex-1 flex-col gap-2 p-4">
      <component :is="tingkat" class="text-lg leading-snug font-bold">
        <!-- Tautan meliputi seluruh kartu (::after); tombol WhatsApp di atasnya (relative z-10) tetap bisa diklik. -->
        <RouterLink
          :to="tautanJasa(jasa.slug)"
          class="rounded-sm after:absolute after:inset-0 after:content-[''] hover:text-primary focus-visible:outline-none focus-visible:after:rounded-xl focus-visible:after:ring-2 focus-visible:after:ring-ring"
        >
          {{ jasa.name }}
        </RouterLink>
      </component>
      <p class="line-clamp-3 text-sm text-muted">{{ jasa.summary }}</p>
      <p class="mt-auto flex items-center gap-1.5 pt-1 font-bold text-success-ink" data-harga-jasa>
        <BadgeCheck class="size-5 shrink-0" aria-hidden="true" />
        {{ LABEL_HARGA_JASA }}
      </p>
      <p v-if="jasa.service_area" class="flex items-center gap-1.5 text-sm text-muted">
        <MapPin class="size-4 shrink-0" aria-hidden="true" />
        {{ jasa.service_area }}
      </p>
      <div class="relative z-10 pt-2">
        <WhatsAppButton block :href="jasaLink(jasa)" label="Konsultasi via WhatsApp" />
      </div>
    </div>
  </article>
</template>

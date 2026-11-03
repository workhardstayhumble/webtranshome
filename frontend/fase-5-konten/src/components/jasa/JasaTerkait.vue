<script setup>
import { ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { Wrench } from 'lucide-vue-next'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { listServices } from '@/services/kontenService'
import { LABEL_HARGA_JASA, tautanJasa } from '@/utils/jasa'

// Detail produk: "Butuh pemasangan?" berisi jasa yang terhubung ke kategori produk ini (diatur di admin › Jasa).
// Tanpa jasa terkait atau bila API gagal, kotak tidak tampil (informasi tambahan, bukan isi utama halaman).
const props = defineProps({
  // Produk dari GET /products/:slug (butuh sku, name, category.slug).
  produk: { type: Object, required: true },
})

const { jasaLink } = useWhatsApp()
const jasa = ref([])

watch(
  () => props.produk?.category?.slug,
  async (slug) => {
    jasa.value = []
    if (!slug) return
    try {
      jasa.value = await listServices({ category: slug, limit: 2 })
    } catch {
      jasa.value = []
    }
  },
  { immediate: true },
)
</script>

<template>
  <section v-if="jasa.length" aria-labelledby="judul-jasa-terkait" class="rounded-xl border border-border bg-subtle p-4 sm:p-5" data-jasa-terkait>
    <h2 id="judul-jasa-terkait" class="flex items-center gap-2 text-lg font-extrabold">
      <Wrench class="size-5 text-primary" aria-hidden="true" />
      Butuh pemasangan?
    </h2>
    <ul class="mt-3 flex flex-col gap-3">
      <li v-for="j in jasa" :key="j.id" class="flex flex-col gap-3 rounded-lg bg-surface p-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex min-w-0 flex-col gap-0.5">
          <RouterLink :to="tautanJasa(j.slug)" class="self-start rounded-sm font-bold hover:text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
            {{ j.name }}
          </RouterLink>
          <span class="text-sm text-muted">{{ j.summary }}</span>
          <span class="text-sm font-semibold text-success-ink">{{ LABEL_HARGA_JASA }}</span>
        </div>
        <WhatsAppButton class="shrink-0" :href="jasaLink(j, produk)" label="Konsultasi jasa" />
      </li>
    </ul>
  </section>
</template>

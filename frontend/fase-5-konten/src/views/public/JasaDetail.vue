<script setup>
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { BadgeCheck, MapPin } from 'lucide-vue-next'
import Breadcrumb from '@/components/ui/Breadcrumb.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import NotFoundState from '@/components/ui/NotFoundState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import IsiArtikel from '@/components/konten/IsiArtikel.vue'
import { remahSeo, urlSitus, useSeo } from '@/composables/useSeo'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { getService } from '@/services/kontenService'
import { pesanError } from '@/services/errors'
import { LABEL_HARGA_JASA } from '@/utils/jasa'

// /jasa/:slug (KEPUTUSAN #20): penjelasan jasa + kotak konsultasi WhatsApp (menempel di samping pada desktop).
const route = useRoute()
const { jasaLink } = useWhatsApp()

const status = ref('loading') // loading | ready | tidak-ada | error
const jasa = ref(null)
const pesan = ref('')

async function muat() {
  status.value = 'loading'
  try {
    jasa.value = await getService(route.params.slug)
    status.value = 'ready'
  } catch (error) {
    if (error?.response?.status === 404) {
      status.value = 'tidak-ada'
      return
    }
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}
watch(
  () => route.params.slug,
  (slug) => {
    if (slug) muat()
  },
  { immediate: true },
)

const remah = computed(() => [{ label: 'Beranda', to: '/' }, { label: 'Jasa', to: '/jasa' }, { label: jasa.value?.name ?? 'Jasa' }])

useSeo({
  title: () => (status.value === 'ready' ? `${jasa.value.name} di Sleman` : status.value === 'tidak-ada' ? 'Jasa tidak ditemukan' : null),
  description: () => (status.value === 'ready' ? `${jasa.value.summary} Konsultasi lewat WhatsApp dengan tim Transhome.` : null),
  canonical: () => `/jasa/${route.params.slug}`,
  robots: () => (status.value === 'tidak-ada' ? 'noindex' : null),
  image: () => jasa.value?.image_url,
  jsonLd: () => {
    if (status.value !== 'ready') return null
    const j = jasa.value
    const url = urlSitus(`/jasa/${j.slug}`)
    return [
      remahSeo(remah.value.map((r, i) => (i === 2 ? { ...r, to: `/jasa/${j.slug}` } : r))),
      {
        '@type': 'Service',
        name: j.name,
        description: j.summary,
        url,
        provider: { '@type': 'HomeAndConstructionBusiness', name: 'Transhome', url: urlSitus('/') },
        ...(j.service_area ? { areaServed: j.service_area } : {}),
      },
    ]
  },
})
</script>

<template>
  <NotFoundState
    v-if="status === 'tidak-ada'"
    title="Jasa tidak ditemukan"
    description="Jasa yang kamu buka tidak ada atau sedang tidak tersedia."
    action-label="Lihat semua jasa"
    action-to="/jasa"
  />

  <div v-else class="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:py-6">
    <Breadcrumb :items="remah" />

    <div v-if="status === 'loading'" class="mt-4 grid gap-8 lg:grid-cols-[minmax(0,1fr)_20rem]" aria-busy="true">
      <span class="sr-only">Memuat jasa…</span>
      <div class="flex flex-col gap-4">
        <Skeleton class="aspect-video rounded-xl" />
        <Skeleton class="h-9 w-2/3" />
        <Skeleton class="h-24" />
      </div>
      <Skeleton class="h-56 rounded-xl" />
    </div>
    <ErrorState v-else-if="status === 'error'" class="mt-6" :message="pesan" @retry="muat" />

    <article v-else class="mt-4 grid items-start gap-8 lg:grid-cols-[minmax(0,1fr)_20rem]">
      <header class="flex flex-col gap-4 lg:col-start-1">
        <img v-if="jasa.image_url" :src="jasa.image_url" alt="" width="800" height="450" class="aspect-video w-full rounded-xl object-cover" />
        <h1 class="text-3xl leading-tight font-extrabold sm:text-4xl">{{ jasa.name }}</h1>
        <p class="text-lg text-muted">{{ jasa.summary }}</p>
      </header>

      <!-- Kotak konsultasi: di HP tepat setelah judul, di desktop kolom kanan yang menempel saat digulir. -->
      <aside
        aria-label="Konsultasi jasa"
        class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-5 lg:sticky lg:top-[calc(var(--tinggi-header,0px)+1.5rem)] lg:col-start-2 lg:row-span-2 lg:row-start-1"
        data-konsultasi-jasa
      >
        <p class="flex items-center gap-2 text-2xl font-extrabold text-success-ink" data-harga-jasa>
          <BadgeCheck class="size-6 shrink-0" aria-hidden="true" />
          {{ LABEL_HARGA_JASA }}
        </p>
        <p v-if="jasa.service_area" class="flex items-center gap-1.5 text-sm text-muted">
          <MapPin class="size-4 shrink-0" aria-hidden="true" />
          Area layanan: {{ jasa.service_area }}
        </p>
        <WhatsAppButton block size="lg" :href="jasaLink(jasa)" label="Konsultasi via WhatsApp" />
        <p class="text-sm text-muted">Kirim foto, ukuran, dan lokasi pekerjaan supaya tim kami bisa menghitung biaya dan jadwalnya.</p>
      </aside>

      <div class="flex flex-col gap-8 lg:col-start-1">
        <IsiArtikel v-if="jasa.description" :isi="jasa.description" />
        <section v-if="jasa.categories?.length" aria-labelledby="judul-produk-jasa" class="flex flex-col gap-3">
          <h2 id="judul-produk-jasa" class="text-xl font-extrabold">Bahan yang bisa dibeli di Transhome</h2>
          <ul class="flex flex-wrap gap-2">
            <li v-for="k in jasa.categories" :key="k.slug">
              <RouterLink
                :to="`/katalog/${k.slug}`"
                class="inline-flex min-h-11 items-center rounded-full border border-border-strong bg-surface px-4 text-sm font-semibold transition-colors duration-150 hover:border-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-9"
              >
                {{ k.name }}
              </RouterLink>
            </li>
          </ul>
        </section>
        <RouterLink to="/jasa" class="self-start rounded-sm text-sm font-bold text-primary hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
          Lihat jasa lainnya
        </RouterLink>
      </div>
    </article>
  </div>
</template>

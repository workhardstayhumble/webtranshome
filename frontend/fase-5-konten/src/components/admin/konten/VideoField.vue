<script setup>
import { ref, useId } from 'vue'
import { Film, X } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { errorField, pesanError } from '@/services/errors'
import { ATURAN_UNGGAH, cekFile, uploadFile } from '@/services/uploadService'

// Satu video opsional (hero Beranda). Pola sama dengan GambarField: unggah, ganti, hapus, pratinjau.
const url = defineModel({ type: String, default: '' })

defineProps({
  label: { type: String, required: true },
  error: { type: String, default: '' },
  hint: { type: String, default: '' },
})

const id = useId()
const input = ref(null)
const sedang = ref(false)
const galat = ref('')

async function onPilih(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  galat.value = cekFile(file, 'video')
  if (galat.value) return
  sedang.value = true
  try {
    url.value = (await uploadFile(file, 'video')).url
  } catch (error) {
    galat.value = errorField(error).file ?? pesanError(error)
  } finally {
    sedang.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-1.5" role="group" :aria-labelledby="`${id}-label`">
    <span :id="`${id}-label`" class="text-sm font-semibold">{{ label }}</span>
    <div class="flex flex-wrap items-start gap-3">
      <div class="aspect-12/5 w-64 max-w-full shrink-0 overflow-hidden rounded-lg border bg-subtle" :class="galat || error ? 'border-danger' : 'border-border'">
        <!-- Pratinjau tanpa suara dengan kontrol, supaya admin bisa memeriksa isinya. -->
        <video v-if="url" :src="url" class="size-full object-cover" controls muted playsinline preload="metadata" data-pratinjau-video />
        <span v-else class="flex size-full items-center justify-center text-faint"><Film class="size-6" aria-hidden="true" /></span>
      </div>
      <div class="flex min-w-0 flex-1 basis-48 flex-col items-start gap-1.5">
        <Button variant="outline" size="sm" :loading="sedang" @click="input?.click()">{{ url ? 'Ganti video' : 'Unggah video' }}</Button>
        <button
          v-if="url"
          type="button"
          class="inline-flex min-h-10 items-center gap-1 rounded-md px-1 text-sm font-semibold text-muted transition-colors duration-150 hover:text-danger focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          @click="url = ''"
        >
          <X class="size-4" aria-hidden="true" />
          Hapus video
        </button>
        <p class="text-xs text-muted">{{ hint || `MP4 atau WebM, ${ATURAN_UNGGAH.video.pesan.replace('harus MP4 atau WebM, ', '')}.` }}</p>
      </div>
    </div>
    <input ref="input" type="file" class="sr-only" :accept="ATURAN_UNGGAH.video.accept" tabindex="-1" aria-hidden="true" data-input-video @change="onPilih" />
    <p v-if="galat || error" role="alert" class="text-sm font-medium text-danger-ink" data-galat>{{ galat || error }}</p>
  </div>
</template>

import http from '@/services/http'

// Video hero Beranda (30 Okt 2026): server sungguhan usulan maks 20 MB. Mode mock menyimpan file sebagai data URL
// di localStorage (kuota ±5 MB), jadi dibatasi 3 MB supaya data contoh lain tetap tersimpan.
const MAKS_VIDEO_MB = import.meta.env.VITE_USE_MOCK === 'true' ? 3 : 20

// Aturan yang sama dicek di browser (pesan cepat) dan di backend (penentu akhir).
export const ATURAN_UNGGAH = {
  video: { tipe: ['video/mp4', 'video/webm'], maks: MAKS_VIDEO_MB * 1024 * 1024, accept: 'video/mp4,video/webm', pesan: `harus MP4 atau WebM, maksimal ${MAKS_VIDEO_MB} MB` },
  gambar: { tipe: ['image/jpeg', 'image/png', 'image/webp'], maks: 2 * 1024 * 1024, accept: 'image/jpeg,image/png,image/webp', pesan: 'harus JPG, PNG, atau WebP, maksimal 2 MB' },
  dokumen: { tipe: ['application/pdf'], maks: 10 * 1024 * 1024, accept: 'application/pdf', pesan: 'harus PDF, maksimal 10 MB' },
}

/** Pesan galat bila file tidak sesuai aturan, atau '' bila boleh diunggah. */
export function cekFile(file, jenis = 'gambar') {
  const aturan = ATURAN_UNGGAH[jenis]
  return aturan.tipe.includes(file.type) && file.size <= aturan.maks ? '' : `${file.name} ${aturan.pesan}`
}

/** POST /admin/uploads (multipart). Hasil: { url, name, size, type }. */
export async function uploadFile(file, jenis = 'gambar') {
  const form = new FormData()
  form.append('file', file)
  form.append('jenis', jenis)
  const { data } = await http.post('/admin/uploads', form, { tanpaToast: true })
  const image = data.data
  if (image?.url?.startsWith('/')) image.url = new URL(image.url, import.meta.env.VITE_API_BASE_URL).href
  return image
}
import { onScopeDispose, toValue, watchEffect } from 'vue'
import { useRoute } from 'vue-router'

// Alamat situs produksi untuk canonical, og:url, dan data terstruktur. Kosong (development) → alamat yang sedang dibuka.
const SITUS = (import.meta.env.VITE_SITE_URL || window.location.origin).replace(/\/+$/, '')
let nomorPemilik = 0

/** URL absolut di situs: urlSitus('/brand') → https://…/brand. URL yang sudah absolut dikembalikan apa adanya. */
export const urlSitus = (path = '/') => (/^https?:\/\//.test(path) ? path : `${SITUS}${path.startsWith('/') ? '' : '/'}${path}`)

/** Data terstruktur BreadcrumbList dari [{ label, to }] (urutan sama dengan Breadcrumb di halaman). */
export const remahSeo = (items) => ({
  '@type': 'BreadcrumbList',
  itemListElement: items.map((item, i) => ({ '@type': 'ListItem', position: i + 1, name: item.label, item: urlSitus(item.to) })),
})

// Tag di <head> yang diatur halaman. Setiap tag ditandai pemiliknya (data-seo), supaya halaman yang ditutup hanya
// menghapus tag miliknya sendiri, bukan tag halaman baru yang sudah terpasang lebih dulu.
const TAG = {
  canonical: ['link[rel="canonical"]', () => Object.assign(document.createElement('link'), { rel: 'canonical' }), 'href'],
  robots: ['meta[name="robots"]', () => meta('name', 'robots'), 'content'],
  ogTitle: ['meta[property="og:title"]', () => meta('property', 'og:title'), 'content'],
  ogDescription: ['meta[property="og:description"]', () => meta('property', 'og:description'), 'content'],
  ogUrl: ['meta[property="og:url"]', () => meta('property', 'og:url'), 'content'],
  ogImage: ['meta[property="og:image"]', () => meta('property', 'og:image'), 'content'],
  jsonLd: ['script[type="application/ld+json"][data-seo]', () => Object.assign(document.createElement('script'), { type: 'application/ld+json' }), 'text'],
}

function meta(atribut, nama) {
  const el = document.createElement('meta')
  el.setAttribute(atribut, nama)
  return el
}

function pasang(pemilik, kunci, isi) {
  const [selektor, buat, kolom] = TAG[kunci]
  let el = document.head.querySelector(selektor)
  if (isi == null || isi === '') {
    if (el?.dataset.seo === pemilik) el.remove()
    return
  }
  if (!el) {
    el = buat()
    document.head.append(el)
  }
  el.dataset.seo = pemilik
  if (kolom === 'text') el.textContent = isi
  else el.setAttribute(kolom, isi)
}

/**
 * SEO per halaman: title dan meta description (pakai meta_title / meta_description jika ada), plus opsional
 * canonical (path), robots (mis. 'noindex, follow' untuk hasil cari/filter), image, dan jsonLd (objek/array schema.org).
 * Semua menerima nilai biasa, ref, atau getter: useSeo({ title: () => produk.value?.name }).
 * Catatan SPA: Google membaca tag ini setelah JavaScript berjalan; pratinjau link di WhatsApp/Facebook tetap
 * memakai meta bawaan di index.html sampai ada prerender/SSR.
 */
export function useSeo({ title, description, canonical, robots, image, jsonLd } = {}) {
  const route = useRoute()
  // Halaman yang sedang ditinggalkan tidak boleh menimpa SEO halaman tujuan.
  const rekam = route?.matched.at(-1)
  const pemilik = `seo-${++nomorPemilik}`

  watchEffect(() => {
    if (route && route.matched.at(-1) !== rekam) return
    const judul = toValue(title)
    const deskripsi = toValue(description)
    if (judul) document.title = `${judul} | Transhome`
    if (deskripsi) document.querySelector('meta[name="description"]')?.setAttribute('content', deskripsi)

    const kanon = toValue(canonical)
    const gambar = toValue(image)
    const data = toValue(jsonLd)
    pasang(pemilik, 'ogTitle', judul ? `${judul} | Transhome` : null)
    pasang(pemilik, 'ogDescription', deskripsi)
    pasang(pemilik, 'canonical', kanon ? urlSitus(kanon) : null)
    pasang(pemilik, 'ogUrl', kanon ? urlSitus(kanon) : null)
    pasang(pemilik, 'robots', toValue(robots))
    // Data URL (mode mock) tidak bisa dipakai crawler.
    pasang(pemilik, 'ogImage', gambar && /^(https?:\/\/|\/)/.test(gambar) ? urlSitus(gambar) : null)
    // "<" di-escape supaya isi dari admin tidak bisa menutup tag <script>.
    pasang(pemilik, 'jsonLd', data ? JSON.stringify({ '@context': 'https://schema.org', '@graph': [data].flat() }).replace(/</g, '\\u003c') : null)
  })

  onScopeDispose(() => {
    for (const kunci of Object.keys(TAG)) pasang(pemilik, kunci, null)
  })
}

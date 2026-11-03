// Artikel, promo, event, dan halaman (tabel articles). Dipakai halaman publik dan editor admin.

export const TIPE_ARTIKEL = {
  ARTIKEL: { label: 'Artikel', kelas: 'bg-subtle text-foreground' },
  PROMO: { label: 'Promo', kelas: 'bg-primary-soft text-primary-hover' },
  EVENT: { label: 'Event', kelas: 'bg-info-soft text-info-ink' },
  HALAMAN: { label: 'Halaman', kelas: 'bg-subtle text-muted' },
}

/**
 * Isi artikel → blok untuk ditampilkan (tanpa HTML mentah, jadi aman dari skrip sisipan).
 * Paragraf dipisah baris kosong; baris berawalan "## " jadi subjudul; blok yang semua barisnya "- " jadi daftar.
 * Subjudul yang langsung diikuti baris lain tanpa baris kosong (salah ketik umum di CMS) tetap dipisah:
 * hanya baris "## " itu yang jadi subjudul, sisanya diproses sebagai blok berikutnya.
 */
export function blokIsi(isi) {
  const satuBlok = (blok) => {
    const baris = blok.split('\n').map((b) => b.trim())
    if (baris[0].startsWith('## ')) {
      const subjudul = { jenis: 'h2', teks: baris[0].slice(3).replace(/\s+/g, ' ').trim() }
      const sisa = baris.slice(1).join('\n').trim()
      return sisa ? [subjudul, ...satuBlok(sisa)] : [subjudul]
    }
    if (baris.every((b) => b.startsWith('- '))) return [{ jenis: 'ul', butir: baris.map((b) => b.slice(2).trim()) }]
    return [{ jenis: 'p', teks: baris.join(' ') }]
  }
  return String(isi ?? '')
    .replace(/\r\n/g, '\n')
    .split(/\n\s*\n/)
    .map((blok) => blok.trim())
    .filter(Boolean)
    .flatMap(satuBlok)
}

/** Alamat publik artikel; token pratinjau hanya untuk draf (dibuka dari panel admin). */
export const tautanArtikel = (slug, token = null) => `/artikel/${slug}${token ? `?preview=${encodeURIComponent(token)}` : ''}`

/** Halaman daftar asal sebuah artikel: promo & event di /promo, sisanya di /artikel. */
export const daftarAsal = (tipe) => (tipe === 'PROMO' || tipe === 'EVENT' ? { label: 'Promo dan event', to: '/promo' } : { label: 'Artikel', to: '/artikel' })

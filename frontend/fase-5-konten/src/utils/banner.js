// Banner Beranda (tabel banners). Status resmi dihitung backend dari is_active + periode (waktu WIB).

// `rasio`: pratinjau gambar di form, sesuai ukuran yang disarankan.
// HOME_HERO: hero paling atas Beranda, gambar atau video (pengganti ilustrasi rumah 3D, 30 Okt 2026).
export const POSISI_BANNER = {
  HOME_HERO: { label: 'Hero Beranda', ukuran: '1920 × 800 px', rasio: 'aspect-12/5' },
  HOME_SLIDER: { label: 'Slider utama', ukuran: '1600 × 640 px', rasio: 'aspect-5/2' },
  HOME_SIDE: { label: 'Samping slider', ukuran: '800 × 520 px', rasio: 'aspect-20/13' },
}

const tanggalWib = new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Jakarta', year: 'numeric', month: '2-digit', day: '2-digit' })

/** Tanggal hari ini di WIB, "YYYY-MM-DD" (sama dengan format input tanggal). */
export const hariIniWib = (kini = new Date()) => tanggalWib.format(kini)

/** "2026-10-07T00:00:00" → "2026-10-07"; kosong → ''. */
export const tanggalSaja = (waktu) => (waktu ? String(waktu).slice(0, 10) : '')

/**
 * Perkiraan status untuk isian form (sebelum disimpan): mulai pukul 00.00, selesai pukul 23.59 WIB.
 * Sama dengan aturan backend, jadi admin langsung tahu apakah banner akan tayang.
 */
export function perkiraanStatus({ is_active, start_date, end_date }, hariIni = hariIniWib()) {
  if (!is_active) return 'NONAKTIF'
  if (start_date && start_date > hariIni) return 'TERJADWAL'
  if (end_date && end_date < hariIni) return 'BERAKHIR'
  return 'TAYANG'
}

/** Tautan banner: halaman Transhome diawali "/", situs lain harus https://. */
export const tautanBannerValid = (link) => !link || link.startsWith('/') || link.startsWith('https://')

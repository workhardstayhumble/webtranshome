// Data contoh banner dan artikel. Kolom mengikuti docs/transhome_postgres.sql. TODO: ganti ke API.
// Gambar berupa SVG contoh berlabel (bukan foto asli). Tanggal dibuat relatif terhadap hari ini (WIB),
// supaya status banner (Tayang, Terjadwal, Berakhir) dan "banner hampir berakhir" di dashboard selalu bisa dicoba.

const tanggalRelatif = (hari) => new Date(Date.now() + 7 * 3600 * 1000 + hari * 86400 * 1000).toISOString().slice(0, 10)

// Teks di dalam SVG harus di-escape: "&" mentah membuat gambar rusak.
const xml = (teks) => String(teks).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

// `judul` kosong: banner (judulnya sudah ditampilkan di atas gambar oleh slider).
function gambarContoh(judul, keterangan, latar, { lebar = 1600, tinggi = 640 } = {}) {
  const svg =
    `<svg xmlns="http://www.w3.org/2000/svg" width="${lebar}" height="${tinggi}" viewBox="0 0 ${lebar} ${tinggi}">` +
    `<rect width="${lebar}" height="${tinggi}" fill="${latar}"/>` +
    (judul ? `<text x="${lebar / 2}" y="${tinggi / 2 - 10}" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 9)}" font-weight="700" fill="#1C2B3A">${xml(judul)}</text>` : '') +
    `<text x="${lebar / 2}" y="${tinggi / 2 + tinggi / 9}" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 16)}" fill="#57534E">${xml(keterangan)}</text>` +
    `<text x="${lebar - 40}" y="${tinggi - 32}" text-anchor="end" font-family="sans-serif" font-size="${Math.round(tinggi / 26)}" fill="#57534E">CONTOH GAMBAR</text>` +
    `</svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

const UKURAN_CONTOH = { HOME_SIDE: { lebar: 800, tinggi: 520 }, HOME_HERO: { lebar: 1920, tinggi: 800 } }

// subtitle, video_url, cta_label, cta2_label, cta2_url hanya dipakai banner HOME_HERO (hero Beranda, 30 Okt 2026).
const banner = (id, title, latar, link_url, position, mulai, selesai, sort_order, is_active = true, hero = {}) => ({
  id, title, image_url: gambarContoh('', position === 'HOME_HERO' ? 'Foto / video hero' : 'Gambar banner promo', latar, UKURAN_CONTOH[position] ?? {}),
  link_url, position, start_at: `${mulai}T00:00:00`, end_at: selesai ? `${selesai}T23:59:59` : null, sort_order, is_active,
  subtitle: hero.subtitle ?? null, video_url: hero.video_url ?? null, cta_label: hero.cta_label ?? null,
  cta2_label: hero.cta2_label ?? null, cta2_url: hero.cta2_url ?? null,
  created_at: '2026-09-01T09:00:00', updated_at: '2026-09-01T09:00:00',
})

// position: HOME_HERO (hero paling atas Beranda: gambar atau video), HOME_SLIDER (slider utama),
// HOME_SIDE (kartu di samping slider).
export const banners = [
  banner(8, 'Dari pondasi sampai atap, semua kebutuhan rumah ada di sini.', '#3A4A5A', '/katalog', 'HOME_HERO', tanggalRelatif(-30), null, 1, true, {
    subtitle: 'Cari bahan untuk tiap tahap rumahmu, cek harga dan stok.', cta_label: 'Jelajahi katalog', cta2_label: 'Lihat promo', cta2_url: '/promo',
  }),
  banner(1, 'Promo keramik lantai', '#E7D7C1', '/katalog/keramik-lantai', 'HOME_SLIDER', tanggalRelatif(-6), tanggalRelatif(3), 1),
  banner(2, 'Paket atap baja ringan', '#F3D9C4', '/katalog/atap', 'HOME_SLIDER', tanggalRelatif(-10), tanggalRelatif(30), 2),
  banner(3, 'Lampu LED hemat energi', '#FDF0C2', '/katalog/lampu', 'HOME_SIDE', tanggalRelatif(-3), tanggalRelatif(20), 1),
  banner(4, 'Pekan cat tembok', '#DCE3D2', '/katalog/cat-tembok', 'HOME_SLIDER', tanggalRelatif(20), tanggalRelatif(34), 3),
  banner(5, 'Promo genteng metal', '#E9DFCE', '/katalog/genteng', 'HOME_SLIDER', tanggalRelatif(-14), tanggalRelatif(-1), 4),
  banner(6, 'Promo sanitasi September', '#D4E2EA', '/katalog/sanitasi', 'HOME_SLIDER', '2026-09-01', '2026-09-30', 5),
  banner(7, 'Promo kran dapur', '#D6D3CD', '/katalog/kran', 'HOME_SLIDER', tanggalRelatif(-2), tanggalRelatif(12), 6, false),
]

// Isi artikel: paragraf dipisah baris kosong, "## " untuk subjudul, "- " untuk daftar (lihat utils/artikel.js).
// Isi boleh ditulis menjorok (template literal); indentasi dibuang supaya rapi di editor admin.
// `labelGambar`: teks pendek untuk gambar sampul contoh (judul panjang tidak muat di gambar 800 px).
const rapikan = (teks) => teks.split('\n').map((baris) => baris.trim()).join('\n').trim()
const artikel = (id, type, title, slug, status, hari, latar, content, preview_token, labelGambar = title) => ({
  id, type, title, slug, status, content: rapikan(content), preview_token,
  thumbnail_url: latar ? gambarContoh(labelGambar, type === 'ARTIKEL' ? 'Panduan isi rumah' : 'Promo & event', latar, { lebar: 800, tinggi: 450 }) : null,
  published_at: status === 'PUBLISHED' ? `${tanggalRelatif(hari)}T08:00:00` : null, author_id: 4,
  created_at: `${tanggalRelatif(hari - 1)}T15:00:00`, updated_at: `${tanggalRelatif(hari)}T08:00:00`,
})

export const articles = [
  artikel(1, 'PROMO', 'Promo keramik lantai bulan ini', 'promo-keramik-lantai', 'PUBLISHED', -6, '#E7D7C1',
    'Contoh isi promo. Harga dan syarat promo diisi Admin Konten.\n\n## Berlaku untuk\n- Keramik lantai pilihan\n- Pembelian langsung di toko\n\nTanyakan stok dan harga promo ke tim kami lewat WhatsApp sebelum datang ke toko.',
    'a41f0c9e7b2d4c18'),
  artikel(2, 'EVENT', 'Demo pemasangan bata ringan di toko', 'demo-pemasangan-bata-ringan', 'PUBLISHED', -12, '#DCE3D2',
    'Contoh isi event. Tanggal, jam, dan tempat diisi Admin Konten.\n\n## Yang akan dibahas\n- Menyiapkan perekat bata ringan\n- Memasang dan meratakan susunan bata\n- Menghitung kebutuhan untuk satu dinding\n\nAcara terbuka untuk umum. Tanyakan jadwalnya lewat WhatsApp.',
    'b53e1dae8c3f5d29'),
  artikel(3, 'PROMO', 'Paket atap baja ringan', 'paket-atap-baja-ringan', 'PUBLISHED', -2, '#F3D9C4',
    'Contoh isi promo paket. Isi paket dan harga diisi Admin Konten.\n\nKirim ukuran atap rumahmu lewat WhatsApp, tim kami bantu hitung kebutuhan rangka dan genteng.',
    'c62f2ebf9d406e3a'),
  // Seri panduan "5 Langkah Mudah Mengisi Rumah" (materi pemilik, 2 Nov 2026): 1 artikel utama + 5 artikel pilar.
  // Pilar 5 sengaja DRAF supaya alur pratinjau & terbitkan di admin tetap bisa dicoba.
  artikel(4, 'ARTIKEL', '5 Langkah Mudah Mengisi Rumah: Dari Lantai Granit, Kamar Mandi, hingga Sofa', '5-langkah-mudah-mengisi-rumah', 'PUBLISHED', -1, '#E7D7C1', `
    Mengisi rumah baru atau merenovasi rumah lama sering bikin bingung: mulai dari mana, dan bagaimana supaya hasilnya serasi? Urutan paling aman adalah dari bawah ke atas, dari yang menempel permanen ke yang bisa dipindah: lantai dulu, lalu dinding, kamar mandi, lampu dan elektronik, baru terakhir sofa dan perabot.

    Urutan ini juga menghemat biaya. Lantai dan dinding paling repot diganti, jadi diputuskan lebih dulu. Warna lantai dan dinding kemudian jadi patokan saat memilih sofa dan perabot lain, sehingga semuanya serasi.

    ## Langkah 1: Pilih lantai yang lega dan aman di kaki

    Granit atau keramik ukuran besar membuat ruangan terlihat lebih luas dan lebih cepat disapu karena garis natnya sedikit. Untuk teras, dapur, dan kamar mandi, pilih keramik kesat (doff atau bertekstur) supaya tidak licin bagi anak kecil dan orang tua. Panduan lengkapnya ada di artikel "Memilih Keramik & Granit Lantai yang Pas dan Aman di Kaki".

    ## Langkah 2: Percantik dinding tanpa membuatnya rapuh

    Pakai cat yang sesuai tempatnya: cat eksterior untuk dinding luar yang kena panas dan hujan, keramik dinding untuk area basah supaya tidak berjamur. Saat memasang cermin atau rak, pakai mata bor dan sekrup yang tepat agar tembok dan keramik tidak retak. Panduan lengkapnya ada di artikel "Mempercantik Dinding dengan Cat, Keramik, dan Hiasan Tempel".

    ## Langkah 3: Lengkapi kamar mandi dengan perlengkapan yang awet

    Ukur kamar mandi sebelum memilih kloset dan wastafel, lalu pilih kloset dua tombol siram yang hemat air. Kran dan shower yang tahan karat dan mudah dibersihkan dari kerak akan jauh lebih awet. Panduan lengkapnya ada di artikel "Perlengkapan Kamar Mandi (Sanitary) yang Awet & Mudah Dibersihkan".

    ## Langkah 4: Atur lampu dan elektronik supaya terang dan hemat

    Lampu putih terang cocok untuk dapur dan meja kerja, lampu kuning hangat untuk kamar tidur. Sesuaikan ukuran AC, TV, dan pemanas air dengan luas ruangan dan jumlah penghuni supaya listrik tidak boros. Panduan lengkapnya ada di artikel "Menata Lampu dan Elektronik agar Rumah Terang & Hemat Listrik".

    ## Langkah 5: Pilih sofa dan kursi paling akhir

    Setelah lantai dan dinding jadi, ukur ruang tamu sebelum membeli sofa supaya jalan di dalam rumah tetap lega. Pilih bahan yang mudah dilap dan warnanya serasi dengan lantai.

    ## Ringkasan daftar belanja

    - Lantai: granit ukuran besar untuk ruang tamu, keramik kesat untuk area basah.
    - Dinding: cat sesuai area, keramik untuk area basah, mata bor yang tepat untuk memasang barang.
    - Kamar mandi: ukur ruangan dulu, pilih kloset hemat air dan kran yang mudah dibersihkan.
    - Lampu & elektronik: putih untuk bekerja, kuning untuk beristirahat, ukuran AC dan TV sesuai ruangan.
    - Sofa: ukur ruangan dan pintu, pilih bahan yang mudah dirawat.

    Belum yakin harus mulai dari mana? Kirim foto dan ukuran ruanganmu ke tim kami lewat WhatsApp, kami bantu susun daftar belanjanya.
  `, 'd7304fc0ae517f4b', '5 langkah mengisi rumah'),

  artikel(5, 'ARTIKEL', 'Memilih Keramik & Granit Lantai yang Pas dan Aman di Kaki', 'memilih-keramik-granit-lantai', 'PUBLISHED', -2, '#DCE3D2', `
    Lantai dipakai setiap hari dan paling repot diganti. Karena itu, pilih dengan menjawab tiga pertanyaan sederhana: ruangannya untuk apa, seberapa luas, dan siapa yang paling sering lewat.

    ## Ukuran besar, ruangan terasa lega

    Granit atau keramik ukuran besar, misalnya 60 × 60 cm atau 60 × 120 cm, membuat ruangan terlihat lebih luas karena garis natnya lebih sedikit. Garis nat yang sedikit juga berarti lebih sedikit celah tempat debu menempel, jadi lantai lebih cepat disapu dan dipel.

    - Ruang tamu dan ruang keluarga: 60 × 60 cm atau lebih besar.
    - Kamar tidur kecil: 50 × 50 cm atau 60 × 60 cm tetap terlihat rapi.
    - Kamar mandi: ukuran kecil seperti 20 × 20 cm atau 30 × 30 cm lebih mudah dibuat miring ke arah lubang pembuangan.

    Warna terang seperti krem, putih tulang, atau abu muda membuat ruangan tampak lebih luas. Motif kayu memberi kesan hangat tanpa perlu merawat kayu asli.

    ## Mengkilap atau kesat? Sesuaikan dengan ruangannya

    Permukaan mengkilap terlihat mewah dan mudah dilap, cocok untuk ruang tamu dan kamar yang selalu kering. Untuk area yang sering basah, pilih permukaan kesat (doff, matte, atau bertekstur) supaya tidak licin.

    - Teras dan carport: kesat atau bertekstur, tahan panas dan hujan.
    - Dapur: doff atau semi-kesat, karena sering terkena cipratan air dan minyak.
    - Kamar mandi: wajib kesat.

    Kalau di rumah ada anak kecil atau orang tua, keamanan lebih penting daripada kilap. Di toko, coba usap permukaan keramik dengan telapak tangan yang sedikit basah: keramik yang kesat terasa menahan, tidak licin.

    ## Sebelum membeli

    - Ukur panjang dan lebar lantai, lalu tambahkan cadangan sekitar 5–10% untuk potongan dan keramik yang pecah.
    - Beli sekaligus dalam satu kali pembelian supaya warna dan kode produksinya seragam.
    - Simpan beberapa keping sisa untuk perbaikan di kemudian hari.

    Kirim ukuran lantai rumahmu ke tim kami lewat WhatsApp, kami bantu hitung jumlah dus yang dibutuhkan.
  `, 'e8415ad1bf62805c', 'Lantai & granit'),

  artikel(6, 'ARTIKEL', 'Mempercantik Dinding dengan Cat, Keramik, dan Hiasan Tempel', 'mempercantik-dinding-cat-keramik-hiasan', 'PUBLISHED', -3, '#FDF0C2', `
    Dinding yang polos bisa dibuat lebih cantik tanpa mengorbankan kekuatannya. Kuncinya ada dua: pilih pelapis dinding yang sesuai tempatnya, lalu pasang barang tempel dengan cara yang benar.

    ## Pilih pelindung dinding sesuai tempatnya

    - Dinding luar: pakai cat eksterior yang tertulis tahan cuaca, karena setiap hari terkena panas matahari dan hujan. Warna terang membantu rumah tidak cepat panas.
    - Dinding dalam: cat interior sudah cukup. Untuk dapur dan kamar anak, pilih cat yang bisa dilap supaya noda mudah dibersihkan.
    - Kamar mandi dan dinding dapur dekat kompor: pasang keramik dinding di area yang sering terkena cipratan, supaya tembok tidak lembap dan berjamur. Pastikan nat terisi rapat; nat anti jamur membuatnya lebih awet.
    - Hiasan dinding: pelapis dinding motif kayu (wall panel) memberi kesan hangat dan rapi tanpa perawatan kayu asli. Untuk area yang lembap, pilih panel yang tahan air.

    Supaya cat awet dan warnanya rata, pastikan dinding sudah kering dan bersih, gunakan cat dasar lebih dulu, lalu cat warna dua lapis. Kalau ada rembesan air, perbaiki dulu sumbernya; cat atau keramik baru tidak akan bertahan di dinding yang basah.

    ## Memasang cermin, rak, dan kabinet dengan aman

    - Cek dulu jalur pipa dan kabel listrik di dalam tembok sebelum mengebor. Biasanya jalurnya lurus ke atas atau ke samping dari stopkontak, sakelar, dan kran.
    - Tempel lakban kertas di titik yang akan dibor pada keramik, supaya mata bor tidak meleset dan permukaan tidak tergores.
    - Pakai mata bor khusus keramik saat menembus keramik dan matikan mode palu (hammer) agar keramik tidak retak. Mode palu boleh dinyalakan lagi setelah mata bor sampai ke bata atau beton.
    - Sesuaikan fischer dan sekrup dengan berat barang. Kabinet gantung dan cermin besar butuh baut atau braket yang lebih kuat.
    - Untuk dinding gipsum, pakai jangkar khusus gipsum, atau kencangkan sekrup ke rangkanya untuk barang yang berat.

    Ragu memilih cat atau alat pasang yang cocok? Tanyakan ke tim kami lewat WhatsApp, sertakan foto dindingnya.
  `, 'f9526be2c073916d', 'Dinding'),

  artikel(7, 'ARTIKEL', 'Perlengkapan Kamar Mandi (Sanitary) yang Awet & Mudah Dibersihkan', 'perlengkapan-kamar-mandi-awet', 'PUBLISHED', -4, '#D4E2EA', `
    Perlengkapan kamar mandi, atau sering disebut sanitary, adalah semua barang yang berhubungan dengan air: kloset, wastafel, kran, dan shower. Karena dipakai setiap hari, pilih yang nyaman, hemat air, dan mudah dibersihkan.

    ## Nyaman dipakai dan hemat air

    - Ukur kamar mandi lebih dulu. Sisakan ruang bebas sekitar 60 cm di depan kloset supaya leluasa saat duduk dan berdiri.
    - Pilih kloset duduk dengan dua tombol siram: tombol kecil untuk buang air kecil, tombol besar untuk buang air besar. Pemakaian air jadi lebih hemat.
    - Kloset monoblok (tangki menyatu dengan dudukan) lebih mudah dibersihkan karena sambungannya sedikit.
    - Cocokkan jarak lubang pembuangan ke dinding dengan ukuran kloset sebelum membeli, supaya tidak perlu membongkar lantai.
    - Wastafel kecil atau wastafel sudut cocok untuk kamar mandi sempit. Tinggi pemasangannya umumnya sekitar 80–85 cm dari lantai.
    - Untuk orang tua, pertimbangkan pegangan di dekat kloset dan area mandi.

    ## Bebas karat, kerak, dan mampet

    - Pilih kran dan shower berbahan stainless steel atau kuningan berlapis krom. Bahan ini lebih tahan karat daripada besi biasa.
    - Air sumur yang berkapur membuat kerak cepat menumpuk di mulut kran dan kepala shower. Kepala shower dengan lubang karet lebih mudah dibersihkan, cukup diusap dengan jari.
    - Bersihkan mulut kran dan kepala shower secara rutin. Rendam di air cuka untuk melarutkan kerak yang membandel.
    - Pasang saringan rambut di lubang pembuangan lantai dan hindari membuang minyak ke saluran air supaya tidak mampet.
    - Kran yang menetes biasanya karena karet di dalamnya aus. Ganti segera supaya air tidak terbuang.

    Butuh bantuan memilih kloset atau kran yang pas dengan ukuran kamar mandimu? Kirim ukurannya ke tim kami lewat WhatsApp.
  `, '0a637cf3d184a27e', 'Kamar mandi'),

  artikel(8, 'ARTIKEL', 'Menata Lampu dan Elektronik agar Rumah Terang & Hemat Listrik', 'lampu-elektronik-terang-hemat-listrik', 'PUBLISHED', -5, '#E9DFCE', `
    Lampu dan elektronik yang tepat membuat rumah nyaman di mata sekaligus ramah di tagihan listrik bulanan. Kuncinya: warna lampu sesuai kegiatan, dan ukuran perangkat sesuai luas ruangan.

    ## Warna lampu yang tepat

    - Putih terang (cool daylight, sekitar 6500K): dapur, meja kerja, kamar mandi, dan garasi, supaya semuanya terlihat jelas.
    - Putih netral (sekitar 4000K): ruang keluarga yang dipakai berkumpul sekaligus beraktivitas.
    - Kuning hangat (warm white, sekitar 3000K): kamar tidur dan ruang santai, supaya mata tidak cepat lelah dan suasana lebih tenang.

    Pakai lampu LED karena lebih hemat listrik dan lebih awet daripada lampu pijar. Untuk tahu seberapa terang, lihat angka lumen di kemasan, bukan hanya angka watt.

    ## Ukuran elektronik yang pas

    AC yang terlalu kecil bekerja terus-menerus dan boros listrik, sedangkan yang terlalu besar membuang daya. Patokan kasar untuk ruangan dengan tinggi plafon sekitar 3 meter yang tidak terkena matahari langsung:

    - ½ PK: sampai sekitar 12 m².
    - ¾ PK: sekitar 12–16 m².
    - 1 PK: sekitar 16–24 m².
    - 1½ PK: sekitar 24–32 m².
    - 2 PK: sekitar 32–40 m².

    Ruangan yang panas, terkena matahari sore, atau berisi banyak orang biasanya butuh satu tingkat lebih besar.

    Untuk TV, jarak duduk yang nyaman sekitar 1,5–2,5 kali ukuran layar. Contohnya, TV 43 inci (sekitar 109 cm) nyaman ditonton dari jarak sekitar 1,6–2,7 meter.

    Untuk pemanas air listrik, kapasitas sekitar 15 liter umumnya cukup untuk 1–2 orang mandi bergantian, dan sekitar 30 liter untuk keluarga 3–4 orang. Cek juga daya listrik rumah sebelum menambah AC atau pemanas air, supaya listrik tidak sering turun.

    ## Kebiasaan kecil yang menghemat listrik

    - Atur suhu AC di sekitar 24–25°C dan tutup pintu serta jendela saat AC menyala.
    - Cabut charger dan perangkat yang tidak dipakai dari stopkontak.
    - Manfaatkan cahaya matahari di siang hari.

    Bingung memilih ukuran AC atau pemanas air? Kirim ukuran ruangan dan jumlah penghuni ke tim kami lewat WhatsApp.
  `, 'b1748e02c5f93a6d', 'Lampu & elektronik'),

  artikel(9, 'ARTIKEL', 'Memilih Sofa dan Kursi yang Empuk Tanpa Membuat Ruangan Sempit', 'memilih-sofa-kursi-ruangan-lega', 'DRAFT', -1, '#F3D9C4', `
    Sofa adalah perabot yang paling sering dipakai bersantai bersama keluarga, tapi juga yang paling sering membuat ruangan terasa sempit. Pilih sofa setelah lantai dan dinding selesai, supaya warnanya bisa diserasikan.

    ## Ukur sebelum beli

    - Ukur panjang dan lebar ruang tamu, lalu tandai area sofa di lantai dengan lakban kertas untuk membayangkan ukurannya.
    - Sisakan jalur jalan sekitar 75–90 cm, dan pastikan sofa tidak menghalangi pintu maupun jendela.
    - Beri jarak sekitar 40–45 cm antara sofa dan meja tamu, supaya kaki leluasa dan meja tetap mudah dijangkau.
    - Ukur juga lebar pintu, lorong, dan tangga, supaya sofa bisa masuk ke dalam rumah.
    - Untuk ruangan kecil, pilih sofa berkaki ramping dengan kolong yang terlihat, sandaran rendah, atau sofa bentuk L di sudut ruangan.

    ## Bahan yang mudah dirawat dan serasi

    - Kain: terasa adem dan pilihan warnanya banyak. Pilih anyaman yang rapat; sarung yang bisa dilepas lebih mudah dicuci.
    - Kulit atau kulit sintetis: mudah dilap bila terkena tumpahan makanan atau minuman, cocok untuk keluarga dengan anak kecil. Jauhkan dari sinar matahari langsung supaya tidak cepat pudar.
    - Warna: lantai granit terang cocok dengan sofa abu-abu, cokelat, atau warna tanah. Pilih warna sofa yang netral, lalu tambahkan bantal berwarna sebagai aksen.
    - Kenyamanan: coba duduk beberapa menit di toko untuk merasakan empuk tidaknya busa dan tinggi sandarannya.

    Mau dibantu memilih sofa yang pas dengan ukuran ruang tamumu? Kirim ukuran dan foto ruangannya ke tim kami lewat WhatsApp.
  `, 'c2859f13d6a04b7e', 'Sofa & kursi'),
]

// Jasa Transhome (KEPUTUSAN #20, 2 Nov 2026). Usulan tabel `services` + `service_categories` (lihat
// docs/FRONTEND_BUILD_PLAN.md). DATA CONTOH: nama, area, dan isi diubah Admin Konten di /admin/jasa.
// Harga jasa tidak ditampilkan: semua jasa memakai label tetap "Harga pasti murah" (utils/jasa.js).
// category_ids = kategori produk yang menampilkan jasa ini di halaman produk (kategori induk berlaku untuk semua anaknya).
const jasa = (id, name, slug, summary, latar, labelGambar, category_ids, description) => ({
  id, name, slug, summary, description: rapikan(description),
  image_url: gambarContoh(labelGambar, 'Jasa Transhome', latar, { lebar: 800, tinggi: 450 }),
  service_area: 'Sleman dan sekitarnya', category_ids, sort_order: id, is_active: true,
  created_at: '2026-11-02T09:00:00', updated_at: '2026-11-02T09:00:00',
})

export const services = [
  jasa(1, 'Jasa pasang keramik & granit', 'pasang-keramik-granit', 'Pemasangan keramik dan granit untuk lantai maupun dinding, termasuk perataan dan pengisian nat.', '#E7D7C1', 'Pasang keramik', [2], `
    Lantai dan dinding yang rata dan rapi bergantung pada cara pemasangannya. Tim kami membantu memasang keramik atau granit yang kamu beli di Transhome.

    ## Yang dikerjakan

    - Mengecek kerataan lantai atau dinding sebelum dipasang.
    - Memasang keramik atau granit dengan perekat yang sesuai.
    - Mengisi nat dan membersihkan sisa pemasangan.

    ## Siapkan sebelum konsultasi

    - Ukuran ruangan (panjang × lebar) atau luas dinding.
    - Foto kondisi lantai atau dinding sekarang.
    - Jenis dan ukuran keramik yang dipilih.
  `),
  jasa(2, 'Jasa pengecatan', 'pengecatan', 'Pengecatan dinding dalam dan luar rumah, dari persiapan permukaan sampai lapisan akhir.', '#FDF0C2', 'Pengecatan', [5], `
    Cat yang awet dimulai dari permukaan yang disiapkan dengan benar. Tim kami membantu mengecat dinding dalam maupun luar rumah.

    ## Yang dikerjakan

    - Membersihkan dan menambal bagian dinding yang retak ringan.
    - Mengoleskan cat dasar, lalu cat warna.
    - Melindungi lantai dan perabot dari cipratan cat.

    ## Siapkan sebelum konsultasi

    - Luas atau ukuran dinding yang akan dicat.
    - Foto dinding sekarang, terutama bagian yang mengelupas atau lembap.
    - Warna dan jenis cat yang diinginkan, bila sudah ada.
  `),
  jasa(3, 'Jasa pasang perlengkapan kamar mandi', 'pasang-perlengkapan-kamar-mandi', 'Pemasangan kloset, wastafel, kran, shower, dan pemanas air.', '#D4E2EA', 'Pasang sanitary', [3, 81], `
    Kloset, wastafel, dan kran yang terpasang rapat tidak mudah bocor. Tim kami membantu memasang perlengkapan kamar mandi yang kamu beli di Transhome.

    ## Yang dikerjakan

    - Memasang kloset, wastafel, kran, shower, atau pemanas air.
    - Menyambungkan ke saluran air bersih dan pembuangan.
    - Mengecek tidak ada rembesan setelah terpasang.

    ## Siapkan sebelum konsultasi

    - Foto kamar mandi dan posisi saluran air.
    - Ukuran ruangan.
    - Daftar barang yang akan dipasang.
  `),
  jasa(4, 'Jasa pasang AC & elektronik', 'pasang-ac-elektronik', 'Pemasangan AC, lampu, stop kontak, dan perangkat elektronik rumah.', '#DCE3D2', 'Pasang AC & elektronik', [6, 8], `
    Perangkat elektronik bekerja paling baik bila dipasang dengan benar. Tim kami membantu memasang AC, lampu, dan perangkat elektronik rumah.

    ## Yang dikerjakan

    - Memasang unit AC dalam dan luar ruangan.
    - Memasang lampu, sakelar, dan stop kontak.
    - Menguji perangkat setelah terpasang.

    ## Siapkan sebelum konsultasi

    - Foto ruangan dan dinding tempat perangkat akan dipasang.
    - Ukuran ruangan untuk AC.
    - Daya listrik rumah (VA), bila tahu.
  `),
  jasa(5, 'Jasa service & perbaikan', 'service-perbaikan', 'Servis dan perbaikan AC, pemanas air, pompa air, dan perangkat elektronik rumah yang bermasalah.', '#E9DFCE', 'Service & perbaikan', [8], `
    AC yang kurang dingin, pemanas air yang tidak panas, atau pompa air yang lemah bisa diperiksa dan diperbaiki tanpa harus langsung membeli baru. Tim kami membantu servis rutin maupun perbaikan.

    ## Yang dikerjakan

    - Servis rutin AC: membersihkan unit dan mengecek kinerjanya.
    - Memeriksa dan memperbaiki pemanas air, pompa air, dan perangkat elektronik rumah.
    - Memberi saran bila ada komponen yang perlu diganti.

    ## Siapkan sebelum konsultasi

    - Merek dan jenis perangkat yang bermasalah.
    - Keluhan singkat, misalnya "AC menetes" atau "air tidak panas".
    - Foto atau video perangkat saat menyala, bila ada.
  `),
  jasa(6, 'Jasa pasang pintu & jendela', 'pasang-pintu-jendela', 'Pemasangan pintu, jendela, kusen, kunci, dan engsel untuk rumah baru maupun renovasi.', '#F3D9C4', 'Pasang pintu & jendela', [7], `
    Pintu dan jendela yang terpasang lurus mudah dibuka, rapat, dan aman dikunci. Tim kami membantu memasang pintu, jendela, beserta kunci dan engselnya.

    ## Yang dikerjakan

    - Memasang kusen, daun pintu, dan jendela.
    - Memasang engsel, kunci, dan pegangan pintu.
    - Menyetel pintu dan jendela supaya rapat dan tidak seret.

    ## Siapkan sebelum konsultasi

    - Ukuran lubang pintu atau jendela (lebar × tinggi).
    - Foto lokasi pemasangan.
    - Jenis pintu atau jendela yang dipilih, bila sudah ada.
  `),
  jasa(7, 'Jasa antar & rakit perabot', 'antar-rakit-perabot', 'Pengantaran dan perakitan sofa, lemari, meja, kasur, dan perabot rumah lainnya sampai siap dipakai.', '#E7E2F0', 'Antar & rakit perabot', [], `
    Perabot besar seperti sofa, lemari, dan kasur repot dibawa dan dirakit sendiri. Tim kami mengantar perabot ke rumahmu dan merakitnya sampai siap dipakai.

    ## Yang dikerjakan

    - Mengantar perabot dari toko ke rumahmu.
    - Merakit lemari, meja, rak, dan tempat tidur.
    - Menata perabot di ruangan yang kamu tentukan.

    ## Siapkan sebelum konsultasi

    - Daftar perabot yang akan diantar atau dirakit.
    - Alamat lengkap dan patokannya.
    - Lebar pintu, lorong, atau tangga yang akan dilewati perabot besar.
  `),
]

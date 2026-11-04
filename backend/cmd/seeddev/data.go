package main

// Semua isi di sini fiktif: nama, harga, stok, member, dan transaksi tidak mewakili data Transhome.

var (
	namaBrand = []string{
		"Merek A", "Merek B", "Merek C", "Merek D", "Merek E", "Merek F",
		"Alfa Keramik", "Beta Sanitair", "Gama Cat", "Delta Pipa", "Epsilon Baja", "Zeta Lampu", "Eta Kabel", "Teta Atap",
		"Iota Kunci", "Kapa Kaca", "Lambda Granit", "Mu Engsel", "Nu Semen", "Ksi Plafon", "Omikron Pintu", "Pi Perkakas",
		"Rho Elektrik", "Sigma Kran", "Tau Bata", "Upsilon Lantai", "Fi Pompa", "Khi Lem", "Psi Panel", "Omega Tangga",
	}
	warnaLogo = []string{"#C2410C", "#1D4ED8", "#15803D", "#7C3AED", "#B91C1C", "#0F766E", "#B45309", "#1C2B3A", "#BE185D", "#4D7C0F"}
	latarFoto = []string{"#D6D3CD", "#E7DFD3", "#DCE3D2", "#D4E2EA", "#E7D7C1", "#F3D9C4", "#FDF0C2", "#E9DFCE"}
	labelFoto = []string{"Foto utama", "Detail tekstur", "Terpasang di ruangan"}

	pengaturanContoh = [][2]string{
		{"address", "Jl. Raya Solo–Yogyakarta Km 11, Berbah, Sleman"},
		{"opening_hours", "Senin–Sabtu, 08.00–17.00"},
		{"maps_url", "https://www.google.com/maps/search/?api=1&query=Transhome+Berbah+Sleman"},
		{"point_ratio_rupiah", "1000000"},
	}

	// Dipakai bila database belum punya kategori ini; kategori yang sudah ada tidak diubah.
	kategoriContoh = [][2]string{
		{"Keramik & Granit", katKeramik},
		{"Cat & Pelapis", katCat},
		{"Lampu & Pencahayaan", katLampu},
		{"Sanitasi & Kamar Mandi", katSanitasi},
		{"Elektrikal & Kelistrikan", katListrik},
	}

	ruanganContoh = []struct{ nama, latar string }{
		{"Dapur", "#E7D7C1"},
		{"Kamar mandi", "#D4E2EA"},
		{"Fasad", "#D6D3CD"},
		{"Ruang keluarga", "#E9DFCE"},
		{"Teras", "#DCE3D2"},
	}

	// Relasi jasa bawaan migration 000007 dengan kategori produk.
	jasaContoh = []struct {
		slug, label, latar string
		kategori           []string
	}{
		{"pengecatan", "Pengecatan", "#FDF0C2", []string{katCat}},
		{"pasang-keramik-granit", "Pasang keramik", "#E7D7C1", []string{katKeramik}},
		{"pasang-perlengkapan-kamar-mandi", "Pasang sanitary", "#D4E2EA", []string{katSanitasi}},
		{"pasang-ac-elektronik", "Pasang AC & elektronik", "#DCE3D2", []string{katLampu, katListrik}},
		{"service-perbaikan", "Service & perbaikan", "#E9DFCE", []string{katListrik}},
		{"pasang-pintu-jendela", "Pasang pintu & jendela", "#F3D9C4", nil},
		{"antar-rakit-perabot", "Antar & rakit perabot", "#E7E2F0", nil},
	}
)

const (
	katKeramik  = "keramik-granit"
	katCat      = "cat-pelapis"
	katLampu    = "lampu-pencahayaan"
	katSanitasi = "sanitasi-kamar-mandi"
	katListrik  = "elektrikal-kelistrikan"

	stokAda  = "TERSEDIA"
	stokSisa = "SISA_STOK"
	stokPO   = "PRE_ORDER"
	stokNol  = "HABIS"

	rDapur    = "dapur"
	rMandi    = "kamar-mandi"
	rFasad    = "fasad"
	rKeluarga = "ruang-keluarga"
	rTeras    = "teras"
)

type spek = map[string]string

type produkContoh struct {
	sku, nama, kategori, brand string
	harga                      int64
	satuan                     string
	minOrder                   int
	stok, labelStok            string
	unggulan, aktif            bool
	ruangan                    []string
	spesifikasi                spek
}

// Lebih dari 24 produk aktif (ukuran satu halaman katalog), supaya perpindahan halaman ikut teruji.
var produkDaftar = []produkContoh{
	{"CONTOH-KRM-001", "Keramik lantai abu doff 40×40", katKeramik, "alfa-keramik", 89500, "dus", 1, stokAda, "", true, true, []string{rKeluarga, rDapur}, spek{"ukuran": "40 × 40 cm", "finishing": "Doff", "warna": "Abu", "isi_per_dus": "6 pcs (0,96 m²)", "garansi": "5 tahun"}},
	{"CONTOH-KRM-002", "Keramik lantai motif kayu 40×40", katKeramik, "alfa-keramik", 96000, "dus", 1, stokSisa, "3 dus", false, true, []string{rKeluarga, rTeras}, spek{"ukuran": "40 × 40 cm", "finishing": "Doff", "warna": "Cokelat kayu", "isi_per_dus": "6 pcs (0,96 m²)"}},
	{"CONTOH-KRM-003", "Keramik lantai putih glossy 50×50", katKeramik, "upsilon-lantai", 112000, "dus", 1, stokAda, "", false, true, []string{rKeluarga}, spek{"ukuran": "50 × 50 cm", "finishing": "Glossy", "warna": "Putih", "isi_per_dus": "4 pcs (1 m²)"}},
	{"CONTOH-KRM-004", "Keramik lantai teraso 60×60", katKeramik, "upsilon-lantai", 168000, "dus", 1, stokPO, "", false, true, []string{rKeluarga, rTeras}, spek{"ukuran": "60 × 60 cm", "finishing": "Doff", "warna": "Teraso", "isi_per_dus": "4 pcs (1,44 m²)"}},
	{"CONTOH-KRM-005", "Keramik dinding putih 30×60", katKeramik, "merek-a", 104000, "dus", 1, stokAda, "", false, true, []string{rDapur, rMandi}, spek{"ukuran": "30 × 60 cm", "finishing": "Glossy", "warna": "Putih", "isi_per_dus": "8 pcs (1,44 m²)"}},
	{"CONTOH-KRM-006", "Keramik dinding motif batu 25×40", katKeramik, "merek-a", 96000, "dus", 1, stokAda, "", false, true, []string{rMandi, rFasad}, spek{"ukuran": "25 × 40 cm", "finishing": "Glossy", "warna": "Batu alam"}},
	{"CONTOH-KRM-007", "Keramik teras kesat 30×30", katKeramik, "alfa-keramik", 74000, "dus", 1, stokAda, "", true, true, []string{rTeras}, spek{"ukuran": "30 × 30 cm", "finishing": "Bertekstur", "warna": "Abu tua"}},
	{"CONTOH-KRM-008", "Keramik kamar mandi kesat 20×20", katKeramik, "upsilon-lantai", 69000, "dus", 1, stokSisa, "5 dus", false, true, []string{rMandi}, spek{"ukuran": "20 × 20 cm", "finishing": "Kesat", "warna": "Biru muda"}},
	{"CONTOH-GRN-001", "Granit lantai abu muda 60×60", katKeramik, "lambda-granit", 245000, "dus", 1, stokAda, "", true, true, []string{rKeluarga}, spek{"ukuran": "60 × 60 cm", "finishing": "Polished", "warna": "Abu muda"}},
	{"CONTOH-GRN-002", "Granit lantai hitam 60×60", katKeramik, "lambda-granit", 265000, "dus", 1, stokPO, "", false, true, []string{rKeluarga}, spek{"ukuran": "60 × 60 cm", "finishing": "Polished", "warna": "Hitam"}},
	{"CONTOH-GRN-003", "Granit lantai putih 60×120", katKeramik, "lambda-granit", 420000, "dus", 1, stokAda, "", false, true, []string{rKeluarga}, spek{"ukuran": "60 × 120 cm", "finishing": "Polished", "warna": "Putih"}},
	{"CONTOH-PRK-001", "Perekat keramik 25 kg", katKeramik, "khi-lem", 98000, "sak", 1, stokAda, "", false, true, []string{rMandi, rKeluarga}, spek{"berat": "25 kg"}},

	{"CONTOH-CAT-001", "Cat tembok interior putih 5 kg", katCat, "gama-cat", 185000, "galon", 1, stokAda, "", true, true, []string{rKeluarga}, spek{"berat": "5 kg", "finishing": "Matt", "warna": "Putih"}},
	{"CONTOH-CAT-002", "Cat tembok interior krem 5 kg", katCat, "gama-cat", 178000, "galon", 1, stokAda, "", false, true, []string{rKeluarga}, spek{"berat": "5 kg", "finishing": "Matt", "warna": "Krem"}},
	{"CONTOH-CAT-003", "Cat tembok interior hijau sage 5 kg", katCat, "gama-cat", 189000, "galon", 1, stokSisa, "2 galon", false, true, []string{rKeluarga}, spek{"berat": "5 kg", "finishing": "Matt", "warna": "Hijau sage"}},
	{"CONTOH-CAT-004", "Cat eksterior tahan cuaca 20 kg", katCat, "merek-b", 640000, "pail", 1, stokAda, "", true, true, []string{rFasad, rTeras}, spek{"berat": "20 kg", "finishing": "Satin", "warna": "Putih"}},
	{"CONTOH-CAT-005", "Cat kayu dan besi hitam 1 kg", katCat, "merek-b", 72000, "kaleng", 1, stokAda, "", false, true, []string{rFasad}, spek{"berat": "1 kg", "finishing": "Glossy", "warna": "Hitam"}},
	{"CONTOH-CAT-006", "Cat dasar tembok 5 kg", katCat, "gama-cat", 142000, "galon", 1, stokNol, "", false, true, []string{rKeluarga, rFasad}, spek{"berat": "5 kg", "warna": "Putih"}},
	{"CONTOH-PNL-001", "Panel dinding motif kayu 3 meter", katCat, "psi-panel", 215000, "lembar", 2, stokPO, "", false, true, []string{rKeluarga}, spek{"ukuran": "16 × 300 cm", "warna": "Walnut"}},

	{"CONTOH-LED-001", "Lampu LED 5 watt putih", katLampu, "zeta-lampu", 22000, "pcs", 1, stokAda, "", false, true, []string{rMandi}, spek{"daya": "5 watt", "warna": "Putih (6500K)"}},
	{"CONTOH-LED-002", "Lampu LED 12 watt putih", katLampu, "zeta-lampu", 45000, "pcs", 1, stokSisa, "8 pcs", true, true, []string{rDapur, rKeluarga}, spek{"daya": "12 watt", "warna": "Putih (6500K)"}},
	{"CONTOH-LED-003", "Lampu LED 9 watt kuning hangat", katLampu, "zeta-lampu", 34000, "pcs", 1, stokAda, "", false, true, []string{rKeluarga}, spek{"daya": "9 watt", "warna": "Kuning (3000K)"}},
	{"CONTOH-LED-004", "Lampu LED 18 watt putih", katLampu, "merek-c", 61000, "pcs", 1, stokAda, "", false, true, []string{rDapur}, spek{"daya": "18 watt", "warna": "Putih (6500K)"}},
	{"CONTOH-LED-005", "Downlight LED 12 watt", katLampu, "merek-c", 88000, "pcs", 1, stokAda, "", false, true, []string{rKeluarga, rDapur}, spek{"daya": "12 watt", "warna": "Netral (4000K)"}},
	{"CONTOH-LED-006", "Lampu taman tahan air", katLampu, "merek-c", 175000, "pcs", 1, stokPO, "", false, true, []string{rTeras, rFasad}, spek{"daya": "7 watt", "ketahanan": "IP65"}},
	{"CONTOH-LED-007", "Lampu dinding teras minimalis", katLampu, "zeta-lampu", 0, "pcs", 1, stokNol, "", false, true, []string{rTeras, rFasad}, spek{"daya": "6 watt", "warna": "Kuning (3000K)"}},

	{"CONTOH-SAN-001", "Kloset duduk monoblok putih", katSanitasi, "beta-sanitair", 3250000, "unit", 1, stokSisa, "3 unit", true, true, []string{rMandi}, spek{"tipe": "Monoblok", "warna": "Putih", "sistem_bilas": "Dual flush", "garansi": "2 tahun"}},
	{"CONTOH-SAN-002", "Kloset jongkok putih", katSanitasi, "beta-sanitair", 285000, "unit", 1, stokAda, "", false, true, []string{rMandi}, spek{"tipe": "Jongkok", "warna": "Putih"}},
	{"CONTOH-SAN-003", "Wastafel sudut keramik", katSanitasi, "beta-sanitair", 460000, "unit", 1, stokAda, "", false, true, []string{rMandi}, spek{"bahan": "Keramik", "warna": "Putih"}},
	{"CONTOH-SAN-004", "Kran wastafel stainless", katSanitasi, "sigma-kran", 245000, "pcs", 1, stokPO, "", true, true, []string{rDapur, rMandi}, spek{"bahan": "Stainless steel", "ukuran_drat": "1/2 inch"}},
	{"CONTOH-SAN-005", "Kran dapur leher angsa", katSanitasi, "sigma-kran", 320000, "pcs", 1, stokAda, "", false, true, []string{rDapur}, spek{"bahan": "Stainless steel", "ukuran_drat": "1/2 inch"}},
	{"CONTOH-SAN-006", "Shower mandi set", katSanitasi, "sigma-kran", 390000, "set", 1, stokAda, "", false, true, []string{rMandi}, spek{"bahan": "Kuningan lapis krom"}},
	{"CONTOH-SAN-007", "Pipa air 1/2 inch 4 meter", katSanitasi, "delta-pipa", 38000, "batang", 5, stokAda, "", false, true, []string{rMandi, rDapur}, spek{"ukuran": "1/2 inch × 4 m", "bahan": "PVC"}},
	{"CONTOH-SAN-008", "Pompa air rumah tangga 125 watt", katSanitasi, "fi-pompa", 780000, "unit", 1, stokAda, "", false, true, nil, spek{"daya": "125 watt", "garansi": "1 tahun"}},

	{"CONTOH-LST-001", "Kabel listrik 2×1,5 mm 50 meter", katListrik, "eta-kabel", 395000, "rol", 1, stokAda, "", false, true, nil, spek{"ukuran": "2 × 1,5 mm", "panjang": "50 m", "sni": "Ada"}},
	{"CONTOH-LST-002", "Kabel listrik 3×2,5 mm 50 meter", katListrik, "eta-kabel", 780000, "rol", 1, stokSisa, "4 rol", false, true, nil, spek{"ukuran": "3 × 2,5 mm", "panjang": "50 m", "sni": "Ada"}},
	{"CONTOH-LST-003", "Stop kontak dinding tunggal", katListrik, "rho-elektrik", 28000, "pcs", 1, stokAda, "", false, true, []string{rKeluarga, rDapur}, spek{"warna": "Putih"}},
	{"CONTOH-LST-004", "Saklar ganda dinding", katListrik, "rho-elektrik", 32000, "pcs", 1, stokAda, "", false, true, []string{rKeluarga}, spek{"warna": "Putih"}},
	{"CONTOH-LST-005", "Pemanas air listrik 15 liter", katListrik, "merek-e", 1450000, "unit", 1, stokAda, "", false, true, []string{rMandi}, spek{"kapasitas": "15 liter", "daya": "350 watt", "garansi": "1 tahun"}},
	// Nonaktif: tidak tampil di katalog publik, hanya di admin.
	{"CONTOH-LST-006", "Kipas angin dinding lama", katListrik, "merek-e", 310000, "unit", 1, stokNol, "", false, false, nil, spek{"ukuran": "16 inch"}},
}

// Harga tidak tampil di website (KEPUTUSAN #15). Harga 0 di daftar produk berarti harganya kosong sama sekali.
var hargaDisembunyikan = map[string]bool{"CONTOH-LST-005": true, "CONTOH-LED-007": true}

var heroContoh = struct{ subjudul, tombol, tombol2, tautan2 string }{
	"Cari bahan untuk tiap tahap rumahmu, cek harga dan stok.", "Jelajahi katalog", "Lihat promo", "/promo",
}

// mulai dan selesai adalah selisih hari dari hari ini; selesai nil berarti tanpa batas akhir.
var bannerDaftar = []struct {
	berkas, judul, keterangan, latar, tautan, posisi string
	lebar, tinggi                                    int
	mulai                                            int
	selesai                                          any
	urutan                                           int
	aktif                                            bool
}{
	{"hero", "Dari pondasi sampai atap, semua kebutuhan rumah ada di sini.", "", "#3A4A5A", "/katalog", "HOME_HERO", 1920, 800, -30, nil, 1, true},
	{"slider-keramik", "Promo keramik lantai", "Gambar banner promo", "#E7D7C1", "/katalog/" + katKeramik, "HOME_SLIDER", 1600, 640, -6, nil, 1, true},
	{"slider-cat", "Pekan cat tembok", "Gambar banner promo", "#DCE3D2", "/katalog/" + katCat, "HOME_SLIDER", 1600, 640, -10, nil, 2, true},
	{"slider-sanitasi", "Perlengkapan kamar mandi pilihan", "Gambar banner promo", "#D4E2EA", "/katalog/" + katSanitasi, "HOME_SLIDER", 1600, 640, -3, nil, 3, true},
	{"samping-lampu", "Lampu LED hemat energi", "Gambar banner promo", "#FDF0C2", "/katalog/" + katLampu, "HOME_SIDE", 800, 520, -3, nil, 1, true},
	{"samping-listrik", "Kelengkapan listrik rumah", "Gambar banner promo", "#F3D9C4", "/katalog/" + katListrik, "HOME_SIDE", 800, 520, -5, nil, 2, true},
	// Tiga banner ini tidak tampil di Beranda; gunanya menguji status di halaman admin.
	{"terjadwal", "Promo akhir tahun", "Gambar banner promo", "#E9DFCE", "/promo", "HOME_SLIDER", 1600, 640, 20, 34, 4, true},
	{"berakhir", "Promo sanitasi lama", "Gambar banner promo", "#D6D3CD", "/katalog/" + katSanitasi, "HOME_SLIDER", 1600, 640, -40, -10, 5, true},
	{"nonaktif", "Promo kran dapur", "Gambar banner promo", "#E7E2F0", "/katalog/" + katSanitasi, "HOME_SLIDER", 1600, 640, -2, nil, 6, false},
}

// Format isi mengikuti utils/artikel.js: paragraf dipisah baris kosong, "## " subjudul, "- " daftar.
var artikelDaftar = []struct {
	tipe, judul, slug, label, latar, status string
	hari                                    int
	isi                                     string
}{
	{"PROMO", "Promo keramik lantai bulan ini", "contoh-promo-keramik-lantai", "Promo keramik", "#E7D7C1", "PUBLISHED", 6, "Ini contoh promo untuk pengujian tampilan, bukan penawaran sungguhan.\n\n## Berlaku untuk\n- Keramik lantai pilihan\n- Pembelian langsung di toko\n\nTanyakan stok dan harga lewat WhatsApp sebelum datang ke toko."},
	{"PROMO", "Paket cat tembok untuk satu rumah", "contoh-paket-cat-tembok", "Paket cat", "#FDF0C2", "PUBLISHED", 2, "Ini contoh promo paket untuk pengujian tampilan, bukan penawaran sungguhan.\n\nKirim ukuran dinding rumahmu lewat WhatsApp, tim kami bantu hitung kebutuhan catnya."},
	{"EVENT", "Demo pemasangan keramik di toko", "contoh-demo-pemasangan-keramik", "Demo keramik", "#DCE3D2", "PUBLISHED", 12, "Ini contoh event untuk pengujian tampilan, bukan acara sungguhan.\n\n## Yang akan dibahas\n- Menyiapkan perekat keramik\n- Meratakan susunan keramik\n- Menghitung kebutuhan satu ruangan"},
	{"ARTIKEL", "Memilih keramik lantai yang aman di kaki", "contoh-memilih-keramik-lantai", "Memilih keramik", "#E9DFCE", "PUBLISHED", 1, "Ini contoh artikel untuk pengujian tampilan. Lantai dipakai setiap hari dan paling repot diganti, jadi pilih sesuai fungsi ruangannya.\n\n## Sesuaikan dengan ruangan\n- Ruang keluarga: ukuran besar agar terasa lega\n- Teras: permukaan kesat agar tidak licin\n- Kamar mandi: wajib kesat\n\nTambahkan cadangan 5–10% untuk potongan dan keramik yang pecah."},
	{"ARTIKEL", "Warna lampu yang pas untuk tiap ruangan", "contoh-warna-lampu-tiap-ruangan", "Warna lampu", "#F3D9C4", "PUBLISHED", 3, "Ini contoh artikel untuk pengujian tampilan. Warna lampu yang tepat membuat ruangan nyaman di mata.\n\n## Panduan singkat\n- Putih terang: dapur dan meja kerja\n- Putih netral: ruang keluarga\n- Kuning hangat: kamar tidur"},
	{"ARTIKEL", "Menyiapkan dinding sebelum dicat", "contoh-menyiapkan-dinding", "Menyiapkan dinding", "#D6D3CD", "PUBLISHED", 5, "Ini contoh artikel untuk pengujian tampilan. Cat awet bila dindingnya kering dan bersih.\n\n## Langkahnya\n- Perbaiki rembesan air lebih dulu\n- Bersihkan dan tambal bagian yang retak\n- Oleskan cat dasar, lalu dua lapis cat warna"},
	{"ARTIKEL", "Perlengkapan kamar mandi yang mudah dibersihkan", "contoh-perlengkapan-kamar-mandi", "Kamar mandi", "#D4E2EA", "PUBLISHED", 8, "Ini contoh artikel untuk pengujian tampilan. Perlengkapan yang tepat membuat kamar mandi lebih awet.\n\n## Yang perlu diperhatikan\n- Ukur ruangan sebelum memilih kloset\n- Pilih kran yang tahan karat\n- Pasang saringan rambut di lubang pembuangan"},
	// Draf: hanya tampil di admin, untuk menguji alur terbitkan.
	{"ARTIKEL", "Memilih sofa untuk ruang keluarga kecil", "contoh-memilih-sofa", "Memilih sofa", "#E7E2F0", "DRAFT", 1, "Ini contoh artikel draf untuk pengujian alur terbit di admin."},
}

type transaksiContoh struct {
	tipe    string
	poin    int64 // EARN dihitung dari nilai dengan rasio point_ratio_rupiah; lainnya bertanda sesuai saldo.
	nilai   int64
	catatan string
	hari    int
}

// Tidak punya kata sandi, jadi akun-akun ini tidak bisa dipakai masuk.
var memberDaftar = []struct {
	nama, lokal                  string
	setuju, aktif, terverifikasi bool
	sinkron                      string
	hari                         int
	transaksi                    []transaksiContoh
}{
	{"Budi Santoso", "budi.santoso", true, true, true, "SYNCED", 150, []transaksiContoh{{"EARN", 0, 4480000, "", 140}, {"EARN", 0, 6900000, "", 110}, {"ADJUST", -1, 0, "Koreksi nota ganda", 105}, {"REDEEM", -5, 0, "Tukar potongan harga di kasir", 80}, {"EARN", 0, 12800000, "", 40}}},
	{"Siti Rahma", "siti.rahma", true, true, true, "SYNCED", 130, []transaksiContoh{{"EARN", 0, 3600000, "", 30}}},
	{"Andi Wijaya", "andi.wijaya", false, true, true, "SYNCED", 120, nil},
	{"Agus Prasetyo", "agus.prasetyo", true, true, true, "SYNCED", 110, []transaksiContoh{{"EARN", 0, 5300000, "", 100}, {"EARN", 0, 2750000, "", 60}}},
	{"Fitri Handayani", "fitri.handayani", true, true, true, "SYNCED", 100, []transaksiContoh{{"EARN", 0, 1100000, "", 90}}},
	{"Hendra Gunawan", "hendra.gunawan", true, false, true, "SYNCED", 90, []transaksiContoh{{"EARN", 0, 8450000, "", 85}, {"REDEEM", -5, 0, "Tukar voucher di kasir", 50}}},
	{"Ratna Sari", "ratna.sari", false, true, true, "SYNCED", 80, nil},
	{"Joko Susilo", "joko.susilo", true, true, true, "SYNCED", 70, []transaksiContoh{{"EARN", 0, 15200000, "", 65}, {"ADJUST", 1, 0, "Koreksi pembulatan", 60}}},
	{"Maya Puspita", "maya.puspita", true, true, true, "SYNCED", 60, []transaksiContoh{{"EARN", 0, 3050000, "", 55}}},
	{"Fajar Nugroho", "fajar.nugroho", true, true, true, "FAILED", 30, []transaksiContoh{{"EARN", 0, 2200000, "", 20}}},
	{"Dimas Pratama", "dimas.pratama", false, true, true, "FAILED", 20, nil},
	{"Yoga Firmansyah", "yoga.firmansyah", false, true, false, "PENDING", 3, nil},
	{"Wulan Septiani", "wulan.septiani", true, true, false, "PENDING", 1, nil},
}

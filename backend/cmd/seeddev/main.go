// Command seeddev mengisi database LOKAL dengan data fiktif supaya semua halaman publik dan admin bisa dicoba:
// brand, ruangan, produk, jasa, banner, artikel, member, poin, dan kontak Qontak. Bukan migration, jadi tidak ikut ke server.
//
//	cd backend
//	go run ./cmd/seeddev          isi data contoh (aman dijalankan ulang)
//	go run ./cmd/seeddev -hapus   hapus data contoh beserta file gambarnya
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	folderLogo   = "../frontend/public/contoh-logo"
	folderGambar = "../frontend/public/contoh-gambar"
	urlGambar    = "/contoh-gambar/"
	// Foto produk harus lewat alamat unggahan: form admin menolak alamat lain saat produk disimpan ulang.
	folderFoto = "storage/uploads/products"
	urlFoto    = "/api/v1/uploads/"
	// Domain .test dicadangkan (RFC 2606), jadi tidak mungkin milik pelanggan sungguhan.
	domainMember = "@contoh.test"
	// jumlahBerlogo mengikuti data mock: hanya brand pertama yang punya logo, sisanya tampil sebagai nama.
	jumlahBerlogo = 10
	// Produk pertama punya galeri foto.
	jumlahBergaleri = 4
	rasioPoin       = 1000000
)

var (
	bukanSlug    = regexp.MustCompile(`[^a-z0-9]+`)
	penggantiXML = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
)

func main() {
	hapus := flag.Bool("hapus", false, "hapus data contoh beserta file gambarnya")
	flag.Parse()

	dsn, err := ambilDatabaseURL()
	if err != nil {
		log.Fatal(err)
	}
	if err := pastikanLokal(dsn); err != nil {
		log.Fatal(err)
	}

	logx.Disable() // Log SQL go-zero per baris tidak berguna untuk seeder.
	ctx := context.Background()
	conn := postgres.New(dsn)
	if *hapus {
		err = hapusContoh(ctx, conn)
	} else {
		err = isiContoh(ctx, conn)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func ambilDatabaseURL() (string, error) {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn, nil
	}
	berkas, err := os.Open(".env")
	if err != nil {
		return "", errors.New("DATABASE_URL tidak ada; jalankan dari folder backend yang berisi .env")
	}
	defer berkas.Close()
	baris := bufio.NewScanner(berkas)
	for baris.Scan() {
		if nilai, ada := strings.CutPrefix(strings.TrimSpace(baris.Text()), "DATABASE_URL="); ada {
			return strings.Trim(nilai, `"'`), nil
		}
	}
	return "", errors.New("DATABASE_URL tidak ada di .env")
}

// pastikanLokal mencegah data fiktif masuk ke database selain milik komputer ini.
// Pesan galat sengaja tidak memuat alamat database karena bisa berisi kredensial.
func pastikanLokal(dsn string) error {
	alamat, err := url.Parse(dsn)
	if err != nil {
		return errors.New("DATABASE_URL tidak bisa dibaca; seeder dibatalkan")
	}
	switch alamat.Hostname() {
	case "localhost", "127.0.0.1", "::1":
		return nil
	}
	return errors.New("seeder hanya boleh dijalankan ke database lokal; dibatalkan")
}

func isiContoh(ctx context.Context, conn sqlx.SqlConn) error {
	if err := os.MkdirAll(folderGambar, 0o755); err != nil {
		return err
	}
	langkah := []struct {
		nama string
		isi  func(context.Context, sqlx.SqlConn) error
	}{
		{"brand", isiBrand}, {"pengaturan", isiPengaturan}, {"kategori", isiKategori}, {"ruangan", isiRuangan},
		{"produk", isiProduk}, {"jasa", isiJasa}, {"banner", isiBanner}, {"artikel", isiArtikel}, {"member", isiMember},
	}
	for _, l := range langkah {
		if err := l.isi(ctx, conn); err != nil {
			return fmt.Errorf("gagal mengisi %s: %w", l.nama, err)
		}
	}
	fmt.Printf("Selesai: %d brand, %d ruangan, %d produk, %d banner, %d artikel, %d member contoh ada di database lokal.\n",
		len(namaBrand), len(ruanganContoh), len(produkDaftar), len(bannerDaftar), len(artikelDaftar), len(memberDaftar))
	return nil
}

func isiBrand(ctx context.Context, conn sqlx.SqlConn) error {
	if err := os.MkdirAll(folderLogo, 0o755); err != nil {
		return err
	}
	for i, nama := range namaBrand {
		kode := slug(nama)
		var logo any
		if i < jumlahBerlogo {
			if err := os.WriteFile(filepath.Join(folderLogo, kode+".svg"), []byte(logoContoh(nama, i)), 0o644); err != nil {
				return err
			}
			logo = "/contoh-logo/" + kode + ".svg"
		}
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.brands (name,slug,logo_url) VALUES ($1,$2,$3)
			ON CONFLICT (slug) DO NOTHING`, nama, kode, logo); err != nil {
			return err
		}
	}
	return nil
}

// isiPengaturan hanya mengisi yang belum ada atau masih kosong; nilai yang sudah diisi admin tidak tertimpa.
func isiPengaturan(ctx context.Context, conn sqlx.SqlConn) error {
	for _, isian := range pengaturanContoh {
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.site_settings (setting_key,setting_value,updated_at)
			VALUES ($1,$2,NOW()) ON CONFLICT (setting_key) DO UPDATE
			SET setting_value=EXCLUDED.setting_value,updated_at=NOW()
			WHERE COALESCE(crm_schema.site_settings.setting_value,'')=''`, isian[0], isian[1]); err != nil {
			return err
		}
	}
	return nil
}

func isiKategori(ctx context.Context, conn sqlx.SqlConn) error {
	for i, k := range kategoriContoh {
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.categories (name,slug,sort_order) VALUES ($1,$2,$3)
			ON CONFLICT (slug) DO NOTHING`, k[0], k[1], i+1); err != nil {
			return err
		}
	}
	return nil
}

func isiRuangan(ctx context.Context, conn sqlx.SqlConn) error {
	for i, r := range ruanganContoh {
		kode := slug(r.nama)
		sampul, err := tulisGambar("ruangan-"+kode, r.nama, "Inspirasi ruangan", r.latar, 800, 800)
		if err != nil {
			return err
		}
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.rooms (name,slug,image_cover,sort_order) VALUES ($1,$2,$3,$4)
			ON CONFLICT (slug) DO NOTHING`, r.nama, kode, sampul, i+1); err != nil {
			return err
		}
	}
	return nil
}

func isiProduk(ctx context.Context, conn sqlx.SqlConn) error {
	for i, p := range produkDaftar {
		spesifikasi, err := json.Marshal(p.spesifikasi)
		if err != nil {
			return err
		}
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.products
			(sku,name,slug,category_id,brand_id,description,specifications,price_general,unit_sale,min_order,stock_status,stock_qty_label,is_featured,is_active,show_price)
			VALUES ($1,$2,$3,(SELECT id FROM crm_schema.categories WHERE slug=$4),(SELECT id FROM crm_schema.brands WHERE slug=$5),$6,$7::jsonb,$8,$9,$10,$11,$12,$13,$14,$15)
			ON CONFLICT (sku) DO NOTHING`,
			p.sku, p.nama, slug(p.nama), p.kategori, p.brand,
			p.nama+". Ini data contoh untuk pengujian tampilan, bukan produk yang dijual.\nHarga dan stok bersifat fiktif.",
			string(spesifikasi), hargaAtauNil(p.harga), p.satuan, p.minOrder, p.stok, kosongJadiNil(p.labelStok), p.unggulan, p.aktif, !hargaDisembunyikan[p.sku]); err != nil {
			return fmt.Errorf("%s: %w", p.sku, err)
		}
		var id int64
		if err := conn.QueryRowCtx(ctx, &id, `SELECT id FROM crm_schema.products WHERE sku=$1`, p.sku); err != nil {
			return fmt.Errorf("%s: %w", p.sku, err)
		}
		for urutan := 0; urutan < jumlahFoto(i); urutan++ {
			foto, err := tulisFoto(fmt.Sprintf("%s-%d", strings.ToLower(p.sku), urutan+1), latarFoto[(i+urutan)%len(latarFoto)])
			if err != nil {
				return err
			}
			if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.product_images (product_id,image_url,is_primary,sort_order)
				SELECT $1::bigint,$2::varchar,$3::boolean,$4::int
				WHERE NOT EXISTS (SELECT 1 FROM crm_schema.product_images WHERE product_id=$1::bigint AND image_url=$2::varchar)`,
				id, foto, urutan == 0, urutan); err != nil {
				return err
			}
		}
		for _, ruangan := range p.ruangan {
			if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.product_rooms (product_id,room_id)
				SELECT $1::bigint,id FROM crm_schema.rooms WHERE slug=$2 ON CONFLICT DO NOTHING`, id, ruangan); err != nil {
				return err
			}
		}
	}
	return nil
}

// jumlahFoto menyebar beberapa produk tanpa foto untuk menguji placeholder. Posisinya dipilih agar tidak
// mengenai produk unggulan dan tidak berderet di awal katalog (urutan bawaan: terbaru dulu).
func jumlahFoto(i int) int {
	switch {
	case i < jumlahBergaleri:
		return len(labelFoto)
	case i%10 == 7:
		return 0
	}
	return 1
}

// isiJasa melengkapi 7 jasa bawaan migration: gambar yang masih kosong dan relasi ke kategori produk.
func isiJasa(ctx context.Context, conn sqlx.SqlConn) error {
	for _, j := range jasaContoh {
		gambar, err := tulisGambar("jasa-"+j.slug, j.label, "Jasa Transhome", j.latar, 800, 450)
		if err != nil {
			return err
		}
		if _, err := conn.ExecCtx(ctx, `UPDATE crm_schema.services SET image_url=$2
			WHERE slug=$1 AND COALESCE(image_url,'')=''`, j.slug, gambar); err != nil {
			return err
		}
		for _, kategori := range j.kategori {
			if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.service_categories (service_id,category_id)
				SELECT s.id,c.id FROM crm_schema.services s, crm_schema.categories c WHERE s.slug=$1 AND c.slug=$2
				ON CONFLICT DO NOTHING`, j.slug, kategori); err != nil {
				return err
			}
		}
	}
	return nil
}

func isiBanner(ctx context.Context, conn sqlx.SqlConn) error {
	for _, b := range bannerDaftar {
		// Judul tidak digambar: komponen banner dan hero sudah menampilkannya di atas gambar.
		gambar, err := tulisGambar("banner-"+b.berkas, "", b.keterangan, b.latar, b.lebar, b.tinggi)
		if err != nil {
			return err
		}
		hero := []any{nil, nil, nil, nil}
		if b.posisi == "HOME_HERO" {
			hero = []any{heroContoh.subjudul, heroContoh.tombol, heroContoh.tombol2, heroContoh.tautan2}
		}
		// Waktu banner disimpan dalam WIB, sama dengan pembanding di bannersHandler.
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.banners
			(title,image_url,link_url,position,start_at,end_at,sort_order,is_active,subtitle,cta_label,cta2_label,cta2_url)
			SELECT $1::varchar,$2::varchar,$3::varchar,$4::varchar,
				timezone('Asia/Jakarta',now()) + $5::int * INTERVAL '1 day',
				timezone('Asia/Jakarta',now()) + $6::int * INTERVAL '1 day',
				$7::int,$8::boolean,$9::varchar,$10::varchar,$11::varchar,$12::varchar
			WHERE NOT EXISTS (SELECT 1 FROM crm_schema.banners WHERE image_url=$2::varchar)`,
			append([]any{b.judul, gambar, b.tautan, b.posisi, b.mulai, b.selesai, b.urutan, b.aktif}, hero...)...); err != nil {
			return fmt.Errorf("%s: %w", b.berkas, err)
		}
	}
	return nil
}

func isiArtikel(ctx context.Context, conn sqlx.SqlConn) error {
	for _, a := range artikelDaftar {
		keterangan := "Promo & event"
		if a.tipe == "ARTIKEL" {
			keterangan = "Panduan isi rumah"
		}
		gambar, err := tulisGambar("artikel-"+a.slug, a.label, keterangan, a.latar, 800, 450)
		if err != nil {
			return err
		}
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.articles
			(type,title,slug,thumbnail_url,content,status,published_at,created_at,updated_at)
			VALUES ($1,$2,$3,$4,$5,$6::varchar,
				CASE WHEN $6::varchar='PUBLISHED' THEN NOW() - $7::int * INTERVAL '1 day' END,
				NOW() - $7::int * INTERVAL '1 day',NOW() - $7::int * INTERVAL '1 day')
			ON CONFLICT (slug) DO NOTHING`, a.tipe, a.judul, a.slug, gambar, a.isi, a.status, a.hari); err != nil {
			return fmt.Errorf("%s: %w", a.slug, err)
		}
	}
	return nil
}

func isiMember(ctx context.Context, conn sqlx.SqlConn) error {
	for i, m := range memberDaftar {
		nomor := i + 1
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.users (email,phone_number,is_active,email_verified_at,phone_verified_at,created_at)
			SELECT $1::varchar,$2::varchar,$3::boolean,
				CASE WHEN $4::boolean THEN NOW() - $5::int * INTERVAL '1 day' END,
				CASE WHEN $4::boolean THEN NOW() - $5::int * INTERVAL '1 day' END,
				NOW() - $5::int * INTERVAL '1 day'
			ON CONFLICT (email) DO NOTHING`,
			m.lokal+domainMember, fmt.Sprintf("628000000%04d", nomor), m.aktif, m.terverifikasi, m.hari); err != nil {
			return fmt.Errorf("%s: %w", m.lokal, err)
		}
		var idUser int64
		if err := conn.QueryRowCtx(ctx, &idUser, `SELECT id FROM crm_schema.users WHERE email=$1`, m.lokal+domainMember); err != nil {
			return fmt.Errorf("%s: %w", m.lokal, err)
		}
		kode := fmt.Sprintf("TF-9%05d", nomor)
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.member_profiles (user_id,member_code,full_name,communication_consent,consent_at,created_at)
			SELECT $1::bigint,$2::varchar,$3::varchar,$4::boolean,
				CASE WHEN $4::boolean THEN NOW() - $5::int * INTERVAL '1 day' END,
				NOW() - $5::int * INTERVAL '1 day'
			ON CONFLICT (user_id) DO NOTHING`, idUser, kode, m.nama, m.setuju, m.hari); err != nil {
			return fmt.Errorf("%s: %w", m.lokal, err)
		}
		if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.user_roles (user_id,role_id)
			SELECT $1::bigint,id FROM crm_schema.roles WHERE code='CUSTOMER' ON CONFLICT DO NOTHING`, idUser); err != nil {
			return err
		}
		if err := isiQontak(ctx, conn, idUser, nomor, kode, m.sinkron, m.hari); err != nil {
			return fmt.Errorf("%s: %w", m.lokal, err)
		}
		var idProfil int64
		if err := conn.QueryRowCtx(ctx, &idProfil, `SELECT id FROM crm_schema.member_profiles WHERE user_id=$1`, idUser); err != nil {
			return fmt.Errorf("%s: %w", m.lokal, err)
		}
		for j, t := range m.transaksi {
			poin, nilai := t.poin, any(nil)
			if t.tipe == "EARN" {
				poin, nilai = t.nilai/rasioPoin, t.nilai
			}
			// Nomor referensi juga diisi untuk REDEEM/ADJUST: kunci unik (type, reference_no) membuat seeder aman diulang.
			if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.point_transactions (member_id,type,points,reference_no,purchase_amount,note,created_at)
				VALUES ($1,$2,$3,$4,$5,$6,NOW() - $7::int * INTERVAL '1 day') ON CONFLICT (type,reference_no) DO NOTHING`,
				idProfil, t.tipe, poin, fmt.Sprintf("CONTOH-%02d%02d", nomor, j+1), nilai, kosongJadiNil(t.catatan), t.hari); err != nil {
				return fmt.Errorf("%s: %w", m.lokal, err)
			}
		}
	}
	return nil
}

func isiQontak(ctx context.Context, conn sqlx.SqlConn, idUser int64, nomor int, kode, status string, hari int) error {
	var idKontak any
	if status == "SYNCED" {
		idKontak = fmt.Sprintf("CONTOH-QTK-%04d", nomor)
	}
	if _, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.qontak_contacts (user_id,qontak_contact_id,sync_status,last_synced_at,created_at)
		SELECT $1::bigint,$2::varchar,$3::varchar,
			CASE WHEN $3::varchar='SYNCED' THEN NOW() - $4::int * INTERVAL '1 day' END,
			NOW() - $4::int * INTERVAL '1 day'
		ON CONFLICT (user_id) DO NOTHING`, idUser, idKontak, status, hari); err != nil {
		return err
	}
	if status == "PENDING" {
		return nil // Belum pernah dikirim, jadi belum ada log.
	}
	kodeStatus, percobaan, jawaban := 200, 0, `{"status":"ok"}`
	if status == "FAILED" {
		kodeStatus, percobaan, jawaban = 503, 3, `{"error":"Contoh kegagalan: layanan tidak tersedia"}`
	}
	_, err := conn.ExecCtx(ctx, `INSERT INTO crm_schema.sync_logs (user_id,event_type,request_payload,response_payload,status_code,retry_count,created_at)
		SELECT $1::bigint,'CONTACT_UPSERT',$2::jsonb,$3::jsonb,$4::int,$5::int,NOW() - $6::int * INTERVAL '1 day'
		WHERE NOT EXISTS (SELECT 1 FROM crm_schema.sync_logs WHERE user_id=$1::bigint AND event_type='CONTACT_UPSERT')`,
		idUser, fmt.Sprintf(`{"member_code":%q}`, kode), jawaban, kodeStatus, percobaan, hari)
	return err
}

// hapusContoh hanya menghapus baris yang masih bisa dikenali sebagai contoh (SKU, domain email, atau gambar contoh).
// Kategori, pengaturan, dan relasi jasa–kategori tidak disentuh karena bisa saja sudah dipakai sungguhan.
func hapusContoh(ctx context.Context, conn sqlx.SqlConn) error {
	var produk, ruangan, banner, artikel, member int64
	for _, p := range produkDaftar {
		n, err := jalankan(ctx, conn, `DELETE FROM crm_schema.products WHERE sku=$1`, p.sku)
		if err != nil {
			return err
		}
		produk += n
	}
	for _, r := range ruanganContoh {
		n, err := jalankan(ctx, conn, `DELETE FROM crm_schema.rooms WHERE slug=$1 AND image_cover LIKE $2`, slug(r.nama), urlGambar+"%")
		if err != nil {
			return err
		}
		ruangan += n
	}
	banner, err := jalankan(ctx, conn, `DELETE FROM crm_schema.banners WHERE image_url LIKE $1`, urlGambar+"%")
	if err != nil {
		return err
	}
	for _, a := range artikelDaftar {
		n, err := jalankan(ctx, conn, `DELETE FROM crm_schema.articles WHERE slug=$1 AND thumbnail_url LIKE $2`, a.slug, urlGambar+"%")
		if err != nil {
			return err
		}
		artikel += n
	}
	if _, err := jalankan(ctx, conn, `UPDATE crm_schema.services SET image_url=NULL WHERE image_url LIKE $1`, urlGambar+"%"); err != nil {
		return err
	}
	// sync_logs.user_id memakai ON DELETE SET NULL, jadi lognya harus dihapus sebelum usernya.
	if _, err := jalankan(ctx, conn, `DELETE FROM crm_schema.sync_logs WHERE user_id IN
		(SELECT id FROM crm_schema.users WHERE email LIKE $1)`, "%"+domainMember); err != nil {
		return err
	}
	if member, err = jalankan(ctx, conn, `DELETE FROM crm_schema.users WHERE email LIKE $1`, "%"+domainMember); err != nil {
		return err
	}
	var brand int64
	for i, nama := range namaBrand {
		kode := slug(nama)
		n, err := jalankan(ctx, conn, `DELETE FROM crm_schema.brands b WHERE b.slug=$1
			AND NOT EXISTS (SELECT 1 FROM crm_schema.products p WHERE p.brand_id=b.id)`, kode)
		if err != nil {
			return err
		}
		brand += n
		if i < jumlahBerlogo && n > 0 {
			if err := os.Remove(filepath.Join(folderLogo, kode+".svg")); err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
	}
	if err := os.RemoveAll(folderGambar); err != nil {
		return err
	}
	fotoContoh, err := filepath.Glob(filepath.Join(folderFoto, "contoh-*.png"))
	if err != nil {
		return err
	}
	for _, foto := range fotoContoh {
		if err := os.Remove(foto); err != nil {
			return err
		}
	}
	fmt.Printf("Selesai: dihapus %d produk, %d ruangan, %d banner, %d artikel, %d member, %d brand contoh.\n", produk, ruangan, banner, artikel, member, brand)
	fmt.Println("Kategori, pengaturan toko, dan relasi jasa–kategori tidak diubah; ganti lewat dashboard admin.")
	return nil
}

func jalankan(ctx context.Context, conn sqlx.SqlConn, kueri string, argumen ...any) (int64, error) {
	hasil, err := conn.ExecCtx(ctx, kueri, argumen...)
	if err != nil {
		return 0, err
	}
	return hasil.RowsAffected()
}

func hargaAtauNil(harga int64) any {
	if harga == 0 {
		return nil
	}
	return harga
}

func kosongJadiNil(teks string) any {
	if teks == "" {
		return nil
	}
	return teks
}

func slug(nama string) string {
	return strings.Trim(bukanSlug.ReplaceAllString(strings.ToLower(nama), "-"), "-")
}

// tulisGambar membuat SVG berlabel "CONTOH" dan mengembalikan alamat yang disimpan ke database.
func tulisGambar(nama, judul, keterangan, latar string, lebar, tinggi int) (string, error) {
	var isi strings.Builder
	yKeterangan := tinggi / 3 // Tanpa judul: teks di bagian atas, karena komponen banner menaruh judul di bagian bawah.
	if judul != "" {
		yKeterangan = tinggi/2 + tinggi/9
		fmt.Fprintf(&isi, `<text x="%d" y="%d" text-anchor="middle" font-family="sans-serif" font-size="%d" font-weight="700" fill="#1C2B3A">%s</text>`,
			lebar/2, tinggi/2-10, tinggi/12, penggantiXML.Replace(judul))
	}
	if keterangan != "" {
		fmt.Fprintf(&isi, `<text x="%d" y="%d" text-anchor="middle" font-family="sans-serif" font-size="%d" fill="#57534E">%s</text>`+
			`<text x="%d" y="%d" text-anchor="end" font-family="sans-serif" font-size="%d" fill="#57534E">CONTOH GAMBAR</text>`,
			lebar/2, yKeterangan, tinggi/20, penggantiXML.Replace(keterangan), lebar-32, tinggi-28, tinggi/28)
	}
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`+
		`<rect width="%d" height="%d" fill="%s"/>%s</svg>`,
		lebar, tinggi, lebar, tinggi, lebar, tinggi, latar, isi.String())
	if err := os.WriteFile(filepath.Join(folderGambar, nama+".svg"), []byte(svg), 0o644); err != nil {
		return "", err
	}
	return urlGambar + nama + ".svg", nil
}

// tulisFoto membuat PNG polos, satu-satunya format di sini yang lolos aturan unggahan (JPG, PNG, WebP).
func tulisFoto(nama, latar string) (string, error) {
	var merah, hijau, biru uint8
	if _, err := fmt.Sscanf(latar, "#%02x%02x%02x", &merah, &hijau, &biru); err != nil {
		return "", err
	}
	if err := os.MkdirAll(folderFoto, 0o750); err != nil {
		return "", err
	}
	foto := image.NewPaletted(image.Rect(0, 0, 800, 800), color.Palette{color.RGBA{merah, hijau, biru, 255}})
	berkas, err := os.Create(filepath.Join(folderFoto, nama+".png"))
	if err != nil {
		return "", err
	}
	defer berkas.Close()
	if err := png.Encode(berkas, foto); err != nil {
		return "", err
	}
	return urlFoto + nama + ".png", nil
}

// logoContoh membuat SVG yang sama dengan logoContoh() di frontend/fase-0-fondasi/src/services/mock/data/catalog.js.
func logoContoh(nama string, i int) string {
	warna := warnaLogo[i%len(warnaLogo)]
	var huruf strings.Builder
	for _, kata := range strings.Fields(nama) {
		if huruf.Len() < 2 {
			huruf.WriteString(strings.ToUpper(kata[:1]))
		}
	}
	sudut := 18
	if i%2 == 1 {
		sudut = 42
	}
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200" viewBox="0 0 200 200">`+
		`<rect x="58" y="28" width="84" height="84" rx="%d" fill="%s"/>`+
		`<text x="100" y="84" text-anchor="middle" font-family="sans-serif" font-size="38" font-weight="800" fill="#fff">%s</text>`+
		`<text x="100" y="152" text-anchor="middle" font-family="sans-serif" font-size="24" font-weight="800" fill="%s">%s</text>`+
		`</svg>`, sudut, warna, huruf.String(), warna, strings.ToUpper(nama))
}

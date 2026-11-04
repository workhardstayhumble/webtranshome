package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"webtranshome/internal/svc"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// Filter dan faset GET /products mengikuti "Detail kontrak katalog" di frontend/docs/FRONTEND_BUILD_PLAN.md.
var (
	statusStokKatalog = []string{"TERSEDIA", "SISA_STOK", "PRE_ORDER", "HABIS"}
	// Daftar resmi dari pemilik masih ditunggu (frontend/CLAUDE.md §12); ini sama dengan yang dipakai mock.
	atributFilterKatalog = []string{"ukuran", "finishing", "warna", "bahan", "daya"}
)

type produkKatalog struct {
	row         publicProductRow
	spesifikasi map[string]string
	ruangan     map[string]bool
}

func (p produkKatalog) hargaTerlihat() bool { return p.row.ShowPrice && p.row.PriceGeneral.Valid }

type opsiKatalog struct {
	Name string `db:"name"`
	Slug string `db:"slug"`
}

// filterKatalog berisi pilihan pengguna yang angkanya dihitung disjungtif.
type filterKatalog struct {
	brand, stock, room map[string]bool
	attr               map[string]map[string]bool
}

func muatProdukKatalog(ctx context.Context, s *svc.ServiceContext) ([]produkKatalog, error) {
	rows, err := fetchPublicProducts(ctx, s)
	if err != nil {
		return nil, err
	}
	var tautan []struct {
		ProductID int64  `db:"product_id"`
		Slug      string `db:"slug"`
	}
	err = s.DB.QueryRowsCtx(ctx, &tautan, `SELECT pr.product_id,room.slug FROM crm_schema.product_rooms pr
		JOIN crm_schema.rooms room ON room.id=pr.room_id WHERE room.is_active=TRUE`)
	if err != nil {
		return nil, err
	}
	ruangan := map[int64]map[string]bool{}
	for _, t := range tautan {
		if ruangan[t.ProductID] == nil {
			ruangan[t.ProductID] = map[string]bool{}
		}
		ruangan[t.ProductID][t.Slug] = true
	}
	hasil := make([]produkKatalog, 0, len(rows))
	for _, row := range rows {
		hasil = append(hasil, produkKatalog{row: row, spesifikasi: bacaAtributFilter(row.Specifications), ruangan: ruangan[row.ID]})
	}
	return hasil, nil
}

func muatOpsiFaset(ctx context.Context, s *svc.ServiceContext) (merek, ruangan []opsiKatalog, err error) {
	if err = s.DB.QueryRowsCtx(ctx, &merek, `SELECT name,slug FROM crm_schema.brands WHERE is_active=TRUE ORDER BY name,id`); err != nil {
		return nil, nil, err
	}
	err = s.DB.QueryRowsCtx(ctx, &ruangan, `SELECT name,slug FROM crm_schema.rooms WHERE is_active=TRUE ORDER BY sort_order,name,id`)
	return merek, ruangan, err
}

func bacaAtributFilter(raw sql.NullString) map[string]string {
	if !raw.Valid {
		return nil
	}
	var isi map[string]any
	if json.Unmarshal([]byte(raw.String), &isi) != nil {
		return nil
	}
	hasil := map[string]string{}
	for _, kunci := range atributFilterKatalog {
		if nilai, ok := isi[kunci].(string); ok && strings.TrimSpace(nilai) != "" {
			hasil[kunci] = strings.TrimSpace(nilai)
		}
	}
	return hasil
}

func daftarKoma(raw string) map[string]bool {
	hasil := map[string]bool{}
	for _, nilai := range strings.Split(raw, ",") {
		if nilai = strings.TrimSpace(nilai); nilai != "" {
			hasil[nilai] = true
		}
	}
	return hasil
}

func bacaFilterKatalog(q url.Values, semua []produkKatalog) filterKatalog {
	f := filterKatalog{brand: daftarKoma(q.Get("brand")), stock: daftarKoma(q.Get("stock")), room: daftarKoma(q.Get("room")), attr: map[string]map[string]bool{}}
	for _, kunci := range atributFilterKatalog {
		if raw := q.Get("attr[" + kunci + "]"); raw != "" {
			f.attr[kunci] = pilihanAtribut(raw, nilaiAtribut(semua, kunci))
		}
	}
	return f
}

func nilaiAtribut(produk []produkKatalog, kunci string) []string {
	unik := map[string]bool{}
	for _, p := range produk {
		if nilai := p.spesifikasi[kunci]; nilai != "" {
			unik[nilai] = true
		}
	}
	hasil := make([]string, 0, len(unik))
	for nilai := range unik {
		hasil = append(hasil, nilai)
	}
	return hasil
}

// pilihanAtribut memecah parameter yang dipisah koma dengan mencocokkan nilai yang dikenal, terpanjang dulu,
// karena nilai atributnya sendiri boleh berisi koma desimal ("2 × 1,5 mm").
func pilihanAtribut(raw string, dikenal []string) map[string]bool {
	urut := append([]string(nil), dikenal...)
	sort.Slice(urut, func(i, j int) bool { return len(urut[i]) > len(urut[j]) })
	hasil := map[string]bool{}
	for raw != "" {
		cocok := false
		for _, nilai := range urut {
			if raw == nilai || strings.HasPrefix(raw, nilai+",") {
				hasil[nilai] = true
				raw = strings.TrimPrefix(strings.TrimPrefix(raw, nilai), ",")
				cocok = true
				break
			}
		}
		if !cocok { // Nilai tak dikenal: tetap dianggap pilihan supaya hasilnya kosong, bukan mengabaikan filter.
			nilai, sisa, _ := strings.Cut(raw, ",")
			hasil[nilai] = true
			raw = sisa
		}
	}
	return hasil
}

// lolos memeriksa semua filter kecuali grup `kecuali` (faset disjungtif: angka opsi mengabaikan grupnya sendiri).
func (f filterKatalog) lolos(p produkKatalog, kecuali string) bool {
	if kecuali != "brand" && len(f.brand) > 0 && !f.brand[p.row.BrandSlug.String] {
		return false
	}
	if kecuali != "stock" && len(f.stock) > 0 && !f.stock[p.row.StockStatus] {
		return false
	}
	if kecuali != "room" && len(f.room) > 0 && !adaIrisan(f.room, p.ruangan) {
		return false
	}
	for kunci, pilihan := range f.attr {
		if kecuali != "attr:"+kunci && len(pilihan) > 0 && !pilihan[p.spesifikasi[kunci]] {
			return false
		}
	}
	return true
}

func adaIrisan(a, b map[string]bool) bool {
	for nilai := range a {
		if b[nilai] {
			return true
		}
	}
	return false
}

// saringDasarKatalog menerapkan pencarian, kategori, exclude, dan unggulan (tanpaHarga), lalu rentang harga (dasar).
// Produk yang harganya disembunyikan tidak ikut filter harga karena pelanggan tidak mengetahui harganya.
func saringDasarKatalog(semua []produkKatalog, q url.Values) (tanpaHarga, dasar []produkKatalog) {
	kata := strings.ToLower(strings.TrimSpace(q.Get("q")))
	kategori := q.Get("category")
	for _, p := range semua {
		r := p.row
		if kata != "" && !strings.Contains(strings.ToLower(r.Name+" "+r.SKU+" "+r.BrandName.String), kata) {
			continue
		}
		if kategori != "" && r.CategorySlug != kategori && r.CategoryParentSlug.String != kategori {
			continue
		}
		if q.Get("exclude") == strconv.FormatInt(r.ID, 10) || (q.Get("featured") == "1" && !r.IsFeatured) {
			continue
		}
		tanpaHarga = append(tanpaHarga, p)
	}
	minimum, adaMin := angkaQuery(q, "min")
	maksimum, adaMaks := angkaQuery(q, "max")
	if !adaMin && !adaMaks {
		return tanpaHarga, tanpaHarga
	}
	for _, p := range tanpaHarga {
		if !p.hargaTerlihat() || (adaMin && p.row.PriceGeneral.Float64 < minimum) || (adaMaks && p.row.PriceGeneral.Float64 > maksimum) {
			continue
		}
		dasar = append(dasar, p)
	}
	return tanpaHarga, dasar
}

func angkaQuery(q url.Values, kunci string) (float64, bool) {
	nilai, err := strconv.ParseFloat(q.Get(kunci), 64)
	return nilai, err == nil
}

// urutkanKatalog: produk tanpa harga yang terlihat selalu di akhir urutan harga. Selain itu tetap terbaru dulu.
func urutkanKatalog(produk []produkKatalog, urutan string) {
	kolasi := collate.New(language.Indonesian, collate.IgnoreCase, collate.Numeric)
	switch urutan {
	case "nama":
		sort.SliceStable(produk, func(i, j int) bool { return kolasi.CompareString(produk[i].row.Name, produk[j].row.Name) < 0 })
	case "harga_terendah", "harga_tertinggi":
		turun := urutan == "harga_tertinggi"
		sort.SliceStable(produk, func(i, j int) bool {
			a, b := produk[i], produk[j]
			// Sesama harga tersembunyi tidak dibandingkan: urutannya tidak boleh membocorkan harga.
			if !a.hargaTerlihat() || !b.hargaTerlihat() {
				return a.hargaTerlihat() && !b.hargaTerlihat()
			}
			if turun {
				return a.row.PriceGeneral.Float64 > b.row.PriceGeneral.Float64
			}
			return a.row.PriceGeneral.Float64 < b.row.PriceGeneral.Float64
		})
	}
}

func hitungFasetKatalog(dasar, tanpaHarga []produkKatalog, f filterKatalog, merek, ruangan []opsiKatalog) map[string]any {
	jumlah := func(kecuali string, cocok func(produkKatalog) bool) int {
		n := 0
		for _, p := range dasar {
			if f.lolos(p, kecuali) && cocok(p) {
				n++
			}
		}
		return n
	}

	fasetMerek := make([]map[string]any, 0, len(merek))
	for _, m := range merek {
		n := jumlah("brand", func(p produkKatalog) bool { return p.row.BrandSlug.String == m.Slug })
		if n > 0 || f.brand[m.Slug] {
			fasetMerek = append(fasetMerek, map[string]any{"slug": m.Slug, "name": m.Name, "count": n})
		}
	}

	fasetStok := make([]map[string]any, 0, len(statusStokKatalog))
	for _, status := range statusStokKatalog {
		n := jumlah("stock", func(p produkKatalog) bool { return p.row.StockStatus == status })
		fasetStok = append(fasetStok, map[string]any{"value": status, "count": n})
	}

	fasetRuangan := make([]map[string]any, 0, len(ruangan))
	for _, r := range ruangan {
		n := jumlah("room", func(p produkKatalog) bool { return p.ruangan[r.Slug] })
		if n > 0 || f.room[r.Slug] {
			fasetRuangan = append(fasetRuangan, map[string]any{"slug": r.Slug, "name": r.Name, "count": n})
		}
	}

	kolasi := collate.New(language.Indonesian, collate.IgnoreCase, collate.Numeric)
	fasetAtribut := make([]map[string]any, 0, len(atributFilterKatalog))
	for _, kunci := range atributFilterKatalog {
		nilaiUnik := nilaiAtribut(dasar, kunci)
		sort.Slice(nilaiUnik, func(i, j int) bool { return kolasi.CompareString(nilaiUnik[i], nilaiUnik[j]) < 0 })
		daftarNilai := make([]map[string]any, 0, len(nilaiUnik))
		for _, nilai := range nilaiUnik {
			n := jumlah("attr:"+kunci, func(p produkKatalog) bool { return p.spesifikasi[kunci] == nilai })
			if n > 0 || f.attr[kunci][nilai] {
				daftarNilai = append(daftarNilai, map[string]any{"value": nilai, "count": n})
			}
		}
		// Satu nilai saja tidak bisa dipakai menyaring, kecuali sedang dipilih (agar bisa dilepas).
		if len(daftarNilai) >= 2 || len(f.attr[kunci]) > 0 {
			fasetAtribut = append(fasetAtribut, map[string]any{"key": kunci, "values": daftarNilai})
		}
	}

	// Dihitung tanpa filter min/maks, supaya rentangnya tidak menyempit oleh filter harga itu sendiri.
	var rentang any
	var terendah, tertinggi float64
	ada := false
	for _, p := range tanpaHarga {
		if !p.hargaTerlihat() || !f.lolos(p, "") {
			continue
		}
		harga := p.row.PriceGeneral.Float64
		if !ada || harga < terendah {
			terendah = harga
		}
		if !ada || harga > tertinggi {
			tertinggi = harga
		}
		ada = true
	}
	if ada {
		rentang = map[string]any{"min": terendah, "max": tertinggi}
	}

	return map[string]any{"brands": fasetMerek, "stock": fasetStok, "rooms": fasetRuangan, "attributes": fasetAtribut, "price": rentang}
}

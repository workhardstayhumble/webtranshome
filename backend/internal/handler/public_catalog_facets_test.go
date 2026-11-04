package handler

import (
	"database/sql"
	"net/url"
	"reflect"
	"testing"
)

func produkUji(id int64, nama, merek, stok string, harga float64, tampil bool, spesifikasi map[string]string, ruangan ...string) produkKatalog {
	row := publicProductRow{ID: id, Name: nama, SKU: nama, StockStatus: stok, ShowPrice: tampil, CategorySlug: "keramik"}
	row.BrandSlug = sql.NullString{String: merek, Valid: merek != ""}
	row.BrandName = row.BrandSlug
	row.PriceGeneral = sql.NullFloat64{Float64: harga, Valid: harga > 0}
	rooms := map[string]bool{}
	for _, r := range ruangan {
		rooms[r] = true
	}
	return produkKatalog{row: row, spesifikasi: spesifikasi, ruangan: rooms}
}

func katalogUji() []produkKatalog {
	return []produkKatalog{
		produkUji(1, "Keramik A", "merek-a", "TERSEDIA", 100, true, map[string]string{"ukuran": "5 mm"}, "dapur"),
		produkUji(2, "Keramik B", "merek-a", "HABIS", 300, true, map[string]string{"ukuran": "1,5 mm"}, "teras"),
		produkUji(3, "Keramik C", "merek-b", "TERSEDIA", 200, true, map[string]string{"ukuran": "12 mm"}),
		produkUji(4, "Keramik D", "merek-b", "TERSEDIA", 999, false, map[string]string{"ukuran": "12 mm"}),
		produkUji(5, "Keramik E", "", "TERSEDIA", 0, false, nil),
	}
}

func fasetUji(q url.Values) map[string]any {
	semua := katalogUji()
	tanpaHarga, dasar := saringDasarKatalog(semua, q)
	merek := []opsiKatalog{{"Merek A", "merek-a"}, {"Merek B", "merek-b"}, {"Merek C", "merek-c"}}
	ruangan := []opsiKatalog{{"Dapur", "dapur"}, {"Teras", "teras"}}
	return hitungFasetKatalog(dasar, tanpaHarga, bacaFilterKatalog(q, semua), merek, ruangan)
}

func hitungOpsi(t *testing.T, faset map[string]any, grup, kunci string) map[string]int {
	t.Helper()
	hasil := map[string]int{}
	for _, opsi := range faset[grup].([]map[string]any) {
		hasil[opsi[kunci].(string)] = opsi["count"].(int)
	}
	return hasil
}

func TestFasetDisjungtif(t *testing.T) {
	faset := fasetUji(url.Values{"brand": {"merek-a"}})
	// Memilih Merek A tidak membuat angka Merek B menjadi 0, dan brand tanpa produk tidak dikirim.
	if got, want := hitungOpsi(t, faset, "brands", "slug"), (map[string]int{"merek-a": 2, "merek-b": 2}); !reflect.DeepEqual(got, want) {
		t.Fatalf("brands = %v, want %v", got, want)
	}
	// Grup lain mengikuti filter brand.
	if got, want := hitungOpsi(t, faset, "stock", "value"), (map[string]int{"TERSEDIA": 1, "SISA_STOK": 0, "PRE_ORDER": 0, "HABIS": 1}); !reflect.DeepEqual(got, want) {
		t.Fatalf("stock = %v, want %v", got, want)
	}
}

func TestFasetOpsiTerpilihTetapDikirim(t *testing.T) {
	faset := fasetUji(url.Values{"brand": {"merek-c"}})
	if got := hitungOpsi(t, faset, "brands", "slug"); got["merek-c"] != 0 || len(got) != 3 {
		t.Fatalf("brands = %v", got)
	}
}

func TestFasetAtributUrutAngkaDanBerkoma(t *testing.T) {
	faset := fasetUji(url.Values{})
	grup := faset["attributes"].([]map[string]any)
	if len(grup) != 1 || grup[0]["key"] != "ukuran" {
		t.Fatalf("attributes = %v", grup)
	}
	urutan := []string{}
	for _, nilai := range grup[0]["values"].([]map[string]any) {
		urutan = append(urutan, nilai["value"].(string))
	}
	if want := []string{"1,5 mm", "5 mm", "12 mm"}; !reflect.DeepEqual(urutan, want) {
		t.Fatalf("urutan = %v, want %v", urutan, want)
	}
}

func TestPilihanAtributMengenaliKomaDesimal(t *testing.T) {
	dikenal := []string{"5 mm", "1,5 mm", "12 mm"}
	tests := map[string][]string{
		"1,5 mm":       {"1,5 mm"},
		"5 mm":         {"5 mm"},
		"1,5 mm,12 mm": {"1,5 mm", "12 mm"},
		"12 mm,1,5 mm": {"12 mm", "1,5 mm"},
		"99 mm":        {"99 mm"},
	}
	for raw, want := range tests {
		got := pilihanAtribut(raw, dikenal)
		if len(got) != len(want) {
			t.Fatalf("%q: got %v, want %v", raw, got, want)
		}
		for _, nilai := range want {
			if !got[nilai] {
				t.Fatalf("%q: got %v, want %v", raw, got, want)
			}
		}
	}
}

func TestFilterAtributNilaiTakDikenalMenghasilkanKosong(t *testing.T) {
	semua := katalogUji()
	f := bacaFilterKatalog(url.Values{"attr[ukuran]": {"99 mm"}}, semua)
	for _, p := range semua {
		if f.lolos(p, "") {
			t.Fatalf("produk %d tidak boleh lolos", p.row.ID)
		}
	}
}

func TestRentangHargaMengabaikanHargaTersembunyiDanFilterHarga(t *testing.T) {
	faset := fasetUji(url.Values{"min": {"150"}})
	if got, want := faset["price"], (map[string]any{"min": 100.0, "max": 300.0}); !reflect.DeepEqual(got, want) {
		t.Fatalf("price = %v, want %v", got, want)
	}
}

func TestFilterHargaTidakMenyertakanHargaTersembunyi(t *testing.T) {
	_, dasar := saringDasarKatalog(katalogUji(), url.Values{"max": {"5000"}})
	for _, p := range dasar {
		if !p.hargaTerlihat() {
			t.Fatalf("produk %d berharga tersembunyi ikut filter harga", p.row.ID)
		}
	}
	if len(dasar) != 3 {
		t.Fatalf("jumlah = %d, want 3", len(dasar))
	}
}

func TestUrutHargaMenaruhHargaTersembunyiDiAkhir(t *testing.T) {
	for urutan, want := range map[string][]int64{
		"harga_terendah":  {1, 3, 2, 4, 5},
		"harga_tertinggi": {2, 3, 1, 4, 5},
	} {
		produk := katalogUji()
		urutkanKatalog(produk, urutan)
		got := []int64{}
		for _, p := range produk {
			got = append(got, p.row.ID)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s = %v, want %v", urutan, got, want)
		}
	}
}

func TestPublicItemTidakMembocorkanHargaTersembunyi(t *testing.T) {
	item := publicItem(produkUji(4, "Keramik D", "merek-b", "TERSEDIA", 999, false, nil).row)
	if item["price_general"] != nil || item["show_price"] != false {
		t.Fatalf("item = %v", item)
	}
	item = publicItem(produkUji(1, "Keramik A", "merek-a", "TERSEDIA", 100, true, nil).row)
	if item["price_general"] != 100.0 || item["show_price"] != true {
		t.Fatalf("item = %v", item)
	}
}

func TestPencarianKategoriIndukDanBrand(t *testing.T) {
	semua := katalogUji()
	semua[0].row.CategoryParentSlug = sql.NullString{String: "lantai", Valid: true}
	if _, dasar := saringDasarKatalog(semua, url.Values{"category": {"lantai"}}); len(dasar) != 1 || dasar[0].row.ID != 1 {
		t.Fatalf("kategori induk = %v", dasar)
	}
	if _, dasar := saringDasarKatalog(semua, url.Values{"q": {"MEREK-B"}}); len(dasar) != 2 {
		t.Fatalf("pencarian brand = %d, want 2", len(dasar))
	}
}

package handler

import (
	"strings"
	"testing"
)

func heroUji(ubah func(*bannerInput)) bannerInput {
	p := bannerInput{Title: "Hero", ImageURL: "/api/v1/uploads/a.png", Position: "HOME_HERO", StartDate: "2026-11-01"}
	if ubah != nil {
		ubah(&p)
	}
	return p
}

func TestValidateBannerHero(t *testing.T) {
	tests := []struct {
		name string
		ubah func(*bannerInput)
		want bool
	}{
		{"tanpa isian hero", nil, true},
		{"lengkap", func(p *bannerInput) {
			p.Subtitle, p.CTALabel, p.CTA2Label, p.CTA2URL, p.VideoURL = "Sub", "Katalog", "Promo", "/promo", "/api/v1/uploads/a.mp4"
		}, true},
		{"subjudul terlalu panjang", func(p *bannerInput) { p.Subtitle = strings.Repeat("a", 161) }, false},
		{"subjudul 160 karakter multibyte", func(p *bannerInput) { p.Subtitle = strings.Repeat("é", 160) }, true},
		{"label tombol terlalu panjang", func(p *bannerInput) { p.CTALabel = strings.Repeat("a", 31) }, false},
		{"label kedua tanpa tautan", func(p *bannerInput) { p.CTA2Label = "Promo" }, false},
		{"tautan kedua tanpa label", func(p *bannerInput) { p.CTA2URL = "/promo" }, false},
		{"tautan kedua http", func(p *bannerInput) { p.CTA2Label, p.CTA2URL = "Promo", "http://a.test" }, false},
		{"tautan kedua skema relatif", func(p *bannerInput) { p.CTA2Label, p.CTA2URL = "Promo", "//a.test" }, false},
		{"video javascript", func(p *bannerInput) { p.VideoURL = "javascript:alert(1)" }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateBanner(heroUji(tt.ubah)); got != tt.want {
				t.Fatalf("validateBanner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeroValuesDiabaikanSelainHero(t *testing.T) {
	p := heroUji(func(p *bannerInput) { p.Subtitle, p.CTALabel = "Sub", "Katalog" })
	if got := heroValues(p); got[0] != "Sub" || got[2] != "Katalog" || got[1] != nil {
		t.Fatalf("hero = %v", got)
	}
	p.Position = "HOME_SLIDER"
	for _, nilai := range heroValues(p) {
		if nilai != nil {
			t.Fatalf("posisi slider harus menyimpan NULL, got %v", nilai)
		}
	}
	// Isian hero yang tidak valid tidak menggagalkan posisi lain karena memang diabaikan.
	p.CTA2Label = "Tanpa tautan"
	if !validateBanner(p) {
		t.Fatal("posisi slider harus mengabaikan isian hero")
	}
}

func TestTokenPratinjauTidakIkutDiResponsPublik(t *testing.T) {
	a := articleRow{ID: 1, Type: "ARTIKEL", Title: "A", Slug: "a", Content: "isi", Status: "DRAFT", PreviewToken: "rahasia"}
	for _, detail := range []bool{true, false} {
		v := a.value(detail)
		if _, ada := v["preview_token"]; ada {
			t.Fatalf("value(%v) membocorkan preview_token: %v", detail, v)
		}
		for _, isi := range v {
			if isi == "rahasia" {
				t.Fatalf("value(%v) membocorkan token: %v", detail, v)
			}
		}
	}
}

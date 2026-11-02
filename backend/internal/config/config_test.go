package config

import (
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestWhatsAppNumberLoadedFromServerEnvironment(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test:unused@localhost/test")
	t.Setenv("JWT_SECRET", "test-secret-not-for-production")
	t.Setenv("QONTAK_BASE_URL", "")
	t.Setenv("QONTAK_API_TOKEN", "")
	for _, tc := range []struct {
		name    string
		number  string
		wantErr bool
	}{
		{name: "configured", number: "6281234567890"},
		{name: "minimum length", number: "6281234567"},
		{name: "maximum length", number: "62812345678901"},
		{name: "unset disables WhatsApp"},
		{name: "local prefix", number: "081234567890", wantErr: true},
		{name: "plus prefix", number: "+6281234567890", wantErr: true},
		{name: "spaces", number: "62812 34567890", wantErr: true},
		{name: "too short", number: "628123456", wantErr: true},
		{name: "too long", number: "628123456789012", wantErr: true},
		{name: "URL instead of number", number: "https://wa.me/6281234567890", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("WHATSAPP_NUMBER", tc.number)
			var c Config
			err := conf.Load("../../etc/transhome-api.yaml", &c, conf.UseEnv())
			if (err != nil) != tc.wantErr {
				t.Fatalf("Load() error=%v, wantErr=%t", err, tc.wantErr)
			}
			if err == nil && c.WhatsApp.Number != tc.number {
				t.Fatalf("WhatsApp.Number=%q, want %q", c.WhatsApp.Number, tc.number)
			}
		})
	}
}

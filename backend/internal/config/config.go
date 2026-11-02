package config

import (
	"errors"
	"regexp"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Database struct {
		DataSource string
	}

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	// Nomor dikelola operator server, bukan melalui database/API admin.
	WhatsApp struct {
		Number string `json:",optional"`
	} `json:",optional"`
}

var whatsAppNumberPattern = regexp.MustCompile(`^628[0-9]{7,11}$`)

func (c Config) Validate() error {
	// Kosong menonaktifkan tautan WhatsApp; tidak pernah kembali ke nomor di database.
	if c.WhatsApp.Number != "" && !whatsAppNumberPattern.MatchString(c.WhatsApp.Number) {
		return errors.New("WHATSAPP_NUMBER harus berformat 628 diikuti 7-11 digit, tanpa spasi atau tanda +")
	}
	return nil
}

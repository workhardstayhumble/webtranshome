# Backend Transhome

Backend website Transhome Muka 1 (Katalog & Membership).
Go + go-zero, PostgreSQL, arsitektur DDD modular monolith.

Penjelasan struktur dan aturan arsitektur: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Kebutuhan

- Go 1.24+
- PostgreSQL 14+
- goctl: `go install github.com/zeromicro/go-zero/tools/goctl@latest`
- golang-migrate: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

## Menjalankan pertama kali

```bash
cd backend
cp .env.example .env        # lalu isi DATABASE_URL dan lainnya
make migrate-up             # buat tabel di database transhome
make gen                    # generate handler, logic, types, dan transhome.go
go mod tidy                 # unduh dependency go-zero
make run                    # server jalan di http://localhost:8888
```

Cek: `GET http://localhost:8888/api/v1/health`

## Nomor WhatsApp dikelola server

- Sumber nomor toko adalah `WHATSAPP_NUMBER` di environment backend, dimuat melalui
  `etc/transhome-api.yaml`. Isi `backend/.env` di server dengan format `628` diikuti
  7–11 digit, tanpa `+` atau spasi. `make run` dan `run.ps1` memuat file `.env`;
  bila menjalankan binary langsung, environment harus dipasang oleh operator/service manager.
- Hanya Tim Cyber/operator yang mempunyai akses server yang mengganti konfigurasi ini.
  Restart backend setelah perubahan. Batasi akses file `.env`/secret manager dan jangan commit nilainya.
- Dashboard hanya menampilkan nomor. `PUT /api/v1/admin/settings` **tidak menerima**
  `wa_number`, termasuk dari Super Admin. Payload yang memuat field tersebut ditolak (422);
  tidak ada endpoint penggantian nomor. Alamat, jam buka, peta, dan izin rasio poin tetap berlaku.
- Nomor lama di `site_settings.wa_number` tidak lagi dibaca atau ditulis oleh aplikasi;
  tidak ada fallback ke database dan tidak perlu migration/penghapusan data lama.
- Jika `WHATSAPP_NUMBER` kosong, nomor pada respons settings kosong dan tombol WhatsApp
  nonaktif sampai operator mengisinya. Format nomor yang tidak valid ditolak saat startup.
- Nomor ini adalah **kontak publik toko** yang dipakai tombol WhatsApp pelanggan,
  bukan secret atau implementasi OTP/device whitelisting. Nomor keamanan internal
  untuk alert/OTP harus dipisahkan dan tidak dipublikasikan melalui settings publik.

## Perintah lain

| Perintah | Fungsi |
|---|---|
| `make worker` | Jalankan proses background |
| `make build` | Compile semua package |
| `make test` | Jalankan test |
| `make migrate-down` | Batalkan 1 migration terakhir |
| `make migrate-new name=xxx` | Buat file migration baru |

## Menambah endpoint baru

1. Tambahkan endpoint di `api/transhome.api` dengan `group` sesuai nama modul.
2. Jalankan `make gen`. File yang sudah ada tidak akan ditimpa.
3. Tulis use case di `internal/modules/<modul>/application/`.
4. Isi `internal/logic/<modul>/...` hanya untuk memanggil use case tersebut.

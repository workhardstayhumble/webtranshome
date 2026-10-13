
# ==========================================
# WebTranshome - Run Backend Server
# ==========================================

# 1. Pindah ke folder backend
Set-Location $PSScriptRoot

# 2. Periksa file .env
if (-not (Test-Path ".env")) {
    Write-Host "ERROR: File .env tidak ditemukan!" -ForegroundColor Red
    exit 1
}

# 3. Baca konfigurasi dari .env
Get-Content ".env" | ForEach-Object {
    $baris = $_.Trim()

    # Lewati baris kosong dan komentar
    if ($baris -eq "" -or $baris.StartsWith("#")) {
        return
    }

    if ($baris -match '^([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)$') {
        $nama = $matches[1]
        $nilai = $matches[2].Trim().Trim('"', "'")

        [Environment]::SetEnvironmentVariable(
            $nama,
            $nilai,
            "Process"
        )
    }
}

# 4. Pastikan DATABASE_URL tersedia
if (-not $env:DATABASE_URL) {
    Write-Host "ERROR: DATABASE_URL belum tersedia!" -ForegroundColor Red
    exit 1
}

# 5. Pastikan JWT_SECRET tersedia
if (-not $env:JWT_SECRET) {
    Write-Host "ERROR: JWT_SECRET belum tersedia!" -ForegroundColor Red
    exit 1
}

# 6. Jalankan backend Go
Write-Host ""
Write-Host "==========================================" -ForegroundColor Blue
Write-Host "   WebTranshome Backend Server" -ForegroundColor White
Write-Host "==========================================" -ForegroundColor Blue
Write-Host "Environment berhasil dimuat." -ForegroundColor Green
Write-Host "Menjalankan server pada localhost:8888..." -ForegroundColor Cyan
Write-Host "Tekan Ctrl + C untuk menghentikan server." -ForegroundColor Yellow
Write-Host ""

go run transhome.go -f etc/transhome-api.yaml

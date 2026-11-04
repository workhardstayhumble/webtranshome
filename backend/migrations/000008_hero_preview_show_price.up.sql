-- =====================================================================
-- Migration 000008: kolom yang dipakai kontrak frontend tetapi belum ada di migration.
-- - articles.preview_token: tautan pratinjau draf.
-- - banners.subtitle, video_url, cta_label, cta2_label, cta2_url: hero Beranda (HOME_HERO).
-- - products.show_price: harga boleh disembunyikan dan boleh kosong (KEPUTUSAN #15).
-- Sebagian database lokal sudah punya kolom artikel dan banner dari skrip manual,
-- jadi semua perintah dibuat aman dijalankan ulang.
-- =====================================================================
BEGIN;

ALTER TABLE crm_schema.articles ADD COLUMN IF NOT EXISTS preview_token VARCHAR(64);
UPDATE crm_schema.articles SET preview_token = replace(gen_random_uuid()::text, '-', '') WHERE preview_token IS NULL;
ALTER TABLE crm_schema.articles ALTER COLUMN preview_token SET DEFAULT replace(gen_random_uuid()::text, '-', '');
ALTER TABLE crm_schema.articles ALTER COLUMN preview_token SET NOT NULL;

ALTER TABLE crm_schema.banners
    ADD COLUMN IF NOT EXISTS subtitle VARCHAR(160),
    ADD COLUMN IF NOT EXISTS video_url VARCHAR(2048),
    ADD COLUMN IF NOT EXISTS cta_label VARCHAR(30),
    ADD COLUMN IF NOT EXISTS cta2_label VARCHAR(30),
    ADD COLUMN IF NOT EXISTS cta2_url VARCHAR(2048);

ALTER TABLE crm_schema.products ADD COLUMN IF NOT EXISTS show_price BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE crm_schema.products ALTER COLUMN price_general DROP NOT NULL;
ALTER TABLE crm_schema.products DROP CONSTRAINT IF EXISTS chk_products_shown_price_required;
ALTER TABLE crm_schema.products ADD CONSTRAINT chk_products_shown_price_required
    CHECK (NOT show_price OR price_general IS NOT NULL);

COMMIT;

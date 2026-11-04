-- Membatalkan migration 000008. Produk tanpa harga menjadi harga 0 karena kolomnya kembali wajib diisi.
BEGIN;

ALTER TABLE crm_schema.products DROP CONSTRAINT IF EXISTS chk_products_shown_price_required;
UPDATE crm_schema.products SET price_general = 0 WHERE price_general IS NULL;
ALTER TABLE crm_schema.products ALTER COLUMN price_general SET NOT NULL;
ALTER TABLE crm_schema.products DROP COLUMN IF EXISTS show_price;

ALTER TABLE crm_schema.banners
    DROP COLUMN IF EXISTS subtitle,
    DROP COLUMN IF EXISTS video_url,
    DROP COLUMN IF EXISTS cta_label,
    DROP COLUMN IF EXISTS cta2_label,
    DROP COLUMN IF EXISTS cta2_url;

ALTER TABLE crm_schema.articles DROP COLUMN IF EXISTS preview_token;

COMMIT;

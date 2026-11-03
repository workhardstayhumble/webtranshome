-- Rollback 000007: hapus izin, menu, dan tabel Jasa (data jasa ikut terhapus).
BEGIN;
DELETE FROM crm_schema.permissions WHERE code IN ('content.service.read','content.service.manage');
DELETE FROM crm_schema.menus WHERE code = 'content.service' AND NOT EXISTS (SELECT 1 FROM crm_schema.permissions p WHERE p.menu_id = crm_schema.menus.id);
DROP TABLE IF EXISTS crm_schema.service_categories;
DROP TABLE IF EXISTS crm_schema.services;
COMMIT;

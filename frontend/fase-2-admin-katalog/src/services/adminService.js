import http from '@/services/http'

/** Menu sidebar yang boleh dilihat user (sudah difilter backend sesuai izin): [{ id, parent_id, name, code, route, icon, sort_order }]. */
export async function getMenus() {
  if (import.meta.env.VITE_USE_MOCK === 'true') {
    const { data } = await http.get('/admin/menus')
    return data.data
  }

  // Backend sudah memfilter menu berdasarkan RBAC. Rute Go != rute Vue.
  const { data } = await http.get('/auth/me/access')
  const ruteVue = {
    '/admin/catalog/products': '/admin/produk',
    '/admin/catalog/categories': '/admin/kategori',
    '/admin/catalog/brands': '/admin/brand',
    '/admin/catalog/rooms': '/admin/ruangan',
    '/admin/produk': '/admin/produk',
    '/admin/kategori': '/admin/kategori',
    '/admin/brand': '/admin/brand',
    '/admin/ruangan': '/admin/ruangan',
    '/admin/member': '/admin/member',
    '/admin/banner': '/admin/banner',
    '/admin/artikel': '/admin/artikel',
    '/admin/sync-qontak': '/admin/sync-qontak',
    '/admin/log-aktivitas': '/admin/log-aktivitas',
    '/admin/role': '/admin/role',
  }
  return (data.menus ?? [])
    .map((menu) => ({
      id: menu.id,
      parent_id: menu.parentId || null,
      name: menu.name,
      code: menu.code,
      route: ruteVue[menu.route] ?? null,
      icon: menu.icon,
      sort_order: menu.sortOrder,
      is_active: menu.isActive,
    }))
    .filter((menu) => menu.is_active && (menu.route || menu.parent_id == null))
}

/** Ringkasan dashboard: { summary, attention, attention_total, logs, logs_scope }. */
export async function getDashboard() {
  const { data } = await http.get('/admin/dashboard')
  return data.data
}
import http from '@/services/http'

const MOCK = import.meta.env.VITE_USE_MOCK === 'true'
let katalogPromise = null

// Backend: GET /admin/access-control/catalog -> { roles, menus, permissions }
async function katalogAkses() {
  katalogPromise ??= http.get('/admin/access-control/catalog')
    .then(({ data }) => data)
    .catch((err) => { katalogPromise = null; throw err })
  return katalogPromise
}

export function refreshKatalogAkses() { katalogPromise = null }

export async function listRoles({ signal } = {}) {
  if (MOCK) {
    const { data } = await http.get('/admin/roles', { signal })
    return data.data
  }
  const catalog = await katalogAkses()
  const petaIzin = new Map(catalog.permissions.map((p) => [p.id, p.code]))
  return Promise.all(catalog.roles.map(async (role) => {
    const { data } = await http.get(`/admin/roles/${role.id}/permissions`, { signal })
    return {
      ...role,
      description: role.description ?? '',
      is_locked: role.code === 'SUPER_ADMIN',
      user_count: null, // Backend belum menyediakan agregat jumlah user per role.
      permissions: (data.permissionIds ?? []).map((id) => petaIzin.get(id)).filter(Boolean),
    }
  }))
}

function grupDariKode(code) {
  if (code.startsWith('catalog.')) return 'Katalog'
  if (code.startsWith('membership.')) return 'Membership'
  if (code.startsWith('member.')) return 'Portal Member'
  if (code.startsWith('content.')) return 'Konten'
  if (code.startsWith('integration.')) return 'Integrasi'
  return 'Sistem'
}

export async function listPermissions({ signal } = {}) {
  if (MOCK) {
    const { data } = await http.get('/admin/permissions', { signal })
    return data.data
  }
  const catalog = await katalogAkses()
  return catalog.menus
    .map((m) => {
      const permissions = catalog.permissions
        .filter((p) => p.menuId === m.id)
        .map((p) => ({ ...p, action: p.action === 'read' ? 'view' : p.action }))
      if (!permissions.length) return null
      const read = permissions.find((p) => p.action === 'view')
      return {
        menu_id: m.id,
        menu: m.name,
        group: grupDariKode(permissions[0].code),
        menu_permission: read?.code ?? null,
        permissions,
      }
    })
    .filter(Boolean)
}

export async function saveRolePermissions(roleId, codes) {
  if (MOCK) {
    const { data } = await http.put(`/admin/roles/${roleId}/permissions`, { permissions: codes })
    return data.data
  }
  const catalog = await katalogAkses()
  const peta = new Map(catalog.permissions.map((p) => [p.code, p.id]))
  if (codes.some((code) => !peta.has(code))) throw new Error('Ada permission yang tidak dikenal backend.')
  await http.put(`/admin/roles/${roleId}/permissions`, { permissionIds: codes.map((c) => peta.get(c)) })
  const role = catalog.roles.find((r) => r.id === roleId)
  return { ...role, name: role?.name ?? '', is_locked: role?.code === 'SUPER_ADMIN', permissions: codes }
}

export async function updateAdminRoles(id, roles) {
  if (MOCK) {
    const { data } = await http.put(`/admin/users/${id}/roles`, { roles })
    return data.data
  }
  const catalog = await katalogAkses()
  const peta = new Map(catalog.roles.map((r) => [r.code, r.id]))
  if (roles.some((role) => !peta.has(role))) throw new Error('Role tidak dikenal backend.')
  await http.put(`/admin/users/${id}/roles`, { roleIds: roles.map((r) => peta.get(r)) })
  return { id, roles: roles.map((code) => catalog.roles.find((r) => r.code === code)) }
}

// API-asli untuk CRUD user dan agregat user admin belum ada dalam ZIP.
// Jangan menampilkan operasi sukses palsu. Mode mock tetap tersedia untuk demo UI.
export async function listAdmins(params = {}, { signal } = {}) {
  if (!MOCK) throw new Error('Daftar admin belum tersedia pada API Go. Tersedia pencarian per email saja.')
  const { data } = await http.get('/admin/users', { params, signal })
  return data
}
export async function createAdmin(data) {
  if (!MOCK) throw new Error('Buat akun admin belum tersedia pada API Go.')
  const response = await http.post('/admin/users', data)
  return response.data.data
}
export async function setAdminActive(id, isActive) {
  if (!MOCK) throw new Error('Ubah status akun admin belum tersedia pada API Go.')
  const { data } = await http.patch(`/admin/users/${id}`, { is_active: isActive })
  return data.data
}
export async function resetAdminPassword(id) {
  if (!MOCK) throw new Error('Reset password admin belum tersedia pada API Go.')
  const { data } = await http.post(`/admin/users/${id}/reset-password`)
  return data.message
}
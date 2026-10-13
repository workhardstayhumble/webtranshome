import { useAuthStore } from '@/stores/auth'

// Translasi hanya untuk alias UI lama. Di mode API asli semua otorisasi berbasis
// kode `permissions.code` dari PostgreSQL; tidak ada bypass SUPER_ADMIN.
const ALIAS = {
  "category.view": "catalog.category.read",
  "category.create": "catalog.category.create",
  "category.update": "catalog.category.update",
  "brand.view": "catalog.brand.read",
  "brand.create": "catalog.brand.create",
  "brand.update": "catalog.brand.update",
  "room.view": "catalog.room.read",
  "room.create": "catalog.room.create",
  "room.update": "catalog.room.update",
  "product.view": "catalog.product.read",
  "product.create": "catalog.product.create",
  "product.update": "catalog.product.update",
  "product.price.update": "catalog.product.price.update",
  "product.stock.update": "catalog.product.stock.update",
  "product.featured": "catalog.featured.manage",
  "member.view": "membership.member.read",
  "member.create": "membership.member.create",
  "member.update": "membership.member.update",
  "point.view": "membership.points.read",
  "point.create": "membership.points.input",
  "point.redeem": "membership.points.redeem",
  "point.void": "membership.points.void",
  "point.rules": "membership.points.rules",
  "reward.manage": "membership.rewards.manage",
  "lead.view": "membership.leads.read",
  "banner.view": "content.banner.read",
  "banner.create": "content.banner.manage",
  "banner.update": "content.banner.manage",
  "article.view": "content.article.read",
  "article.create": "content.article.manage",
  "article.update": "content.article.manage",
  "review.moderate": "content.reviews.moderate",
  "promo.manage": "content.promo.manage",
  "setting.update": "content.pages.manage",
  "user.view": "rbac.manage",
  "user.create": "rbac.manage",
  "user.update": "rbac.manage",
  "role.view": "rbac.manage",
  "role.create": "rbac.manage",
  "role.update": "rbac.manage",
  "log.view": "audit.read",
  "qontak.retry": "integration.qontak.retry",
  "qontak.view": "integration.qontak.read",
  "qontak.configure": "integration.qontak.configure",
  "wishlist.manage": "member.wishlist.manage",
  "member.points": "member.points.read",
  "material.calculate": "member.material.calculate",
  "stock.alert": "member.stock_alert.manage",
  "review.create": "member.review.create"
}

export function kodeIzin(input) {
  return ALIAS[input] ?? input
}

export function boleh(input) {
  const auth = useAuthStore()
  if (!input || !auth.isLoggedIn) return false
  const kode = kodeIzin(input)
  // Backward compatibility bagi data mock; tidak berlaku pada backend asli.
  return auth.can(kode) || (import.meta.env.VITE_USE_MOCK === 'true' && (auth.can(input) || Object.entries(ALIAS).some(([lama, baru]) => baru === kode && auth.can(lama))))
}

export function usePermission() {
  return { can: boleh, kodeIzin }
}

function perbarui(el, binding) {
  el.hidden = !boleh(binding.value)
}
export const vCan = { mounted: perbarui, updated: perbarui }
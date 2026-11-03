// Data contoh akun, poin, dan pengaturan toko. Kolom mengikuti docs/transhome_postgres.sql. TODO: ganti ke API.
// Akun contoh HANYA untuk mode mock (VITE_USE_MOCK=true). Semua password: rahasia123.

export const MOCK_PASSWORD = 'rahasia123'

// password_hash: tidak ada = memakai MOCK_PASSWORD (akun contoh); null = belum punya password (masuk lewat Google);
// berisi = hash SHA-256 dari akun yang daftar/ganti password di mode mock (lihat routesAkun.js).
export const users = [
  {
    id: 101, email: 'member@example.com', phone_number: '6281200000001', is_active: true, roles: ['CUSTOMER'],
    email_verified_at: '2026-05-30T10:05:00', google_id: null, last_login_at: '2026-10-06T19:22:00',
    profile: { id: 1, member_code: 'TF-000123', full_name: 'Budi Santoso', tier: 'TAHAP_1', communication_consent: true, consent_at: '2026-05-30T10:00:00', created_at: '2026-05-30T10:00:00' },
  },
  {
    id: 102, email: 'siti.rahma@example.com', phone_number: '6281200000002', is_active: true, roles: ['CUSTOMER'],
    email_verified_at: '2026-10-05T19:20:00', google_id: null, last_login_at: '2026-10-05T19:25:00',
    profile: { id: 2, member_code: 'TF-000124', full_name: 'Siti Rahma', tier: 'TAHAP_1', communication_consent: true, consent_at: '2026-10-05T19:12:00', created_at: '2026-10-05T19:12:00' },
  },
  {
    // Member tanpa nomor HP (HP opsional, KEPUTUSAN #11).
    id: 103, email: 'andi.wijaya@example.com', phone_number: null, is_active: true, roles: ['CUSTOMER'],
    email_verified_at: '2026-10-06T08:45:00', google_id: null, last_login_at: '2026-10-06T08:50:00',
    profile: { id: 3, member_code: 'TF-000125', full_name: 'Andi Wijaya', tier: 'TAHAP_1', communication_consent: false, consent_at: null, created_at: '2026-10-06T08:40:00' },
  },
  { id: 1, email: 'superadmin@example.com', phone_number: null, is_active: true, roles: ['SUPER_ADMIN'], email_verified_at: '2026-09-01T08:00:00', google_id: null, profile: { full_name: 'Super Admin Contoh' } },
  { id: 2, email: 'katalog@example.com', phone_number: null, is_active: true, roles: ['ADMIN_KATALOG'], email_verified_at: '2026-09-01T08:00:00', google_id: null, profile: { full_name: 'Admin Katalog Contoh' } },
  { id: 3, email: 'membership@example.com', phone_number: null, is_active: true, roles: ['ADMIN_MEMBERSHIP'], email_verified_at: '2026-09-01T08:00:00', google_id: null, profile: { full_name: 'Admin Membership Contoh' } },
  { id: 4, email: 'konten@example.com', phone_number: null, is_active: true, roles: ['ADMIN_KONTEN'], email_verified_at: '2026-09-01T08:00:00', google_id: null, profile: { full_name: 'Admin Konten Contoh' } },
]

// Member contoh tambahan (id 104–125, TF-000126 dst.) supaya daftar member admin, filter, dan paginasi bisa dicoba.
// [nama, email, punya HP, setuju info promo, terdaftar, status sync Qontak, { aktif, verifikasi }]
export const MEMBER_TAMBAHAN = [
  ['Agus Prasetyo', 'agus.prasetyo', true, true, '2026-06-02T09:10:00', 'SYNCED'],
  ['Fitri Handayani', 'fitri.handayani', true, true, '2026-06-09T13:25:00', 'SYNCED'],
  ['Bambang Sutrisno', 'bambang.sutrisno', true, false, '2026-06-18T10:40:00', 'SYNCED'],
  ['Yuliana Kartika', 'yuliana.kartika', false, true, '2026-06-27T16:05:00', 'SYNCED'],
  ['Eko Saputra', 'eko.saputra', true, true, '2026-07-04T08:55:00', 'SYNCED'],
  ['Nur Aini', 'nur.aini', true, false, '2026-07-11T11:30:00', 'SYNCED'],
  ['Hendra Gunawan', 'hendra.gunawan', true, true, '2026-07-19T14:15:00', 'SYNCED', { aktif: false }],
  ['Ratna Sari', 'ratna.sari', false, false, '2026-07-26T09:45:00', 'SYNCED'],
  ['Joko Susilo', 'joko.susilo', true, true, '2026-08-02T15:20:00', 'SYNCED'],
  ['Maya Puspita', 'maya.puspita', true, true, '2026-08-09T10:10:00', 'SYNCED'],
  ['Wahyu Hidayat', 'wahyu.hidayat', false, true, '2026-08-16T13:00:00', 'SYNCED'],
  ['Indah Permatasari', 'indah.permatasari', true, false, '2026-08-23T17:35:00', 'SYNCED'],
  ['Teguh Santoso', 'teguh.santoso', true, true, '2026-08-30T08:20:00', 'SYNCED'],
  ['Lina Marlina', 'lina.marlina', true, true, '2026-09-06T12:45:00', 'SYNCED'],
  ['Rudi Hartono', 'rudi.hartono', false, false, '2026-09-13T09:05:00', 'SYNCED'],
  ['Sri Wahyuni', 'sri.wahyuni', true, true, '2026-09-20T14:50:00', 'SYNCED'],
  ['Dimas Pratama', 'dimas.pratama', true, false, '2026-09-24T10:30:00', 'SYNCED'],
  ['Ani Rahmawati', 'ani.rahmawati', true, true, '2026-09-27T16:40:00', 'SYNCED'],
  ['Fajar Nugroho', 'fajar.nugroho', true, true, '2026-10-01T11:15:00', 'FAILED'],
  ['Putri Ayuningtyas', 'putri.ayuningtyas', false, true, '2026-10-03T19:00:00', 'SYNCED'],
  ['Yoga Firmansyah', 'yoga.firmansyah', true, false, '2026-10-06T20:10:00', 'PENDING', { verifikasi: false }],
  ['Wulan Septiani', 'wulan.septiani', false, true, '2026-10-07T07:30:00', 'PENDING', { verifikasi: false }],
]

MEMBER_TAMBAHAN.forEach(([nama, email, punyaHp, setuju, dibuat, , { aktif = true, verifikasi = true } = {}], i) => {
  users.push({
    id: 104 + i, email: `${email}@example.com`, phone_number: punyaHp ? `62857123${String(4001 + i).padStart(5, '0')}` : null,
    is_active: aktif, roles: ['CUSTOMER'], email_verified_at: verifikasi ? dibuat : null, google_id: null, last_login_at: verifikasi ? dibuat : null,
    profile: {
      id: 4 + i, member_code: `TF-${String(126 + i).padStart(6, '0')}`, full_name: nama, tier: 'TAHAP_1',
      communication_consent: setuju, consent_at: setuju ? dibuat : null, created_at: dibuat,
    },
  })
})

const izinKatalog = ['product', 'category', 'brand', 'room'].flatMap((m) => [`${m}.view`, `${m}.create`, `${m}.update`])
const izinMembership = ['member.view', 'member.update', 'point.create', 'qontak.retry']
const izinKonten = [
  'banner.view', 'banner.create', 'banner.update', 'article.view', 'article.create', 'article.update', 'article.publish',
  'service.view', 'service.create', 'service.update', 'setting.update',
]

// Usulan pemetaan role → izin (lampiran docs/FRONTEND_BUILD_PLAN.md). Backend yang menentukan versi final.
export const rolePermissions = {
  SUPER_ADMIN: [...izinKatalog, 'brand.delete', ...izinMembership, 'point.adjust', ...izinKonten, 'setting.point_ratio', 'user.view', 'user.create', 'user.update', 'role.view', 'role.update', 'log.view'],
  ADMIN_KATALOG: [...izinKatalog, 'brand.delete'],
  ADMIN_MEMBERSHIP: izinMembership,
  ADMIN_KONTEN: izinKonten,
  CUSTOMER: [],
}

// Riwayat poin Budi (13 transaksi, saldo 24) supaya filter dan "Muat lebih banyak" bisa dicoba.
const riwayatBudi = [
  ['2026-06-05T10:12:00', 'EARN', 4, 'NT-2026-06-0117', '4480000.00', null],
  ['2026-06-14T15:30:00', 'EARN', 1, 'NT-2026-06-0301', '1250000.00', null],
  ['2026-06-20T09:00:00', 'ADJUST', -1, null, null, 'Koreksi nota ganda'],
  ['2026-06-28T13:45:00', 'EARN', 6, 'NT-2026-06-0533', '6900000.00', null],
  ['2026-07-12T11:20:00', 'EARN', 2, 'NT-2026-07-0215', '2600000.00', null],
  ['2026-07-20T16:05:00', 'REDEEM', -5, null, null, 'Tukar potongan harga di kasir'],
  ['2026-07-27T10:40:00', 'EARN', 5, 'NT-2026-07-0412', '5150000.00', null],
  ['2026-08-03T14:10:00', 'EARN', 3, 'NT-2026-08-0098', '3400000.00', null],
  ['2026-08-15T09:30:00', 'EARN', 2, 'NT-2026-08-0127', '2100000.00', null],
  ['2026-09-02T11:15:00', 'REDEEM', -10, null, null, 'Tukar voucher di kasir'],
  ['2026-09-21T16:40:00', 'EARN', 12, 'NT-2026-09-0311', '12800000.00', null],
  ['2026-09-28T10:05:00', 'ADJUST', 1, null, null, 'Koreksi pembulatan'],
  ['2026-10-05T14:20:00', 'EARN', 4, 'NT-2026-10-0456', '4250000.00', null],
]

// Transaksi member lain: [member_profiles.id, tanggal, jenis, poin, nota, total belanja, catatan].
const riwayatLain = [
  [2, '2026-10-05T19:40:00', 'EARN', 3, 'NT-2026-10-0431', '3600000.00', null],
  [4, '2026-06-10T11:05:00', 'EARN', 5, 'NT-2026-06-0188', '5300000.00', null],
  [4, '2026-08-21T10:20:00', 'EARN', 2, 'NT-2026-08-0164', '2750000.00', null],
  [5, '2026-07-02T14:35:00', 'EARN', 1, 'NT-2026-07-0027', '1100000.00', null],
  [7, '2026-07-15T09:50:00', 'EARN', 8, 'NT-2026-07-0298', '8450000.00', null],
  [7, '2026-09-10T15:10:00', 'REDEEM', -5, null, null, 'Tukar potongan harga di kasir'],
  [8, '2026-07-30T13:25:00', 'EARN', 2, 'NT-2026-07-0466', '2400000.00', null],
  [12, '2026-08-12T16:00:00', 'EARN', 15, 'NT-2026-08-0109', '15200000.00', null],
  [12, '2026-09-18T11:40:00', 'EARN', 4, 'NT-2026-09-0255', '4900000.00', null],
  [13, '2026-08-25T10:00:00', 'EARN', 3, 'NT-2026-08-0201', '3050000.00', null],
  [16, '2026-09-03T12:30:00', 'EARN', 6, 'NT-2026-09-0042', '6700000.00', null],
  [17, '2026-09-08T09:15:00', 'EARN', 2, 'NT-2026-09-0137', '2200000.00', null],
  [19, '2026-09-22T17:05:00', 'EARN', 9, 'NT-2026-09-0322', '9800000.00', null],
  [21, '2026-09-29T10:45:00', 'EARN', 1, 'NT-2026-09-0398', '1650000.00', null],
  [21, '2026-10-04T14:00:00', 'EARN', 3, 'NT-2026-10-0219', '3300000.00', null],
]

export const pointTransactions = [
  ...riwayatBudi.map(([created_at, type, points, reference_no, purchase_amount, note]) => ({ member_id: 1, created_at, type, points, reference_no, purchase_amount, note })),
  ...riwayatLain.map(([member_id, created_at, type, points, reference_no, purchase_amount, note]) => ({ member_id, created_at, type, points, reference_no, purchase_amount, note })),
].map((t, i) => ({ id: i + 1, ...t, created_by: t.type === 'ADJUST' ? 1 : 3 }))

// Mode mock tidak mengirim email sungguhan: email verifikasi & reset password masuk ke sini
// dan bisa dibuka di /dev/email. Token disimpan supaya link hanya bisa dipakai sekali.
export const emailKeluar = []
export const tokenVerifikasi = []
export const tokenReset = []

// site_settings yang boleh dibaca publik. Nomor WhatsApp di sini nomor contoh, bukan nomor toko.
export const publicSettings = {
  wa_number: '6280000000000',
  address: 'Jl. Raya Solo–Yogyakarta Km 11, Berbah, Sleman',
  opening_hours: 'Senin–Sabtu, 08.00–17.00',
  maps_url: 'https://www.google.com/maps/search/?api=1&query=Transhome+Berbah+Sleman',
  point_ratio_rupiah: 1000000,
}

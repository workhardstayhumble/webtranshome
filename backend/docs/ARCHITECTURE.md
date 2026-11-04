# Arsitektur Backend Transhome

Backend memakai **DDD modular monolith** di atas framework **go-zero**.
Satu aplikasi, satu database PostgreSQL (`transhome`, schema `crm_schema`),
dibagi menjadi beberapa modul (bounded context) yang tidak saling mengakses
tabel milik modul lain.

## Struktur folder

```
backend/
├── api/transhome.api          Definisi endpoint HTTP (sumber untuk goctl)
├── etc/transhome-api.yaml     Konfigurasi server
├── migrations/                File SQL golang-migrate (000001_init_schema, 000002_admin_refresh_sessions, 000003_rbac_access_management, 000004_catalog_permissions, 000006_role_capabilities, 000007_services, 000008_hero_preview_show_price)
├── cmd/worker/                Proses background (sync Qontak, retry, notifikasi)
├── transhome.go               Entry point API (dibuat oleh `make gen`)
├── docs/ARCHITECTURE.md       Dokumen ini
└── internal/
    ├── config/                Struct konfigurasi (dibuat goctl)
    ├── svc/                   ServiceContext: tempat merakit repository & service
    ├── middleware/            Adapter HTTP untuk autentikasi JWT, otorisasi, logging
    ├── types/                 Request/response HTTP (dibuat goctl dari .api)
    ├── handler/<modul>/       Handler HTTP (dibuat goctl)
    ├── logic/<modul>/         Penghubung HTTP -> use case (TIPIS, tanpa aturan bisnis)
    ├── shared/
    │   ├── domain/            Shared kernel: base error, domain event, value object umum
    │   ├── eventbus/          Bus event antar modul (in-process)
    │   ├── tx/                Transaksi / unit of work
    │   ├── mailer/            Kirim email
    │   └── storage/           Simpan foto produk & datasheet PDF
    └── modules/
        ├── identityaccess/    users, password_resets, admin_refresh_sessions, roles, menus, permissions, user_roles, role_permissions
        ├── membership/        member_profiles, point_transactions
        ├── catalog/           categories, brands, products (+ product_images, product_rooms), rooms
        ├── content/           articles, banners, site_settings
        ├── audit/             activity_logs
        └── integration/
            └── qontak/        qontak_contacts, sync_logs
```

Nama folder di `handler/` dan `logic/` **sama persis** dengan nama modul di
`modules/`, supaya mudah dicari.

## Lapisan di dalam setiap modul

```
modules/<modul>/
├── domain/
│   ├── <aggregate>/   Entity, value object, aturan bisnis, INTERFACE repository
│   └── service/       Domain service (aturan yang melibatkan lebih dari 1 aggregate)
├── application/
│   ├── command/       Use case yang mengubah data (create, update, assign)
│   ├── query/         Use case yang hanya membaca data
│   ├── dto/           Input/output use case
│   ├── listener/      Menangani domain event dari modul lain
│   └── port/          Interface ke modul lain / layanan luar yang dibutuhkan modul ini
└── infrastructure/
    ├── postgres/      Implementasi repository (GORM/SQL) + model tabel
    └── <adapter>/     Implementasi port: google, qontakapi, notifier, identityaccess
```

Folder `listener/` dan `port/` hanya dibuat di modul yang memang membutuhkannya.

## Alur request

```
Request publik:
HTTP -> handler -> logic -> application use case -> domain/repository port
                                                   ^                 |
                                                   |                 v
                                    infrastructure adapter <- PostgreSQL

Request terlindungi:
HTTP -> authentication middleware -> authorization middleware -> handler -> logic
             |                               |
             v                               v
       verifikasi JWT              application query/policy
                                             |
                                             v
                                    permission repository
```

## Authentication dan authorization

Authentication dan authorization adalah tanggung jawab berbeda:

- **Authentication** membuktikan identitas pemanggil. Login admin mencari akun,
  memverifikasi password terhadap hash bcrypt, lalu menerbitkan access token JWT.
- **Authorization** memutuskan apakah identitas yang sudah terautentikasi boleh
  menjalankan aksi tertentu. Keputusan dibuat per route/aksi berdasarkan permission,
  bukan hanya karena pengguna berhasil login atau memiliki token.

### Alur yang direncanakan untuk request terlindungi

1. Route menyatakan permission code yang diwajibkan, misalnya permission untuk
   membaca katalog. Permission code adalah identifier stabil dari konfigurasi/data,
   bukan teks pesan yang ditampilkan kepada pengguna.
2. Authentication middleware mengambil bearer token dan memverifikasi signature
   dengan algoritma yang diizinkan, expiration, serta klaim identitas. Middleware
   menaruh principal terverifikasi ke request context. Klaim dari token yang belum
   diverifikasi tidak boleh dipakai.
3. Authorization middleware mengambil principal dan permission yang diwajibkan
   route, lalu meminta keputusan kepada application query/policy Identity & Access.
4. Application layer memakai repository interface milik domain. Adapter PostgreSQL
   Identity & Access membaca relasi `user_roles`, `roles`, `role_permissions`, dan
   `permissions`, termasuk status aktif yang relevan.
5. Middleware melanjutkan request hanya jika permission diberikan; jika tidak,
   request ditolak. Ketiadaan/kegagalan autentikasi dibedakan dari penolakan akses.

Middleware adalah adapter HTTP: ia boleh membaca request dan menulis response,
tetapi tidak memuat query SQL atau aturan role/permission. Handler dan logic tetap
tipis. Aturan keputusan berada di application/domain; akses tabel hanya berada pada
adapter PostgreSQL dalam modul Identity & Access. Permission tidak dipercaya dari
body request atau klaim role yang dikirim client. Jika permission di-cache kelak,
aturan invalidasi cache harus memastikan perubahan role/permission berlaku sesuai
kebijakan keamanan.

### Kondisi implementasi saat ini

- Skema `users`, `roles`, `menus`, `permissions`, `user_roles`, dan
  `role_permissions` menjadi sumber data RBAC. Query aplikasi `CheckPermission`
  memakai kontrak repository domain; adapter PostgreSQL memeriksa user, role,
  menu, dan permission yang terhubung, termasuk status aktif user/role/menu.
- Authentication middleware memverifikasi bearer access token JWT HS256 dan
  menaruh user ID terverifikasi ke request context. Permission middleware
  memeriksa kode permission yang ditetapkan pada route; keputusan tidak dibuat
  dari role di middleware. `GET /api/v1/auth/me/access` mengembalikan menu dan
  permission efektif milik principal untuk kebutuhan UI; endpoint itu tidak
  menggantikan pemeriksaan permission pada setiap route. Endpoint katalog dan
  assignment RBAC berikut dilindungi permission `rbac.manage`:
  `GET /api/v1/admin/access-control/catalog`,
  `GET /api/v1/admin/users?email=...`,
  `POST /api/v1/admin/roles`, `PUT /api/v1/admin/roles/{roleId}`,
  `GET|PUT /api/v1/admin/users/{userId}/roles`, dan
  `GET|PUT /api/v1/admin/roles/{roleId}/permissions`.
- Migration `000003` menambahkan menu/permission pengelolaan akses dan memberi
  `rbac.manage` kepada role `SUPER_ADMIN` melalui `role_permissions`. API hanya
  menampilkan katalog menu/permission; penambahan katalog harus dilakukan
  melalui migration bersama route yang menerapkannya. API membuat role baru,
  mengubah nama/deskripsi/status role tanpa mengubah code, dan mengganti
  assignment user-role serta role-permission secara atomik. Role dinonaktifkan
  sebagai pengganti penghapusan; role `SUPER_ADMIN` tidak dapat dinonaktifkan
  melalui API.
- Kode role baru di-trim dan dinormalisasi ke huruf besar; format yang diterima
  adalah `^[A-Z][A-Z0-9_]*$` dengan panjang maksimal 50 karakter.
- Category Admin API menyediakan daftar, detail, pembuatan, dan partial update
  kategori melalui `/api/v1/admin/catalog/categories`. Route membutuhkan
  authentication dan permission `catalog.category.manage`, yang diberikan
  migration `000004` kepada `SUPER_ADMIN` dan `ADMIN_KATALOG`. Migration hanya
  menyemai data menu/permission dan relasi role; tidak mengubah struktur schema.
  Deaktivasi memakai `isActive: false`; `parentId: 0` pada PATCH menghapus
  relasi parent, `imageUrl: ""` menghapus URL gambar, sedangkan field yang
  dihilangkan tidak diubah.
- Brand Admin API menyediakan daftar, detail, pembuatan, dan partial update
  melalui `/api/v1/admin/catalog/brands`. Route membutuhkan authentication dan
  permission `catalog.brand.manage`, yang diberikan migration `000004` kepada
  `SUPER_ADMIN` dan `ADMIN_KATALOG`. Migration hanya menyemai data
  menu/permission dan relasi role; tidak mengubah struktur schema.
  Deaktivasi memakai `isActive: false`; `logoUrl: ""` menghapus URL logo,
  sedangkan field yang dihilangkan tidak diubah.
- Room Admin API menyediakan daftar, detail, pembuatan, dan partial update
  melalui `/api/v1/admin/catalog/rooms`. Route membutuhkan authentication dan
  permission `catalog.room.manage`, yang diberikan migration `000004` kepada
  `SUPER_ADMIN` dan `ADMIN_KATALOG`. Migration hanya menyemai data
  menu/permission dan relasi role; tidak mengubah struktur schema.
  Deaktivasi memakai `isActive: false`; `imageCover` wajib selalu memiliki
  nilai sesuai constraint schema, sedangkan field lain yang dihilangkan tidak
  diubah.
- Penggantian permission pada role `SUPER_ADMIN` ditolak jika akan menghapus
  permission `rbac.manage`. Ini adalah invariant untuk mencegah role bootstrap
  kehilangan izin, bukan mekanisme authorization untuk route.
- Login memeriksa email/password, status aktif, dan role staf/admin aktif
  (`SUPER_ADMIN`, `ADMIN_KATALOG`, `ADMIN_MEMBERSHIP`, atau `ADMIN_KONTEN`)
  sebelum menerbitkan token admin. Gerbang ini hanya menentukan kelayakan masuk
  ke aplikasi administratif; tidak memilih API atau aksi yang boleh diakses.
- Kebijakan token yang disepakati: access token JWT HS256 berlaku 15 menit dan
  refresh session admin berlaku absolut 7 hari. Migration `000002` menyiapkan
  tabel untuk menyimpan hash refresh token, bukan token mentah, serta mendukung
  rotasi dan revocation; runtime lifecycle belum diimplementasikan.
- Cookie refresh token, endpoint refresh/logout, dan penerbitan pasangan token
  belum diimplementasikan. Konfigurasi access token yang berjalan masih 24 jam.
- Health check dan login tetap publik. Endpoint fitur lain harus didaftarkan
  dengan authentication middleware dan permission spesifik; jangan menambah
  pengecualian Super Admin berdasarkan role di middleware.
- Jika akses RBAC terkunci, recovery dilakukan sebagai operasi manual terkontrol
  oleh operator database. Pastikan email admin tepercaya sebelum menjalankan
  SQL berikut, lalu ganti nilai placeholder:

  ```sql
  BEGIN;

  INSERT INTO crm_schema.role_permissions (role_id, permission_id)
  SELECT r.id, p.id
  FROM crm_schema.roles r
  CROSS JOIN crm_schema.permissions p
  WHERE r.code = 'SUPER_ADMIN'
    AND p.code = 'rbac.manage'
  ON CONFLICT (role_id, permission_id) DO NOTHING;

  INSERT INTO crm_schema.user_roles (user_id, role_id)
  SELECT u.id, r.id
  FROM crm_schema.users u
  CROSS JOIN crm_schema.roles r
  WHERE u.email = 'REPLACE_WITH_TRUSTED_ADMIN_EMAIL'
    AND u.is_active = TRUE
    AND r.code = 'SUPER_ADMIN'
    AND r.is_active = TRUE
  ON CONFLICT (user_id, role_id) DO NOTHING;

  COMMIT;
  ```

  Setelah operasi, verifikasi user-role dan role-permission yang terbentuk
  sebelum mengembalikan akses melalui API. SQL recovery ini tidak diekspos
  sebagai endpoint.

## 5 aturan wajib

1. **`domain/` tidak boleh meng-import** GORM, go-zero, `database/sql`, atau package
   `infrastructure` apa pun. Domain hanya berisi Go murni.
2. **Interface repository** ada di `domain/<aggregate>/`, **implementasinya** ada di
   `infrastructure/postgres/`.
3. **Aturan bisnis** (nota tidak boleh dobel, saldo poin harus cukup, role nonaktif
   tidak boleh di-assign) ditulis di entity atau domain service. Bukan di handler,
   logic, atau query SQL.
4. **`logic/` harus tipis**: ubah `types` -> `dto`, panggil use case, ubah hasilnya
   kembali ke `types`. Tidak ada `if` aturan bisnis di sini.
5. **Modul tidak boleh mengakses tabel milik modul lain.** Gunakan `port/` + adapter
   (contoh: `membership/infrastructure/identityaccess`) atau domain event lewat
   `shared/eventbus`.

## `types/` vs `application/dto/`

| | `internal/types/` | `modules/*/application/dto/` |
|---|---|---|
| Dibuat oleh | goctl dari `transhome.api` | Ditulis manual |
| Isi | Bentuk JSON request/response | Input/output use case |
| Boleh tahu HTTP? | Ya (tag `json`, `form`, `path`) | Tidak |

Dengan begitu use case tetap bisa dipakai dari worker atau test tanpa HTTP.

## Contoh alur antar modul (event)

```
identityaccess: RegisterUser  -> publish UserRegistered
membership:     listener      -> buat member_profiles (member_code, tier TAHAP_1)
qontak:         listener      -> qontak_contacts PENDING -> worker kirim ke API Qontak
audit:          listener      -> tulis activity_logs (untuk aksi admin)
```

## Kenapa `integration/` punya sub-folder?

Muka 2 (e-commerce) akan menambah integrasi lain, misalnya payment gateway dan
kurir. Semuanya diletakkan sejajar dengan `qontak/` di dalam `integration/`.

# Backend Coding Guidelines (Golang Clean Architecture)

Dokumen ini berisi aturan dan standar coding yang harus diikuti dalam pengembangan project ini.

## 1. General Rules
- **No `any` or `interface{}`**: Hindari penggunaan `any` kecuali benar-benar diperlukan (misalnya pada generic response). Gunakan struct yang jelas untuk setiap data.
- **Strict Typing**: Gunakan tipe data yang spesifik. Jangan gunakan string jika data tersebut adalah enum (gunakan `type MyType string` dan `const`).
- **Clean Architecture**: Ikuti pembagian layer:
    - `handler`: Hanya mengurus HTTP (parsing request, kirim response).
    - `service`: Tempat business logic. Tidak boleh tahu tentang HTTP context.
    - `repository`: Tempat query database. Tidak boleh ada business logic.
    - `dto`: Objek untuk transfer data antar layer atau ke client.

## 2. Error Handling
- Jangan pernah mengabaikan error (`_ = someFunc()`). Selalu handle error.
- Gunakan error wrapping jika ingin menambahkan context pada error.
- Response error ke client harus konsisten menggunakan `dto.WebResponse`.

## 3. Database
- Gunakan GORM untuk operasi database.
- Hindari logic di dalam repository yang terlalu kompleks. Jika butuh logic, pindahkan ke service.
- Gunakan `DeletedAt` (Soft Delete) untuk data-data penting.

## 4. Security
- Password wajib di-hash menggunakan `bcrypt`.
- Gunakan JWT untuk autentikasi. Token tidak boleh menyimpan data sensitif (hanya `user_id` dan `role`).
- Implementasikan Role-Based Access Control (RBAC) di middleware.

## 5. Documentation
- Setiap endpoint baru harus didaftarkan di Postman collection.
- Gunakan komentar pada fungsi yang memiliki logic kompleks.

## 6. Tech Stack
- Framework: Fiber v3
- Database: PostgreSQL (Local: SQLite)
- Auth: JWT
- Hashing: Bcrypt

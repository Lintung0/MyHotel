# 🏨 MyHotel API Documentation

## 📋 Daftar Isi
1. [Authentication](#authentication)
2. [Rooms](#rooms)
3. [Bookings](#bookings)
4. [Reviews](#reviews)
5. [Admin Management](#admin-management)

---

## 🔐 Authentication

### Register (Pendaftaran)
- **Endpoint:** `POST /api/auth/register`
- **Access:** Public
- **Request Body:**
```json
{
  "username": "john_doe",
  "password": "password123",
  "email": "john@example.com",
  "full_name": "John Doe"
}
```
- **Response Success (200):**
```json
{
  "success": true,
  "message": "Pendaftaran berhasil",
  "data": null
}
```

### Login (Masuk)
- **Endpoint:** `POST /api/auth/login`
- **Access:** Public
- **Request Body:**
```json
{
  "username": "john_doe",
  "password": "password123"
}
```
- **Response Success (200):**
```json
{
  "success": true,
  "message": "Login Berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "member"
    }
  }
}
```

---

## 🛏️ Rooms (Kamar)

### Get All Rooms (Lihat Semua Kamar)
- **Endpoint:** `GET /api/rooms`
- **Access:** Public
- **Query Parameters:**
  - `page` (optional, default: 1)
  - `limit` (optional, default: 10)
  - `sort` (optional, default: "created_at desc")
- **Response Success (200):**
```json
{
  "success": true,
  "message": "Berhasil mengambil data kamar",
  "data": {
    "rooms": [
      {
        "id": 1,
        "room_number": "101",
        "type": "Suite",
        "price": 500000,
        "description": "Kamar mewah dengan pemandangan laut",
        "status": "available",
        "max_occupancy": 4,
        "images": []
      }
    ],
    "page": 1,
    "limit": 10
  }
}
```

### Get Room Detail (Detail Kamar)
- **Endpoint:** `GET /api/rooms/:id`
- **Access:** Public
- **Response Success (200):** (sama seperti Get All Rooms, tapi untuk 1 kamar)

### Get Available Rooms (Kamar Tersedia)
- **Endpoint:** `POST /api/rooms/available`
- **Access:** Public
- **Request Body:**
```json
{
  "check_in_date": "2025-12-20",
  "check_out_date": "2025-12-25"
}
```
- **Query Parameters:**
  - `page` (optional)
  - `limit` (optional)
  - `sort` (optional)

---

## 📅 Bookings (Pemesanan)

### Create Booking (Buat Pemesanan)
- **Endpoint:** `POST /api/member/bookings`
- **Access:** Member Only
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "room_id": 1,
  "check_in_date": "2025-12-20",
  "check_out_date": "2025-12-25",
  "payment_method": "credit_card"
}
```
- **Response Success (201):**
```json
{
  "success": true,
  "message": "Pemesanan berhasil dibuat",
  "data": {
    "id": 1,
    "user_id": 1,
    "room_id": 1,
    "check_in_date": "2025-12-20",
    "check_out_date": "2025-12-25",
    "total_price": 2500000,
    "payment_method": "credit_card",
    "payment_status": "pending",
    "booking_status": "confirmed"
  }
}
```

### Get My Bookings (Lihat Pemesanan Saya)
- **Endpoint:** `GET /api/member/bookings`
- **Access:** Member Only
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `page` (optional)
  - `limit` (optional)
  - `sort` (optional)

### Cancel Booking (Batalkan Pemesanan)
- **Endpoint:** `DELETE /api/member/bookings/:id`
- **Access:** Member Only
- **Headers:** `Authorization: Bearer <token>`
- **Response Success (200):**
```json
{
  "success": true,
  "message": "Pemesanan berhasil dibatalkan",
  "data": null
}
```

---

## ⭐ Reviews (Ulasan)

### Create Review (Buat Ulasan)
- **Endpoint:** `POST /api/member/reviews`
- **Access:** Member Only
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "booking_id": 1,
  "rating": 5,
  "comment": "Kamar sangat bersih dan nyaman!"
}
```
- **Response Success (201):**
```json
{
  "success": true,
  "message": "Ulasan berhasil dibuat",
  "data": {
    "id": 1,
    "booking_id": 1,
    "user_id": 1,
    "rating": 5,
    "comment": "Kamar sangat bersih dan nyaman!"
  }
}
```

### Get Room Reviews (Lihat Ulasan Kamar)
- **Endpoint:** `GET /api/reviews/room/:roomId`
- **Access:** Public
- **Query Parameters:**
  - `page` (optional)
  - `limit` (optional)
  - `sort` (optional)

### Get Review Detail (Detail Ulasan)
- **Endpoint:** `GET /api/reviews/:id`
- **Access:** Public

---

## 👨‍💼 Admin Management

### Create Room (Buat Kamar Baru)
- **Endpoint:** `POST /api/admin/rooms`
- **Access:** Admin Only
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "room_number": "102",
  "type": "Suite",
  "price": 500000,
  "description": "Kamar mewah dengan pemandangan laut",
  "max_occupancy": 4
}
```

### Update Room (Ubah Data Kamar)
- **Endpoint:** `PUT /api/admin/rooms/:id`
- **Access:** Admin Only
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:** (semua field optional)
```json
{
  "room_number": "102",
  "type": "Suite",
  "price": 600000,
  "status": "available"
}
```

### Delete Room (Hapus Kamar)
- **Endpoint:** `DELETE /api/admin/rooms/:id`
- **Access:** Admin Only
- **Headers:** `Authorization: Bearer <token>`

### Add Room Image (Tambah Gambar Kamar)
- **Endpoint:** `POST /api/admin/rooms/:id/images`
- **Access:** Admin Only
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "image_url": "https://example.com/image.jpg",
  "is_primary": true
}
```

### Delete Room Image (Hapus Gambar Kamar)
- **Endpoint:** `DELETE /api/admin/rooms/:id/images/:imageId`
- **Access:** Admin Only
- **Headers:** `Authorization: Bearer <token>`

### Get All Bookings (Lihat Semua Pemesanan)
- **Endpoint:** `GET /api/admin/bookings`
- **Access:** Admin Only
- **Headers:** `Authorization: Bearer <token>`
- **Query Parameters:**
  - `page` (optional)
  - `limit` (optional)
  - `sort` (optional)

### Update Payment Status (Ubah Status Pembayaran)
- **Endpoint:** `PUT /api/admin/bookings/:id/payment-status`
- **Access:** Admin Only
- **Headers:** `Authorization: Bearer <token>`
- **Request Body:**
```json
{
  "payment_status": "paid"
}
```
- **Allowed Values:** `pending`, `paid`, `failed`

### Delete Review (Hapus Ulasan)
- **Endpoint:** `DELETE /api/admin/reviews/:id`
- **Access:** Admin Only
- **Headers:** `Authorization: Bearer <token>`

---

## 🔄 Status Booking & Payment

### Booking Status
- `confirmed` - Pemesanan dikonfirmasi
- `cancelled` - Pemesanan dibatalkan
- `completed` - Pemesanan selesai

### Payment Status
- `pending` - Menunggu pembayaran
- `paid` - Sudah dibayar
- `failed` - Pembayaran gagal

### Room Status
- `available` - Tersedia
- `booked` - Sudah dibooking
- `maintenance` - Dalam perbaikan

---

## 📝 Error Response

### Format Error Response
```json
{
  "success": false,
  "message": "Pesan error",
  "data": null
}
```

### Common HTTP Status Codes
- `200` - OK / Success
- `201` - Created / Resource berhasil dibuat
- `400` - Bad Request / Request tidak valid
- `401` - Unauthorized / Token tidak valid
- `403` - Forbidden / Tidak memiliki akses
- `404` - Not Found / Resource tidak ditemukan
- `409` - Conflict / Data sudah ada
- `500` - Internal Server Error / Error server

---

## 🚀 Cara Menjalankan Project

### Prerequisites
- Go 1.25+
- MySQL 8.0+

### Setup
1. Clone repository:
   ```bash
   git clone <repository-url>
   cd Backend
   ```

2. Setup database:
   ```bash
   mysql -u root -p < database-schema.sql
   ```

3. Create `.env` file (copy dari `.env.example`):
   ```bash
   cp .env.example .env
   ```

4. Update `.env` dengan konfigurasi lokal Anda

5. Run project:
   ```bash
   go run cmd/main.go
   ```

6. Atau build dan jalankan:
   ```bash
   go build -o myhotel cmd/main.go
   ./myhotel
   ```

Server akan berjalan di `http://localhost:8080`

---

## 📚 Technology Stack
- **Framework:** Fiber (Go web framework)
- **Database:** MySQL with GORM ORM
- **Authentication:** JWT (JSON Web Tokens)
- **Password Hashing:** bcrypt
- **Port Parsing:** dateparse

---

## ✅ File Structure Summary

```
Backend/
├── cmd/
│   └── main.go                 # Entry point aplikasi
├── internal/
│   ├── app/
│   │   └── services/          # Business logic
│   ├── domain/
│   │   ├── models/            # Data models
│   │   └── repositories/      # Repository interfaces
│   ├── infra/
│   │   ├── database/          # Database connection
│   │   ├── gorm/repositories/ # Repository implementation
│   │   └── http/
│   │       ├── handlers/      # API handlers
│   │       └── routes/        # Route definitions
│   └── config/                # Configuration
├── pkg/
│   └── utils/                 # Helper functions
└── go.mod                      # Go module definition
```

---

## 📌 Catatan Penting

1. **JWT Token:** Disertakan di setiap request protected menggunakan header:
   ```
   Authorization: Bearer <token>
   ```

2. **Role-based Access:** Ada 2 role:
   - `admin` - Akses ke semua management endpoints
   - `member` - Akses member features (booking, review)

3. **Database Models:** Semua model akan otomatis di-create saat app startup via AutoMigrate

4. **Pagination:** Gunakan query parameters `page`, `limit`, dan `sort`

---

Generated with ❤️

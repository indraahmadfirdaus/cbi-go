
# Go Todo List API - Vercel Edition

Aplikasi Todo List berbasis REST API yang dibangun dengan Go dan di-deploy di Vercel, dilengkapi dengan sistem autentikasi JWT dan manajemen pengguna.

## Fitur

- ✅ Sistem autentikasi dengan JWT
- ✅ Manajemen pengguna (registrasi, login)
- ✅ CRUD operasi untuk todo list
- ✅ Middleware autentikasi
- ✅ Validasi input
- ✅ Error handling yang komprehensif
- ✅ Serverless deployment di Vercel

## Teknologi yang Digunakan

- Go - Bahasa pemrograman utama
- JWT - Untuk autentikasi dan otorisasi
- bcrypt - Untuk hashing password
- net/http - HTTP server bawaan Go
- Vercel - Platform deployment serverless

## Live Demo

API ini tersedia di: `https://go-auth-h5cmnkod8-indrafrds-projects.vercel.app/`

## Struktur Project

```
go-auth-cbi/
├── api/
│   └── index.go     # Serverless function handler
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── vercel.json      # Vercel deployment config
└── README.md        # Dokumentasi ini
```

## API Documentation

### Base URL
```
https://go-auth-h5cmnkod8-indrafrds-projects.vercel.app/
```

### Authentication

Semua endpoint todo memerlukan autentikasi. Sertakan JWT token di header:
```
Authorization: Bearer <your-jwt-token>
```

---

## Endpoints

### 🔐 Authentication

#### Register User
```http
POST /api/register
```

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Login User
```http
POST /api/login
```

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

---

### ✅ Todo Management

#### Get All Todos
```http
GET /api/todos
```
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": 1,
    "title": "Complete project documentation",
    "description": "Write comprehensive README and API docs",
    "completed": false,
    "user_id": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

#### Create Todo
```http
POST /api/todos
```
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "New Todo Item",
  "description": "Description of the todo item"
}
```

**Response (201):**
```json
{
  "message": "Todo created successfully",
  "todo": {
    "id": 2,
    "title": "New Todo Item",
    "description": "Description of the todo item",
    "completed": false,
    "user_id": 1,
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

#### Get Todo by ID
```http
GET /api/todos/{id}
```
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": 1,
  "title": "Complete project documentation",
  "description": "Write comprehensive README and API docs",
  "completed": false,
  "user_id": 1,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

#### Update Todo
```http
PUT /api/todos/{id}
```
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "Updated Todo Title",
  "description": "Updated description",
  "completed": true
}
```

**Response (200):**
```json
{
  "message": "Todo updated successfully",
  "todo": {
    "id": 1,
    "title": "Updated Todo Title",
    "description": "Updated description",
    "completed": true,
    "user_id": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:30:00Z"
  }
}
```

#### Delete Todo
```http
DELETE /api/todos/{id}
```
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "Todo deleted successfully"
}
```

---

## Error Responses

Semua error response menggunakan format berikut:

```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes

- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

## Testing dengan cURL

### 1. Test API Status
```bash
curl https://go-auth-h5cmnkod8-indrafrds-projects.vercel.app/
```

### 2. Register User
```bash
curl -X POST https://go-auth-h5cmnkod8-indrafrds-projects.vercel.app/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Login
```bash
curl -X POST https://go-auth-h5cmnkod8-indrafrds-projects.vercel.app/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 4. Create Todo (dengan token)
```bash
curl -X POST https://go-auth-h5cmnkod8-indrafrds-projects.vercel.app/api/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "My First Todo",
    "description": "This is my first todo item"
  }'
```

### 5. Get All Todos
```bash
curl -X GET https://go-auth-h5cmnkod8-indrafrds-projects.vercel.app/api/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Default Test Account

Untuk testing, tersedia akun default:
- **Email:** `admin@example.com`
- **Password:** `password123`

## Vercel Configuration

Proyek ini menggunakan konfigurasi Vercel dengan:
- **Runtime:** Go 1.x serverless functions
- **Handler:** `/api/index.go`
- **Environment Variables:** `JWT_SECRET`

## Kontribusi

1. Fork repository
2. Buat feature branch (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

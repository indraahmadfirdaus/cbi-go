# Go Todo List API

Aplikasi Todo List berbasis REST API yang dibangun dengan Go, dilengkapi dengan sistem autentikasi JWT dan manajemen pengguna.

## Fitur

- ✅ Sistem autentikasi dengan JWT
- ✅ Manajemen pengguna (registrasi, login)
- ✅ CRUD operasi untuk todo list
- ✅ Middleware autentikasi
- ✅ Validasi input
- ✅ Error handling yang komprehensif

## Teknologi yang Digunakan

- **Go** - Bahasa pemrograman utama
- **JWT** - Untuk autentikasi dan otorisasi
- **bcrypt** - Untuk hashing password
- **net/http** - HTTP server bawaan Go

## Instalasi

### Prasyarat

- Go 1.19 atau lebih baru
- Git

### Langkah Instalasi

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd go-auth-cbi
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit file `.env` dan sesuaikan konfigurasi:
   ```env
   JWT_SECRET=your-super-secret-jwt-key-here
   PORT=8080
   ```

4. **Jalankan aplikasi**
   ```bash
   go run .
   ```

   Server akan berjalan di `http://localhost:8080`

## Struktur Project

```
go-auth-cbi/
├── main.go          # Entry point aplikasi
├── models.go        # Data models dan structs
├── handlers.go      # HTTP handlers
├── routes.go        # Route definitions
├── middleware.go    # Authentication middleware
├── utils.go         # Utility functions
├── .env.example     # Environment variables template
├── vercel.json      # Vercel deployment config
└── README.md        # Dokumentasi ini
```

## API Documentation

### Base URL
```
http://localhost:8080
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

### 👥 User Management

#### Get All Users
```http
GET /api/users
```
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

#### Get User by ID
```http
GET /api/users/{id}
```
**Headers:** `Authorization: Bearer <token>`

#### Update User
```http
PUT /api/users/{id}
```
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "username": "john_updated",
  "email": "john_updated@example.com"
}
```

#### Delete User
```http
DELETE /api/users/{id}
```
**Headers:** `Authorization: Bearer <token>`

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

### 1. Register User
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Create Todo (dengan token)
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "My First Todo",
    "description": "This is my first todo item"
  }'
```

### 4. Get All Todos
```bash
curl -X GET http://localhost:8080/api/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Deployment

### Vercel

Proyek ini sudah dikonfigurasi untuk deployment di Vercel:

1. Install Vercel CLI:
   ```bash
   npm i -g vercel
   ```

2. Deploy:
   ```bash
   vercel
   ```

3. Set environment variables di Vercel dashboard:
   - `JWT_SECRET`
   - `PORT` (optional, default: 8080)

## Kontribusi

1. Fork repository
2. Buat feature branch (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

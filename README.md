# Go Todo List API - Vercel Edition

A REST API-based Todo List application built with Go and deployed on Vercel, complete with JWT authentication system and user management.

## Features

- ✅ JWT authentication system
- ✅ User management (registration, login)
- ✅ CRUD operations for todo list
- ✅ Authentication middleware
- ✅ Input validation
- ✅ Comprehensive error handling
- ✅ Serverless deployment on Vercel

## Technologies Used

- Go - Main programming language
- JWT - For authentication and authorization
- bcrypt - For password hashing
- net/http - Go's built-in HTTP server
- Vercel - Serverless deployment platform

## Live Demo

This API is available at: `https://go-auth-cbi.vercel.app/`

## Project Structure

```
go-auth-cbi/
├── api/
│   └── index.go     # Serverless function handler
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── vercel.json      # Vercel deployment config
└── README.md        # This documentation
```

## API Documentation

### Base URL
```
https://go-auth-cbi.vercel.app/
```

### Authentication

All todo endpoints require authentication. Include JWT token in header:
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

All error responses use the following format:

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

## Testing with cURL

### 1. Test API Status
```bash
curl https://go-auth-cbi.vercel.app/
```

### 2. Register User
```bash
curl -X POST https://go-auth-cbi.vercel.app/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Login
```bash
curl -X POST https://go-auth-cbi.vercel.app/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 4. Create Todo (with token)
```bash
curl -X POST https://go-auth-cbi.vercel.app/api/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "My First Todo",
    "description": "This is my first todo item"
  }'
```

### 5. Get All Todos
```bash
curl -X GET https://go-auth-cbi.vercel.app/api/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Default Test Account

For testing purposes, a default account is available:
- **Email:** `admin@example.com`
- **Password:** `password123`

## Vercel Configuration

This project uses Vercel configuration with:
- **Runtime:** Go 1.x serverless functions
- **Handler:** `/api/index.go`
- **Environment Variables:** `JWT_SECRET`

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

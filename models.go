package main

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Todo structs
type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	UserID      int    `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}

var (
	users      = make(map[int]User)
	usersMux   = sync.RWMutex{}
	nextID     = 1
	todos      = make(map[int]Todo)
	todosMux   = sync.RWMutex{}
	nextTodoID = 1
	jwtKey     = []byte("your-secret-key")
)

func InitData() {
	usersMux.Lock()
	defer usersMux.Unlock()

	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	users[1] = User{ID: 1, Name: "John Doe", Email: "john@example.com", Password: string(hashedPassword1)}
	users[2] = User{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Password: string(hashedPassword2)}
	nextID = 3

	// Initialize sample todos
	todosMux.Lock()
	defer todosMux.Unlock()

	now := time.Now().Format(time.RFC3339)
	todos[1] = Todo{
		ID:          1,
		Title:       "Complete project documentation",
		Description: "Write comprehensive API documentation",
		Completed:   false,
		UserID:      1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	todos[2] = Todo{
		ID:          2,
		Title:       "Review code changes",
		Description: "Review pull requests from team members",
		Completed:   true,
		UserID:      1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	nextTodoID = 3
}

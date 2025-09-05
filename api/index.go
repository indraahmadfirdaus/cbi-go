package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// Todo struct
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Request structs
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// Global variables (in production, use database)
var (
	users     []User
	todos     []Todo
	nextUserID = 1
	nextTodoID = 1
	jwtSecret  = []byte("your-secret-key")
)

// Initialize data
func initData() {
	if len(users) == 0 {
		// Create sample users
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		users = append(users, User{
			ID:       nextUserID,
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
		})
		nextUserID++

		// Create sample todos
		todos = append(todos, Todo{
			ID:          nextTodoID,
			Title:       "Sample Todo",
			Description: "This is a sample todo",
			Completed:   false,
			UserID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		nextTodoID++
	}

	// Get JWT secret from environment
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = []byte(secret)
	}
}

// JWT functions
func generateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func validateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}

// CORS helper
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// Auth middleware
func requireAuth(w http.ResponseWriter, r *http.Request) (int, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return 0, false
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	userID, err := validateJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return 0, false
	}

	return userID, true
}

// Main handler function
func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize data
	initData()
	
	// Enable CORS
	enableCORS(w)
	
	// Handle OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	path := r.URL.Path
	method := r.Method
	
	// Route handling
	switch {
	case path == "/" && method == "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Todo List API is running!",
			"version": "1.0.0",
		})
		
	case path == "/api/register" && method == "POST":
		handleRegister(w, r)
		
	case path == "/api/login" && method == "POST":
		handleLogin(w, r)
		
	case path == "/api/todos" && method == "GET":
		handleGetTodos(w, r)
		
	case path == "/api/todos" && method == "POST":
		handleCreateTodo(w, r)
		
	case strings.HasPrefix(path, "/api/todos/") && method == "GET":
		handleGetTodoByID(w, r)
		
	case strings.HasPrefix(path, "/api/todos/") && method == "PUT":
		handleUpdateTodo(w, r)
		
	case strings.HasPrefix(path, "/api/todos/") && method == "DELETE":
		handleDeleteTodo(w, r)
		
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// Handler functions
func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Check if user exists
	for _, user := range users {
		if user.Email == req.Email {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	
	// Create user
	user := User{
		ID:       nextUserID,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	users = append(users, user)
	nextUserID++
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Find user
	var user *User
	for i := range users {
		if users[i].Email == req.Email {
			user = &users[i]
			break
		}
	}
	
	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	// Generate JWT
	token, err := generateJWT(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	
	var userTodos []Todo
	for _, todo := range todos {
		if todo.UserID == userID {
			userTodos = append(userTodos, todo)
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userTodos)
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	todo := Todo{
		ID:          nextTodoID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	todos = append(todos, todo)
	nextTodoID++
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func handleGetTodoByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	
	for _, todo := range todos {
		if todo.ID == id && todo.UserID == userID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todo)
			return
		}
	}
	
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	
	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	for i := range todos {
		if todos[i].ID == id && todos[i].UserID == userID {
			todos[i].Title = req.Title
			todos[i].Description = req.Description
			todos[i].Completed = req.Completed
			todos[i].UpdatedAt = time.Now()
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}
	
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	
	for i, todo := range todos {
		if todo.ID == id && todo.UserID == userID {
			todos = append(todos[:i], todos[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Todo deleted successfully",
			})
			return
		}
	}
	
	http.Error(w, "Todo not found", http.StatusNotFound)
}
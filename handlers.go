package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var registerReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	usersMux.Lock()
	defer usersMux.Unlock()

	// Check if email already exists
	for _, existingUser := range users {
		if existingUser.Email == registerReq.Email {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Create new user with auto-generated ID and hashed password
	user := User{
		ID:       nextID,
		Name:     registerReq.Name,
		Email:    registerReq.Email,
		Password: string(hashedPassword),
	}
	nextID++
	users[user.ID] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	usersMux.RLock()
	defer usersMux.RUnlock()

	var foundUser User
	var userFound bool

	// Check credentials with bcrypt
	for _, user := range users {
		if user.Email == loginReq.Email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
			if err == nil {
				foundUser = user
				userFound = true
				break
			}
		}
	}

	if !userFound {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: foundUser.ID,
		Email:  foundUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Password automatically excluded from JSON response due to json:"-" tag
	response := LoginResponse{
		Token: tokenString,
		User:  foundUser,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	usersMux.RLock()
	defer usersMux.RUnlock()

	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

// CreateUserHandler creates a new user
// @Summary Create a new user
// @Description Create a new user (requires authentication)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body RegisterRequest true "User data (without ID)"
// @Success 201 {object} User "User created successfully (password excluded)"
// @Failure 400 {object} map[string]string "Invalid JSON"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 409 {object} map[string]string "Email already exists"
// @Router /users [post]
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var registerReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	usersMux.Lock()
	defer usersMux.Unlock()

	// Check if email already exists
	for _, existingUser := range users {
		if existingUser.Email == registerReq.Email {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	user := User{
		ID:       nextID,
		Name:     registerReq.Name,
		Email:    registerReq.Email,
		Password: string(hashedPassword),
	}
	nextID++
	users[user.ID] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updateReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	usersMux.Lock()
	defer usersMux.Unlock()

	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Hash password if provided
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateReq.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	updatedUser := User{
		ID:       id,
		Name:     updateReq.Name,
		Email:    updateReq.Email,
		Password: string(hashedPassword),
	}
	users[id] = updatedUser

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	usersMux.RLock()
	defer usersMux.RUnlock()

	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	usersMux.Lock()
	defer usersMux.Unlock()

	if _, exists := users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(users, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromToken(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todosMux.RLock()
	defer todosMux.RUnlock()

	var userTodos []Todo
	for _, todo := range todos {
		if todo.UserID == userID {
			userTodos = append(userTodos, todo)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userTodos)
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromToken(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var createReq CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if createReq.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todosMux.Lock()
	defer todosMux.Unlock()

	now := time.Now().Format(time.RFC3339)
	newTodo := Todo{
		ID:          nextTodoID,
		Title:       createReq.Title,
		Description: createReq.Description,
		Completed:   false,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	todos[nextTodoID] = newTodo
	nextTodoID++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func GetTodoByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromToken(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todosMux.RLock()
	defer todosMux.RUnlock()

	todo, exists := todos[id]
	if !exists || todo.UserID != userID {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromToken(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var updateReq UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todosMux.Lock()
	defer todosMux.Unlock()

	todo, exists := todos[id]
	if !exists || todo.UserID != userID {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	if updateReq.Title != "" {
		todo.Title = updateReq.Title
	}
	if updateReq.Description != "" {
		todo.Description = updateReq.Description
	}
	if updateReq.Completed != nil {
		todo.Completed = *updateReq.Completed
	}
	todo.UpdatedAt = time.Now().Format(time.RFC3339)

	todos[id] = todo

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromToken(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todosMux.Lock()
	defer todosMux.Unlock()

	todo, exists := todos[id]
	if !exists || todo.UserID != userID {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	delete(todos, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// Helper function to get user ID from JWT token
func getUserIDFromToken(r *http.Request) int {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims.UserID
	}

	return 0
}

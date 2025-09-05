package main

import (
	"fmt"
	"net/http"
	"strings"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes() {
	http.HandleFunc("/", RouteHandler)
}

// RouteHandler handles all HTTP routes
func RouteHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	// Auth routes
	case path == "/api/register":
		RegisterHandler(w, r)
	case path == "/api/login":
		LoginHandler(w, r)

	// Todo routes
	case path == "/api/todos" && r.Method == "GET":
		AuthMiddleware(GetTodosHandler)(w, r)
	case path == "/api/todos" && r.Method == "POST":
		AuthMiddleware(CreateTodoHandler)(w, r)
	case strings.HasPrefix(path, "/api/todos/") && r.Method == "GET":
		AuthMiddleware(GetTodoByIDHandler)(w, r)
	case strings.HasPrefix(path, "/api/todos/") && r.Method == "PUT":
		AuthMiddleware(UpdateTodoHandler)(w, r)
	case strings.HasPrefix(path, "/api/todos/") && r.Method == "DELETE":
		AuthMiddleware(DeleteTodoHandler)(w, r)

	// Default route
	case path == "/":
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Todo List API is running!", "version": "1.0.0"}`)
	default:
		EnableCORS(w, r)
		if r.Method == "OPTIONS" {
			return
		}
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

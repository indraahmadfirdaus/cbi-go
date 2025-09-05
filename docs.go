package main

// @title Backend Authentication API
// @version 1.0
// @description This is a sample authentication server with JWT tokens.
// @description It provides user registration, login, and CRUD operations for users.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type \"Bearer\" followed by a space and JWT token.
package main

import "os"

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetPort() string {
	return GetEnv("PORT", "8080")
}

func GetJWTSecret() []byte {
	secret := GetEnv("JWT_SECRET", "your-secret-key")
	return []byte(secret)
}

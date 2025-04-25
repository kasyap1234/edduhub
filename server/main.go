// File: cmd/main.go
package main

import (
	"eduhub/server/api/app"
	"log" // Use standard log for fatal startup errors before custom logger is ready

	"github.com/joho/godotenv"
)

// @title           EduHub API
// @version         1.0
// @description     API for the EduHub platform.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080 // Change if needed
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization // Or X-Session-Token depending on your auth mechanism

func main() {
	// Load .env file FIRST
	err := godotenv.Load()
	if err != nil {
		// Log a warning, don't necessarily stop if it might run with existing env vars
		log.Println("Warning: Error loading .env file:", err)
	}

	// Create the app instance (which loads config, logger, db, etc.)
	setup := app.New()

	// Start the application
	err = setup.Start()
	if err != nil {
		// Use log.Fatalf for fatal errors during startup
		log.Fatalf("Failed to start server: %v", err)
	}
	// If Start() returns successfully (e.g., graceful shutdown), log it.
	log.Println("Server stopped gracefully.")
}

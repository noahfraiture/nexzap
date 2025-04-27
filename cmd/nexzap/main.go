package main

import (
	"fmt"
	"log"
	"net/http"
	"nexzap/internal/db"
	"nexzap/internal/handlers"
	"nexzap/internal/services"
	"os"
)

func main() {
	os.Setenv("APP_ENV", "dev")

	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Check database health
	if err := database.HealthCheck(); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	// Initialize services
	exerciseService, err := services.NewExerciseService()
	if err != nil {
		log.Fatalf("Failed to initialize exercise service: %v", err)
	}
	sheetService := services.NewSheetService(database)
	markdownService := services.NewMarkdownParser()
	importService := services.NewImportService(database)

	app := &handlers.App{
		Database:        database,
		ExerciseService: exerciseService,
		MarkdownService: markdownService,
		SheetService:    sheetService,
		ImportService:   importService,
	}

	// Nuke and populate the database (only in development)
	if err := database.NukeDatabase(); err != nil {
		log.Fatalf("Failed to nuke database: %v", err)
	}
	if err := database.Populate(); err != nil {
		log.Fatalf("Failed to populate database: %v", err)
	}
	if err := importService.RefreshTutorials(); err != nil {
		log.Fatalf("Failed to refresh tutorials: %v", err)
	}

	// Set up the router
	handlers.SetupRouter(app)

	// Start the server
	port := "8080"
	fmt.Println("Server running on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

package handlers

import (
	"net/http"
	"nexzap/internal/db"
	"nexzap/internal/services"
)

type App struct {
	Database        *db.Database
	ExerciseService *services.ExerciseService
	MarkdownService *services.MarkdownParser
	SheetService    *services.SheetService
	ImportService   *services.ImportService
	HistoryService  *services.HistoryService
}

func NewApp(
	database *db.Database,
	exerciseService *services.ExerciseService,
	markdownService *services.MarkdownParser,
	sheetService *services.SheetService,
	importService *services.ImportService,
	historyService *services.HistoryService,

) *App {
	return &App{
		Database:        database,
		ExerciseService: exerciseService,
		MarkdownService: markdownService,
		SheetService:    sheetService,
		ImportService:   importService,
		HistoryService:  historyService,
	}
}

// SetupRouter configures the HTTP router with handlers
func SetupRouter(app *App) {
	// Static file serving
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/images/favicon.ico")
	})

	http.HandleFunc("/", app.HomeHandler)
	http.HandleFunc("/sheet", app.SheetHandler)
	http.HandleFunc("/submit", app.SubmitHandler)

}

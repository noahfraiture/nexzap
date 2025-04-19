package main

import (
	"fmt"
	"log"
	"net/http"
	"nexzap/internal/db"
	"nexzap/internal/handlers"
	"nexzap/internal/services"
)

func router(s *services.Service) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/images/favicon.ico")
	})

	http.HandleFunc("/", handlers.HomeHandler())
	http.HandleFunc("GET /sheet", handlers.SheetHandler())
	http.HandleFunc("POST /submit", handlers.SubmitHandler(s))
	http.HandleFunc("POST /import", handlers.ImportHandler())
	http.HandleFunc("GET /refresh", handlers.RefreshHandler())
}

func main() {
	services.InitMarkdown()
	exerciseService, err := services.NewService()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = db.NukeDatabase()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Populate()
	if err != nil {
		log.Fatal(err)
	}

	router(exerciseService)
	port := "8080"
	fmt.Println("Server running on port " + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

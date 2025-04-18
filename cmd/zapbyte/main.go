package main

import (
	"fmt"
	"log"
	"net/http"
	"nexzap/internal/db"
	"nexzap/internal/handlers"
	"nexzap/internal/services"
)

func router() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/images/favicon.ico")
	})

	http.HandleFunc("/", handlers.HomeHandler())
	http.HandleFunc("/sheet", handlers.SheetHandler())
}

func main() {
	services.InitMarkdown()
	err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	router()
	port := "8080"
	fmt.Println("Server running on port " + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"zapbyte/internal/db"
	"zapbyte/internal/handlers"
	"zapbyte/internal/services"
)

func router() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/images/favicon.ico")
	})

	http.HandleFunc("/", handlers.HomeHandler())
	http.HandleFunc("/sheet", handlers.SheetHandler())

	// TODO : remove or protect
	http.HandleFunc("/api/tutorials", handlers.InsertTutorialHandler())
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

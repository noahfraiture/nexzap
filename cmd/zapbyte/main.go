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
	http.HandleFunc("/", handlers.HomeHandler())
	http.HandleFunc("/next", handlers.PageHandler())
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/images/favicon.ico")
	})
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

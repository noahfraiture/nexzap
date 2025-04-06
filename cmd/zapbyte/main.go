package main

import (
	"net/http"
	"zapbyte/internal/handlers"
)

func fromHtmx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { handlers.HomeHandler(fromHtmx(r), w, r) })
	http.HandleFunc("/tutorials", func(w http.ResponseWriter, r *http.Request) { handlers.TutorialHandler(fromHtmx(r), w, r) })

	port := "8080"
	println("Server running on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

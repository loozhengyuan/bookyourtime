package main

import (
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "title"}
	renderTemplate(w, "index.html", p)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	eventName := r.FormValue("eventName")

	// Generate Event
	err := createEvent(eventName)
	if err != nil {
		log.Fatal("Event creation failed")
	}

	// Output status
	w.Write([]byte("Success!"))
}

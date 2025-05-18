package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	handler := NewHandler()

	if err := handler.registerRoutes(mux); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const HOST = "localhost:7000"

func main() {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	server := &http.Server{
		Addr:    HOST,
		Handler: mux,
	}

	fmt.Printf("Listening on http://%v", HOST)

	serverErr := server.ListenAndServe()

	if serverErr != nil {
		log.Fatal("cannot start server", serverErr)
	}
}

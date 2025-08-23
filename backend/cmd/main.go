package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		curr := time.Now()
		fileName := curr.Format("2006-01-02-150405")
		w.Write([]byte(fileName))
	})

	// conn := connect()

	http.ListenAndServe(":8080", r)

}

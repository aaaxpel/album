package cmd

import (
	"net/http"

	images "github.com/aaaxpel/album/internal/routes/images"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/api/upload", images.UploadHandler)

	http.ListenAndServe(":8080", r)
}

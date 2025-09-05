package cmd

import (
	"log"
	"net/http"

	images "github.com/aaaxpel/album/internal/routes/images"
	"github.com/aaaxpel/album/internal/routes/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/debug", middleware.Profiler())

	r.Get("/api/image/:uuid", images.GetOneHandler)
	r.Post("/api/upload", images.UploadHandler)

	r.Post("/api/register", users.Register)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

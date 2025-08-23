package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func Router() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse form: %v\n", err)
			return
		}

		files := r.MultipartForm.File["file"]

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "Error retrieving file", http.StatusBadRequest)
				return
			}

			defer file.Close()

			// destination file
			name, _ := uuid.NewV7()
			dst, err := os.Create(filepath.Join("uploads", name.String()+filepath.Ext(fileHeader.Filename)))
			if err != nil {
				http.Error(w, "Error creating the file", http.StatusInternalServerError)
				fmt.Fprintf(os.Stderr, "Error creating the file: %v\n", err)
				return
			}

			defer dst.Close()

			// copy file contents
			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, "Error saving file", http.StatusInternalServerError)
				return
			}
		}

		w.Write([]byte("Files uploaded successfully!"))
	})

	http.ListenAndServe(":8080", r)
}

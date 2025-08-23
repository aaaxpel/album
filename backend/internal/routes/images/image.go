package images

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	db "github.com/aaaxpel/album/internal/db"
	"github.com/google/uuid"
)

func GetOneHandler(w http.ResponseWriter, r *http.Request) {

}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB RAM allocation
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

		// fmt.Println(fileHeader.Size)

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

	saveToDB()

	w.Write([]byte("Files uploaded successfully!"))
}

func saveToDB() {
	conn := db.Connect()
	conn.Ping(context.Background())
}

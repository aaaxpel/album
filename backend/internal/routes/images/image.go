package images

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/chai2010/webp"
	"github.com/google/uuid"
)

type FileError struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

func GetOneHandler(w http.ResponseWriter, r *http.Request) {

}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB RAM allocation
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse form: %v\n", err)
		return
	}

	files := r.MultipartForm.File["file"]

	const fileCount = 10

	jobs := make(chan *multipart.FileHeader, fileCount)
	errors := make(chan FileError, fileCount)

	var wg sync.WaitGroup
	for range fileCount {
		wg.Add(1)
		go worker(jobs, errors, &wg)
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(errors)
	}()

	var fileErrors []FileError

	for err := range errors {
		fileErrors = append(fileErrors, err)
	}

	w.Header().Set("Content-Type", "application/json")

	switch len(fileErrors) {
	case 0:
		w.WriteHeader(200) // Success
	case len(files):
		w.WriteHeader(400) // Bad request
	default:
		w.WriteHeader(207) // Partial success
	}

	json.NewEncoder(w).Encode(fileErrors)

	saveToDB()
}

func worker(jobs <-chan *multipart.FileHeader, errors chan<- FileError, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		err := processFile(job)
		if err != nil {
			errors <- FileError{job.Filename, err.Error()}
		}
	}
}

func processFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("error retrieving file")
	}

	defer file.Close()

	name, _ := uuid.NewV7()

	// Preview
	decodedFile, err := decodeImage(file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		switch err.Error() {
		case "invalid type":
			return fmt.Errorf("invalid file type")
		default:
			return fmt.Errorf("error decoding image: %v", err.Error())
		}
	}

	encodingErr := make(chan error)

	go func() {
		encodingErr <- encodeImage(name, decodedFile)
	}()

	if err := <-encodingErr; err != nil {
		return fmt.Errorf("error encoding image: %v", err)
	}

	// This is here just so I don't forget about database
	// fmt.Println(fileHeader.Size)

	// Resetting reader position to the beginning
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	// Original destination file
	output, err := os.Create(filepath.Join("uploads", "original", name.String()+filepath.Ext(fileHeader.Filename)))
	if err != nil {
		return fmt.Errorf("error creating the file: %v", err.Error())
	}

	defer output.Close()

	// Saving original file / Copying contents to output
	_, err = io.Copy(output, file)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err.Error())
	}

	return nil
}

func decodeImage(img multipart.File, contentType string) (image.Image, error) {
	var decodedImg image.Image
	var err error

	switch contentType {
	case "image/jpeg":
		decodedImg, err = jpeg.Decode(img)
	case "image/png":
		decodedImg, err = png.Decode(img)
	case "image/gif":
		decodedImg, err = gif.Decode(img)
	default:
		err := fmt.Errorf("invalid type")
		return nil, err
	}

	if err != nil {
		// Error decoding images
		return nil, err
	}

	return decodedImg, nil
}

func encodeImage(name uuid.UUID, img image.Image) error {
	var buf bytes.Buffer

	if err := webp.Encode(&buf, img, &webp.Options{Quality: 0.70}); err != nil {
		err := fmt.Errorf("encoding error: %v", err.Error())
		return err
	}

	output, err := os.Create(filepath.Join("uploads", "preview", name.String()+"_preview"+".webp"))
	if err != nil {
		err := fmt.Errorf("error creating the file: %v", err.Error())
		return err
	}

	// Saving the file
	_, err = io.Copy(output, &buf)
	if err != nil {
		log.Printf("Error saving file: %v", err.Error())
		return fmt.Errorf("error saving file: %v", err.Error())
	}

	return nil
}

func saveToDB() {
	// conn := db.Connect()
	// conn.Ping(context.Background())
}

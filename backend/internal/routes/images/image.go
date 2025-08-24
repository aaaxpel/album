package images

import (
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	db "github.com/aaaxpel/album/internal/db"
	"github.com/google/uuid"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
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

	const fileCount = 10

	jobs := make(chan *multipart.FileHeader, fileCount)
	errors := make(chan error, fileCount)

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

	// for _, fileHeader := range files {
	// 	go processFile(w, fileHeader)
	// }

	saveToDB()

	w.Write([]byte("Files uploaded successfully!"))
}

func worker(jobs <-chan *multipart.FileHeader, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		err := processFile(job)
		errors <- err
	}
}

func processFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		// error handling
		// return fmt.Errorf("Error retrieving file")
	}

	defer file.Close()

	name, _ := uuid.NewV7()

	// Preview
	decodedFile, err := decodeImage(file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		switch err.Error() {
		case "invalid type":
			// error handling
			// return fmt.Errorf("Invalid file type")
		default:
			fmt.Println(err.Error())
			// error handling
			// return fmt.Errorf("Error decoding image")
		}
	}

	go func() {
		err = encodeImage(name, decodedFile)
		if err != nil {
			// error handling
			return
		}
	}()

	// fmt.Println(fileHeader.Size)

	// Resetting reader position to the beginning
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	// Original destination file
	output, err := os.Create(filepath.Join("uploads", "original", name.String()+filepath.Ext(fileHeader.Filename)))
	if err != nil {
		// error handling
		fmt.Fprintf(os.Stderr, "Error creating the file: %v\n", err)
		// return
	}

	defer output.Close()

	// Saving original file / Copying contents to output
	_, err = io.Copy(output, file)
	if err != nil {
		// error handling
		// http.Error(w, "Error saving file", http.StatusInternalServerError)
		// return
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
		err = fmt.Errorf("error decoding image: %v", err.Error())
		return nil, err
	}

	return decodedImg, nil
}

func encodeImage(name uuid.UUID, img image.Image) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 70)
	if err != nil {
		err := fmt.Errorf("encoding error: %v", err.Error())
		return err
	}

	output, err := os.Create(filepath.Join("uploads", "preview", name.String()+"_preview"+".webp"))
	if err != nil {
		err := fmt.Errorf("error creating the file: %v", err.Error())
		return err
	}

	if err := webp.Encode(output, img, options); err != nil {
		err := fmt.Errorf("error encoding the file: %v", err.Error())
		return err
	}

	return nil
}

func saveToDB() {
	conn := db.Connect()
	conn.Ping(context.Background())
}

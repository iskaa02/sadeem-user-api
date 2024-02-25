package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/iskaa02/sadeem-user-api/api_error"
)

// return filename and error
func uploadFile(r *http.Request, id string) (string, error) {
	r.ParseMultipartForm(4 >> 20)
	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Printf("error opening form file: %s", err)
		return "", err
	}
	defer file.Close()

	// check if file is a png
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	Imagetype := http.DetectContentType(bytes)
	fmt.Println(Imagetype)
	if Imagetype != "image/png" {
		return "", api_error.NewBadRequestError("image_type_png_only", errors.New(""))
	}
	// save to images/
	filename := filepath.Join("images", id+".png")
	dst, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error creating destination file: %s", err)
		return "", err
	}
	defer dst.Close()
	_, err = dst.Write(bytes)
	if err != nil {
		fmt.Printf("error copying file: %s", err)
		return "", err
	}
	return filename, err
}

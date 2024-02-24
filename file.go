package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	if _, err = file.Read(buff); err != nil {
		fmt.Printf("error checking for file type: %s", err)
		return "", err
	}
	Imagetype := http.DetectContentType(buff)
	if Imagetype != "image/png" {
		return "", err
	}

	// save to images/
	filename := filepath.Join("images", id+".png")
	dst, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error creating destination file: %s", err)
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Printf("error copying file: %s", err)
		return "", err
	}
	return filename, err
}

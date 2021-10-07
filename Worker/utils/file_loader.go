package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(srcUrl, dstPath string) error {
	resp, err := http.Get(srcUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func UploadFile(srcPath, dstUrl string) (string, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, _ := w.CreateFormFile("file", filepath.Base(file.Name())) // TODO: export fieldname f.e. "file" to parameter
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	w.Close()
	resp, err := http.Post(dstUrl, w.FormDataContentType(), body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// return location header value if exists
	return resp.Header.Get("Location"), nil
}

func UploadImage(img image.Image, name, dstUrl string) (string, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, _ := w.CreateFormFile("file", name+".jpeg")
	err := jpeg.Encode(part, img, nil)
	if err != nil {
		return "", err
	}
	w.Close()
	resp, err := http.Post(dstUrl, w.FormDataContentType(), body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return resp.Header.Get("Location"), nil
}

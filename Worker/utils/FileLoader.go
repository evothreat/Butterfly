package utils

import (
	"bytes"
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
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name())) // TODO: export fieldname to parameter
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	writer.Close()
	resp, err := http.Post(dstUrl, writer.FormDataContentType(), body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// return location header value if exists
	return resp.Header.Get("Location"), nil
}

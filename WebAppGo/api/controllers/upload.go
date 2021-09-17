package controllers

import (
	"WebAppGo/api"
	"WebAppGo/api/models"
	"WebAppGo/utils"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TODO: before creating and copying the file check whether worker exists!

func CreateUpload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil || !utils.IsValidFilename(file.Filename) {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, fileExt), time.Now().Unix(), fileExt)
	filePath := filepath.Join(api.UPLOADS_DIR, c.Param("wid"), fileName)
	fileType := "NONE"
	if fileExt != "" {
		fileType = strings.ToUpper(strings.TrimPrefix(fileExt, "."))
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	fileSize, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO uploads(filename,type,size,created,worker_id) VALUES(?,?,?,?,?)",
		fileName, fileType, fileSize, time.Now(), c.Param("wid"))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func GetUpload(c echo.Context) error {
	uploadId, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	var upload models.Upload
	row := db.QueryRow("SELECT * FROM uploads WHERE worker_id=? AND id=?", c.Param("wid"), uploadId)
	if err := upload.Scan(row); err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	filePath := filepath.Join(api.UPLOADS_DIR, c.Param("wid"), upload.Filename) // TODO: handle file not found error
	if _, ok := c.QueryParams()["attach"]; ok {
		return c.Attachment(filePath, upload.Filename)
	}
	return c.File(filePath)
}

func DeleteUpload(c echo.Context) error {
	uploadId, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	var upload models.Upload
	row := db.QueryRow("SELECT * FROM uploads WHERE worker_id=? AND id=?", c.Param("wid"), uploadId)
	if err := upload.Scan(row); err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	_, err = db.Exec("DELETE FROM uploads WHERE worker_id=? AND id=?", c.Param("wid"), uploadId)
	if err != nil {
		return err
	}
	os.Remove(filepath.Join(api.UPLOADS_DIR, c.Param("wid"), upload.Filename)) // TODO: check for error?
	return c.NoContent(http.StatusOK)
}

func GetUploadInfo(c echo.Context) error {
	uploadId, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	if uploadId == 0 {
		rows, err := db.Query("SELECT * FROM uploads WHERE worker_id=?", c.Param("wid"))
		if err != nil {
			return err
		}
		uploads := make([]*models.Upload, 0, 10) // TODO: examine the number of rows first!
		for rows.Next() {
			u := &models.Upload{}
			if err := u.Scan(rows); err != nil {
				return err
			}
			uploads = append(uploads, u)
		}
		return c.JSON(http.StatusOK, &uploads)
	}
	var upload models.Upload
	row := db.QueryRow("SELECT * FROM uploads WHERE worker_id=? AND id=?", c.Param("wid"), uploadId)
	if err := upload.Scan(row); err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	return c.JSON(http.StatusOK, &upload)
}

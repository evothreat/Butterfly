package api

import (
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

func CreateUpload(c echo.Context) error {
	workerId := c.Param("wid")
	file, err := c.FormFile("file")
	if err != nil || !utils.IsValidFilename(file.Filename) {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, fileExt), time.Now().Unix(), fileExt)
	filePath := filepath.Join(UPLOADS_DIR, workerId, fileName)
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
	res, err := db.Exec("INSERT INTO uploads(filename,type,size,created,worker_id) VALUES(?,?,?,?,?)",
		fileName, fileType, fileSize, time.Now(), workerId)
	if err != nil {
		if isNoReferencedRowErr(err) {
			// would be good if we check this at begin
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	uploadId, _ := res.LastInsertId()
	// TODO: avoid raw url strings
	c.Response().Header().Set(echo.HeaderLocation, fmt.Sprintf("/api/workers/%s/uploads/%d", workerId, uploadId))
	return c.NoContent(http.StatusCreated)
}

func GetUpload(c echo.Context) error {
	workerId := c.Param("wid")
	uploadId, _ := strconv.Atoi(c.Param("uid"))
	var upload models.Upload
	row := db.QueryRow("SELECT * FROM uploads WHERE worker_id=? AND id=?", workerId, uploadId)
	if err := upload.Scan(row); err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	filePath := filepath.Join(UPLOADS_DIR, workerId, upload.Filename) // handle file not found error?
	if _, ok := c.QueryParams()["attach"]; ok {
		return c.Attachment(filePath, upload.Filename)
	}
	return c.File(filePath)
}

func DeleteUpload(c echo.Context) error {
	workerId := c.Param("wid")
	uploadId, _ := strconv.Atoi(c.Param("uid"))
	var upload models.Upload
	row := db.QueryRow("SELECT * FROM uploads WHERE worker_id=? AND id=?", workerId, uploadId)
	err := upload.Scan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}
	_, err = db.Exec("DELETE FROM uploads WHERE worker_id=? AND id=?", workerId, uploadId)
	if err != nil {
		return err
	}
	os.Remove(filepath.Join(UPLOADS_DIR, workerId, upload.Filename)) // TODO: check for error?
	return c.NoContent(http.StatusOK)
}

func GetUploadInfo(c echo.Context) error {
	uploadId, _ := strconv.Atoi(c.Param("uid"))
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

func GetAllUploadsInfo(c echo.Context) error {
	rows, err := db.Query("SELECT * FROM uploads WHERE worker_id=?", c.Param("wid"))
	if err != nil {
		return err
	}
	uploads := make([]*models.Upload, 0, MIN_LIST_CAP)
	for rows.Next() {
		u := &models.Upload{}
		if err := u.Scan(rows); err != nil {
			return err
		}
		uploads = append(uploads, u)
	}
	return c.JSON(http.StatusOK, &uploads)
}

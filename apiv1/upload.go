package apiv1

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v6"
	"github.com/mndyu/localchat-server/config"
)

type uploadResult struct {
	FileName string
	Url      string
}

func saveFile(fileName string, src io.Reader) (string, error) {
	var buf bytes.Buffer
	r := io.TeeReader(src, &buf)

	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	fileHash := fmt.Sprintf("%x", h.Sum(nil))
	ext := filepath.Ext(fileName)
	hashedFileName := fileHash + ext
	fileDir := config.PublicDirectory
	filePath := filepath.Join(fileDir, hashedFileName)

	// Destination
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, &buf); err != nil {
		return "", err
	}

	urlPath := path.Join(config.PublicPrefix, hashedFileName)
	return urlPath, nil
}

func uploadFile(fileName string, src io.Reader, size int64) error {
	endpoint := "localhost:9000"
	accessKeyID := "minio"
	secretAccessKey := "minio123"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	bucketName := "mymusic"
	location := "us-east-1"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			return err
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the zip file
	objectName := fileName
	// filePath := "static/index.html"
	contentType := "text/html"

	// Upload the zip file with FPutObject
	// n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	n, err := minioClient.PutObject(bucketName, objectName, src, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
	fmt.Printf("Successfully uploaded %s of size %d\n", objectName, n)

	return nil
}

func PostUpload(context *Context, c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("upload failed %s", err.Error()))
	}
	files := form.File["files"]

	var result []uploadResult
	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("upload failed %s", err.Error()))
		}
		defer src.Close()

		fmt.Println("upload:", file.Filename)
		urlPath, err := saveFile(file.Filename, src)
		// urlPath, err = uploadFile(file.Filename, src, file.Size)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("upload failed %s", err.Error()))
		}
		result = append(result, uploadResult{Url: urlPath})
	}
	if len(result) == 1 {
		return c.JSON(http.StatusOK, result[0])
	} else {
		return c.JSON(http.StatusOK, result)
	}
}

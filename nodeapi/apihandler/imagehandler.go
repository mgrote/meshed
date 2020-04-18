package apihandler

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const maxUploadSize = 200 * 1024 * 1024 // 200 mb
const uploadPath = "./tmp"

func UploadFileHandler(writer http.ResponseWriter, request *http.Request) {
	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		request.Body = http.MaxBytesReader(writer, request.Body, maxUploadSize)
		if err := request.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(writer, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		// parse and validate file and post parameters
		file, _, err := request.FormFile("uploadFile")
		if err != nil {
			renderError(writer, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(writer, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		detectedFileType := http.DetectContentType(fileBytes)
		switch detectedFileType {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
		case "application/pdf":
			break
		default:
			renderError(writer, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileName := request.FormValue("filename")
		if len(fileName) == 0 { // empty string, delivered if no form value 'filename'
			fileName = randToken(12)
		}
		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		if err != nil {
			renderError(writer, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
		newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
		fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(writer, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(writer, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		renderSucces(writer, "SUCCESS")
	//})
}

func renderSucces(writer http.ResponseWriter, message string) {
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte(message))
	if err != nil {
		log.Fatal("Error while writing respose")
	}
}

func renderError(writer http.ResponseWriter, message string, statusCode int) {
	writer.WriteHeader(statusCode)
	_, err := writer.Write([]byte(message))
	if err != nil {
		log.Fatal("Error while writing respose")
	}
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

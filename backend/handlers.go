package main

import (
	"encoding/json"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ansel1/merry"
	uuid "github.com/satori/go.uuid"
)

type Handlers struct {
	FolderPath   string
	FileInfoByID map[string]*FileInfo
	Journal      *Journal
}

type FileInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
}

func (handlers *Handlers) UploadFiles(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(merry.Details(err))
			http.Error(responseWriter, err.Error(), merry.HTTPCode(err))
			return
		}
	}()
	mediaType, params, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil {
		err = merry.WithHTTPCode(err, http.StatusBadRequest)
		return
	}
	if mediaType != "multipart/form-data" {
		err = merry.New("Content-Type is not multipart/form-data").WithHTTPCode(http.StatusBadRequest)
		return
	}
	reader := multipart.NewReader(request.Body, params["boundary"])
	var fileInfos []*FileInfo
	for {
		var part *multipart.Part
		part, err = reader.NextPart()
		switch err {
		case nil:
		case io.EOF:
			json.NewEncoder(responseWriter).Encode(fileInfos)
			err = nil
			return
		default:
			err = merry.WithHTTPCode(err, http.StatusBadRequest)
			return
		}
		err = handlers.CheckPartIsFile(part)
		if err != nil {
			return
		}
		var fileInfo *FileInfo
		fileInfo, err = handlers.SavePartToFileSystem(part)
		if err != nil {
			return
		}
		fileInfos = append(fileInfos, fileInfo)
		part.Close()
	}
}

func (handlers *Handlers) CheckPartIsFile(part *multipart.Part) (err error) {
	if part.FileName() == "" {
		err = merry.New("part is not a file")
		err = merry.WithHTTPCode(err, http.StatusBadRequest)
		return
	}
	return
}

func (handlers *Handlers) SavePartToFileSystem(part *multipart.Part) (fileInfo *FileInfo, err error) {
	fileID := uuid.NewV4().String()
	name := filepath.Join(handlers.FolderPath, fileID)
	file, err := os.Create(name)
	if err != nil {
		err = merry.WithHTTPCode(err, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	_, err = io.Copy(file, part)
	if err != nil {
		err = merry.WithHTTPCode(err, http.StatusInternalServerError)
		return
	}
	fileInfo = &FileInfo{
		ID:       fileID,
		Name:     part.FileName(),
		MimeType: part.Header.Get("Content-Type"),
	}
	handlers.Journal.AddFile(fileInfo)
	return
}

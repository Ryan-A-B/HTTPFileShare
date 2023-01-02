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
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Handlers struct {
	FolderPath string
	FileInfos  []*FileInfo
	Journal    *Journal
}

type FileInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
}

func (handlers *Handlers) AddFile(responseWriter http.ResponseWriter, request *http.Request) {
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
	part, err := reader.NextPart()
	if err != nil {
		err = merry.WithHTTPCode(err, http.StatusBadRequest)
		return
	}
	err = handlers.CheckPartIsFile(part)
	if err != nil {
		return
	}
	defer part.Close()
	var fileInfo *FileInfo
	fileInfo, err = handlers.SavePartToFileSystem(part)
	if err != nil {
		return
	}
	handlers.Journal.AddFile(fileInfo)
	handlers.FileInfos = append(handlers.FileInfos, fileInfo)
	_, err = reader.NextPart()
	switch err {
	case nil:
		err = merry.New("too many parts").WithHTTPCode(http.StatusBadRequest)
		return
	case io.EOF:
		err = nil
	default:
		err = merry.WithHTTPCode(err, http.StatusBadRequest)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(responseWriter).Encode(fileInfo)
	if err != nil {
		err = merry.WithHTTPCode(err, http.StatusInternalServerError)
		return
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
	return
}

type ListFilesOutput struct {
	Items []*FileInfo `json:"items"`
}

func (handlers *Handlers) ListFiles(responseWriter http.ResponseWriter, request *http.Request) {
	output := &ListFilesOutput{
		Items: handlers.FileInfos,
	}
	err := json.NewEncoder(responseWriter).Encode(output)
	if err != nil {
		log.Println(merry.Details(err))
		http.Error(responseWriter, err.Error(), merry.HTTPCode(err))
		return
	}
}

func (handlers *Handlers) DownloadFile(responseWriter http.ResponseWriter, request *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(merry.Details(err))
			http.Error(responseWriter, err.Error(), merry.HTTPCode(err))
			return
		}
	}()
	vars := mux.Vars(request)
	fileID := vars["file_id"]
	fileInfo, err := handlers.GetFileInfo(fileID)
	if err != nil {
		return
	}
	file, err := os.Open(filepath.Join(handlers.FolderPath, fileInfo.ID))
	if err != nil {
		err = merry.WithHTTPCode(err, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	responseWriter.Header().Set("Content-Type", fileInfo.MimeType)
	responseWriter.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name)
	_, err = io.Copy(responseWriter, file)
	if err != nil {
		err = merry.WithHTTPCode(err, http.StatusInternalServerError)
		return
	}
}

func (handlers *Handlers) GetFileInfo(fileID string) (fileInfo *FileInfo, err error) {
	for _, fileInfo = range handlers.FileInfos {
		if fileInfo.ID == fileID {
			return
		}
	}
	err = merry.New("file not found").WithHTTPCode(http.StatusNotFound)
	return
}

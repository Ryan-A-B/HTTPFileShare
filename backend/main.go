package main

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	file, err := os.OpenFile("journal.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileInfoByID, err := BuildFileInfoByIDFromJournal(file)
	if err != nil {
		panic(err)
	}
	router := GetRouter(&Handlers{
		FolderPath:   "files",
		FileInfoByID: fileInfoByID,
		Journal:      NewJournal(file),
	})
	err = http.ListenAndServe(":8080", router)
	panic(err)
}

func GetRouter(handlers *Handlers) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/upload", handlers.UploadFiles).Methods("POST")
	return
}

func BuildFileInfoByIDFromJournal(reader io.Reader) (fileInfoByID map[string]*FileInfo, err error) {
	fileInfoByID = make(map[string]*FileInfo)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var event Event
		err = json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			return
		}
		switch event.Type {
		case "add_file":
			fileInfoByID[event.AddFile.ID] = event.AddFile
		}
	}
	err = scanner.Err()
	if err != nil {
		return
	}
	return
}

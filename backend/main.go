package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
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
	fileInfos, err := BuildFileInfosFromJournal(file)
	if err != nil {
		panic(err)
	}
	router := GetRouter(&Handlers{
		FolderPath: "files",
		FileInfos:  fileInfos,
		Journal:    NewJournal(file),
	})
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(CORSAllowOriginMiddleware)
	log.Println("Listening on port 9000")
	err = http.ListenAndServe(":9000", router)
	panic(err)
}

func CORSAllowOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func GetRouter(handlers *Handlers) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/files", handlers.AddFile).Methods("POST").Name("AddFile")
	router.HandleFunc("/files", handlers.ListFiles).Methods("GET").Name("ListFiles")
	router.HandleFunc("/files/{file_id}", handlers.DownloadFile).Methods("GET").Name("DownloadFile")
	return
}

func BuildFileInfosFromJournal(reader io.Reader) (fileInfos []*FileInfo, err error) {
	scanner := bufio.NewScanner(reader)
	fileInfos = make([]*FileInfo, 0)
	for scanner.Scan() {
		var event Event
		err = json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			return
		}
		switch event.Type {
		case "add_file":
			fileInfos = append(fileInfos, event.AddFile)
		}
	}
	err = scanner.Err()
	if err != nil {
		return
	}
	return
}

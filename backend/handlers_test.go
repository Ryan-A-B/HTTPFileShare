package main

import (
	"bytes"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	os.MkdirAll("testdata", 0777)
	code := m.Run()
	os.RemoveAll("testdata")
	os.Exit(code)
}

func TestHandlers(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	Convey("TestHandlers", t, func() {
		journalBuffer := bytes.NewBuffer(nil)
		router := GetRouter(&Handlers{
			FolderPath:   "testdata",
			FileInfoByID: make(map[string]*FileInfo),
			Journal:      NewJournal(journalBuffer),
		})
		server := httptest.NewServer(router)
		Convey("TestUploadFiles", func() {
			target := server.URL + "/upload"
			Convey("success", func() {
				Convey("one file", func() {
					buffer := bytes.NewBuffer(nil)
					names := []string{
						uuid.NewV4().String(),
					}
					writer := multipart.NewWriter(buffer)
					for _, name := range names {
						part, err := writer.CreateFormFile("file", name)
						So(err, ShouldBeNil)
						_, err = part.Write([]byte("test"))
						So(err, ShouldBeNil)
					}
					err := writer.Close()
					So(err, ShouldBeNil)
					request, err := http.NewRequest(http.MethodPost, target, buffer)
					So(err, ShouldBeNil)
					request.Header.Set("Content-Type", writer.FormDataContentType())
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					defer response.Body.Close()
					So(response.StatusCode, ShouldEqual, http.StatusOK)
					fileInfoByID, err := BuildFileInfoByIDFromJournal(journalBuffer)
					So(err, ShouldBeNil)
					for _, name := range names {
						found := false
						for _, fileInfo := range fileInfoByID {
							if fileInfo.Name == name {
								delete(fileInfoByID, fileInfo.ID)
								found = true
								break
							}
						}
						So(found, ShouldBeTrue)
					}
					So(len(fileInfoByID), ShouldEqual, 0)
				})
				Convey("two files", func() {
					buffer := bytes.NewBuffer(nil)
					names := []string{
						uuid.NewV4().String(),
						uuid.NewV4().String(),
					}
					writer := multipart.NewWriter(buffer)
					for _, name := range names {
						part, err := writer.CreateFormFile("file", name)
						So(err, ShouldBeNil)
						_, err = part.Write([]byte("test"))
						So(err, ShouldBeNil)
					}
					err := writer.Close()
					So(err, ShouldBeNil)
					request, err := http.NewRequest(http.MethodPost, target, buffer)
					So(err, ShouldBeNil)
					request.Header.Set("Content-Type", writer.FormDataContentType())
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusOK)
					fileInfoByID, err := BuildFileInfoByIDFromJournal(journalBuffer)
					So(err, ShouldBeNil)
					for _, name := range names {
						found := false
						for _, fileInfo := range fileInfoByID {
							if fileInfo.Name == name {
								delete(fileInfoByID, fileInfo.ID)
								found = true
								break
							}
						}
						So(found, ShouldBeTrue)
					}
					So(len(fileInfoByID), ShouldEqual, 0)
				})
				Convey("n files", func() {
					n := rand.Intn(100) + 1
					var names []string
					buffer := bytes.NewBuffer(nil)
					writer := multipart.NewWriter(buffer)
					for i := 0; i < n; i++ {
						name := uuid.NewV4().String()
						part, err := writer.CreateFormFile("file", name)
						So(err, ShouldBeNil)
						_, err = part.Write(uuid.NewV4().Bytes())
						So(err, ShouldBeNil)
						names = append(names, name)
					}
					err := writer.Close()
					So(err, ShouldBeNil)
					request, err := http.NewRequest(http.MethodPost, target, buffer)
					So(err, ShouldBeNil)
					request.Header.Set("Content-Type", writer.FormDataContentType())
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusOK)
					fileInfoByID, err := BuildFileInfoByIDFromJournal(journalBuffer)
					So(err, ShouldBeNil)
					for _, name := range names {
						found := false
						for _, fileInfo := range fileInfoByID {
							if fileInfo.Name == name {
								delete(fileInfoByID, fileInfo.ID)
								found = true
								break
							}
						}
						So(found, ShouldBeTrue)
					}
					So(len(fileInfoByID), ShouldEqual, 0)
				})
			})
			Convey("failure", func() {
				Convey("wrong content type", func() {
					request, err := http.NewRequest(http.MethodPost, target, nil)
					So(err, ShouldBeNil)
					request.Header.Set("Content-Type", "application/json")
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusBadRequest)
				})
				Convey("part is not a file", func() {
					buffer := bytes.NewBuffer(nil)
					writer := multipart.NewWriter(buffer)
					part, err := writer.CreateFormField("test")
					So(err, ShouldBeNil)
					_, err = part.Write([]byte("test"))
					So(err, ShouldBeNil)
					err = writer.Close()
					So(err, ShouldBeNil)
					request, err := http.NewRequest(http.MethodPost, target, buffer)
					So(err, ShouldBeNil)
					request.Header.Set("Content-Type", writer.FormDataContentType())
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusBadRequest)
				})
				Convey("invalid boundary", func() {
					buffer := bytes.NewBuffer(nil)
					writer := multipart.NewWriter(buffer)
					part, err := writer.CreateFormFile("file", "test")
					So(err, ShouldBeNil)
					_, err = part.Write([]byte("test"))
					So(err, ShouldBeNil)
					err = writer.Close()
					So(err, ShouldBeNil)
					request, err := http.NewRequest(http.MethodPost, target, buffer)
					So(err, ShouldBeNil)
					request.Header.Set("Content-Type", "multipart/form-data; boundary=invalid")
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusBadRequest)
				})
			})
		})
	})
}

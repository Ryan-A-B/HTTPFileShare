package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
			FolderPath: "testdata",
			FileInfos:  nil,
			Journal:    NewJournal(journalBuffer),
		})
		server := httptest.NewServer(router)
		Convey("UploadFiles", func() {
			target := server.URL + "/files"
			Convey("success", func() {
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
				fileInfos, err := BuildFileInfosFromJournal(journalBuffer)
				So(err, ShouldBeNil)
				So(len(fileInfos), ShouldEqual, 1)
				for i, name := range names {
					So(fileInfos[i].Name, ShouldEqual, name)
				}
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
				Convey("multiple files", func() {
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
					So(response.StatusCode, ShouldEqual, http.StatusBadRequest)
				})
			})
		})
		Convey("ListFiles", func() {
			target := server.URL + "/files"
			Convey("success", func() {
				Convey("no files", func() {
					request, err := http.NewRequest(http.MethodGet, target, nil)
					So(err, ShouldBeNil)
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusOK)
					var output ListFilesOutput
					err = json.NewDecoder(response.Body).Decode(&output)
					So(err, ShouldBeNil)
					So(len(output.Items), ShouldEqual, 0)
				})
				Convey("1 file", func() {
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
					request.Header.Set("Content-Type", writer.FormDataContentType())
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusOK)
					request, err = http.NewRequest(http.MethodGet, target, nil)
					So(err, ShouldBeNil)
					response, err = http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusOK)
					fileInfos, err := BuildFileInfosFromJournal(journalBuffer)
					So(err, ShouldBeNil)
					So(len(fileInfos), ShouldEqual, 1)
					So(fileInfos[0].Name, ShouldEqual, "test")
				})
				Convey("n files", func() {
					n := rand.Intn(100) + 1
					var names []string
					for i := 0; i < n; i++ {
						buffer := bytes.NewBuffer(nil)
						writer := multipart.NewWriter(buffer)
						name := uuid.NewV4().String()
						part, err := writer.CreateFormFile("file", name)
						So(err, ShouldBeNil)
						_, err = part.Write(uuid.NewV4().Bytes())
						So(err, ShouldBeNil)
						err = writer.Close()
						So(err, ShouldBeNil)
						request, err := http.NewRequest(http.MethodPost, target, buffer)
						So(err, ShouldBeNil)
						request.Header.Set("Content-Type", writer.FormDataContentType())
						response, err := http.DefaultClient.Do(request)
						So(err, ShouldBeNil)
						So(response.StatusCode, ShouldEqual, http.StatusOK)
						names = append(names, name)
					}
					request, err := http.NewRequest(http.MethodGet, target, nil)
					So(err, ShouldBeNil)
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusOK)
					var output ListFilesOutput
					err = json.NewDecoder(response.Body).Decode(&output)
					So(err, ShouldBeNil)
					So(len(output.Items), ShouldEqual, n)
					for i, name := range names {
						So(output.Items[i].Name, ShouldEqual, name)
					}
				})
			})
			Convey("DownloadFile", func() {
				Convey("success", func() {
					buffer := bytes.NewBuffer(nil)
					writer := multipart.NewWriter(buffer)
					name := uuid.NewV4().String()
					part, err := writer.CreateFormFile("file", name)
					So(err, ShouldBeNil)
					content := uuid.NewV4().Bytes()
					_, err = part.Write(content)
					So(err, ShouldBeNil)
					err = writer.Close()
					So(err, ShouldBeNil)
					target := server.URL + "/files"
					request, err := http.NewRequest(http.MethodPost, target, buffer)
					So(err, ShouldBeNil)
					request.Header.Set("Content-Type", writer.FormDataContentType())
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusOK)
					var fileInfo FileInfo
					err = json.NewDecoder(response.Body).Decode(&fileInfo)
					So(err, ShouldBeNil)
					So(fileInfo.Name, ShouldEqual, name)
					target = server.URL + "/files/" + fileInfo.ID
					request, err = http.NewRequest(http.MethodGet, target, nil)
					So(err, ShouldBeNil)
					response, err = http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusOK)
					So(response.Header.Get("Content-Type"), ShouldEqual, "application/octet-stream")
					So(response.Header.Get("Content-Disposition"), ShouldEqual, "attachment; filename="+name)
					body, err := ioutil.ReadAll(response.Body)
					So(err, ShouldBeNil)
					So(body, ShouldResemble, content)
				})
				Convey("file not found", func() {
					target := server.URL + "/files/" + uuid.NewV4().String()
					request, err := http.NewRequest(http.MethodGet, target, nil)
					So(err, ShouldBeNil)
					response, err := http.DefaultClient.Do(request)
					So(err, ShouldBeNil)
					So(response.StatusCode, ShouldEqual, http.StatusNotFound)
				})
			})
		})
	})
}

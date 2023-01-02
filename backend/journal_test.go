package main

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestJournal(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	Convey("TestJournal", t, func() {
		buffer := bytes.NewBuffer(nil)
		journal := NewJournal(buffer)
		Convey("AddFile once", func() {
			journal.AddFile(&FileInfo{
				ID:       "id",
				Name:     "name",
				MimeType: "mime_type",
			})
			So(buffer.String(), ShouldEqual, "{\"type\":\"add_file\",\"add_file\":{\"id\":\"id\",\"name\":\"name\",\"mime_type\":\"mime_type\"}}\n")
		})
		Convey("AddFile twice", func() {
			journal.AddFile(&FileInfo{
				ID:       "id1",
				Name:     "name1",
				MimeType: "mime_type1",
			})
			journal.AddFile(&FileInfo{
				ID:       "id2",
				Name:     "name2",
				MimeType: "mime_type2",
			})
			So(buffer.String(), ShouldEqual, "{\"type\":\"add_file\",\"add_file\":{\"id\":\"id1\",\"name\":\"name1\",\"mime_type\":\"mime_type1\"}}\n{\"type\":\"add_file\",\"add_file\":{\"id\":\"id2\",\"name\":\"name2\",\"mime_type\":\"mime_type2\"}}\n")
		})
		Convey("AddFile n times", func() {
			n := rand.Intn(1000)
			var expectedFileInfos []*FileInfo
			for i := 0; i < n; i++ {
				fileInfo := &FileInfo{
					ID:       uuid.NewV4().String(),
					Name:     uuid.NewV4().String(),
					MimeType: uuid.NewV4().String(),
				}
				journal.AddFile(fileInfo)
				expectedFileInfos = append(expectedFileInfos, fileInfo)
			}
			actualFileInfos, err := BuildFileInfosFromJournal(buffer)
			So(err, ShouldBeNil)
			So(actualFileInfos, ShouldResemble, expectedFileInfos)
		})
	})
}

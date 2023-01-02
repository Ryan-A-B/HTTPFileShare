package main

import (
	"encoding/json"
	"io"
)

type Journal struct {
	JSONEncoder *json.Encoder
}

func NewJournal(writer io.Writer) *Journal {
	return &Journal{
		JSONEncoder: json.NewEncoder(writer),
	}
}

type Event struct {
	Type    string    `json:"type"`
	AddFile *FileInfo `json:"add_file"`
}

func (journal *Journal) AddFile(fileInfo *FileInfo) {
	event := Event{
		Type:    "add_file",
		AddFile: fileInfo,
	}
	journal.JSONEncoder.Encode(event)
}

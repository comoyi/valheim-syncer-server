package server

import (
	"net/http"
	"time"
)

type ServerFileInfo struct {
	Files []*FileInfo `json:"files"`
}

type FileInfo struct {
	//Name string `json:"name"`
	Path string `json:"path"`
	Hash string `json:"hash"`
}

var serverFileInfo *ServerFileInfo = &ServerFileInfo{}
var baseDir string

func Start() {

	go func() {
		refreshFileInfo()
	}()

	http.HandleFunc("/sync", sync)
	http.HandleFunc("/files", files)
	err := http.ListenAndServe(":10123", nil)
	if err != nil {
		return
	}
}

func refreshFileInfo() {
	for {
		select {
		case <-time.After(5 * time.Second):
			files := make([]*FileInfo, 0)

			file := &FileInfo{
				//Name: "",
				Path: "",
				Hash: "",
			}
			files = append(files, file)

			serverFileInfo.Files = files
		}
	}
}

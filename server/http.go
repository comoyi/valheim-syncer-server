package server

import (
	"encoding/json"
	"github.com/comoyi/valheim-syncer-server/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func files(writer http.ResponseWriter, request *http.Request) {
	bytes, err := json.Marshal(serverFileInfo)
	if err != nil {
		log.Debugf("json.Marshal failed, err: %s\n", err)
		return
	}

	j := string(bytes)
	log.Debugf("json: %s\n", j)
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = writer.Write(bytes)
	if err != nil {
		log.Debugf("write failed, err: %s\n", err)
		return
	}
}

func sync(writer http.ResponseWriter, request *http.Request) {

	relativePath := request.FormValue("file")

	filePath := filepath.Join(baseDir, relativePath)
	if !strings.HasPrefix(filePath, baseDir) {
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = io.Copy(writer, f)
	if err != nil {
		return
	}
}

func announcement(writer http.ResponseWriter, request *http.Request) {
	content := announcementContent
	announcement := &Announcement{Content: content}

	bytes, err := json.Marshal(announcement)
	if err != nil {
		log.Debugf("json.Marshal failed, err: %s\n", err)
		return
	}

	j := string(bytes)
	log.Debugf("json: %s\n", j)
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = writer.Write(bytes)
	if err != nil {
		log.Debugf("write failed, err: %s\n", err)
		return
	}
}

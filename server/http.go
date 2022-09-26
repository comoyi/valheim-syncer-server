package server

import (
	"encoding/json"
	"fmt"
	"github.com/comoyi/valheim-syncer-server/log"
	"mime"
	"net/http"
)

func files(writer http.ResponseWriter, request *http.Request) {
	bytes, err := json.Marshal(serverFileInfo)
	if err != nil {
		log.Debugf("json.Marshal failed, err: %s\n", err)
		return
	}

	j := string(bytes)
	log.Debugf("json: %s\n", j)
	writer.Header().Set("Content-Type", mime.TypeByExtension("json"))
	_, err = writer.Write(bytes)
	if err != nil {
		log.Debugf("write failed, err: %s\n", err)
		return
	}
}

func sync(writer http.ResponseWriter, request *http.Request) {

	file := request.FormValue("file")
	fmt.Printf("file: %s\n", file)
	_, err := writer.Write([]byte("Hi\nHello\nASD\n"))
	if err != nil {
		return
	}
}

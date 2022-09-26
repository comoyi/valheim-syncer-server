package server

import (
	"crypto/md5"
	"fmt"
	"github.com/comoyi/valheim-syncer-server/log"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ServerFileInfo struct {
	Files []*FileInfo `json:"files"`
}

type FileInfo struct {
	Path string `json:"path"`
	Type int8   `json:"type"`
	Hash string `json:"hash"`
}

const (
	TypeFile int8 = 1
	TypeDir  int8 = 2
)

var serverFileInfo *ServerFileInfo = &ServerFileInfo{}
var baseDir string

func Start() {

	baseDir = "/tmp/vvv"

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
	doRefreshFileInfo()
	for {
		select {
		case <-time.After(10 * time.Second):
			doRefreshFileInfo()
		}
	}
}

func doRefreshFileInfo() {
	log.Debugf("refresh files info\n")
	files := make([]*FileInfo, 0)

	err := filepath.Walk(baseDir, walkFun(&files))
	if err != nil {
		log.Debugf("refresh files info failed\n")
		return
	}

	serverFileInfo.Files = files
}

func walkFun(files *[]*FileInfo) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		pathRelative := strings.TrimPrefix(path, baseDir)
		var file *FileInfo
		if info.IsDir() {
			file = &FileInfo{
				Path: pathRelative,
				Type: TypeDir,
				Hash: "",
			}
		} else {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			bytes, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			hashSumRaw := md5.Sum(bytes)
			hashSum := fmt.Sprintf("%x", hashSumRaw)
			log.Debugf("file: %s, hashSum: %s\n", path, hashSum)
			file = &FileInfo{
				Path: pathRelative,
				Type: TypeFile,
				Hash: hashSum,
			}
		}
		*files = append(*files, file)
		return nil
	}
}

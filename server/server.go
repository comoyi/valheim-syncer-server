package server

import (
	"crypto/md5"
	"fmt"
	"github.com/comoyi/valheim-syncer-server/config"
	"github.com/comoyi/valheim-syncer-server/log"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var serverFileInfo *ServerFileInfo = &ServerFileInfo{
	ScanStatus: ScanStatusWait,
	Files:      make([]*FileInfo, 0),
}

var baseDir string

var appName = "Valheim Syncer Server"
var versionText = "0.0.1"

func Start() {

	baseDir = config.Conf.BaseDir
	if baseDir == "" {
		fmt.Printf("baseDir invalid\n")
		log.Errorf("baseDir invalid\n")
		return
	}
	baseDir = strings.TrimSuffix(baseDir, string(os.PathSeparator))

	go func() {
		refreshFileInfo()
	}()

	http.HandleFunc("/sync", sync)
	http.HandleFunc("/files", files)
	http.HandleFunc("/announcement", announcement)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Conf.Port), nil)
	if err != nil {
		fmt.Printf("server start failed err: %v\n", err)
		log.Errorf("server start failed err: %v\n", err)
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

	serverFileInfo.ScanStatus = ScanStatusScanning

	err := filepath.Walk(baseDir, walkFun(&files))
	if err != nil {
		log.Debugf("refresh files info failed\n")
		serverFileInfo.ScanStatus = ScanStatusFailed
		return
	}

	serverFileInfo.Files = files
	serverFileInfo.ScanStatus = ScanStatusCompleted
}

func walkFun(files *[]*FileInfo) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if !strings.HasPrefix(path, baseDir) {
			log.Warnf("path not excepted, path: %s\n", path)
			return nil
		}
		pathRelative := strings.TrimPrefix(path, baseDir)
		if pathRelative == "" {
			return nil
		}
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
			log.Tracef("file: %s, hashSum: %s\n", path, hashSum)
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

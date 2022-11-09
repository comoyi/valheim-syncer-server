package server

import (
	"fmt"
	"github.com/comoyi/valheim-syncer-server/config"
	"github.com/comoyi/valheim-syncer-server/log"
	"github.com/comoyi/valheim-syncer-server/util/cryptoutil/md5util"
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

var appName = "Valheim Syncer Server"
var versionText = "1.0.6"

var baseDir string

func Start() {

	baseDir = config.Conf.Dir
	if baseDir != "" {
		baseDir = filepath.Clean(baseDir)
	}

	setAnnouncement(config.Conf.Announcement)

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
	interval := config.Conf.Interval
	doRefreshFileInfo()
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Second):
			doRefreshFileInfo()
		}
	}
}

func doRefreshFileInfo() {
	log.Debugf("refresh files info\n")

	if baseDir == "" {
		log.Errorf("baseDir invalid\n")
		setDirStatusLedRed()
		return
	}

	files := make([]*FileInfo, 0)

	serverFileInfo.ScanStatus = ScanStatusScanning

	err := filepath.Walk(baseDir, walkFun(&files))
	if err != nil {
		log.Debugf("refresh files info failed, err: %v\n", err)
		serverFileInfo.ScanStatus = ScanStatusFailed
		setDirStatusLedRed()
		return
	}

	serverFileInfo.Files = files
	serverFileInfo.ScanStatus = ScanStatusCompleted
	setDirStatusLedGreen()
}

func walkFun(files *[]*FileInfo) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == baseDir {
			return nil
		}
		if !strings.HasPrefix(path, baseDir) {
			log.Warnf("path not expected, baseDir: %s, path: %s\n", baseDir, path)
			return fmt.Errorf("path not expected, baseDir: %s, path: %s\n", baseDir, path)
		}
		relativePath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}
		if relativePath == "." ||
			relativePath == ".." ||
			strings.HasPrefix(relativePath, "./") ||
			strings.HasPrefix(relativePath, ".\\") ||
			strings.HasPrefix(relativePath, "../") ||
			strings.HasPrefix(relativePath, "..\\") {
			return fmt.Errorf("relativePath not expected, baseDir: %s, path: %s, relativePath: %s\n", baseDir, path, relativePath)
		}
		if relativePath == "" {
			return nil
		}
		var file *FileInfo
		if info.IsDir() {
			log.Tracef("dir:  %s\n", path)
			file = &FileInfo{
				RelativePath: relativePath,
				Type:         TypeDir,
				Hash:         "",
			}
		} else if info.Mode()&os.ModeSymlink != 0 {
			log.Tracef("symlink:  %s\n", path)
			file = &FileInfo{
				RelativePath: relativePath,
				Type:         TypeSymlink,
				Hash:         "",
			}
		} else if info.Mode().IsRegular() {
			hashSum, err := md5util.SumFile(path)
			if err != nil {
				return err
			}
			log.Tracef("file: %s, hashSum: %s\n", path, hashSum)
			file = &FileInfo{
				RelativePath: relativePath,
				Type:         TypeFile,
				Hash:         hashSum,
			}
		} else {
			log.Tracef("unhandled file type, filepath:  %s\n", path)
			return nil
		}
		*files = append(*files, file)
		return nil
	}
}

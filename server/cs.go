package server

type ScanStatus int8

const (
	ScanStatusWait      ScanStatus = 10
	ScanStatusScanning  ScanStatus = 20
	ScanStatusFailed    ScanStatus = 30
	ScanStatusCompleted ScanStatus = 40
)

type FileType int8

const (
	TypeFile FileType = 1
	TypeDir  FileType = 2
)

type ServerFileInfo struct {
	ScanStatus ScanStatus  `json:"status"`
	Files      []*FileInfo `json:"files"`
}

type FileInfo struct {
	RelativePath string   `json:"relative_path"`
	Type         FileType `json:"type"`
	Hash         string   `json:"hash"`
}

type Announcement struct {
	Content string `json:"content"`
	Hash    string `json:"hash"`
}

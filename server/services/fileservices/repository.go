package fileservices

import "os"

type FileRepo interface {
	DownloadFile(fileName string) ([]byte, error)
	UploadFile(fileName string, data []byte) error
	GetListFiles() ([]os.DirEntry, error)
}

package repo

import (
	"os"
	"path/filepath"
)

type FileRepo struct {
	storage string
}

func NewFileRepo(storage string) *FileRepo {
	return &FileRepo{storage: storage}
}

func (f FileRepo) DownloadFile(fileName string) ([]byte, error) {

	filePath := filepath.Join("server", f.storage, fileName)
	return os.ReadFile(filePath)
}

func (f FileRepo) UploadFile(fileName string, data []byte) error {

	path := filepath.Join("server", f.storage, fileName)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (f FileRepo) GetListFiles() ([]os.DirEntry, error) {
	path := filepath.Join("server", f.storage)
	return os.ReadDir(path)
}

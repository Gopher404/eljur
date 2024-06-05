package file

import (
	"eljur/internal/domain/models"
	"eljur/pkg/tr"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FilesStorage struct {
	storagePath string
}

func NewFilesStorage(storagePath string) *FilesStorage {
	return &FilesStorage{
		storagePath: storagePath,
	}
}

func (f *FilesStorage) Get(path string) (*models.File, error) {
	path = filepath.Join(f.storagePath, path)
	file, err := os.Open(path)
	if err != nil {
		return nil, tr.Trace(err)
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, tr.Trace(err)
	}

	return &models.File{
		Name:    fileInfo.Name(),
		ModTime: fileInfo.ModTime(),
		Data:    file,
	}, nil
}

func (f *FilesStorage) ListDir(path string) ([]models.FileInfo, error) {
	path = filepath.Join(f.storagePath, path)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, tr.Trace(err)
	}
	var filesInfo []models.FileInfo
	for _, file := range files {
		i, _ := file.Info()
		i.ModTime()
		filesInfo = append(filesInfo, models.FileInfo{
			Name:    file.Name(),
			IsDir:   file.IsDir(),
			ModTime: getTimeFromFile(file),
		})
	}
	return filesInfo, nil
}

func (f *FilesStorage) Save(path string, data io.Reader) error {
	path = filepath.Join(f.storagePath, path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		return tr.Trace(err)
	}

	if _, err := io.Copy(file, data); err != nil {
		return tr.Trace(err)
	}

	return nil
}

func (f *FilesStorage) CreateDir(path string) error {
	path = filepath.Join(f.storagePath, path)
	if err := os.Mkdir(path, 666); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func getTimeFromFile(file os.DirEntry) string {
	i, err := file.Info()
	if err != nil {
		return ""
	}
	return i.ModTime().Format(time.DateTime)
}

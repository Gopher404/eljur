package files

import (
	"eljur/internal/domain/models"
	"eljur/internal/storage"
	"eljur/pkg/tr"
	"io"
)

func New(s storage.Files) *FileService {
	return &FileService{
		s: s,
	}
}

type FileService struct {
	s storage.Files
}

func (f *FileService) Get(path string) (*models.File, error) {
	file, err := f.s.Get(path)
	if err != nil {
		return nil, tr.Trace(err)
	}
	return file, nil
}

func (f *FileService) ListDir(path string) ([]models.FileInfo, error) {
	dir, err := f.s.ListDir(path)
	if err != nil {
		return nil, tr.Trace(err)
	}
	return dir, nil
}

func (f *FileService) SaveFile(path string, data io.Reader) error {
	if err := f.s.Save(path, data); err != nil {
		return tr.Trace(err)
	}
	return nil
}

func (f *FileService) CreateDir(path string) error {
	if err := f.s.CreateDir(path); err != nil {
		return tr.Trace(err)
	}
	return nil
}

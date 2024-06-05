package models

import (
	"io"
	"time"
)

type FileInfo struct {
	Name    string `json:"name"`
	IsDir   bool   `json:"is_dir"`
	ModTime string `json:"mod_time"`
}

type File struct {
	Name    string
	Data    io.ReadSeeker
	ModTime time.Time
}

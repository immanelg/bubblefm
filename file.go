package main

import (
	"io/fs"
	"time"
)

type File struct {
	Name     string
	Path     string
	IsDir    bool
	Size     int64
	Modified time.Time
	Mode     fs.FileMode
}

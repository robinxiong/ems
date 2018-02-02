package filesystem

import (
	"path/filepath"
	"log"
)

type FileSystem struct {
	Base string //用于文件存储的根目录
}

func New(base string) *FileSystem {
	absbase, err := filepath.Abs(base)
	if err != nil {
		log.Println("FileSystem storage's directory haven't been initialized")
	}

	return &FileSystem{Base: absbase}
}




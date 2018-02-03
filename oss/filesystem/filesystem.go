package filesystem

import (
	"path/filepath"
	"log"
	"strings"
	"os"
	"io"
	"ems/oss"
)

type FileSystem struct {
	Base string //用于文件存储的根目录
}

func New(base string) *FileSystem {
	//如果是相对目录，则是相对于当前应用程序
	absbase, err := filepath.Abs(base)
	if err != nil {
		log.Println("FileSystem storage's directory haven't been initialized")
	}

	return &FileSystem{Base: absbase}
}


// GetEndpoint get endpoint, FileSystem's endpoint is /
func (fileSystem FileSystem) GetEndpoint() string {
	return "/"
}

func (fileSystem FileSystem) GetURL(path string) (url string, err error) {
	return
}
func (fileSystem FileSystem) Get(path string) (*os.File, error) {
	return os.Open(fileSystem.GetFullPath(path))
}


// Put store a reader into given path
func (fileSystem FileSystem) Put(path string, reader io.Reader) (*oss.Object, error) {
	var (
		fullpath = fileSystem.GetFullPath(path)
		err      = os.MkdirAll(filepath.Dir(fullpath), os.ModePerm)
	)

	if err != nil {
		return nil, err
	}

	dst, err := os.Create(fullpath)

	if err == nil {
		if seeker, ok := reader.(io.ReadSeeker); ok {
			seeker.Seek(0, 0)
		}
		_, err = io.Copy(dst, reader)
	}

	return &oss.Object{Path: path, Name: filepath.Base(path), StorageInterface: fileSystem}, err
}


func (fileSystem FileSystem) Delete(path string) error {
	return os.Remove(fileSystem.GetFullPath(path))
}

// List list all objects under current path
func (fileSystem FileSystem) List(path string) ([]*oss.Object, error) {
	var (
		objects  []*oss.Object
		fullpath = fileSystem.GetFullPath(path)
	)

	filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if path == fullpath {
			return nil
		}

		if err == nil && !info.IsDir() {
			modTime := info.ModTime()
			objects = append(objects, &oss.Object{
				Path:             strings.TrimPrefix(path, fileSystem.Base),
				Name:             info.Name(),
				LastModified:     &modTime,
				StorageInterface: fileSystem,
			})
		}
		return nil
	})

	return objects, nil
}
// GetFullPath get full path from absolute/relative path, 不是StorageInterface
func (fileSystem FileSystem) GetFullPath(path string) string {
	fullpath := path
	if !strings.HasPrefix(path, fileSystem.Base) {
		fullpath, _ = filepath.Abs(filepath.Join(fileSystem.Base, path))
	}
	return fullpath
}


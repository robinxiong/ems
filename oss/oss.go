package oss
import (
	"io"
	"os"
	"time"
)

// StorageInterface define common API to operate storage
type StorageInterface interface {
	GetURL(path string) (string, error)
	Get(path string) (*os.File, error)
	Put(path string, reader io.Reader) (*Object, error)
	Delete(path string) error
	List(path string) ([]*Object, error)
	GetEndpoint() string
}

// Object content object
type Object struct {
	Path             string
	Name             string
	LastModified     *time.Time
	StorageInterface StorageInterface
}

// 返回一个对像的内容
func (object Object) Get() (*os.File, error) {
	return object.StorageInterface.Get(object.Path)
}

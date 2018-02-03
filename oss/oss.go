package oss
import (
	"io"
	"os"
	"time"
)

// StorageInterface define common API to operate storage
//创建一个保存资源的平台，可以是本地的文件系统,也可以是s3， 它用于获取对像，文件，保存对像， 删除指定路径文件
type StorageInterface interface {
	GetURL(path string) (string, error)
	Get(path string) (*os.File, error)
	Put(path string, reader io.Reader) (*Object, error)
	Delete(path string) error
	List(path string) ([]*Object, error) //列出指定目录下的文件，以Object对像返回
	GetEndpoint() string
}

// Object content object
type Object struct {
	Path             string
	Name             string
	LastModified     *time.Time
	StorageInterface StorageInterface  //平台的引用，文件或者s3的引用
}

// 返回一个对像的内容
func (object Object) Get() (*os.File, error) {
	return object.StorageInterface.Get(object.Path)
}

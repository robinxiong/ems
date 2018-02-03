package media

import (
	"io"
	"mime/multipart"
)

// Base defined a base struct for storages
// 所以media的基础类，它实现了Media接口，其它类filesystem和oss（以对像保存文件）都继承此类
// model struct保存以聚合的方式，oss Oss字段，保存它相关的资源, 例如
/*
	type Product struct {
	  gorm.Model
	  Image oss.OSS
	}
 */
type Base struct {
	FileName    string
	Url         string
	CropOptions map[string]*CropOption `json:",omitempty"`
	Delete      bool                   `json:"-"`
	Crop        bool                   `json:"-"`
	FileHeader  FileHeader             `json:"-"`
	Reader      io.Reader              `json:"-"`
	cropped     bool
}


// CropOption includes crop options
type CropOption struct {
	X, Y, Width, Height int
}

// FileHeader is an interface, for matched values, when call its `Open` method will return `multipart.File`
type FileHeader interface {
	Open() (multipart.File, error)
}
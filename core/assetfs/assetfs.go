package assetfs

import (
	"fmt"
	"runtime/debug"
)

type Interface interface{
	PrependPath(path string) error     //将路径添加到其它路径的前面
	RegisterPath(path string) error   //注册一个路径
	Asset(name string) ([]byte, error)  //精确查找文件
	Glob(pattern string) (matches []string, err error) //模糊查找path下的文件pattern, 如果找到, 返回相对于路径下的文件路径
	Compile() error
	NameSpace(nameSpace string) Interface  //指定子AssetFileSystem
}

// AssetFS default assetfs
var assetFS Interface = &AssetFileSystem{}
var used bool



// AssetFS get AssetFS
func AssetFS() Interface {
	used = true
	return assetFS
}

// SetAssetFS set assetfs
func SetAssetFS(fs Interface) {
	if used {
		fmt.Println("WARNING: AssetFS is used before overwrite it!")
		debug.PrintStack()
	}

	assetFS = fs
}
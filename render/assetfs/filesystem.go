package assetfs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//AssetFileSystem AssetFS base on FileSystem

type AssetFileSystem struct {
	paths        []string             //用来保存所有的asset路径
	nameSpacedFS map[string]Interface //保存Sub AssetFileSystem
}

func (fs *AssetFileSystem) RegisterPath(pth string) error {
	if _, err := os.Stat(pth); !os.IsNotExist(err) {
		var existing bool
		for _, p := range fs.paths {
			if p == pth {
				existing = true
				break
			}
		}
		if !existing {
			fs.paths = append(fs.paths, pth)
		}
		return nil
	}
	return errors.New("path not found")
}

// PrependPath prepend path to view paths, 将第一个参数添加到paths最前面
func (fs *AssetFileSystem) PrependPath(pth string) error {
	if _, err := os.Stat(pth); !os.IsNotExist(err) {
		var existing bool
		for _, p := range fs.paths {
			if p == pth {
				existing = true
				break
			}
		}
		if !existing {
			fs.paths = append([]string{pth}, fs.paths...)
		}
		return nil
	}
	return errors.New("not found")
}

// Asset get content with name from assetfs
// 查找fs.paths下的所有文件为name的资源
func (fs *AssetFileSystem) Asset(name string) ([]byte, error) {
	for _, pth := range fs.paths {
		filePath := filepath.Join(pth, name)
		if _, err := os.Stat(filePath); err == nil {
			return ioutil.ReadFile(filePath)
		}
	}
	return []byte{}, fmt.Errorf("%v not found", name)
}

// Glob list matched files from assetfs
func (fs *AssetFileSystem) Glob(pattern string) (matches []string, err error) {
	for _, pth := range fs.paths {
		if results, err := filepath.Glob(filepath.Join(pth, pattern)); err == nil {
			for _, result := range results {
				matches = append(matches, strings.TrimPrefix(result, pth))
			}
		}
	}
	return
}

func (fs *AssetFileSystem) Compile() error {
	return nil
}

// NameSpace return namespaced filesystem
func (fs *AssetFileSystem) NameSpace(nameSpace string) Interface {
	if fs.nameSpacedFS == nil {
		fs.nameSpacedFS = map[string]Interface{}
	}
	fs.nameSpacedFS[nameSpace] = &AssetFileSystem{}
	return fs.nameSpacedFS[nameSpace]
}

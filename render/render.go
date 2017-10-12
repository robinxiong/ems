package render

import (
	"net/http"
	"html/template"
	"ems/render/assetfs"
	"path/filepath"
	"ems/core/utils"
	"strings"
	"os"
)

//DefaultLayout default layout name
const DefaultLayout = "application"

//DefaultViewPath default view path
const DefaultViewPath = "app/views"

// Render the render struct.
type Render struct {
	*Config
	funcMaps template.FuncMap
}
// Config render config
type Config struct {
	ViewPaths []string
	DefaultLayout string
	FuncMapMaker func(render *Render, request *http.Request, writer http.ResponseWriter) template.FuncMap
	AssetFileSystem assetfs.Interface
}

//New initialize the render struct
func New(config *Config, viewPaths ...string) *Render{
	if config == nil {
		config = &Config{}
	}

	if config.DefaultLayout == "" {
		config.DefaultLayout = DefaultLayout
	}

	if config.AssetFileSystem == nil {
		config.AssetFileSystem = assetfs.AssetFS().NameSpace("views")
	}

	config.ViewPaths = append(append(config.ViewPaths, viewPaths...), DefaultViewPath)

	r := &Render{
		Config: config,
		funcMaps:map[string]interface{},
	}

	for _, viewPath := range config.ViewPaths {
		r.RegisterViewPath(viewPath)
	}
	return r
}

//RegisterViewPath register view path
func (render *Render) RegisterViewPath(paths ...string) {
	for _, pth := range paths {
		if filepath.IsAbs(pth) {
			render.ViewPaths = append(render.ViewPaths, pth)
			render.AssetFileSystem.RegisterPath(pth) //将视图路径也添加到AssetFileSystem
		} else {
			//将当前路径添加到pth前，生成绝对路径, 如果不存在，则调用utils.AppRoot中的路径
			if absPath, err := filepath.Abs(pth); err == nil && isExistingDir(pth){
				render.ViewPaths = append(render.ViewPaths, absPath)
				render.AssetFileSystem.RegisterPath(absPath)
			} else if isExistingDir(filepath.Join(utils.AppRoot, "vendor", pth)) {
				render.AssetFileSystem.RegisterPath(filepath.Join(utils.AppRoot, "vendor", pth))
			} else {
				for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
					if p:= filepath.Join(gopath, "src", pth); isExistingDir(p) {
						render.ViewPaths = append(render.ViewPaths, p)
						render.AssetFileSystem.RegisterPath(p)
					}
				}
			}
		}
	}
}

// PrependViewPath prepend view path
func (render *Render) PrependViewPath(paths ...string) {
	for _, pth := range paths {
		if filepath.IsAbs(pth) {
			render.ViewPaths = append([]string{pth}, render.ViewPaths...)
			render.AssetFileSystem.PrependPath(pth)
		} else {
			if absPath, err := filepath.Abs(pth); err == nil && isExistingDir(absPath) {
				render.ViewPaths = append([]string{absPath}, render.ViewPaths...)
				render.AssetFileSystem.PrependPath(absPath)
			} else if isExistingDir(filepath.Join(utils.AppRoot, "vendor", pth)) {
				render.AssetFileSystem.PrependPath(filepath.Join(utils.AppRoot, "vendor", pth))
			} else {
				for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
					if p := filepath.Join(gopath, "src", pth); isExistingDir(p) {
						render.ViewPaths = append([]string{p}, render.ViewPaths...)
						render.AssetFileSystem.PrependPath(p)
					}
				}
			}
		}
	}
}

// SetAssetFS set asset fs for render
func (render *Render) SetAssetFS(assetFS assetfs.Interface) {
	for _, viewPath := range render.ViewPaths {
		assetFS.RegisterPath(viewPath)
	}

	render.AssetFileSystem = assetFS
}

// Layout set layout for template.
func (render *Render) Layout(name string) *Te

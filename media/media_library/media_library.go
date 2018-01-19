package media_library

import (
	"encoding/json"
	"ems/media/oss"
	"ems/media"
	"github.com/jinzhu/gorm"
)

type File struct {
	ID          json.Number
	Url         string
	VideoLink   string
	FileName    string
	Description string
}

type MediaBox struct {
	Values string `json:"-" gorm:"size:4294967295;"`
	Files  []File `json:",omitempty"`
}


type MediaLibraryStorage struct {
	oss.OSS
	Sizes        map[string]*media.Size `json:",omitempty"`
	Video        string
	SelectedType string
	Description  string
}

type MediaLibrary struct {
	gorm.Model
	SelectedType string
	File         MediaLibraryStorage `sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
}
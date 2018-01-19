package help

import "github.com/jinzhu/gorm"

type QorHelpEntry struct {
	gorm.Model
	Title      string
	Categories Categories
	Body       string `gorm:"size:65532"`
}

type Categories struct {
	RawValue   string
	Categories []string
}
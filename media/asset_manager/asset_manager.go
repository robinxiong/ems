package asset_manager

import (
	"github.com/jinzhu/gorm"
	"ems/media/oss"
)

type AssetManager struct {
	gorm.Model
	File oss.OSS `media_library:"URL:/system/assets/{{primary_key}}/{{filename_with_hash}}"`
}
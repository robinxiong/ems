package banner_editor

import (
	"github.com/jinzhu/gorm"
	"ems/serializable_meta"
)

// QorBannerEditorSetting default setting model
type QorBannerEditorSetting struct {
	gorm.Model
	serializable_meta.SerializableMeta
}
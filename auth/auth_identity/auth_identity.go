package auth_identity

import (
	"github.com/jinzhu/gorm"
	"time"
)

type AuthIdentity struct{
	gorm.Model
	Basic
	SignLogs
}

type Basic struct {
	Provider string // phone, email, wechat, github...
	UID  string `gorm:"column:uid"`
	EncryptedPassword string
	UserID            string
	ConfirmedAt       *time.Time  //验证的时间
}


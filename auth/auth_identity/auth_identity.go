package auth_identity

import (
	"github.com/jinzhu/gorm"
	"time"
	"ems/auth/claims"
)
//用于认证的表，它表含用户名和加密的密码
type AuthIdentity struct{
	gorm.Model
	Basic
	SignLogs
}
//参考providers/password/handlers/DefaultAuthorizeHandler
type Basic struct {
	Provider string // phone, email, wechat, github...
	UID  string `gorm:"column:uid"` //邮箱帐号名
	EncryptedPassword string
	UserID            string
	ConfirmedAt       *time.Time  //验证的时间
}

func (basic Basic) ToClaims() *claims.Claims {
	claims := claims.Claims{}
	claims.Provider = basic.Provider
	claims.Id = basic.UID
	claims.UserID = basic.UserID
	return &claims
}

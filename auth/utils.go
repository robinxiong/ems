package auth

import (
	"ems/auth/claims"
	"ems/core/utils"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

const CurrentUser utils.ContextKey = "current_user"

func (auth *Auth) GetCurrentUser(req *http.Request) interface{} {
	//req context下是否有current_user
	if currentUser := req.Context().Value(CurrentUser); currentUser != nil {
		return currentUser
	}

	claims, err := auth.SessionStorer.Get(req)

	if err == nil {

		context := &Context{Auth: auth, Claims: claims, Request: req}
		if user, err := auth.UserStorer.Get(claims, context); err == nil {

			return user
		}
	}

	return nil
}

//providers/handlers/DefaultAuthorizeHandler
func (auth *Auth) GetDB(request *http.Request) *gorm.DB {
	db := request.Context().Value(utils.ContextDBName) //site/config/routes/routes.go中设置
	if tx, ok := db.(*gorm.DB); ok {
		return tx
	}
	return auth.Config.DB
}

//登录成功后，更新claims中的LastLoginAt为当前时间，并且更新session
func (auth *Auth) Login(w http.ResponseWriter, req *http.Request, claimer claims.ClaimerInterface) error {
	claims := claimer.ToClaims()
	now := time.Now()
	claims.LastLoginAt = &now

	return auth.SessionStorer.Update(w, req, claims)
}

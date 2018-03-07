package auth

import (
	"ems/auth/auth_identity"
	"ems/auth/claims"
	"ems/core/utils"
	"fmt"
	"reflect"

	"github.com/jinzhu/copier"
)

type UserStorerInterface interface {
	//第一个参数为用户信息表
	Save(schema *Schema, context *Context) (user interface{}, userId string, err error)
	Get(claims *claims.Claims, context *Context) (user interface{}, err error)
}

type UserStorer struct {
}

func (*UserStorer) Save(schema *Schema, context *Context) (user interface{}, userId string, err error) {

	var tx = context.Auth.GetDB(context.Request)

	if context.Auth.Config.UserModel != nil {
		currentUser := reflect.New(utils.ModelType(context.Auth.Config.UserModel)).Interface()
		//从schema中复制currentUser中相同的字段
		copier.Copy(currentUser, schema)
		err = tx.Create(currentUser).Error
		return currentUser, fmt.Sprint(tx.NewScope(currentUser).PrimaryKeyValue()), err
	}
	return nil, "", nil
}

//根据claims中的User id去查找用户信息
func (*UserStorer) Get(Claims *claims.Claims, context *Context) (user interface{}, err error) {
	var tx = context.Auth.GetDB(context.Request)
	//如果给auth配置了UserModel, 则根据claims中保存的UserID去查找数据库中的用户信息
	if context.Auth.Config.UserModel != nil {

		if Claims.UserID != "" {
			currentUser := reflect.New(utils.ModelType(context.Auth.Config.UserModel)).Interface()
			if err = tx.First(currentUser, Claims.UserID).Error; err == nil {

				return currentUser, nil
			}
			return nil, ErrInvalidAccount
		}
	}

	//没有指定UserModel
	var (
		authIdentity = reflect.New(utils.ModelType(context.Auth.Config.AuthIdentityModel)).Interface()
		authInfo     = auth_identity.Basic{
			Provider: Claims.Provider,
			UID:      Claims.Id,
		}
	)

	//首先查找auth_identity表, 如果找到记录
	//查找usermodel,
	if !tx.Where(authInfo).First(&authIdentity).RecordNotFound() {
		if context.Auth.Config.UserModel != nil {
			if authBasicInfo, ok := authIdentity.(interface {
				ToClaims() *claims.Claims
			}); ok {
				currentUser := reflect.New(utils.ModelType(context.Auth.Config.UserModel)).Interface()
				if err = tx.First(currentUser, authBasicInfo.ToClaims().UserID).Error; err == nil {
					return currentUser, nil
				}
				return nil, ErrInvalidAccount
			}
		}
		return authIdentity, nil
	}
	return nil, ErrInvalidAccount

}

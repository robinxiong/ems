package password

import (
	"ems/auth"
	"ems/auth/auth_identity"
	"ems/auth/claims"
	"ems/core/utils"
	"ems/session"
	"reflect"
	"strings"
)

func DefaultAuthorizeHandler(context *auth.Context) (*claims.Claims, error) {
	var (
		authInfo auth_identity.AuthIdentity //用于记录认证信息
		req      = context.Request
		provider = context.Provider.(*Provider)
		tx       = context.Auth.GetDB(req)
	)

	req.ParseForm()
	authInfo.Provider = provider.GetName()
	authInfo.UID = strings.TrimSpace(req.Form.Get("login")) //邮箱帐号名

	//从auth_identity表中查找用户名和provider, 如果此用户名不存在，则返回错误
	if tx.Model(context.Auth.AuthIdentityModel).Where(authInfo).Scan(&authInfo).RecordNotFound() {
		return nil, auth.ErrInvalidAccount
	}

	//Confirmable在auth_theme/clean注册provider是指定为true
	//发送验证邮件, 并向用户返回当前帐号没有验证的错误
	if provider.Config.Confirmable && authInfo.ConfirmedAt == nil {
		currentUser, _ := context.Auth.UserStorer.Get(authInfo.ToClaims(), context)
		provider.Config.ConfirmMailer(authInfo.UID, context, authInfo.ToClaims(), currentUser)
		return nil, ErrUnconfirmed
	}

	if err := provider.Encryptor.Compare(authInfo.EncryptedPassword, strings.TrimSpace(req.Form.Get("password"))); err == nil {
		return authInfo.ToClaims(), err
	}
	return nil, auth.ErrInvalidPassword

}
func DefaultRegisterHandler(context *auth.Context) (*claims.Claims, error) {
	var (
		err         error
		req         = context.Request
		schema      auth.Schema
		authInfo    auth_identity.Basic
		provider, _ = context.Provider.(*Provider)
		tx          = context.Auth.GetDB(req)
		currentUser interface{} //保存后的用户信息
	)
	req.ParseForm()

	//获取邮箱或者手机
	uid := req.Form.Get("login")
	password := req.Form.Get("password")

	if uid == "" {
		return nil, auth.ErrInvalidAccount
	}

	if password == "" {
		return nil, auth.ErrInvalidPassword
	}

	authInfo.Provider = provider.GetName()
	authInfo.UID = strings.TrimSpace(uid)

	//根据UID和provider， 查找auth_iidentity表是否注册了用户名, 如果找到, 则返回错误
	if !tx.Model(context.Auth.AuthIdentityModel).Where(authInfo).Scan(&authInfo).RecordNotFound() {
		return nil, auth.ErrInvalidAccount
	}

	//将明文密码加密
	if authInfo.EncryptedPassword, err = provider.Encryptor.Digest(strings.TrimSpace(password)); err == nil {
		//保存用户信息到auth.Schema

		schema.Provider = authInfo.Provider
		schema.UID = authInfo.UID
		schema.Email = authInfo.UID
		schema.RawInfo = req

		//调用auth.UserStorer来保存用户信息
		//context是用来获取request, UserModel
		currentUser, authInfo.UserID, err = context.Auth.UserStorer.Save(&schema, context)

		if err != nil {
			return nil, err
		}

		authIdentity := reflect.New(utils.ModelType(context.Auth.Config.AuthIdentityModel)).Interface()
		//在auth_identity表中创建一个记录, 它保存用户uid和加密后的password
		if err = tx.Where(authInfo).FirstOrCreate(authIdentity).Error; err == nil {
			if provider.Config.Confirmable {
				context.SessionStorer.Flash(context.Writer, req, session.Message{Message: ConfirmFlashMessage, Type: "success"})
				err = provider.Config.ConfirmMailer(schema.Email, context, authInfo.ToClaims(), currentUser)
			}
			return authInfo.ToClaims(), err
		}
	}

	return nil, err
}

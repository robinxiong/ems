package auth

import (
	"ems/auth/claims"
	"ems/session"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// SessionStoreInterface用于将auth.Claims中的信息保存到Session中，或者从session中获取登录信息
// 同时输出错误信息, 这个接口会被auth嵌入, 所以auth也将拥有这些方法
type SessionStorerInterface interface {
	//从request Authorization头中获取token,或者从session中获取到token
	Get(req *http.Request) (*claims.Claims, error)
	//将claims(认证信息)，加密为token, 然后添加到session中
	Update(w http.ResponseWriter, req *http.Request, claims *claims.Claims) error
	// 从session中删除auth_session
	Delete(w http.ResponseWriter, req *http.Request) error
	//向session中添加flash message, 这些message会保存为_flashs, 模板中可以通过context的Flashs方法找回
	//，context也是调用SessionStorer的Flashes方法
	Flash(w http.ResponseWriter, req *http.Request, message session.Message) error
	// 返回Session中的_flashes的message
	Flashes(w http.ResponseWriter, req *http.Request) []session.Message
	// claims转变为token
	SignedToken(claims *claims.Claims) string
	//校验token
	ValidateClaims(tokenString string) (*claims.Claims, error)
}

type SessionStorer struct {
	SessionName    string                   //session的名称，通过这个sessionManager通过这个名称，获取值或者保存session
	SigningMethod  jwt.SigningMethod        //加签的方法 HMAC SHA, RSA, RSA-PSS, and ECDSA
	SignedString   string                   //jwt的公钥
	SessionManager session.ManagerInterface //session module
}

// Get 从request的Authorization头中,获取到token, 然后返回一个Claims
func (sessionStorer *SessionStorer) Get(req *http.Request) (*claims.Claims, error) {
	tokenString := req.Header.Get("Authorization")
	// 或者从cookie中获取
	if tokenString == "" {
		tokenString = sessionStorer.SessionManager.Get(req, sessionStorer.SessionName)
	}
	return sessionStorer.ValidateClaims(tokenString)
}

// Update update claims with session manager
func (sessionStorer *SessionStorer) Update(w http.ResponseWriter, req *http.Request, claims *claims.Claims) error {
	token := sessionStorer.SignedToken(claims)
	return sessionStorer.SessionManager.Add(w, req, sessionStorer.SessionName, token)
}

// Delete delete claims from session manager
func (sessionStorer *SessionStorer) Delete(w http.ResponseWriter, req *http.Request) error {
	sessionStorer.SessionManager.Pop(w, req, sessionStorer.SessionName)
	return nil
}

// Flash add flash message to session data
func (sessionStorer *SessionStorer) Flash(w http.ResponseWriter, req *http.Request, message session.Message) error {

	return sessionStorer.SessionManager.Flash(w, req, message)
}

// Flashes returns a slice of flash messages from session data
func (sessionStorer *SessionStorer) Flashes(w http.ResponseWriter, req *http.Request) []session.Message {

	return sessionStorer.SessionManager.Flashes(w, req)
}

//使用sessionStorer.SignedString对claims加签和解签
// SignedString可以是私钥和单独的公钥
// SignedToken generate signed token with Claims
func (sessionStorer *SessionStorer) SignedToken(claims *claims.Claims) string {
	token := jwt.NewWithClaims(sessionStorer.SigningMethod, claims)
	signedToken, _ := token.SignedString([]byte(sessionStorer.SignedString))

	return signedToken
}

func (sessionStorer *SessionStorer) ValidateClaims(tokenString string) (*claims.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != sessionStorer.SigningMethod {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(sessionStorer.SignedString), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*claims.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

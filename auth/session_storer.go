package auth

import (
	"net/http"
	"ems/auth/claims"
	"github.com/dgrijalva/jwt-go"
	"ems/session"
	"fmt"
	"errors"
)

// SessionStoreInterface用于认证时，获取token中的信息, 它会将认证信息保存到session中
type SessionStoreInterface interface {
	Get(req *http.Request) (*claims.Claims, error)
}

type SessionStorer struct {
	SessionName string  //session的名称，通过这个sessionManager通过这个名称，获取值或者保存session
	SigningMethod jwt.SigningMethod //加签的方法 HMAC SHA, RSA, RSA-PSS, and ECDSA
	SignedString string //jwt的公钥
	SessionManager  session.ManagerInterface//session module
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

func (sessionStorer *SessionStorer) ValidateClaims(tokenString string) (*claims.Claims, error){
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

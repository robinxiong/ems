package auth

import (
	"ems/auth/claims"
	"net/http"
)

// Context context
type Context struct {
	*Auth
	Claims   *claims.Claims
	Provider Provider  //controller/serveMux中，当路径中有/auth/password/login时，找到password provider，然后保存到context
	Request  *http.Request
	Writer   http.ResponseWriter
}

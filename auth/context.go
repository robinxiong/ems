package auth

import (
	"ems/auth/claims"
	"net/http"
	"ems/session"
)

// Context context
type Context struct {
	*Auth
	Claims   *claims.Claims
	Provider Provider  //controller/serveMux中，当路径中有/auth/password/login时，找到password provider，然后保存到context
	Request  *http.Request
	Writer   http.ResponseWriter
}


// Flashes get flash messages
func (context Context) Flashes() []session.Message {
	return context.Auth.SessionStorer.Flashes(context.Writer, context.Request)
}

// FormValue get form value with name
func (context Context) FormValue(name string) string {
	return context.Request.Form.Get(name)
}

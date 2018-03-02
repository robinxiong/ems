package claims

import (
	"time"
	"github.com/dgrijalva/jwt-go"
)

//在jwt将claims的内容加密为token, 在验证token时，又可以解析回claims对像

type Claims struct {
	Provider                         string         `json:"provider,omitempty"`
	UserID                           string         `json:"userid,omitempty"`
	LastLoginAt                      *time.Time     `json:"last_login,omitempty"`
	LastActiveAt                     *time.Time     `json:"last_active,omitempty"`
	LongestDistractionSinceLastLogin *time.Duration `json:"distraction_time,omitempty"`
	jwt.StandardClaims
}



// ToClaims implement ClaimerInterface
func (claims *Claims) ToClaims() *Claims {
	return claims
}

// ClaimerInterface claimer interface
type ClaimerInterface interface {
	ToClaims() *Claims
}

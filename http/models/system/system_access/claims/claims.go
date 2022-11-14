package claims

import "github.com/dgrijalva/jwt-go"

// Claims - токен.
type Claims struct {
	jwt.StandardClaims
	Login string `json:"login" bson:"login" yaml:"login" form:"login" description:"Логин"` // Логин
}

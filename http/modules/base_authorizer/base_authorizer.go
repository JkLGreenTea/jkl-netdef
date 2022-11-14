package base_authorizer

import (
	http_api_server_cfg "JkLNetDef/engine/config/http_api_server"
	"JkLNetDef/engine/http/models/system/system_access/claims"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/models/user"
	"context"
	"crypto/sha1"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Authorizer - авторизатор
type Authorizer struct {
	Utils   *Utils
	Modules *Modules
}

// Utils - утилиты авторизатора.
type Utils struct {
	Config *http_api_server_cfg.Server
}

// Modules - модули авторизатора.
type Modules struct {
	Logger interfacies.Logger
}

// ModelAuth - модель авторизации.
type ModelAuth struct {
	Login    string `json:"login" form:"login" yaml:"login" description:"Логин" validate:"required,min=4,max=32"`           // Логин
	Password string `json:"password" form:"password" yaml:"password" description:"Пароль" validate:"required,min=4,max=32"` // Пароль
}

// New - создание авторизатора
func New(cfg *http_api_server_cfg.Server, log interfacies.Logger) *Authorizer {
	return &Authorizer{
		Utils: &Utils{
			Config: cfg,
		},
		Modules: &Modules{
			Logger: log,
		},
	}
}

// SignIn - подпись токена.
func (authorizer *Authorizer) SignIn(ctx context.Context, us *user.User, tm int64) (*claims.Claims, string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(us.Password))
	pwd.Write([]byte(authorizer.Utils.Config.SystemAccess.Token.HasSalt))

	claims_ := &claims.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tm,
			IssuedAt:  time.Now().Unix(),
		},
		Login: us.Login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims_)
	signedToken, err := token.SignedString([]byte(authorizer.Utils.Config.SystemAccess.Token.SignedKey))

	return claims_, signedToken, err
}

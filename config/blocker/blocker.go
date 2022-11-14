package blocker

import "time"

// Blocker - конфигурация блокировщика.
type Blocker struct {
	// Путь к файлу с IP-адресами по странам (по умолчанию 'system/data/ipv4.json').
	IPv4ByCountry string `json:"ip_v4_by_country" bson:"ip_v4_by_country" yaml:"ip_v4_by_country" form:"ip_v4_by_country" description:"Путь к файлу с IP-адресами по странам (по умолчанию 'system/data/ipv4.json')."`

	// TokenName имя токена. Используется как ключ для записи cookie
	// (по умолчанию 'proxy_token').
	TokenName string `json:"token_name" bson:"token_name" yaml:"token_name" form:"token_name" description:"Имя токена. Используется как ключ для записи cookie (по умолчанию 'proxy_token')."`

	// TokenSecretKey секретный ключ токена.
	TokenSecretKey string `json:"token_secret_key" bson:"token_secret_key" yaml:"token_secret_key" form:"token_secret_key" description:"Секретный ключ токена."`

	// TokenExpire время жизни токена
	// (по умолчанию time.Hour * 4).
	TokenExpire time.Duration `json:"token_expire" bson:"token_expire" yaml:"token_expire" form:"token_expire" description:"Время жизни токена (по умолчанию time.Hour * 4)."`

}
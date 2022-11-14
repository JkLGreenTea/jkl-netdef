package token

import "go.mongodb.org/mongo-driver/bson/primitive"

// Token - токен.
type Token struct {
	ID        primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                                        // ID
	Host      string             `json:"host" bson:"host" yaml:"host" form:"host" description:"IP-адрес"`                           // IP-адрес
	Location  string             `json:"location" bson:"location" yaml:"location" form:"location" description:"Местоположение"`     // Местоположение
	UserAgent string             `json:"user_agent" bson:"user_agent" yaml:"user_agent" form:"user_agent" description:"User-Agent"` // User-Agent

	Date   int64  `json:"date" bson:"date" yaml:"date" form:"date" description:"Дата начала действия"`        // Дата начала действия
	Expire int64  `json:"expire" bson:"expire" yaml:"expire" form:"expire" description:"Дата конца действия"` // Дата конца действия
	Data   string `json:"data" bson:"data" yaml:"data" form:"data" description:"Data информация токена"`      // Data информация токена
}

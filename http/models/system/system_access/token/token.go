package token

import (
	"JkLNetDef/engine/http/models/system/meta_data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token - токен.
type Token struct {
	ID      primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID" `                          // ID
	Owner   primitive.ObjectID `json:"owner" bson:"owner" yaml:"owner" form:"owner" description:"Владелец" `         // Владелец
	Created int64              `json:"created" bson:"created" yaml:"created" form:"created" description:"Создан" `   // Создан
	Expire  int64              `json:"expire" bson:"expire" yaml:"expire" form:"expire" description:"Действует до" ` // Действует до
	Data    string             `json:"data" bson:"data" yaml:"data" form:"data" description:"Дата информация" `      // Дата информация

	NoDelete bool                `json:"no_delete" bson:"no_delete" yaml:"no_delete" form:"no_delete" description:"Запрещено ли удаление" ` // Запрещено ли удаление
	Meta     *meta_data.MetaData `json:"meta" bson:"meta" yaml:"meta" form:"meta" description:"Мета данные" `                               // Мета данные
}

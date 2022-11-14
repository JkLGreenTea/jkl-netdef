package session

import (
	"JkLNetDef/engine/http/models/system/meta_data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session - сессия пользователя.
type Session struct {
	ID    primitive.ObjectID     `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                // ID
	Token primitive.ObjectID     `json:"token" bson:"token" yaml:"token" form:"token" description:"Token"`  // Token
	Data  map[string]interface{} `json:"data" bson:"data" yaml:"data" form:"data" description:"Информация"` // Информация

	NoDelete bool                `json:"no_delete" bson:"no_delete" yaml:"no_delete" form:"no_delete" description:"Запрещено ли удаление"` // Запрещено ли удаление
	Meta     *meta_data.MetaData `json:"meta" bson:"meta" yaml:"meta" form:"meta" description:"Мета данные"`                               // Мета данные
}

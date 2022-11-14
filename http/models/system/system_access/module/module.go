package module

import (
	"JkLNetDef/engine/http/models/system/meta_data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Module - модуль.
type Module struct {
	ID           primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                                                                        // ID
	Title        string             `json:"title" bson:"title" yaml:"title" form:"title" description:"Название"`                                                       // Название
	Description  string             `json:"description" bson:"description" yaml:"description" form:"description" description:"Описание"`                               // Описание
	Locked       bool               `json:"locked" bson:"locked" yaml:"locked" form:"locked" description:"Закрыт ли"`                                                  // Закрыт ли
	Authorized   bool               `json:"authorized" bson:"authorized" yaml:"authorized" form:"authorized" description:"Обязательно быть авторизированным"`          // Обязательно быть авторизированным
	HttpRequests []string           `json:"http_requests" bson:"http_requests" yaml:"http_requests" form:"http_requests" description:"Http запросы входящие в модуль"` // Http запросы входящие в модуль

	Meta *meta_data.MetaData `json:"meta" bson:"meta" yaml:"meta" form:"meta" description:"Meta данные"` // Meta данные
}

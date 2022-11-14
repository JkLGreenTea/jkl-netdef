package http_request

import (
	"JkLNetDef/engine/http/models/system/meta_data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Request - http запрос.
type Request struct {
	ID          primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                                                               // ID
	Method      string             `json:"method" bson:"method" yaml:"method" form:"method" description:"Метод"`                                             // Метод
	URL         string             `json:"url" bson:"url" yaml:"url" form:"url" description:"url"`                                                           // url
	Version     string             `json:"version" bson:"version" yaml:"version" form:"version" description:"Версия"`                                        // Версия
	Locked      bool               `json:"locked" bson:"locked" yaml:"locked" form:"locked" description:"Закрыт ли"`                                         // Закрыт ли
	Authorized  bool               `json:"authorized" bson:"authorized" yaml:"authorized" form:"authorized" description:"Обязательно быть авторизированным"` // Обязательно быть авторизированным
	IsStatic    bool               `json:"is_static" bson:"is_static" yaml:"is_static" form:"is_static" description:"Является ли статичным"`                 // Является ли статичным
	IsSystem    bool               `json:"is_system" bson:"is_system" yaml:"is_system" form:"is_system" description:"Закрыт ли"`                             // Закрыт ли
	Info        string             `json:"info" bson:"info" yaml:"info" form:"info" description:"Информация"`                                                // Информация
	Title       string             `json:"title" yaml:"title" form:"title" description:"Наименование"`                                                       // Наименование
	Description string             `json:"description" bson:"description" yaml:"description" form:"description" description:"Описание"`                      // Описание

	Meta *meta_data.MetaData `json:"meta" bson:"meta" yaml:"meta" form:"meta" description:"Meta данные"` // Meta данные
}

package role

import (
	"JkLNetDef/engine/http/models/system/meta_data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role - роль.
type Role struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                                                                // ID
	Title        string               `json:"title" bson:"title" yaml:"title" form:"title" description:"Название"`                                               // Название
	HttpRequests []string             `json:"http_requests" bson:"http_requests" yaml:"http_requests" form:"http_requests" description:"Доступные http запросы"` // Доступные запросы
	Modules      []primitive.ObjectID `json:"modules" bson:"modules" yaml:"modules" form:"modules" description:"Доступные модули"`                               // Доступные модули

	Meta *meta_data.MetaData `json:"meta" bson:"meta" yaml:"meta" form:"meta" description:"Meta данные"` // Meta данные
}

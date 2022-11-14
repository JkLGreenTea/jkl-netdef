package notifications

import (
	"JkLNetDef/engine/http/models/system/meta_data"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Notification - уведомление.
type Notification struct {
	ID        primitive.ObjectID `json:"id" bson:"_id" yaml:"id" form:"id" description:"ID"`                                          // ID
	Title     string             `json:"title" bson:"title" yaml:"title" form:"title" description:"Заголовок"`                        // Заголовок
	Message   string             `json:"message" bson:"message" yaml:"message" form:"message" description:"Сообщение"`                // Сообщение
	Date      time.Time          `json:"date" bson:"date" yaml:"date" form:"date" description:"Дата"`                                 // Дата
	Recipient primitive.ObjectID `json:"recipient" bson:"recipient" yaml:"recipient" form:"recipient" description:"Получатель"`       // Получатель
	Sender    primitive.ObjectID `json:"sender" bson:"sender" yaml:"sender" form:"sender" description:"Отправитель"`                  // Отправитель
	IsReading bool               `json:"is_reading" bson:"is_reading" yaml:"is_reading" form:"is_reading" description:"Прочитано ли"` // Прочитано ли

	Meta *meta_data.MetaData `json:"meta" bson:"meta" yaml:"meta" form:"meta" description:"Meta данные"` // Meta данные
}

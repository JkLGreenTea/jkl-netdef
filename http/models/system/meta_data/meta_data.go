package meta_data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// MetaData - мета данные.
type MetaData struct {
	DateCreated time.Time          `json:"date_created" bson:"date_created" yaml:"date_created" form:"date_created" description:"Когда создан"`  // Когда создан
	Created     primitive.ObjectID `json:"created" bson:"created" yaml:"created" form:"created" description:"Кем создан"`                        // Кем создан
	DateChanged time.Time          `json:"date_changed" bson:"date_changed" yaml:"date_changed" form:"date_changed" description:"Когда изменён"` // Когда изменён
	Changed     primitive.ObjectID `json:"changed" bson:"changed" yaml:"changed" form:"changed" description:"Кем изменён"`                       // Кем изменён
}

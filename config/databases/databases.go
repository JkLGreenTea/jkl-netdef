package databases

// Databases - конфигурация баз данных.
type Databases struct {
	// Конфигурация MongoDB.
	Mongo *MongoDB `json:"mongodb" bson:"mongodb" yaml:"mongodb" form:"mongodb" description:"Конфигурация MongoDB."`
}
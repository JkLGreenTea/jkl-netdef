package blocker

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	"JkLNetDef/engine/interfacies"
	"go.mongodb.org/mongo-driver/mongo"
)

// Blocker - запросы блокировщика.
type Blocker struct {
	Title  string        // Название
	Client *mongo.Client // Клиент бд

	Config *databases_cfg.MongoDB // Конфиг
	Logger interfacies.Logger     // Логгер
}

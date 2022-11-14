package controller_reputation

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	"JkLNetDef/engine/interfacies"
	"go.mongodb.org/mongo-driver/mongo"
)

// ControllerReputation - запросы контроллера репутации.
type ControllerReputation struct {
	Title  string        // Название
	Client *mongo.Client // Клиент бд

	Config *databases_cfg.MongoDB // Конфиг
	Logger interfacies.Logger     // Логгер
}

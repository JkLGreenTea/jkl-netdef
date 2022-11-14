package background_process

import (
	"JkLNetDef/engine/databases"
	"JkLNetDef/engine/http/utils/postman"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/utils/synchronizer"
)

// Process - фоновой процесс.
type Process struct {
	Databases *databases.Databases
	Utils     *Utils
	Modules   *Modules
}

// Utils - утилиты.
type Utils struct {
	Sync    *synchronizer.Synchronizer
	Postman *postman.Postman
}

// Modules - модули.
type Modules struct {
	Logger               interfacies.Logger
	HttpServerLogger     interfacies.HttpServerLogger
	ManagerMetadata      interfacies.ManagerMetaData
	ManagerNotifications interfacies.ManagerNotifications
	SystemAccess         interfacies.SystemAccess
	ManagerSessions      interfacies.ManagerSessions
	Authorizer           interfacies.Authorizer
	Validator            interfacies.Validator
}

type HandlerFunc func(process *Process) error

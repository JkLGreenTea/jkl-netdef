package mongo

import (
	"JkLNetDef/engine/databases/mongo/collections/blocker"
	"JkLNetDef/engine/databases/mongo/collections/controller_reputation"
	"JkLNetDef/engine/databases/mongo/collections/main_collections"
	"JkLNetDef/engine/databases/mongo/collections/system_access_collections"
	"JkLNetDef/engine/databases/mongo/collections/system_collections"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB - MongoDB.
type MongoDB struct {
	Client *mongo.Client

	Blocker              *blocker.Blocker
	ControllerReputation *controller_reputation.ControllerReputation
	SystemAccess         *system_access_collections.SystemAccess
	System               *system_collections.System
	Main                 *main_collections.Main
}

// Disconnect - разрыв соединения с MongoDB.
func (database *MongoDB) Disconnect() error {
	return database.Client.Disconnect(context.Background())
}

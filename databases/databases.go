package databases

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	my_mongodb "JkLNetDef/engine/databases/mongo"
	"JkLNetDef/engine/databases/mongo/collections/blocker"
	"JkLNetDef/engine/databases/mongo/collections/controller_reputation"
	"JkLNetDef/engine/databases/mongo/collections/main_collections"
	"JkLNetDef/engine/databases/mongo/collections/system_access_collections"
	"JkLNetDef/engine/databases/mongo/collections/system_collections"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Databases struct {
	Title   string
	Mongo   *my_mongodb.MongoDB
	Config  *databases_cfg.MongoDB
	Modules *Modules
}

// Modules - модули базы данных.
type Modules struct {
	Logger interfacies.Logger
}

// New - создание баз данных.
func New(title string, cfg *databases_cfg.Databases, log interfacies.Logger) (*Databases, error) {
	databases := &Databases{
		Title: title,
		Modules: &Modules{
			Logger: log,
		},
		Config: cfg.Mongo,
	}

	// MongoDB
	{
		client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin", databases.Config.UserName, databases.Config.Password, databases.Config.Addr)))
		if err != nil {
			return nil, err
		}

		if client != nil {
			err = client.Connect(context.TODO())
			if err != nil {
				return nil, err
			}
		} else {
			err = errors.New("Client is nil. ")
			return nil, err
		}

		databases.Mongo = &my_mongodb.MongoDB{
			Client: client,

			Blocker: &blocker.Blocker{
				Title:  "DB-Blocker",
				Client: client,
				Config: cfg.Mongo,
				Logger: log,
			},
			ControllerReputation: &controller_reputation.ControllerReputation{
				Title:  "DB-ControllerReputation",
				Client: client,
				Config: cfg.Mongo,
				Logger: log,
			},
			SystemAccess: &system_access_collections.SystemAccess{
				HttpRequests: &system_access_collections.SysHttpRequests{
					Title:  "DB-SysHttpRequests",
					Client: client,
					Config: cfg.Mongo,
					Logger: log,
				},
				Modules: &system_access_collections.SysModules{
					Title:  "DB-SysModules",
					Client: client,
					Config: cfg.Mongo,
					Logger: log,
				},
				Tokens: &system_access_collections.SysTokens{
					Title:  "DB-SysTokens",
					Client: client,
					Config: cfg.Mongo,
					Logger: log,
				},
				Roles: &system_access_collections.SysRoles{
					Title:  "DB-SysRoles",
					Client: client,
					Config: cfg.Mongo,
					Logger: log,
				},
			},
			System: &system_collections.System{
				Sessions: &system_collections.SysSessions{
					Title:  "DB-SysSessions",
					Client: client,
					Config: cfg.Mongo,
					Logger: log,
				},
				Notifications: &system_collections.SysNotifications{
					Title:  "DB-SysNotifications",
					Client: client,
					Config: cfg.Mongo,
					Logger: log,
				},
			},
			Main: &main_collections.Main{
				Users: &main_collections.Users{
					Title:  "DB-Users",
					Client: client,
					Config: cfg.Mongo,
					Logger: log,
				},
			},
		}
	}

	return databases, nil
}

// Disconnect - Завершения соединений с базами данных.
func (databases *Databases) Disconnect() error {
	if err := databases.Mongo.Disconnect(); err != nil {
		databases.Modules.Logger.FATAL(base_logger.Message{
			Sender: databases.Title,
			Text:   err.Error(),
		})
	}

	return nil
}

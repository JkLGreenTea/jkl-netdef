package http

import (
	"JkLNetDef/engine/config"
	"JkLNetDef/engine/http/models/system/schema"
	"JkLNetDef/engine/http/modules/base_authorizer"
	"JkLNetDef/engine/http/modules/base_http_api_logger"
	"JkLNetDef/engine/http/modules/base_manager_notifications"
	"JkLNetDef/engine/http/modules/base_manager_sessions"
	"JkLNetDef/engine/http/modules/base_system_access"
	"JkLNetDef/engine/http/modules/base_validator"
	"JkLNetDef/engine/http/utils/postman"
	"JkLNetDef/engine/modules/base_logger"
	"JkLNetDef/engine/modules/base_manager_metadata"
	"JkLNetDef/engine/services"
	"JkLNetDef/engine/services/users_service"
	"JkLNetDef/engine/services/users_service/users_service_footer"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// buildUtils - создание утилит.
func (engine *Engine) buildUtils() error {
	postman_ := postman.New(engine.Config.Engine, engine.Modules.Logger)

	valid, err := base_validator.New(engine.Modules.Logger)
	if err != nil {
		engine.Modules.Logger.FATAL(base_logger.Message{
			Sender: engine.Title,
			Text:   err.Error(),
		})

		return err
	}

	engine.Modules.Validator = valid
	engine.Utils.Postman = postman_

	return nil
}

// buildModules - создание модулей.
func (engine *Engine) buildModules() error {
	managerSessions := new(base_manager_sessions.Manager)
	managerMetadata := new(base_manager_metadata.Manager)
	systemAccess := new(base_system_access.SystemAccess)
	managerNotifications := new(base_manager_notifications.Manager)
	authorizer := new(base_authorizer.Authorizer)

	// CREATE MODULES
	{
		// authorizer
		{
			authorizer.Utils = &base_authorizer.Utils{
				Config: engine.Config.Engine,
			}
		}

		// manager sessions
		{
			managerSessions.Databases = engine.Databases
		}

		// manager metadata
		{
			managerMetadata.Modules = &base_manager_metadata.Modules{
				Logger:          engine.Modules.Logger,
				ManagerSessions: managerSessions,
			}
		}

		// system access
		{
			systemAccess.Databases = engine.Databases
			systemAccess.Utils = &base_system_access.Utils{
				Config: engine.Config.Engine,
			}
		}

		// manager notifications
		{
			managerNotifications.Databases = engine.Databases
		}
	}

	// SET MODULES
	{
		// authorizer
		{
			authorizer.Modules = &base_authorizer.Modules{
				Logger: engine.Modules.Logger,
			}
		}

		// manager sessions
		{
			managerSessions.Modules = &base_manager_sessions.Modules{
				ManagerMetadata:      managerMetadata,
				ManagerNotifications: managerNotifications,
				Logger:               engine.Modules.Logger,
			}
		}

		// manager metadata
		{
			managerMetadata.Modules = &base_manager_metadata.Modules{
				Logger:          engine.Modules.Logger,
				ManagerSessions: managerSessions,
			}
		}

		// system access
		{
			systemAccess.Modules = &base_system_access.Modules{
				Logger:          engine.Modules.Logger,
				ManagerMetadata: managerMetadata,
				ManagerSessions: managerSessions,
				Authorizer:      authorizer,
			}
		}

		// manager notifications
		{
			managerNotifications.Modules = &base_manager_notifications.Modules{
				Logger:          engine.Modules.Logger,
				ManagerMetadata: managerMetadata,
			}
		}
	}

	engine.Modules.Authorizer = authorizer
	engine.Modules.SystemAccess = systemAccess
	engine.Modules.ManagerSessions = managerSessions
	engine.Modules.ManagerNotifications = managerNotifications
	engine.Modules.ManagerMetaData = managerMetadata
	engine.Modules.HttpServerLogger = base_http_api_logger.New("Http-Log", engine.Config.ApiLogger, engine.Modules.Logger, engine.Utils.Synchronizer, managerSessions)

	return nil
}

// buildServices - создание сервисов.
func (engine *Engine) buildServices(cfg *config.GlobalConfig) error {
	engine.Services = &services.Services{
		Users: &users_service.Service{
			Title: "Serv-Users",
			Footer: &users_service_footer.ServiceFooter{
				Title:     "ServF-Users",
				Databases: engine.Databases,
				Config:    cfg,
				Utils: &users_service_footer.Utils{
					Sync: engine.Utils.Synchronizer,
				},
				Modules: &users_service_footer.Modules{
					Logger:          engine.Modules.Logger,
					ManagerMetaData: engine.Modules.ManagerMetaData,
					ManagerSessions: engine.Modules.ManagerSessions,
					Validator:       engine.Modules.Validator,
				},
			},
			Utils: &users_service.Utils{
				Sync: engine.Utils.Synchronizer,
			},
			Modules: &users_service.Modules{
				Logger:          engine.Modules.Logger,
				ManagerMetaData: engine.Modules.ManagerMetaData,
				ManagerSessions: engine.Modules.ManagerSessions,
				Validator:       engine.Modules.Validator,
			},
		},
	}

	return nil
}

// getRequests - подгрузка запросов.
func (engine *Engine) loadRequests(schem *schema.Schema) error {
	requests := schem.Requests

	routes := engine.httpGin.Routes()

	for index, req := range requests {
		for _, route := range routes {
			if req.Method == route.Method {
				if route.Path == req.URL || route.Path == req.URL+"/" {
					req.URL = route.Path
				}
			}
		}

		bdReq, err := engine.Databases.Mongo.SystemAccess.HttpRequests.GetByURLAndMethod(req.Method, req.URL)
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				engine.Modules.Logger.INFO(base_logger.Message{
					Sender: engine.Title,
					Text:   fmt.Sprintf("Запрос '%s %s' не загружен в систему, начата загрузка... ", req.Method, req.URL),
				})
			} else {
				engine.Modules.Logger.WARN(base_logger.Message{
					Sender: engine.Title,
					Text:   err.Error(),
				})

				return err
			}
		}

		if bdReq == nil {
			// Заполнение данных
			{
				req.ID = primitive.NewObjectID()

				// Meta
				{
					meta, _, err := engine.Modules.ManagerMetaData.NewMeta(nil)
					if err != nil {
						engine.Modules.Logger.WARN(base_logger.Message{
							Sender: engine.Title,
							Text:   err.Error(),
						})

						return err
					}

					req.Meta = meta
				}
			}

			err = engine.Databases.Mongo.SystemAccess.HttpRequests.Add(req)
			if err != nil {
				engine.Modules.Logger.WARN(base_logger.Message{
					Sender: engine.Title,
					Text:   err.Error(),
				})
			}

			engine.Modules.Logger.INFO(base_logger.Message{
				Sender: engine.Title,
				Text:   fmt.Sprintf("Запрос '%s %s' загружен. ", req.Method, req.URL),
			})
		} else {
			schem.Requests[index] = bdReq
		}
	}

	for _, group := range schem.Groups {
		err := engine.loadRequests(group)
		if err != nil {
			engine.Modules.Logger.WARN(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})
		}
	}

	return nil
}

package engine

import (
	"JkLNetDef/engine/config"
	"JkLNetDef/engine/databases"
	http_api "JkLNetDef/engine/http/engine"
	"JkLNetDef/engine/http/modules/base_manager_notifications"
	"JkLNetDef/engine/http/modules/base_manager_sessions"
	"JkLNetDef/engine/http/modules/base_validator"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/models/controller_reputation"
	"JkLNetDef/engine/modules/base_blocker"
	"JkLNetDef/engine/modules/base_controller_reputation"
	"JkLNetDef/engine/modules/base_controller_technical_works"
	"JkLNetDef/engine/modules/base_logger"
	"JkLNetDef/engine/modules/base_manager_metadata"
	"JkLNetDef/engine/proxy"
	"JkLNetDef/engine/services"
	"JkLNetDef/engine/services/users_service"
	"JkLNetDef/engine/services/users_service/users_service_footer"
	"JkLNetDef/engine/utils/synchronizer"
	"sync"
)

// Engine - движок.
type Engine struct {
	ProxyStorage *proxy.Storage

	Api       *Api
	Services  *services.Services
	Modules   *Modules
	Utils     *Utils
	Databases *databases.Databases
	Config    *config.GlobalConfig

	states       *States
	currentState State

	chanStop chan struct{}
}

// Modules - модули.
type Modules struct {
	Logger                   interfacies.Logger
	ControllerReputation     interfacies.ControllerReputation
	ControllerTechnicalWorks interfacies.ControllerTechnicalWorks
	Blocker                  interfacies.Blocker
	ManagerMetaData          interfacies.ManagerMetaData
	ManagerSessions          interfacies.ManagerSessions
	ManagerNotifications     interfacies.ManagerNotifications
	Validator                interfacies.Validator
}

// Utils - утилиты.
type Utils struct {
	Synchronizer *synchronizer.Synchronizer
}

// Api - Api для взаимодействия с прокси движком.
type Api struct {
	Http *http_api.Engine
}

// New - создание движка.
func New(cfg *config.GlobalConfig) (*Engine, error) {
	eng := &Engine{
		Modules: new(Modules),
		Utils:   new(Utils),

		Config:   cfg,
		chanStop: make(chan struct{}),
	}

	// Элементы
	{
		// Основа
		{
			// Synchronizer
			{
				eng.Utils.Synchronizer = synchronizer.New()
			}

			// Logger
			{
				log, err := base_logger.New(eng.Config.Loggers.Global, eng.Utils.Synchronizer)
				if err != nil {
					return nil, err
				}
				eng.Modules.Logger = log
			}

			// db
			{
				dbs, err := databases.New("Databases", cfg.Databases, eng.Modules.Logger)
				if err != nil {
					return nil, err
				}

				eng.Databases = dbs
			}

			// BLocker
			{
				blocker, err := base_blocker.New("Blocker", cfg.Blocker, eng.Modules.Logger, eng.Databases)
				if err != nil {
					return nil, err
				}

				eng.Modules.Blocker = blocker
			}

			// ManagerNotifications
			{
				eng.Modules.ManagerNotifications = &base_manager_notifications.Manager{
					Title:     "ManagerNotifications",
					Databases: eng.Databases,
				}
			}

			// ManagerSessions
			{
				eng.Modules.ManagerSessions = &base_manager_sessions.Manager{
					Title:     "ManagerSessions",
					Databases: eng.Databases,
				}
			}

			// ManagerMetaData
			{
				eng.Modules.ManagerMetaData = base_manager_metadata.New("Manager-MetaData", eng.Modules.Logger, eng.Modules.ManagerSessions)
			}

			// ManagerMetaData
			{
				var err error

				eng.Modules.Validator, err = base_validator.New(eng.Modules.Logger)
				if err != nil {
					return nil, err
				}
			}

			// ControllerReputation
			{
				eng.Modules.ControllerReputation = &base_controller_reputation.ControllerReputation{
					Config:            cfg.ControllerReputation,
					Clients:           make([]*controller_reputation.Client, 0),
					Hosts:             make(map[string]*controller_reputation.Host),
					CounterResetScore: make(map[string]uint8),

					CounterResetScoreRWMutex: new(sync.RWMutex),
					ClientsRwMutex:           new(sync.RWMutex),
					HostsRwMutex:             new(sync.RWMutex),

					Synchronizer: eng.Utils.Synchronizer,
				}
			}

			// ControllerTechnicalWorks
			{
				eng.Modules.ControllerTechnicalWorks = base_controller_technical_works.New()
			}
		}

		// Модули
		{
			// ManagerNotifications
			{
				eng.Modules.ManagerNotifications.(*base_manager_notifications.Manager).Modules = &base_manager_notifications.Modules{
					Logger:          eng.Modules.Logger,
					ManagerMetadata: eng.Modules.ManagerMetaData,
				}
			}

			// ManagerSessions
			{
				eng.Modules.ManagerSessions.(*base_manager_sessions.Manager).Modules = &base_manager_sessions.Modules{
					ManagerMetadata:      eng.Modules.ManagerMetaData,
					ManagerNotifications: eng.Modules.ManagerNotifications,
					Logger:               eng.Modules.Logger,
				}
			}
		}
	}

	// Состояния
	{
		eng.states = &States{
			Run:  &engineStateRun{eng},
			Stop: &engineStateStop{eng},
		}

		eng.currentState = eng.states.Stop
	}

	// Хранилище прокси
	{
		proxyStorage, err := proxy.NewStorage()
		if err != nil {
			return nil, err
		}

		eng.ProxyStorage = proxyStorage
	}

	// Сервисы
	{
		eng.Services = &services.Services{
			Users: &users_service.Service{
				Title: "Serv-Users",
				Footer: &users_service_footer.ServiceFooter{
					Title:     "ServF-Users",
					Databases: eng.Databases,
					Config:    eng.Config,
					Utils: &users_service_footer.Utils{
						Sync: eng.Utils.Synchronizer,
					},
					Modules: &users_service_footer.Modules{
						Logger:          eng.Modules.Logger,
						ManagerMetaData: eng.Modules.ManagerMetaData,
						ManagerSessions: eng.Modules.ManagerSessions,
						Validator:       eng.Modules.Validator,
					},
				},
				Utils: &users_service.Utils{
					Sync: eng.Utils.Synchronizer,
				},
				Modules: &users_service.Modules{
					Logger:          eng.Modules.Logger,
					ManagerMetaData: eng.Modules.ManagerMetaData,
					ManagerSessions: eng.Modules.ManagerSessions,
					Validator:       eng.Modules.Validator,
				},
			},
		}
	}

	// Api
	{
		eng.Api = new(Api)

		httpServer, err := http_api.Build(eng.Config, eng.Utils.Synchronizer, eng.Modules.Logger, eng.Databases, eng.Modules.Blocker)
		if err != nil {
			return nil, err
		}

		eng.Api.Http = httpServer

		if err = eng.Api.Http.Init(); err != nil {
			return nil, err
		}
	}

	eng.Modules.Logger.INFO(base_logger.Message{
		Sender: "ENGINE",
		Text:   "Экземпляр движка создан. ",
	})

	return eng, nil
}

// Run - запуск движка.
func (engine *Engine) Run() error {
	return engine.currentState.run()
}

// Stop - остановка движка.
func (engine *Engine) Stop() error {
	return engine.currentState.stop()
}

// NewProxy - создание прокси.
func (engine *Engine) NewProxy(cfg *config.ProxyConfig) (interfacies.HttpProxy, error) {
	return engine.currentState.newProxy(cfg)
}

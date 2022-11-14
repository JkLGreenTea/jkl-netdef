package databases

// MongoDB - конфигурация базы данных Mongo.
type MongoDB struct {
	// UserName имя пользователя базы данных.
	UserName string `json:"user_name" bson:"user_name" yaml:"user_name" form:"user_name" description:"Имя пользователя базы данных."`

	// Password пароль пользователя базы данных.
	Password string `json:"password" bson:"password" yaml:"password" form:"password" description:"Пароль пользователя базы данных."`

	// Addr адрес базы данных.
	Addr string `json:"addr" bson:"addr" yaml:"addr" form:"addr" description:"Адрес базы данных."`

	// Blocker конфигурация таблицы блокировщика.
	Blocker *MongoDBBlocker `json:"blocker" bson:"blocker" yaml:"blocker" form:"blocker" description:"Конфигурация таблицы блокировщика."`

	// ReputationController конфигурация таблицы контроллера репутации клиентов.
	ReputationController *MongoDBReputationController `json:"reputation_controller" bson:"reputation_controller" yaml:"reputation_controller" form:"reputation_controller" description:"Конфигурация таблицы контроллера репутации клиентов."`

	// System конфигурация системной таблицы.
	System *MongodbSystem `json:"system" bson:"system" yaml:"system" form:"system" description:"Конфигурация системной таблицы."`

	// Main конфигурация главной таблицы.
	Main *MongoDBMain `json:"main" bson:"main" yaml:"main" form:"main" description:"Конфигурация главной таблицы."`

	// Dashboard конфигурация таблицы админ панели.
	Dashboard *MongoDBDashboard `json:"dashboard" bson:"dashboard" yaml:"dashboard" form:"dashboard" description:"Конфигурация таблицы админ панели.."`
}

// MongoDBBlocker - конфигурации бд блокировщика.
type MongoDBBlocker struct {
	// DatabaseName название таблицы
	// (по умолчанию 'jkl-netdef_blocker').
	DatabaseName string `json:"database_name" bson:"database_name" yaml:"database_name" form:"database_name" description:"Название таблицы (по умолчанию 'jkl-netdef_blocker')."`

	// CollectionLocationBlackList коллекция черного списка местоположений клиентов
	// (по умолчанию 'location_black_list').
	CollectionLocationBlackList string `json:"collection_location_black_list" bson:"collection_location_black_list" yaml:"collection_location_black_list" form:"collection_location_black_list" description:"Коллекция черного списка местоположений клиентов (по умолчанию 'location_black_list')."`

	// CollectionLocationWhiteList коллекция белого списка местоположений клиентов
	// (по умолчанию 'location_white_list').
	CollectionLocationWhiteList string `json:"collection_location_white_list" bson:"collection_location_white_list" yaml:"collection_location_white_list" form:"collection_location_white_list" description:"Коллекция белого списка местоположений клиентов (по умолчанию 'location_white_list')."`

	// CollectionHostBlackList коллекция черного списка хостов
	// (по умолчанию 'host_black_list').
	CollectionHostBlackList string `json:"collection_host_black_list" bson:"collection_host_black_list" yaml:"collection_host_black_list" form:"collection_host_black_list" description:"Коллекция черного списка хостов (по умолчанию 'host_black_list')."`

	// CollectionHostWhiteList коллекция белого списка хостов
	// (по умолчанию 'host_white_list').
	CollectionHostWhiteList string `json:"collection_host_white_list" bson:"collection_host_white_list" yaml:"collection_host_white_list" form:"collection_host_white_list" description:"Коллекция белого списка хостов (по умолчанию 'host_white_list')."`

	// CollectionHostHardWhiteList коллекция жесткого белого списка хостов
	// (по умолчанию 'host_hard_white_list').
	CollectionHostHardWhiteList string `json:"collection_host_hard_white_list" bson:"collection_host_hard_white_list" yaml:"collection_host_hard_white_list" form:"collection_host_hard_white_list" description:"Коллекция жесткого белого списка хостов (по умолчанию 'host_hard_white_list')."`

	// CollectionClientListOnCaptchaCheck коллекция списка клиентов на каптчу
	// (по умолчанию 'client_list_on_captcha_check').
	CollectionClientListOnCaptchaCheck string `json:"collection_client_list_on_captcha_check" bson:"collection_client_list_on_captcha_check" yaml:"collection_client_list_on_captcha_check" form:"collection_client_list_on_captcha_check" description:"Коллекция списка клиентов на каптчу (по умолчанию 'client_list_on_captcha_check')."`

	// CollectionClientListOnCaptchaCheck коллекция списка хостов в бане
	// (по умолчанию 'host_ban_list').
	CollectionHostBanList string `json:"collection_host_ban_list" bson:"collection_host_ban_list" yaml:"collection_host_ban_list" form:"collection_host_ban_list" description:"Коллекция списка хостов в бане (по умолчанию 'host_ban_list')."`

	// CollectionTokens коллекция токенов.
	// (по умолчанию 'tokens').
	CollectionTokens string `json:"collection_tokens" bson:"collection_tokens" yaml:"collection_tokens" form:"collection_tokens" description:"Коллекция токенов. (по умолчанию 'tokens')."`

	// CollectionUserAgentWhiteList коллекция белого списка User-Agent
	// (по умолчанию 'user-agent_white_list').
	CollectionUserAgentWhiteList string `json:"collection_user_agent_white_list" bson:"collection_user_agent_white_list" yaml:"collection_user_agent_white_list" form:"collection_user_agent_white_list" description:"Коллекция белого списка User-Agent (по умолчанию 'user-agent_white_list')."`

	// CollectionUserAgentBlackList коллекция черного списка User-Agent
	// (по умолчанию 'user-agent_black_list').
	CollectionUserAgentBlackList string `json:"collection_user_agent_black_list" bson:"collection_user_agent_black_list" yaml:"collection_user_agent_black_list" form:"collection_user_agent_black_list" description:"Коллекция черного списка User-Agent (по умолчанию 'user-agent_black_list')."`
}

// MongoDBReputationController - конфигурации бд контроллера репутации.
type MongoDBReputationController struct {
	// DatabaseName название таблицы.
	// (по умолчанию 'jkl-netdef_controller_reputation').
	DatabaseName string `json:"database_name" bson:"database_name" yaml:"database_name" form:"database_name" description:"Название таблицы. (по умолчанию 'jkl-netdef_controller_reputation')."`

	// CollectionClientReputation коллекция репутации клиентов
	// (по умолчанию 'client_reputation').
	CollectionClientReputation string `json:"collection_client_reputation" bson:"collection_client_reputation" yaml:"collection_client_reputation" form:"collection_client_reputation" description:"Коллекция репутации клиентов (по умолчанию 'client_reputation')."`
}

// MongodbSystem - конфигурации бд системы.
type MongodbSystem struct {
	// DatabaseName название таблицы.
	// (по умолчанию 'jkl-netdef_sys').
	DatabaseName string `json:"database_name" bson:"database_name" yaml:"database_name" form:"database_name" description:"Название таблицы (по умолчанию 'jkl-netdef_sys')."`

	// CollectionSessions коллекция сессий
	// (по умолчанию 'sessions').
	CollectionSessions string `json:"collection_sessions" bson:"collection_sessions" yaml:"collection_sessions" form:"collection_sessions" description:"Коллекция сессий (по умолчанию 'sessions')."`

	// CollectionRoles коллекция ролей
	// (по умолчанию 'roles').
	CollectionRoles string `json:"collection_roles" bson:"collection_roles" yaml:"collection_roles" form:"collection_roles" description:"Коллекция ролей (по умолчанию 'roles')."`

	// CollectionHttpRequests коллекция http запросов
	// (по умолчанию 'http_requests').
	CollectionHttpRequests string `json:"collection_http_requests" bson:"collection_http_requests" yaml:"collection_http_requests" form:"collection_http_requests" description:"Коллекция http запросов (по умолчанию 'http_requests')."`

	// CollectionTokens коллекция системных токенов
	// (по умолчанию 'tokens').
	CollectionTokens string `json:"collection_tokens" bson:"collection_tokens" yaml:"collection_tokens" form:"collection_tokens" description:"Коллекция системных токенов (по умолчанию 'tokens')."`

	// CollectionNotifications коллекция уведомлений
	// (по умолчанию 'notifications').
	CollectionNotifications string `json:"collection_notifications" bson:"collection_notifications" yaml:"collection_notifications" form:"collection_notifications" description:"Коллекция уведомлений (по умолчанию 'notifications')."`

	// CollectionModules коллекция модулей
	// (по умолчанию 'modules').
	CollectionModules string `json:"collection_modules" bson:"collection_modules" yaml:"collection_modules" form:"collection_modules" description:"Коллекция модулей (по умолчанию 'modules')."`
}

// MongoDBMain - конфигурации бд основных данных.
type MongoDBMain struct {
	// DatabaseName название таблицы.
	// (по умолчанию 'jkl-netdef_main').
	DatabaseName string `json:"database_name" bson:"database_name" yaml:"database_name" form:"database_name" description:"Название таблицы. (по умолчанию 'jkl-netdef_main')."`

	// CollectionUser коллекция пользователей
	// (по умолчанию 'users').
	CollectionUser string `json:"collection_user" bson:"collection_user" yaml:"collection_user" form:"collection_user" description:"Коллекция пользователей (по умолчанию 'users')."`
}

// MongoDBDashboard - конфигурации бд админ панели.
type MongoDBDashboard struct {
	// DatabaseName название таблицы.
	// (по умолчанию 'jkl-netdef_dashboard').
	DatabaseName string `json:"database_name" bson:"database_name" yaml:"database_name" form:"database_name" description:"Название таблицы. (по умолчанию 'jkl-netdef_dashboard')."`

	// CollectionMainMenu коллекция главного меню админ панели
	// (по умолчанию 'main_menu').
	CollectionMainMenu string `json:"collection_main_menu" bson:"collection_main_menu" yaml:"collection_main_menu" form:"collection_main_menu" description:"CollectionMainMenu коллекция главного меню админ панели (по умолчанию 'main_menu')."`
}

package config

import (
	blocker_cfg "JkLNetDef/engine/config/blocker"
	controller_reputation_cfg "JkLNetDef/engine/config/controller_reputation"
	databases_cfg "JkLNetDef/engine/config/databases"
	http_api_server_cfg "JkLNetDef/engine/config/http_api_server"
	http_client_cfg "JkLNetDef/engine/config/http_client"
	http_proxy_cfg "JkLNetDef/engine/config/http_proxy"
	loggers_cfg "JkLNetDef/engine/config/loggers"
	system_cfg "JkLNetDef/engine/config/system"
	"net/http"
	"time"
)

// GlobalConfig - конфигурация всей системы.
type GlobalConfig struct {
	HttpApiServer *http_api_server_cfg.Server `json:"http_api_server" bson:"http_api_server" yaml:"http_api_server" form:"http_api_server" description:"Http сервера"` // Http Api
	Proxies       map[string]*ProxyConfig `json:"proxies" bson:"proxies" yaml:"proxies" form:"proxies" description:"Proxy сервера"` // Proxy сервера

	Loggers   *loggers_cfg.Loggers     `json:"loggers" bson:"loggers" yaml:"loggers" form:"loggers" description:"Логгеры"`             // Логгеры
	System    *system_cfg.System       `json:"system" bson:"system" yaml:"system" form:"system" description:"Системный конфиг"`        // Системный конфиг
	Databases *databases_cfg.Databases `json:"databases" bson:"databases" yaml:"databases" form:"databases" description:"Базы данных"` // Базы данных

	Blocker              *blocker_cfg.Blocker                            `json:"blocker" bson:"blocker" yaml:"blocker" form:"blocker" description:"Блокировщик"`                                                                  // Блокировщик
	ControllerReputation *controller_reputation_cfg.ControllerReputation `json:"controller_reputation" bson:"controller_reputation" yaml:"controller_reputation" form:"controller_reputation" description:"Контроллер репутации"` // Контроллер репутации
}

// ProxyConfig - конфигурация прокси.
type ProxyConfig struct {
	Engine        *http_proxy_cfg.HTTPProxy   `json:"engine" bson:"engine" yaml:"engine" form:"engine" description:"Http proxy движок"`                                // Http proxy
	HttpApiServer *http_api_server_cfg.Server `json:"http_api_server" bson:"http_api_server" yaml:"http_api_server" form:"http_api_server" description:"Http сервера"` // Http Api

	GlobalLogger    *loggers_cfg.Global    `json:"global_logger" bson:"global_logger" yaml:"global_logger" form:"global_logger" description:"Глобальный логгер"`                 // Глобальный логгер
	HttpProxyLogger *loggers_cfg.HttpProxy `json:"http_proxy_logger" bson:"http_proxy_logger" yaml:"http_proxy_logger" form:"http_proxy_logger" description:"Http proxy логгер"` // Http proxy логгер

	HTTPClient *http_client_cfg.HTTPClient `json:"http_client" bson:"http_client" yaml:"http_client" form:"http_client" description:"Http client (http.Client)"` // Http client
	Databases  *databases_cfg.Databases    `json:"databases" bson:"databases" yaml:"databases" form:"databases" description:"Базы данных"`                       // Базы данных

	Blocker              *blocker_cfg.Blocker                            `json:"blocker" bson:"blocker" yaml:"blocker" form:"blocker" description:"Блокировщик"`                                                                  // Блокировщик
	ControllerReputation *controller_reputation_cfg.ControllerReputation `json:"controller_reputation" bson:"controller_reputation" yaml:"controller_reputation" form:"controller_reputation" description:"Контроллер репутации"` // Контроллер репутации
}

// HttpApiConfig - конфигурация Http Api.
type HttpApiConfig struct {
	Engine *http_api_server_cfg.Server `json:"http_api_server" bson:"http_api_server" yaml:"http_api_server" form:"http_api_server" description:"Http сервера"` // Http Api

	ApiLogger *loggers_cfg.Api `json:"api_logger" bson:"api_logger" yaml:"api_logger" form:"api_logger" description:"Конфигурация api логгера."`
}

// Default - получить дефолтную конфигурацию всей системы.
func Default() *GlobalConfig {
	return &GlobalConfig{
		HttpApiServer: defaultHttpApiServer(),
		Proxies:       make(map[string]*ProxyConfig),

		Loggers:   defaultLoggers(),
		System:    defaultSystem(),
		Databases: defaultDatabases(),

		Blocker:              defaultBlocker(),
		ControllerReputation: defaultControllerReputation(),
	}
}

// DefaultProxyConfig - получить дефолтную конфигурацию движка прокси сервера.
func DefaultProxyConfig() *ProxyConfig {
	return &ProxyConfig{
		HttpApiServer: defaultHttpApiServer(),
		Engine:        defaultHTTPProxy(),

		GlobalLogger:    defaultLoggers().Global,
		HttpProxyLogger: defaultLoggers().HttpProxy,

		HTTPClient: defaultHTTPClient(),
		Databases:  defaultDatabases(),

		Blocker:              defaultBlocker(),
		ControllerReputation: defaultControllerReputation(),
	}
}

// DefaultHttpApiConfig - получить дефолтную конфигурацию движка Http Api сервера.
func DefaultHttpApiConfig() *HttpApiConfig {
	return &HttpApiConfig{
	}
}

// defaultHTTPProxy - получить дефолтную конфигурацию http_proxy_cfg.HTTPProxy.
func defaultHTTPProxy() *http_proxy_cfg.HTTPProxy {
	return &http_proxy_cfg.HTTPProxy{
		ErrorPage:          "system/templates/html/error.html",
		CheckClientPage:    "system/templates/html/check_client.html",
		TechnicalWorksPage: "system/templates/html/technical_works.html",
		CheckCaptchaPage:   "system/templates/html/captcha.html",

		TLSCertificateDir: "system/certs/",
		AllowHosts:        make([]string, 0),
		RemoveHeaders: []string{
			"Connection",
			"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
			"Keep-Alive",
			"Proxy-Authenticate",
			"Proxy-Authorization",
			"Te",      // canonicalized version of "TE"
			"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
			"Transfer-Encoding",
			"Upgrade",
		},
		ClientLimit: 100,

		ValueRequestPerMinuteForNoDataPage: 600,

		EnableHostWhiteList:     true,
		EnableHardHostWhiteList: false,
		EnableHostBlackList:     true,

		EnableLocationWhiteList: false,
		EnableLocationBlackList: false,

		Host:     "",
		Type:     "",
		Protocol: "",
		InAddr:   "",
		OutAddr:  "",

		MethodsUnverifiedByToken: []string{
			http.MethodOptions,
		},

		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  300 * time.Second,
	}
}

// defaultHTTPClient - получить дефолтную конфигурацию http_client_cfg.HTTPClient.
func defaultHTTPClient() *http_client_cfg.HTTPClient {
	return &http_client_cfg.HTTPClient{
					Timeout:       0,
					FlushInterval: 0,

					Proto:      "HTTP/1.1",
					ProtoMajor: 1,
					ProtoMinor: 1,

					Transport: &http_client_cfg.Transport{
						TLSHandshakeTimeout:    15 * time.Second,
						DisableKeepAlives:      true,
						DisableCompression:     false,
						MaxIdleConns:           100,
						MaxIdleConnsPerHost:    100,
						MaxConnsPerHost:        0,
						IdleConnTimeout:        90 * time.Second,
						ResponseHeaderTimeout:  0,
						ExpectContinueTimeout:  1 * time.Second,
						MaxResponseHeaderBytes: 0,
						WriteBufferSize:        0,
						ReadBufferSize:         0,
						ForceAttemptHTTP2:      false,
					},
					Dialer: &http_client_cfg.Dialer{
						Timeout:       30 * time.Second,
						Deadline:      time.Time{},
						FallbackDelay: 0,
						KeepAlive:     30 * time.Second,
					},
				}
}

// defaultHttpApiServer - получить дефолтную конфигурацию http_api_server_cfg.Server.
func defaultHttpApiServer() *http_api_server_cfg.Server {
	return &http_api_server_cfg.Server{
		Domain:                      "",
		Name:                        "",
		Version:                     "",
		GinMode:                     "",
		HttpAddr:                    "",
		HttpsAddr:                   "",

		PprofFileDirectory:          "system/pprof/api/",
		PostmanCollectionsDirectory: "system/postman/api/",
		EnableGenPostmanCollections: true,
		EnablePathRewrite:           false,
		EnablePprof:                 true,
		EnableCors:                  true,
		SystemAccess: &http_api_server_cfg.SystemAccess{
			Title: "SystemAccess",
			Token: &http_api_server_cfg.Token{
				HasSalt:   "SQsad1Qza",
				LifeTime:  36000.0,
				SignedKey: "signed_key",
			},
		},
		Cors: &http_api_server_cfg.Cors{
			AllowOrigins:           []string{"*"},
			AllowMethods:           []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE", "HEAD"},
			AllowHeaders:           []string{"Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Cookie", "Authorization", "RequestURL", "Set-Cookie"},
			AllowCredentials:       true,
			ExposeHeaders:          []string{"Content-Disposition", "Content-Length"},
			MaxAge:                 36000.0,
			AllowAllOrigins:        false,
			AllowFiles:             true,
			AllowBrowserExtensions: false,
			AllowWildcard:          false,
			AllowWebSockets:        false,
		},

		TLSCertificateDir: "system/certs/",
		CaptchaImgDir:     "system/api/http/static/captcha/img/",
		TemplateFiles:               make([]string, 0),
		StaticFiles:                 make([]http_api_server_cfg.StaticFiles, 0),
	}
}

// defaultLoggers - получить дефолтную конфигурацию loggers_cfg.Loggers.
func defaultLoggers() *loggers_cfg.Loggers {
	return &loggers_cfg.Loggers{
		Api: &loggers_cfg.Api{
			Title:       "Api-Log",
			LogLevel:    "DEBUG",
			TimeFormat:  "Monday, 02 Jan 2006 15:04:05",
			LogFilePath: "system/logs/api/",
			EnableCallerPath: &loggers_cfg.EnableCallerPath{
				DEBUG: true,
				INFO:  false,
				WARN:  false,
				ERROR: true,
				FATAL: true,
			},
			EnableOutputFile: &loggers_cfg.LogEnableOutputFile{
				DEBUG: true,
				INFO:  true,
				WARN:  true,
				ERROR: true,
				FATAL: true,
			},
			EnableOutputInTerminal: true,
			Colors: &loggers_cfg.LogColors{
				DEBUG: "Green",
				INFO:  "Cyan",
				WARN:  "Yellow",
				ERROR: "Red",
				FATAL: "Purple",

				ALogColor: "Blue",
				PLogColor: "Cyan",
				GLogColor: "Purple",

				GET:     "Green",
				POST:    "Cyan",
				PUT:     "Purple",
				DELETE:  "Red",
				PATCH:   "Yellow",
				OPTIONS: "Blue",
				HEAD:    "Gray",

				HTTPCode100: "Green",
				HTTPCode200: "Cyan",
				HTTPCode300: "Purple",
				HTTPCode400: "Yellow",
				HTTPCode500: "Red",
			},
		},
		Global: &loggers_cfg.Global{
			Title:       "Global-Log",
			LogLevel:    "DEBUG",
			TimeFormat:  "Monday, 02 Jan 2006 15:04:05",
			LogFilePath: "system/logs/global/",
			EnableCallerPath: &loggers_cfg.EnableCallerPath{
				DEBUG: true,
				INFO:  false,
				WARN:  false,
				ERROR: true,
				FATAL: true,
			},
			EnableOutputFile: &loggers_cfg.LogEnableOutputFile{
				DEBUG: true,
				INFO:  true,
				WARN:  true,
				ERROR: true,
				FATAL: true,
			},
			EnableOutputInTerminal: true,
			Colors: &loggers_cfg.LogColors{
				DEBUG: "Green",
				INFO:  "Cyan",
				WARN:  "Yellow",
				ERROR: "Red",
				FATAL: "Purple",

				ALogColor: "Blue",
				PLogColor: "Cyan",
				GLogColor: "Purple",

				GET:     "Green",
				POST:    "Cyan",
				PUT:     "Purple",
				DELETE:  "Red",
				PATCH:   "Yellow",
				OPTIONS: "Blue",
				HEAD:    "Gray",

				HTTPCode100: "Green",
				HTTPCode200: "Cyan",
				HTTPCode300: "Purple",
				HTTPCode400: "Yellow",
				HTTPCode500: "Red",
			},
		},
		HttpProxy: &loggers_cfg.HttpProxy{
			Title:       "HttpProxy-Log",
			LogLevel:    "DEBUG",
			TimeFormat:  "Monday, 02 Jan 2006 15:04:05",
			LogFilePath: "system/logs/proxy/http/",
			EnableCallerPath: &loggers_cfg.EnableCallerPath{
				DEBUG: false,
				INFO:  false,
				WARN:  false,
				ERROR: false,
				FATAL: false,
			},
			EnableOutputFile: &loggers_cfg.LogEnableOutputFile{
				DEBUG: true,
				INFO:  true,
				WARN:  true,
				ERROR: true,
				FATAL: true,
			},
			EnableOutputInTerminal: true,
			Colors: &loggers_cfg.LogColors{
				DEBUG: "Green",
				INFO:  "Cyan",
				WARN:  "Yellow",
				ERROR: "Red",
				FATAL: "Purple",

				ALogColor: "Blue",
				PLogColor: "Cyan",
				GLogColor: "Purple",

				GET:     "Green",
				POST:    "Cyan",
				PUT:     "Purple",
				DELETE:  "Red",
				PATCH:   "Yellow",
				OPTIONS: "Blue",
				HEAD:    "Gray",

				HTTPCode100: "Green",
				HTTPCode200: "Cyan",
				HTTPCode300: "Purple",
				HTTPCode400: "Yellow",
				HTTPCode500: "Red",
			},
		},
	}
}

// defaultSystem - получить дефолтную конфигурацию system_cfg.System.
func defaultSystem() *system_cfg.System {
	return &system_cfg.System{
		PasswordChars:  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
		PasswordLength: 12,
	}
}

// defaultDatabases - получить дефолтную конфигурацию databases_cfg.Databases.
func defaultDatabases() *databases_cfg.Databases {
	return &databases_cfg.Databases{
		Mongo: &databases_cfg.MongoDB{
			UserName: "nepeople",
			Password: "03Dor14proG",
			Addr:     "176.99.5.230:40000",

			Blocker: &databases_cfg.MongoDBBlocker{
				DatabaseName: "jkl-netdef_blocker",

				CollectionLocationBlackList: "location_black_list",
				CollectionLocationWhiteList: "location_white_list",

				CollectionHostBlackList:     "host_black_list",
				CollectionHostWhiteList:     "host_white_list",
				CollectionHostHardWhiteList: "host_hard_white_list",

				CollectionClientListOnCaptchaCheck: "client_list_on_captcha_check",
				CollectionHostBanList:              "host_ban_list",
				CollectionTokens:                   "tokens",

				CollectionUserAgentWhiteList: "user-agent_white_list",
				CollectionUserAgentBlackList: "user-agent_black_list",
			},
			ReputationController: &databases_cfg.MongoDBReputationController{
				DatabaseName: "jkl-netdef_controller_reputation",

				CollectionClientReputation: "client_reputation",
			},
			System: &databases_cfg.MongodbSystem{
				DatabaseName: "jkl-netdef_sys",

				CollectionSessions:      "sessions",
				CollectionRoles:         "roles",
				CollectionHttpRequests:  "http_requests",
				CollectionTokens:        "tokens",
				CollectionNotifications: "notifications",
				CollectionModules:       "modules",
			},
			Main: &databases_cfg.MongoDBMain{
				DatabaseName: "jkl-netdef_main",

				CollectionUser: "users",
			},
			Dashboard: &databases_cfg.MongoDBDashboard{
				DatabaseName: "jkl-netdef_dashboard",

				CollectionMainMenu: "main_menu",
			},
		},
	}
}

// defaultBlocker - получить дефолтную конфигурацию blocker_cfg.Blocker.
func defaultBlocker() *blocker_cfg.Blocker {
	return &blocker_cfg.Blocker{
		IPv4ByCountry: "system/data/ipv4.json",

		TokenName:      "proxy_token",
		TokenSecretKey: "12321321dfasdasf3<<1!/AS21",
		TokenExpire:    time.Hour * 4,
	}
}

// defaultControllerReputation - получить дефолтную конфигурацию databases_cfg.Databases.
func defaultControllerReputation() *controller_reputation_cfg.ControllerReputation {
	return &controller_reputation_cfg.ControllerReputation{
		MinValueScore:                   0,
		MaxValueRequestsPerMinute:       450,
		MaxValueRequestsPerMinuteByHost: 1000,
		ValueCounterResetScore:          10,
		MaxValueScore:                   100,
		ValueScoreForBan:                -100,
	}
}

//// Default - получить дефолтную конфигурацию прокси сервера.
//func Default() *GlobalConfig {
//	return &GlobalConfig{
//		HTTPProxy:
//		},




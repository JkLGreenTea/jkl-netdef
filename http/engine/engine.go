package http

import (
	"JkLNetDef/engine/config"
	"JkLNetDef/engine/databases"
	"JkLNetDef/engine/http/models/status"
	"JkLNetDef/engine/http/models/system/schema"
	"JkLNetDef/engine/http/models/system/system_access/http_request"
	"JkLNetDef/engine/http/utils/background_process"
	"JkLNetDef/engine/http/utils/postman"
	"JkLNetDef/engine/http/utils/signature"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/models/tls_certificate"
	"JkLNetDef/engine/modules/base_blocker"
	"JkLNetDef/engine/modules/base_logger"
	"JkLNetDef/engine/services"
	"JkLNetDef/engine/utils/synchronizer"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	ginpprof "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

// go tool pprof -png pprof.main.samples.cpu.001.pb.gz > profile1.png
// go tool pprof http://samgk.ru:40006/debug/pprof/heap
// go tool pprof http://samgk.ru:40006/debug/pprof/profile
// go tool pprof http://samgk.ru:40006/debug/pprof/block

// Engine - сервер.
type Engine struct {
	Title string

	httpGin  *gin.Engine
	httpsGin *gin.Engine

	Router *RouterGroup
	Schema *schema.Schema

	Config    *config.HttpApiConfig
	Databases *databases.Databases
	Services  *services.Services
	Utils     *Utils
	Modules   *Modules

	TLSCertificates     map[string]*tls_certificate.TLSCertificate
	backgroundProcesses []background_process.HandlerFunc
}

// Utils - утилиты.
type Utils struct {
	Synchronizer *synchronizer.Synchronizer
	Postman      *postman.Postman
}

// Modules - модули.
type Modules struct {
	Logger               interfacies.Logger
	HttpServerLogger     interfacies.HttpServerLogger
	ManagerMetaData      interfacies.ManagerMetaData
	ManagerNotifications interfacies.ManagerNotifications
	SystemAccess         interfacies.SystemAccess
	ManagerSessions      interfacies.ManagerSessions
	Authorizer           interfacies.Authorizer
	Validator            interfacies.Validator
	Blocker              interfacies.Blocker
}

// Build - создать сервер.
func Build(cfg *config.GlobalConfig, sync_ *synchronizer.Synchronizer, baseLogger interfacies.Logger,
	dbs *databases.Databases, blocker interfacies.Blocker) (*Engine, error) {
	engine := &Engine{
		httpGin:  gin.New(),
		httpsGin: gin.New(),

		Utils:   new(Utils),
		Modules: new(Modules),

		TLSCertificates:     make(map[string]*tls_certificate.TLSCertificate),
		backgroundProcesses: make([]background_process.HandlerFunc, 0),
	}

	ginpprof.Register(engine.httpGin)
	ginpprof.Register(engine.httpsGin)

	// Херня без которой не могут работать другие части системы
	{
		engine.Config = &config.HttpApiConfig{
			Engine:    cfg.HttpApiServer,
			ApiLogger: cfg.Loggers.Api,
		}

		engine.Utils.Synchronizer = sync_

		engine.Modules.Blocker = blocker
		engine.Modules.Logger = baseLogger
		engine.Databases = dbs
	}

	// Схема, роутер, метрика
	{
		engine.Schema = &schema.Schema{
			Requests: make([]*http_request.Request, 0),
			Groups:   make(map[string]*schema.Schema),
		}

		var routerPath string

		if engine.Config.Engine.Name != "" {
			routerPath = fmt.Sprintf("/%s", engine.Config.Engine.Name)
		}

		if engine.Config.Engine.Version != "" {
			routerPath += fmt.Sprintf("/%s", engine.Config.Engine.Version)
		}

		routerPath += "/"

		engine.Router = engine.Group(routerPath, engine.Config.Engine.Version, "", false)
	}

	// Билды
	{
		err := engine.buildUtils()
		if err != nil {
			engine.Modules.Logger.FATAL(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		err = engine.buildModules()
		if err != nil {
			engine.Modules.Logger.FATAL(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			return nil, err
		}
	}

	// Фоновые процессы
	{
		err := engine.loadDefaultBackgroundProcess()
		if err != nil {
			engine.Modules.Logger.FATAL(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			return nil, err
		}
	}

	// gin mode
	{
		ginMode := strings.ToLower(engine.Config.Engine.GinMode)

		switch ginMode {
		case "debug":
			gin.SetMode(gin.DebugMode)
		case "release":
			gin.SetMode(gin.ReleaseMode)
		case "test":
			gin.SetMode(gin.TestMode)
		}
	}

	return engine, nil
}

// New - создать новый сервер.
func New(cfg *config.GlobalConfig) (*Engine, error) {
	eng := &Engine{
		httpGin:  gin.New(),
		httpsGin: gin.New(),

		Utils:   new(Utils),
		Modules: new(Modules),

		TLSCertificates:     make(map[string]*tls_certificate.TLSCertificate),
		backgroundProcesses: make([]background_process.HandlerFunc, 0),
	}

	ginpprof.Register(eng.httpGin)
	ginpprof.Register(eng.httpsGin)

	// Херня без которой не могут работать другие части системы
	{
		eng.Config = &config.HttpApiConfig{
			Engine:    cfg.HttpApiServer,
			ApiLogger: cfg.Loggers.Api,
		}
		eng.Utils.Synchronizer = synchronizer.New()

		baseLogger, err := base_logger.New(cfg.Loggers.Global, eng.Utils.Synchronizer)
		if err != nil {
			log.Fatal(err)

			return nil, err
		}
		eng.Modules.Logger = baseLogger
	}

	// Бд
	{
		dbs, err := databases.New("Databases", cfg.Databases, eng.Modules.Logger)
		if err != nil {
			eng.Modules.Logger.FATAL(base_logger.Message{
				Sender: eng.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		eng.Databases = dbs
	}

	// Blocker
	{
		blocker, err := base_blocker.New("Blocker", cfg.Blocker, eng.Modules.Logger, eng.Databases)
		if err != nil {
			return nil, err
		}
		eng.Modules.Blocker = blocker
	}

	// Схема, роутер, метрика
	{
		eng.Schema = &schema.Schema{
			Requests: make([]*http_request.Request, 0),
			Groups:   make(map[string]*schema.Schema),
		}

		var routerPath string

		if eng.Config.Engine.Name != "" {
			routerPath = fmt.Sprintf("/%s", eng.Config.Engine.Name)
		}

		if eng.Config.Engine.Version != "" {
			routerPath += fmt.Sprintf("/%s", eng.Config.Engine.Version)
		}

		routerPath += "/"

		eng.Router = eng.Group(routerPath, eng.Config.Engine.Version, "", false)
	}

	// Билды
	{
		err := eng.buildUtils()
		if err != nil {
			eng.Modules.Logger.FATAL(base_logger.Message{
				Sender: eng.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		err = eng.buildModules()
		if err != nil {
			eng.Modules.Logger.FATAL(base_logger.Message{
				Sender: eng.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		err = eng.buildServices(cfg)
		if err != nil {
			eng.Modules.Logger.FATAL(base_logger.Message{
				Sender: eng.Title,
				Text:   err.Error(),
			})

			return nil, err
		}
	}

	// Фоновые процессы
	{
		err := eng.loadDefaultBackgroundProcess()
		if err != nil {
			eng.Modules.Logger.FATAL(base_logger.Message{
				Sender: eng.Title,
				Text:   err.Error(),
			})

			return nil, err
		}
	}

	// gin mode
	{
		ginMode := strings.ToLower(eng.Config.Engine.GinMode)

		switch ginMode {
		case "debug":
			gin.SetMode(gin.DebugMode)
		case "release":
			gin.SetMode(gin.ReleaseMode)
		case "test":
			gin.SetMode(gin.TestMode)
		}
	}

	return eng, nil
}

// Init - инициализация.
func (engine *Engine) Init() error {
	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Инициализация сервера... ",
	})

	// Templates
	{
		pwd, err := os.Getwd()
		if err != nil {
			engine.Modules.Logger.FATAL(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			return err
		}

		for _, templatePath := range engine.Config.Engine.TemplateFiles {
			engine.LoadHTMLGlob(path.Join(pwd, templatePath))
		}
	}

	// StaticFiles
	{
		for _, staticDir := range engine.Config.Engine.StaticFiles {
			err := engine.LoadStaticDirFiles(staticDir.Path, path.Join(engine.Router.path, staticDir.Root))
			if err != nil {
				engine.Modules.Logger.FATAL(base_logger.Message{
					Sender: engine.Title,
					Text:   err.Error(),
				})

				return err
			}
		}
	}

	// NoRoute
	{
		engine.NoRoute(func(ctx *gin.Context) {
			// Ответ
			{
				engine.WriteResponse(ctx, http.StatusNotFound, gin.H{
					"status": status.Error,
					"code":   http.StatusNotFound,
					"data":   "",
				})
			}
		})
	}

	// NoMethod
	{
		engine.NoMethod(func(ctx *gin.Context) {
			// Ответ
			{
				engine.WriteResponse(ctx, http.StatusNotFound, gin.H{
					"status": status.Error,
					"code":   http.StatusNotFound,
					"data":   "",
				})
			}
		})
	}

	// Middlewares
	{
		engine.Modules.Logger.INFO(base_logger.Message{
			Sender: engine.Title,
			Text:   "Запуск промежуточно программного обеспечения... ",
		})

		// Система
		{
			// WaitGroup middleware
			{
				engine.enableRequestWaitGroupMiddleware()
			}

			// Path rewrite middleware
			{
				if engine.Config.Engine.EnablePathRewrite {
					engine.enablePathRewriteMiddleware()
				}
			}

			// Cors
			{
				if engine.Config.Engine.EnableCors {
					engine.enableCorsMiddleware()
				}
			}
		}

		// Request logger middleware
		{
			if engine.Config.ApiLogger.EnableOutputInTerminal {
				engine.enableRequestLoggerMiddleware()
			}
		}

		// Утилиты
		{
			// Pprof middleware
			{
				if engine.Config.Engine.EnablePprof {
					engine.enablePprofMiddleware()
				}
			}
		}

		engine.Modules.Logger.INFO(base_logger.Message{
			Sender: engine.Title,
			Text:   "Всё промежуточное программное обеспечение запущено. ",
		})
	}

	// Requests
	{
		// favicon.ico
		{
			engine.Router.GET("/favicon.ico", "Иконка", "", "",
				false, false, false,
				func(ctx *gin.Context) {
					//ctx.File(path.Join(engine.Config.HttpApiServer.EnablePprof.Path, "favicon.ico"))
				})
		}

		// ping
		{
			engine.Router.GET("/ping", "", "Пинг сервера. ", "",
				true, false, false, engine.Modules.SystemAccess.HttpMiddleware(false),
				func(ctx *gin.Context) {
					type RequestArgs struct{}

					type ResponseArgs struct {
						Message string `json:"message" yaml:"message" form:"message" description:"Сообщение" validate:"required"`
					}

					requestArgs := new(RequestArgs)
					responseArgs := new(ResponseArgs)

					// Системная фигня
					{
						if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
							ctx.Abort()
							return
						}
					}

					responseArgs.Message = "pong"

					// Ответ
					{
						engine.WriteResponse(ctx, http.StatusOK, gin.H{
							"status": status.Success,
							"code":   http.StatusOK,
							"data":   responseArgs,
						})
					}
				})
		}

		// health
		{
			engine.Router.GET("/health", "", "Проверить сервер. ", "",
				true, false, false, engine.Modules.SystemAccess.HttpMiddleware(false),
				func(ctx *gin.Context) {
					type RequestArgs struct{}

					type ResponseArgs struct{}

					requestArgs := new(RequestArgs)
					responseArgs := new(ResponseArgs)

					// Системная фигня
					{
						if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
							ctx.Abort()
							return
						}
					}

					// Ответ
					{
						engine.WriteResponse(ctx, http.StatusOK, gin.H{
							"status": status.Success,
							"code":   http.StatusOK,
							"data":   responseArgs,
						})
					}
				})
		}

		// Schema
		{
			engine.Router.GET("/schema", "", "Просмотр схемы запросов. ", "",
				true, false, false, engine.Modules.SystemAccess.HttpMiddleware(false),
				func(ctx *gin.Context) {
					type RequestArgs struct {
					}

					type ResponseArgs struct {
						Schema *schema.Schema `json:"schema" yaml:"schema" form:"schema" description:"Схема запросов" validate:"required"`
					}

					requestArgs := new(RequestArgs)
					responseArgs := new(ResponseArgs)

					// Системная фигня
					{
						if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
							ctx.Abort()
							return
						}
					}

					// Запись ответа
					{
						responseArgs.Schema = engine.Schema
					}

					// Ответ
					{
						engine.WriteResponse(ctx, http.StatusOK, gin.H{
							"status": status.Success,
							"code":   http.StatusOK,
							"data":   responseArgs,
						})
					}
				})
		}

		// Documentation
		{
			engine.Router.GET("/docs", "", "Просмотр OpenAPI документации. ", "",
				true, false, false, engine.Modules.SystemAccess.HttpMiddleware(false),
				func(ctx *gin.Context) {
					fContent, err := ioutil.ReadFile("system/api/http/web/docs/index.html")
					if err != nil {
						engine.Modules.Logger.WARN(base_logger.Message{
							Sender: engine.Title,
							Text:   err.Error(),
						})

						engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
							"status": status.Failed,
							"code":   http.StatusBadRequest,
							"error": map[string]interface{}{
								"message": "Ошибка чтения HTML файла. ",
								"err":     err.Error(),
								"fields":  make(map[string]interface{}),
							},
						})
					}

					ctx.Data(http.StatusOK, "text/html; charset=utf-8", fContent)
				})
		}
	}

	return nil
}

// LoadTLSCertificate - загрузить TLS сертификат.
func (engine *Engine) LoadTLSCertificate(domain, cert, key string) error {
	engine.TLSCertificates[domain] = &tls_certificate.TLSCertificate{
		Domain: domain,
		Cert:   path.Join(engine.Config.Engine.TLSCertificateDir, cert),
		Key:    path.Join(engine.Config.Engine.TLSCertificateDir, key),
	}

	return nil
}

// Run - запуск Http сервера.
func (engine *Engine) Run() error {
	err := engine.LoadRequests()
	if err != nil {
		engine.Modules.Logger.FATAL(base_logger.Message{
			Sender: engine.Title,
			Text:   err.Error(),
		})

		return err
	}

	err = engine.GenerationPostmanCollections()
	if err != nil {
		engine.Modules.Logger.FATAL(base_logger.Message{
			Sender: engine.Title,
			Text:   err.Error(),
		})

		return err
	}

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Сервер запущен. ",
	})

	// Запуск фоновых процессов
	{
		process := &background_process.Process{
			Databases: engine.Databases,
			Utils: &background_process.Utils{
				Sync:    engine.Utils.Synchronizer,
				Postman: engine.Utils.Postman,
			},
			Modules: &background_process.Modules{
				Logger:               engine.Modules.Logger,
				HttpServerLogger:     engine.Modules.HttpServerLogger,
				ManagerMetadata:      engine.Modules.ManagerMetaData,
				ManagerNotifications: engine.Modules.ManagerNotifications,
				SystemAccess:         engine.Modules.SystemAccess,
				ManagerSessions:      engine.Modules.ManagerSessions,
				Authorizer:           engine.Modules.Authorizer,
				Validator:            engine.Modules.Validator,
			},
		}

		for _, handle := range engine.backgroundProcesses {
			go func(handle background_process.HandlerFunc, process *background_process.Process) {
				if err := handle(process); err != nil {
					engine.Modules.Logger.FATAL(base_logger.Message{
						Sender: engine.Title,
						Text:   err.Error(),
					})
				}
			}(handle, process)
		}
	}

	// Http
	{
		if engine.Config.Engine.HttpAddr != "" {
			engine.Utils.Synchronizer.HttpServer.WaitGroup.Add(1)

			go func() {
				defer func() {
					engine.Utils.Synchronizer.HttpServer.WaitGroup.Done()
				}()

				err = engine.httpGin.Run(engine.Config.Engine.HttpAddr)
				if err != nil {
					engine.Modules.Logger.FATAL(base_logger.Message{
						Sender: engine.Title,
						Text:   err.Error(),
					})
				}
			}()
		}
	}

	// Https
	{
		if engine.Config.Engine.HttpsAddr != "" {
			var srv *http.Server

			engine.Utils.Synchronizer.HttpServer.WaitGroup.Add(1)

			// Server
			{
				tlsConfig := &tls.Config{}

				// tls
				{
					var err error

					tlsConfig.Certificates = make([]tls.Certificate, len(engine.TLSCertificates), len(engine.TLSCertificates))

					i := 0
					for _, cert := range engine.TLSCertificates {
						tlsConfig.Certificates[i], err = tls.LoadX509KeyPair(cert.Cert, cert.Key)
						if err != nil {
							engine.Modules.Logger.ERROR(err.Error())
							return err
						}

						i++
					}
				}

				srv = &http.Server{
					Addr:         engine.Config.Engine.HttpsAddr,
					Handler:      engine.httpsGin,
					TLSConfig:    tlsConfig,
					TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
				}
			}

			go func() {
				defer func() {
					engine.Utils.Synchronizer.HttpServer.WaitGroup.Done()
				}()

				if err = srv.ListenAndServeTLS("", ""); err != nil {
					engine.Modules.Logger.ERROR(err.Error())
				}
			}()
		}
	}

	engine.Utils.Synchronizer.HttpServer.WaitGroup.Wait()

	return nil
}

// Stop - остановка Http Api.
func (engine *Engine) Stop() error {
	if err := engine.Databases.Disconnect(); err != nil {
		return err
	}

	return nil
}

// LoadBackgroundProcess - загрузка фоновых процессов.
func (engine *Engine) LoadBackgroundProcess(processes ...background_process.HandlerFunc) error {
	engine.backgroundProcesses = append(engine.backgroundProcesses, processes...)

	return nil
}

// SystemHandle - обработка системных процессов.
func (engine *Engine) SystemHandle(ctx *gin.Context, requestArgs, responseArgs interface{}) bool {
	// Возврат аргументов запроса ?signature=
	{
		if ctx.Request.URL.Query().Get("signature") == "true" {
			signat := signature.GetStructSignature(requestArgs, responseArgs)

			engine.WriteResponse(ctx, http.StatusOK, gin.H{
				"status":    status.Success,
				"code":      http.StatusOK,
				"signature": signat,
				"data":      make(map[string]interface{}),
			})
			return false
		}
	}

	// Чтение данных
	{
		err := engine.ReadRequest(ctx, requestArgs)
		if err != nil {
			engine.Modules.Logger.WARN(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
				"status": status.Failed,
				"code":   http.StatusBadRequest,
				"error": map[string]interface{}{
					"message": "Ошибка чтения данных. ",
					"err":     err.Error(),
					"fields":  make(map[string]interface{}),
				},
			})
			return false
		}

		buff, err := json.Marshal(requestArgs)
		if err != nil {
			engine.Modules.Logger.WARN(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})
		}
		ctx.Set("request", string(buff))
	}

	// Валидация данных
	{
		errs := engine.Modules.Validator.Struct(requestArgs)
		if len(errs) > 0 {
			err := errors.New("Request validation error. ")
			engine.Modules.Logger.WARN(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
				"status": status.Failed,
				"code":   http.StatusBadRequest,
				"error": map[string]interface{}{
					"message": "Ошибка валидации данных. ",
					"err":     err.Error(),
					"fields":  errs,
				},
			})
			return false
		}
	}

	return true
}

// LoadRequests - подгрузка запросов.
func (engine *Engine) LoadRequests() error {
	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Подгрузка запросов... ",
	})

	err := engine.loadRequests(engine.Schema)
	if err != nil {
		engine.Modules.Logger.WARN(base_logger.Message{
			Sender: engine.Title,
			Text:   err.Error(),
		})
	}
	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Запросы загружены. ",
	})

	return nil
}

// NoRoute - gin.NoRoute
func (engine *Engine) NoRoute(handlers ...gin.HandlerFunc) {
	engine.httpGin.NoRoute(handlers...)
	engine.httpsGin.NoRoute(handlers...)
}

// NoMethod - gin.NoMethod
func (engine *Engine) NoMethod(handlers ...gin.HandlerFunc) {
	engine.httpGin.NoMethod(handlers...)
	engine.httpsGin.NoMethod(handlers...)
}

// LoadHTMLGlob - gin.LoadHTMLGlob
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.httpGin.LoadHTMLGlob(pattern)
	engine.httpsGin.LoadHTMLGlob(pattern)
}

// GenerationPostmanCollections - генерация postman коллекций.
func (engine *Engine) GenerationPostmanCollections() error {
	if !engine.Config.Engine.EnableGenPostmanCollections {
		return nil
	}

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Генерация postman коллекций... ",
	})

	err := engine.Utils.Postman.Generation(engine.Schema)
	if err != nil {
		engine.Modules.Logger.WARN(base_logger.Message{
			Sender: engine.Title,
			Text:   err.Error(),
		})
	}
	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "генерация postman коллекций завершена. ",
	})

	return nil
}

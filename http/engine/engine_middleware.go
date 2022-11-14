package http

import (
	"JkLNetDef/engine/modules/base_logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/profile"
	"path"
	"time"
)

// System --------------------------

// enableRequestLoggerMiddleware - включение логгер middleware на запросы.
func (engine *Engine) enableRequestLoggerMiddleware() {
	engine.httpGin.Use(engine.Modules.HttpServerLogger.HttpMiddleware())
	engine.httpsGin.Use(engine.Modules.HttpServerLogger.HttpMiddleware())
	engine.Router.Use(engine.Modules.HttpServerLogger.HttpMiddleware())

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Request logger middleware запущен. ",
	})
}

// enablePathRewriteMiddleware - включить middleware с перезаписью путей (вставляет / в конце).
func (engine *Engine) enablePathRewriteMiddleware() {
	handle := func(ctx *gin.Context) {
		if ctx.Request.URL.Path[len(ctx.Request.URL.Path)-1] != '/' {
			ctx.Request.URL.Path += "/"
		}

		engine.httpGin.HandleContext(ctx)
		engine.httpsGin.HandleContext(ctx)
	}

	engine.httpGin.Use(handle)
	engine.httpsGin.Use(handle)
	engine.Router.Use(handle)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Path rewrite middleware запущен. ",
	})
}

// enableRequestWaitGroupMiddleware - включение wait group middleware.
func (engine *Engine) enableRequestWaitGroupMiddleware() {
	waitGroupMiddleware := func(ctx *gin.Context) {

		engine.Utils.Synchronizer.HttpServer.RequestsWg.Add(1)

		defer func() {
			engine.Utils.Synchronizer.HttpServer.RequestsWg.Done()
		}()

		ctx.Next()
	}

	engine.httpGin.Use(waitGroupMiddleware)
	engine.httpsGin.Use(waitGroupMiddleware)
	engine.Router.Use(waitGroupMiddleware)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Request wait group middleware запущен. ",
	})
}

// enableCorsMiddleware - включение cors middleware.
func (engine *Engine) enableCorsMiddleware() {
	middleware_ := cors.New(cors.Config{
		AllowOrigins:           engine.Config.Engine.Cors.AllowOrigins,
		AllowMethods:           engine.Config.Engine.Cors.AllowMethods,
		AllowHeaders:           engine.Config.Engine.Cors.AllowHeaders,
		AllowCredentials:       engine.Config.Engine.Cors.AllowCredentials,
		ExposeHeaders:          engine.Config.Engine.Cors.ExposeHeaders,
		MaxAge:                 time.Duration(engine.Config.Engine.Cors.MaxAge) * time.Second,
		AllowAllOrigins:        engine.Config.Engine.Cors.AllowAllOrigins,
		AllowFiles:             engine.Config.Engine.Cors.AllowFiles,
		AllowBrowserExtensions: engine.Config.Engine.Cors.AllowBrowserExtensions,
		AllowWildcard:          engine.Config.Engine.Cors.AllowWildcard,
		AllowWebSockets:        engine.Config.Engine.Cors.AllowWebSockets,
	})

	engine.httpGin.Use(middleware_)
	engine.httpsGin.Use(middleware_)
	engine.Router.Use(middleware_)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Cross-Origin Resource Sharing middleware запущен. ",
	})
}

// Other --------------------------

// enablePprofMiddleware - включение middleware с профилированием.
// go tool pprof -http=0.0.0.0:40000 ./system/pprof/
// http://samgk.ru:40000/
func (engine *Engine) enablePprofMiddleware() {
	middleware := func(ctx *gin.Context) {
		var typePprof func(*profile.Profile)
		typePprofStr := ctx.Query("pprof")

		if typePprofStr != "" {
			switch typePprofStr {
			case "mem":
				typePprof = profile.MemProfile
			case "mem_heap":
				typePprof = profile.MemProfileHeap
			case "mem_allocs":
				typePprof = profile.MemProfileAllocs
			case "block":
				typePprof = profile.BlockProfile
			case "cpu":
				typePprof = profile.CPUProfile
			case "goroutine":
				typePprof = profile.GoroutineProfile
			case "mutex":
				typePprof = profile.MutexProfile
			case "no_shutdown_hook":
				typePprof = profile.NoShutdownHook
			case "quiet":
				typePprof = profile.Quiet
			case "threadcreation":
				typePprof = profile.ThreadcreationProfile
			case "trace":
				typePprof = profile.TraceProfile
			}

			prof := profile.Start(typePprof, profile.ProfilePath(path.Join(engine.Config.Engine.PprofFileDirectory, engine.Title)))
			ctx.Next()
			prof.Stop()

			return
		}

		ctx.Next()
	}

	engine.httpGin.Use(middleware)
	engine.httpsGin.Use(middleware)
	engine.Router.Use(middleware)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Pprof middleware запущен. ",
	})
}

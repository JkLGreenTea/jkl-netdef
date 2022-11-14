package http

import (
	"JkLNetDef/engine/http/models/system/schema"
	"JkLNetDef/engine/http/models/system/system_access/http_request"
	"JkLNetDef/engine/modules/base_logger"
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"path"
)

// Handle - обертка над gin.Handle.
func (engine *Engine) Handle(httpMethod, relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      httpMethod,
		URL:         path.Join("/", relativePath),
		Version:     engine.Config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	engine.Schema.Requests = append(engine.Schema.Requests, req)

	iRoutes1 := engine.httpGin.Handle(httpMethod, relativePath, handlers...)
	iRoutes2 := engine.httpsGin.Handle(httpMethod, relativePath, handlers...)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   fmt.Sprintf("Запрос \"%s %s\" зарегистрирован. (%d middleware)", httpMethod, relativePath, len(engine.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// GET - обертка над gin.GET.
func (engine *Engine) GET(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "GET",
		URL:         path.Join("/", relativePath),
		Version:     engine.Config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	engine.Schema.Requests = append(engine.Schema.Requests, req)

	iRoutes1 := engine.httpGin.Handle("GET", relativePath, handlers...)
	iRoutes2 := engine.httpsGin.Handle("GET", relativePath, handlers...)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   fmt.Sprintf("Запрос \"GET %s\" зарегистрирован. (%d middleware)", relativePath, len(engine.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// POST - обертка над gin.POST.
func (engine *Engine) POST(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "POST",
		URL:         path.Join("/", relativePath),
		Version:     engine.Config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	engine.Schema.Requests = append(engine.Schema.Requests, req)

	iRoutes1 := engine.httpGin.Handle("POST", relativePath, handlers...)
	iRoutes2 := engine.httpsGin.Handle("POST", relativePath, handlers...)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   fmt.Sprintf("Запрос \"POST %s\" зарегистрирован. (%d middleware)", relativePath, len(engine.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// PUT - обертка над gin.PUT.
func (engine *Engine) PUT(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "PUT",
		URL:         path.Join("/", relativePath),
		Version:     engine.Config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	engine.Schema.Requests = append(engine.Schema.Requests, req)

	iRoutes1 := engine.httpGin.Handle("PUT", relativePath, handlers...)
	iRoutes2 := engine.httpsGin.Handle("PUT", relativePath, handlers...)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   fmt.Sprintf("Запрос \"PUT %s\" зарегистрирован. (%d middleware)", relativePath, len(engine.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// PATCH - обертка над gin.PATCH.
func (engine *Engine) PATCH(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "PATCH",
		URL:         path.Join("/", relativePath),
		Version:     engine.Config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	engine.Schema.Requests = append(engine.Schema.Requests, req)

	iRoutes1 := engine.httpGin.Handle("PATCH", relativePath, handlers...)
	iRoutes2 := engine.httpsGin.Handle("PATCH", relativePath, handlers...)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   fmt.Sprintf("Запрос \"PATCH %s\" зарегистрирован. (%d middleware)", relativePath, len(engine.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// DELETE - обертка над gin.DELETE.
func (engine *Engine) DELETE(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "DELETE",
		URL:         path.Join("/", relativePath),
		Version:     engine.Config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	engine.Schema.Requests = append(engine.Schema.Requests, req)

	iRoutes1 := engine.httpGin.Handle("DELETE", relativePath, handlers...)
	iRoutes2 := engine.httpsGin.Handle("DELETE", relativePath, handlers...)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   fmt.Sprintf("Запрос \"DELETE %s\" зарегистрирован. (%d middleware)", relativePath, len(engine.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// OPTIONS - обертка над gin.OPTIONS.
func (engine *Engine) OPTIONS(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "OPTIONS",
		URL:         path.Join("/", relativePath),
		Version:     engine.Config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	engine.Schema.Requests = append(engine.Schema.Requests, req)

	iRoutes1 := engine.httpGin.Handle("OPTIONS", relativePath, handlers...)
	iRoutes2 := engine.httpsGin.Handle("OPTIONS", relativePath, handlers...)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   fmt.Sprintf("Запрос \"OPTIONS %s\" зарегистрирован. (%d middleware)", relativePath, len(engine.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// HEAD - обертка над gin.HEAD.
func (engine *Engine) HEAD(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "HEAD",
		URL:         path.Join("/", relativePath),
		Version:     engine.Config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	engine.Schema.Requests = append(engine.Schema.Requests, req)

	iRoutes1 := engine.httpGin.Handle("HEAD", relativePath, handlers...)
	iRoutes2 := engine.httpsGin.Handle("HEAD", relativePath, handlers...)

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   fmt.Sprintf("Запрос \"HEAD %s\" зарегистрирован. (%d middleware)", relativePath, len(engine.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// Group - обертка над gin.Group.
func (engine *Engine) Group(relativePath string, title, description string, isSystem bool, handlers ...gin.HandlerFunc) *RouterGroup {
	schem := &schema.Schema{
		Requests:    make([]*http_request.Request, 0),
		Groups:      make(map[string]*schema.Schema),
		Title:       title,
		Description: description,
		IsSystem:    isSystem,
	}

	if engine.Schema.Groups[relativePath] == nil {
		engine.Schema.Groups[relativePath] = schem
	} else {
		schem = engine.Schema.Groups[relativePath]
	}

	return &RouterGroup{
		Schema:   schem,
		httpGin:  engine.httpGin.Group(relativePath, handlers...),
		httpsGin: engine.httpsGin.Group(relativePath, handlers...),

		config:  engine.Config,
		utils:   engine.Utils,
		modules: engine.Modules,
		path:    relativePath,
	}
}

// LoadStaticDirFiles - загрузить директорию со статическими файлами.
func (engine *Engine) LoadStaticDirFiles(dir, root string) error {
	req := &http_request.Request{
		Method:      "GET",
		URL:         path.Join("/", root),
		Version:     engine.Config.Engine.Version,
		Locked:      false,
		Authorized:  false,
		IsSystem:    true,
		IsStatic:    true,
		Info:        "Статичные файлы.",
		Description: "",
		Title:       "Статичные файлы. ",
	}
	engine.Schema.Requests = append(engine.Schema.Requests, req)

	engine.httpGin.Use(static.Serve(root, static.LocalFile(dir, true)))
	engine.httpsGin.Use(static.Serve(root, static.LocalFile(dir, true)))

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   fmt.Sprintf("Статические файлы подгружены path: \"%s\" root: \"%s\". ", dir, root),
	})

	return nil
}

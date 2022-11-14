package http

import (
	"JkLNetDef/engine/config"
	"JkLNetDef/engine/http/models/system/schema"
	"JkLNetDef/engine/http/models/system/system_access/http_request"
	"JkLNetDef/engine/modules/base_logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
)

// RouterGroup - обертка над gin.RouterGroup
type RouterGroup struct {
	Title  string
	Schema *schema.Schema

	httpGin  *gin.RouterGroup
	httpsGin *gin.RouterGroup

	config  *config.HttpApiConfig
	utils   *Utils
	modules *Modules
	path    string
}

// Group - обертка над gin.Group.
func (router *RouterGroup) Group(relativePath, title, description string, isSystem bool, handlers ...gin.HandlerFunc) *RouterGroup {
	schem := &schema.Schema{
		Requests:    make([]*http_request.Request, 0),
		Groups:      make(map[string]*schema.Schema),
		Title:       title,
		Description: description,
		IsSystem:    isSystem,
	}

	if router.Schema.Groups[relativePath] == nil {
		router.Schema.Groups[relativePath] = schem
	} else {
		schem = router.Schema.Groups[relativePath]
	}

	return &RouterGroup{
		Schema:   schem,
		httpGin:  router.httpGin.Group(relativePath, handlers...),
		httpsGin: router.httpsGin.Group(relativePath, handlers...),

		config:  router.config,
		utils:   router.utils,
		modules: router.modules,
		path:    path.Join(router.path, relativePath),
	}
}

// Use - обертка над gin.Use.
func (router *RouterGroup) Use(handlers ...gin.HandlerFunc) {
	router.httpGin.Use(handlers...)
	router.httpsGin.Use(handlers...)
}

// Handle - обертка над gin.Handle.
func (router *RouterGroup) Handle(httpMethod, relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      httpMethod,
		URL:         path.Join(router.path, relativePath),
		Version:     router.config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	router.Schema.Requests = append(router.Schema.Requests, req)

	iRoutes1 := router.httpGin.Handle(httpMethod, relativePath, handlers...)
	iRoutes2 := router.httpsGin.Handle(httpMethod, relativePath, handlers...)

	router.modules.Logger.INFO(base_logger.Message{
		Sender: router.Title,
		Text:   fmt.Sprintf("Запрос \"%s %s\" зарегистрирован. (%d middleware)", req.Method, req.URL, len(router.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// GET - обертка над gin.GET.
func (router *RouterGroup) GET(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "GET",
		URL:         path.Join(router.path, relativePath),
		Version:     router.config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	router.Schema.Requests = append(router.Schema.Requests, req)

	iRoutes1 := router.httpGin.Handle("GET", relativePath, handlers...)
	iRoutes2 := router.httpsGin.Handle("GET", relativePath, handlers...)

	router.modules.Logger.INFO(base_logger.Message{
		Sender: router.Title,
		Text:   fmt.Sprintf("Запрос \"%s %s\" зарегистрирован. (%d middleware)", req.Method, req.URL, len(router.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// POST - обертка над gin.POST.
func (router *RouterGroup) POST(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "POST",
		URL:         path.Join(router.path, relativePath),
		Version:     router.config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	router.Schema.Requests = append(router.Schema.Requests, req)

	iRoutes1 := router.httpGin.Handle("POST", relativePath, handlers...)
	iRoutes2 := router.httpsGin.Handle("POST", relativePath, handlers...)

	router.modules.Logger.INFO(base_logger.Message{
		Sender: router.Title,
		Text:   fmt.Sprintf("Запрос \"%s %s\" зарегистрирован. (%d middleware)", req.Method, req.URL, len(router.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// PUT - обертка над gin.PUT.
func (router *RouterGroup) PUT(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "PUT",
		URL:         path.Join(router.path, relativePath),
		Version:     router.config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	router.Schema.Requests = append(router.Schema.Requests, req)

	iRoutes1 := router.httpGin.Handle("PUT", relativePath, handlers...)
	iRoutes2 := router.httpsGin.Handle("PUT", relativePath, handlers...)

	router.modules.Logger.INFO(base_logger.Message{
		Sender: router.Title,
		Text:   fmt.Sprintf("Запрос \"%s %s\" зарегистрирован. (%d middleware)", req.Method, req.URL, len(router.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// PATCH - обертка над gin.PATCH.
func (router *RouterGroup) PATCH(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "PATCH",
		URL:         path.Join(router.path, relativePath),
		Version:     router.config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	router.Schema.Requests = append(router.Schema.Requests, req)

	iRoutes1 := router.httpGin.Handle("PATCH", relativePath, handlers...)
	iRoutes2 := router.httpsGin.Handle("PATCH", relativePath, handlers...)

	router.modules.Logger.INFO(base_logger.Message{
		Sender: router.Title,
		Text:   fmt.Sprintf("Запрос \"%s %s\" зарегистрирован. (%d middleware)", req.Method, req.URL, len(router.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// DELETE - обертка над gin.DELETE.
func (router *RouterGroup) DELETE(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "DELETE",
		URL:         path.Join(router.path, relativePath),
		Version:     router.config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	router.Schema.Requests = append(router.Schema.Requests, req)

	iRoutes1 := router.httpGin.Handle("DELETE", relativePath, handlers...)
	iRoutes2 := router.httpsGin.Handle("DELETE", relativePath, handlers...)

	router.modules.Logger.INFO(base_logger.Message{
		Sender: router.Title,
		Text:   fmt.Sprintf("Запрос \"%s %s\" зарегистрирован. (%d middleware)", req.Method, req.URL, len(router.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// OPTIONS - обертка над gin.OPTIONS.
func (router *RouterGroup) OPTIONS(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "OPTIONS",
		URL:         path.Join(router.path, relativePath),
		Version:     router.config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	router.Schema.Requests = append(router.Schema.Requests, req)

	iRoutes1 := router.httpGin.Handle("OPTIONS", relativePath, handlers...)
	iRoutes2 := router.httpsGin.Handle("OPTIONS", relativePath, handlers...)

	router.modules.Logger.INFO(base_logger.Message{
		Sender: router.Title,
		Text:   fmt.Sprintf("Запрос \"%s %s\" зарегистрирован. (%d middleware)", req.Method, req.URL, len(router.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

// HEAD - обертка над gin.HEAD.
func (router *RouterGroup) HEAD(relativePath, info, title, description string, isSystem, locked, authorized bool, handlers ...gin.HandlerFunc) (gin.IRoutes, gin.IRoutes) {
	req := &http_request.Request{
		Method:      "HEAD",
		URL:         path.Join(router.path, relativePath),
		Version:     router.config.Engine.Version,
		Locked:      locked,
		Authorized:  authorized,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
	}
	router.Schema.Requests = append(router.Schema.Requests, req)

	iRoutes1 := router.httpGin.Handle("HEAD", relativePath, handlers...)
	iRoutes2 := router.httpsGin.Handle("HEAD", relativePath, handlers...)

	router.modules.Logger.INFO(base_logger.Message{
		Sender: router.Title,
		Text:   fmt.Sprintf("Запрос \"%s %s\" зарегистрирован. (%d middleware)", req.Method, req.URL, len(router.httpGin.Handlers)),
	})

	return iRoutes1, iRoutes2
}

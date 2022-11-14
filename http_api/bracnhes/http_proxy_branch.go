package bracnhes

import (
	http_client_cfg "JkLNetDef/engine/config/http_client"
	http_proxy_cfg "JkLNetDef/engine/config/http_proxy"
	my_http "JkLNetDef/engine/http/engine"
	"JkLNetDef/engine/http/models/status"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"JkLNetDef/engine/proxy"
	http_reverse_proxy "JkLNetDef/engine/proxy/http_reverse"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

// InitHttpProxyBranch - инициализация ветви Http прокси.
func InitHttpProxyBranch(engine *my_http.Engine, proxyStorage *proxy.Storage) error {
	router := engine.Router.Group("/proxy", "Ветка прокси. ", "", true)
	httpProxyRouter := router.Group("/http_reverse", "Ветка Http Reverse прокси. ", "", true)

	// Запросы
	{
		// Получить список Http прокси.
		httpProxyRouter.GET("/list", "", "Получить список Http прокси. ", "",
			true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
			func(ctx *gin.Context) {
				type RequestArgs struct{}

				type ResponseArgs struct {
					Proxies []http_reverse_proxy.Proxy `json:"proxies" yaml:"proxies" form:"proxies" description:"Http reverse прокси. "`
				}

				requestArgs := new(RequestArgs)
				responseArgs := &ResponseArgs{
					Proxies: make([]http_reverse_proxy.Proxy, 0, proxyStorage.Len()),
				}

				// Системная фигня
				{
					if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
						ctx.Abort()
						return
					}
				}

				// Обработка
				{
					prxs, err := proxyStorage.GetAll()
					if err != nil {
						engine.Modules.Logger.WARN(base_logger.Message{
							Sender: engine.Title,
							Text:   err.Error(),
						})
						engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
							"status": status.Failed,
							"code":   http.StatusBadRequest,
							"error": map[string]interface{}{
								"message": "Ошибка получения Http reverse прокси. ",
								"err":     err.Error(),
								"fields":  make(map[string]interface{}),
							},
						})
						ctx.Abort()
						return
					}

					for _, prx_ := range prxs {
						switch prx := prx_.(type) {
						case *http_reverse_proxy.Proxy:
							{
								prx__ := *prx
								prx__.Config = nil
								responseArgs.Proxies = append(responseArgs.Proxies, prx__)
							}
						}
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

		// Получить все Http прокси.
		httpProxyRouter.GET("/all", "", "Получить все Http прокси. ", "",
			true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
			func(ctx *gin.Context) {
				type RequestArgs struct{}

				type ResponseArgs struct {
					Proxies []http_reverse_proxy.Proxy `json:"proxies" yaml:"proxies" form:"proxies" description:"Http reverse прокси. "`
				}

				requestArgs := new(RequestArgs)
				responseArgs := &ResponseArgs{
					Proxies: make([]http_reverse_proxy.Proxy, 0, proxyStorage.Len()),
				}

				// Системная фигня
				{
					if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
						ctx.Abort()
						return
					}
				}

				// Обработка
				{
					prxs, err := proxyStorage.GetAll()
					if err != nil {
						engine.Modules.Logger.WARN(base_logger.Message{
							Sender: engine.Title,
							Text:   err.Error(),
						})
						engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
							"status": status.Failed,
							"code":   http.StatusBadRequest,
							"error": map[string]interface{}{
								"message": "Ошибка получения Http reverse прокси. ",
								"err":     err.Error(),
								"fields":  make(map[string]interface{}),
							},
						})
						ctx.Abort()
						return
					}

					for _, prx_ := range prxs {
						switch prx := prx_.(type) {
						case *http_reverse_proxy.Proxy:
							{
								prx__ := *prx
								prx__.Config = nil
								responseArgs.Proxies = append(responseArgs.Proxies, prx__)
							}
						}
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

		// Получить Http прокси по ID.
		httpProxyRouter.GET("/by_id/:id", "", "Получить Http прокси по ID. ", "",
			true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
			func(ctx *gin.Context) {
				type RequestArgs struct {
					ID string `uri:"id" description:"ID Http reverse прокси. "`
				}

				type ResponseArgs struct {
					Proxy interfacies.HttpProxy `json:"proxy" yaml:"proxy" form:"proxy" description:"Http reverse прокси. "`
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

				// Обработка
				{
					prxs, err := proxyStorage.GetAll()
					if err != nil {
						engine.Modules.Logger.WARN(base_logger.Message{
							Sender: engine.Title,
							Text:   err.Error(),
						})
						engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
							"status": status.Failed,
							"code":   http.StatusBadRequest,
							"error": map[string]interface{}{
								"message": "Ошибка получения Http reverse прокси. ",
								"err":     err.Error(),
								"fields":  make(map[string]interface{}),
							},
						})
						ctx.Abort()
						return
					}

				found:
					for _, prx_ := range prxs {
						switch prx := prx_.(type) {
						case *http_reverse_proxy.Proxy:
							{
								if prx.ID.Hex() == requestArgs.ID {
									responseArgs.Proxy = prx
									break found
								}
							}
						}
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

		// Изменить конфиг Http прокси.
		httpProxyRouter.PUT("/config", "", "Изменить конфиг Http прокси. ", "",
			true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
			func(ctx *gin.Context) {
				type RequestArgs struct {
					ID         string                      `json:"id" yaml:"id" form:"id" description:"ID Http reverse прокси. " validate:"required"`
					Title      string                      `json:"title" yaml:"title" form:"title" description:"Название Http reverse прокси. "`
					Engine     *http_proxy_cfg.HTTPProxy   `json:"engine" yaml:"engine" form:"engine" description:"Http proxy движок"`
					HTTPClient *http_client_cfg.HTTPClient `json:"http_client" yaml:"http_client" form:"http_client" description:"Http client (http.Client)"`
				}

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

				// Обработка
				{
					prxs, err := proxyStorage.GetAll()
					if err != nil {
						engine.Modules.Logger.WARN(base_logger.Message{
							Sender: engine.Title,
							Text:   err.Error(),
						})
						engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
							"status": status.Failed,
							"code":   http.StatusBadRequest,
							"error": map[string]interface{}{
								"message": "Ошибка получения Http reverse прокси. ",
								"err":     err.Error(),
								"fields":  make(map[string]interface{}),
							},
						})
						ctx.Abort()
						return
					}

				found:
					for _, prx_ := range prxs {
						switch prx := prx_.(type) {
						case *http_reverse_proxy.Proxy:
							{
								if prx.ID.Hex() == requestArgs.ID {

									prx.Title = requestArgs.Title

									prx.Config.Engine = requestArgs.Engine
									prx.Config.HTTPClient = requestArgs.HTTPClient

									break found
								}
							}
						}
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

	// Technical works
	{
		technicalWorksHttpProxyRouter := httpProxyRouter.Group("/technical_works", "Ветка тех. работ Http Reverse прокси. ", "", true)

		// Globals
		{
			globalTechnicalWorksHttpProxyRouter := technicalWorksHttpProxyRouter.Group("/global", "Ветка глобальных тех. работ Http Reverse прокси. ", "", true)

			// Запросы
			{
				// Получить доступ к прокси во время глобальных тех. работ.
				globalTechnicalWorksHttpProxyRouter.POST("/access", "", "Получить доступ к прокси во время глобальных тех. работ. ", "",
					true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
					func(ctx *gin.Context) {
						type RequestArgs struct {
							ID string `json:"id" yaml:"id" form:"id" description:"ID Http reverse прокси. " validate:"required"`
						}

						type ResponseArgs struct {}

						requestArgs := new(RequestArgs)
						responseArgs := new(ResponseArgs)

						// Системная фигня
						{
							if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
								ctx.Abort()
								return
							}
						}

						// Обработка
						{
							prxs, err := proxyStorage.GetAll()
							if err != nil {
								engine.Modules.Logger.WARN(base_logger.Message{
									Sender: engine.Title,
									Text:   err.Error(),
								})
								engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
									"status": status.Failed,
									"code":   http.StatusBadRequest,
									"error": map[string]interface{}{
										"message": "Ошибка получения Http reverse прокси. ",
										"err":     err.Error(),
										"fields":  make(map[string]interface{}),
									},
								})
								ctx.Abort()
								return
							}

						found:
							for _, prx_ := range prxs {
								switch prx := prx_.(type) {
								case *http_reverse_proxy.Proxy:
									{
										if prx.ID.Hex() == requestArgs.ID {
											host, _, _ := net.SplitHostPort(ctx.Request.RemoteAddr)

											prx.TechnicalWorks.IssueGlobalAccess(host, ctx.Request.UserAgent())

											break found
										}
									}
								}
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

				// Остановить глобальные тех. работы.
				globalTechnicalWorksHttpProxyRouter.POST("/stop", "", "Остановить глобальные тех. работы. ", "",
					true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
					func(ctx *gin.Context) {
						type RequestArgs struct {
							ID string `json:"id" yaml:"id" form:"id" description:"ID Http reverse прокси. " validate:"required"`
						}

						type ResponseArgs struct {}

						requestArgs := new(RequestArgs)
						responseArgs := new(ResponseArgs)

						// Системная фигня
						{
							if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
								ctx.Abort()
								return
							}
						}

						// Обработка
						{
							prxs, err := proxyStorage.GetAll()
							if err != nil {
								engine.Modules.Logger.WARN(base_logger.Message{
									Sender: engine.Title,
									Text:   err.Error(),
								})
								engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
									"status": status.Failed,
									"code":   http.StatusBadRequest,
									"error": map[string]interface{}{
										"message": "Ошибка получения Http reverse прокси. ",
										"err":     err.Error(),
										"fields":  make(map[string]interface{}),
									},
								})
								ctx.Abort()
								return
							}

						found:
							for _, prx_ := range prxs {
								switch prx := prx_.(type) {
								case *http_reverse_proxy.Proxy:
									{
										if prx.ID.Hex() == requestArgs.ID {
											prx.TechnicalWorks.StopGlobal()

											break found
										}
									}
								}
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

				// Запустить глобальные тех. работы.
				globalTechnicalWorksHttpProxyRouter.POST("/start", "", "Запустить глобальные тех. работы. ", "",
					true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
					func(ctx *gin.Context) {
						type RequestArgs struct {
							ID        string `json:"id" yaml:"id" form:"id" description:"ID Http reverse прокси. " validate:"required"`
							StartTime int64  `json:"start_time" yaml:"start_time" form:"start_time" description:"Время старта тех. работ. " validate:"required"`
							EndTime   int64  `json:"end_time" yaml:"end_time" form:"end_time" description:"Время конца тех. работ. " validate:"required"`
						}

						type ResponseArgs struct {}

						requestArgs := new(RequestArgs)
						responseArgs := new(ResponseArgs)

						// Системная фигня
						{
							if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
								ctx.Abort()
								return
							}
						}

						// Обработка
						{
							prxs, err := proxyStorage.GetAll()
							if err != nil {
								engine.Modules.Logger.WARN(base_logger.Message{
									Sender: engine.Title,
									Text:   err.Error(),
								})
								engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
									"status": status.Failed,
									"code":   http.StatusBadRequest,
									"error": map[string]interface{}{
										"message": "Ошибка получения Http reverse прокси. ",
										"err":     err.Error(),
										"fields":  make(map[string]interface{}),
									},
								})
								ctx.Abort()
								return
							}

						found:
							for _, prx_ := range prxs {
								switch prx := prx_.(type) {
								case *http_reverse_proxy.Proxy:
									{
										if prx.ID.Hex() == requestArgs.ID {
											prx.TechnicalWorks.StartGlobal(requestArgs.StartTime, requestArgs.EndTime)

											break found
										}
									}
								}
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
		}

		// By pages
		{
			byPagesTechnicalWorksHttpProxyRouter := technicalWorksHttpProxyRouter.Group("/by_pages", "Ветка тех. работ Http Reverse прокси по страницам. ", "", true)

			// Запросы
			{
				// Получить доступ к прокси во время тех. работ на странице.
				byPagesTechnicalWorksHttpProxyRouter.POST("/access", "", "Получить доступ к прокси во время тех. работ на странице. ", "",
					true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
					func(ctx *gin.Context) {
						type RequestArgs struct {
							ID  string `json:"id" yaml:"id" form:"id" description:"ID Http reverse прокси " validate:"required"`
							URL string `json:"url" yaml:"url" form:"url" description:"Страница " validate:"required"`
						}

						type ResponseArgs struct {}

						requestArgs := new(RequestArgs)
						responseArgs := new(ResponseArgs)

						// Системная фигня
						{
							if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
								ctx.Abort()
								return
							}
						}

						// Обработка
						{
							prxs, err := proxyStorage.GetAll()
							if err != nil {
								engine.Modules.Logger.WARN(base_logger.Message{
									Sender: engine.Title,
									Text:   err.Error(),
								})
								engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
									"status": status.Failed,
									"code":   http.StatusBadRequest,
									"error": map[string]interface{}{
										"message": "Ошибка получения Http reverse прокси. ",
										"err":     err.Error(),
										"fields":  make(map[string]interface{}),
									},
								})
								ctx.Abort()
								return
							}

						found:
							for _, prx_ := range prxs {
								switch prx := prx_.(type) {
								case *http_reverse_proxy.Proxy:
									{
										if prx.ID.Hex() == requestArgs.ID {
											host, _, _ := net.SplitHostPort(ctx.Request.RemoteAddr)

											prx.TechnicalWorks.IssueByPagesAccess(host, ctx.Request.UserAgent(), requestArgs.URL)

											break found
										}
									}
								}
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

				// Остановить тех. работы на странице.
				byPagesTechnicalWorksHttpProxyRouter.POST("/stop", "", "Остановить тех. работы на странице. ", "",
					true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
					func(ctx *gin.Context) {
						type RequestArgs struct {
							ID string `json:"id" yaml:"id" form:"id" description:"ID Http reverse прокси. " validate:"required"`
							URL string `json:"url" yaml:"url" form:"url" description:"Страница " validate:"required"`
						}

						type ResponseArgs struct {}

						requestArgs := new(RequestArgs)
						responseArgs := new(ResponseArgs)

						// Системная фигня
						{
							if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
								ctx.Abort()
								return
							}
						}

						// Обработка
						{
							prxs, err := proxyStorage.GetAll()
							if err != nil {
								engine.Modules.Logger.WARN(base_logger.Message{
									Sender: engine.Title,
									Text:   err.Error(),
								})
								engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
									"status": status.Failed,
									"code":   http.StatusBadRequest,
									"error": map[string]interface{}{
										"message": "Ошибка получения Http reverse прокси. ",
										"err":     err.Error(),
										"fields":  make(map[string]interface{}),
									},
								})
								ctx.Abort()
								return
							}

						found:
							for _, prx_ := range prxs {
								switch prx := prx_.(type) {
								case *http_reverse_proxy.Proxy:
									{
										if prx.ID.Hex() == requestArgs.ID {
											prx.TechnicalWorks.StopByPage(requestArgs.URL)

											break found
										}
									}
								}
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

				// Запустить тех. работы на странице.
				byPagesTechnicalWorksHttpProxyRouter.POST("/start", "", "Запустить тех. работы на странице. ", "",
					true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
					func(ctx *gin.Context) {
						type RequestArgs struct {
							ID        string `json:"id" yaml:"id" form:"id" description:"ID Http reverse прокси. " validate:"required"`
							URL       string `json:"url" yaml:"url" form:"url" description:"Страница " validate:"required"`
							StartTime int64  `json:"start_time" yaml:"start_time" form:"start_time" description:"Время старта тех. работ. " validate:"required"`
							EndTime   int64  `json:"end_time" yaml:"end_time" form:"end_time" description:"Время конца тех. работ. " validate:"required"`
						}

						type ResponseArgs struct {}

						requestArgs := new(RequestArgs)
						responseArgs := new(ResponseArgs)

						// Системная фигня
						{
							if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
								ctx.Abort()
								return
							}
						}

						// Обработка
						{
							prxs, err := proxyStorage.GetAll()
							if err != nil {
								engine.Modules.Logger.WARN(base_logger.Message{
									Sender: engine.Title,
									Text:   err.Error(),
								})
								engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
									"status": status.Failed,
									"code":   http.StatusBadRequest,
									"error": map[string]interface{}{
										"message": "Ошибка получения Http reverse прокси. ",
										"err":     err.Error(),
										"fields":  make(map[string]interface{}),
									},
								})
								ctx.Abort()
								return
							}

						found:
							for _, prx_ := range prxs {
								switch prx := prx_.(type) {
								case *http_reverse_proxy.Proxy:
									{
										if prx.ID.Hex() == requestArgs.ID {
											prx.TechnicalWorks.StartByPage(requestArgs.StartTime, requestArgs.EndTime, requestArgs.URL)

											break found
										}
									}
								}
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
		}
	}

	// Listener
	{
		listenerHttpProxyRouter := httpProxyRouter.Group("/listener", "Ветка слушателя Http Reverse прокси. ", "", true)

		// Запросы
		{
			// Остановить слушателя Http Reverse прокси.
			listenerHttpProxyRouter.POST("/stop", "", "Остановить слушателя Http Reverse прокси. ", "",
				true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
				func(ctx *gin.Context) {
					type RequestArgs struct {
						ID string `json:"id" yaml:"id" form:"id" description:"ID Http reverse прокси. " validate:"required"`
					}

					type ResponseArgs struct {}

					requestArgs := new(RequestArgs)
					responseArgs := new(ResponseArgs)

					// Системная фигня
					{
						if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
							ctx.Abort()
							return
						}
					}

					// Обработка
					{
						prxs, err := proxyStorage.GetAll()
						if err != nil {
							engine.Modules.Logger.WARN(base_logger.Message{
								Sender: engine.Title,
								Text:   err.Error(),
							})
							engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
								"status": status.Failed,
								"code":   http.StatusBadRequest,
								"error": map[string]interface{}{
									"message": "Ошибка получения Http reverse прокси. ",
									"err":     err.Error(),
									"fields":  make(map[string]interface{}),
								},
							})
							ctx.Abort()
							return
						}

					found:
						for _, prx_ := range prxs {
							switch prx := prx_.(type) {
							case *http_reverse_proxy.Proxy:
								{
									if prx.ID.Hex() == requestArgs.ID {
										if err := prx.Stop(); err != nil {
											engine.Modules.Logger.WARN(base_logger.Message{
												Sender: engine.Title,
												Text:   err.Error(),
											})
											engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
												"status": status.Failed,
												"code":   http.StatusBadRequest,
												"error": map[string]interface{}{
													"message": "Ошибка выключения Http reverse прокси. ",
													"err":     err.Error(),
													"fields":  make(map[string]interface{}),
												},
											})
											ctx.Abort()
											return
										}

										break found
									}
								}
							}
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

			// Запустить слушателя Http Reverse прокси.
			listenerHttpProxyRouter.POST("/start", "", "Запустить слушателя Http Reverse прокси. ", "",
				true, false, true, engine.Modules.SystemAccess.HttpMiddleware(false),
				func(ctx *gin.Context) {
					type RequestArgs struct {
						ID string `json:"id" yaml:"id" form:"id" description:"ID Http reverse прокси. " validate:"required"`
					}

					type ResponseArgs struct {}

					requestArgs := new(RequestArgs)
					responseArgs := new(ResponseArgs)

					// Системная фигня
					{
						if !engine.SystemHandle(ctx, requestArgs, responseArgs) {
							ctx.Abort()
							return
						}
					}

					// Обработка
					{
						prxs, err := proxyStorage.GetAll()
						if err != nil {
							engine.Modules.Logger.WARN(base_logger.Message{
								Sender: engine.Title,
								Text:   err.Error(),
							})
							engine.WriteResponse(ctx, http.StatusBadRequest, gin.H{
								"status": status.Failed,
								"code":   http.StatusBadRequest,
								"error": map[string]interface{}{
									"message": "Ошибка получения Http reverse прокси. ",
									"err":     err.Error(),
									"fields":  make(map[string]interface{}),
								},
							})
							ctx.Abort()
							return
						}

					found:
						for _, prx_ := range prxs {
							switch prx := prx_.(type) {
							case *http_reverse_proxy.Proxy:
								{
									if prx.ID.Hex() == requestArgs.ID {
										engine.Utils.Synchronizer.Proxy.WaitGroup.Add(1)

										go func() {
											defer engine.Utils.Synchronizer.Proxy.WaitGroup.Done()

											if err = prx.Listen(); err != nil {
												engine.Modules.Logger.WARN(base_logger.Message{
													Sender: engine.Title,
													Text:   err.Error(),
												})
												return
											}
										}()

										break found
									}
								}
							}
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
	}

	return nil
}

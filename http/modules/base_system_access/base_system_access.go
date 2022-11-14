package base_system_access

import (
	http_api_server_cfg "JkLNetDef/engine/config/http_api_server"
	"JkLNetDef/engine/databases"
	"JkLNetDef/engine/http/models/status"
	"JkLNetDef/engine/http/models/system/meta_data"
	"JkLNetDef/engine/http/models/system/session"
	"JkLNetDef/engine/http/models/system/system_access/http_request"
	"JkLNetDef/engine/http/models/system/system_access/token"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/models/user"
	"JkLNetDef/engine/modules/base_logger"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

// SystemAccess - система доступа
type SystemAccess struct {
	Title     string
	Databases *databases.Databases

	Utils   *Utils
	Modules *Modules
}

// Utils - утилиты системы доступа.
type Utils struct {
	Config *http_api_server_cfg.Server
}

// Modules - модули системы доступа.
type Modules struct {
	Logger          interfacies.Logger
	ManagerMetadata interfacies.ManagerMetaData
	ManagerSessions interfacies.ManagerSessions
	Authorizer      interfacies.Authorizer
}

// WriteHttpResponse - запись http ответа.
func (system *SystemAccess) WriteHttpResponse(ctx *gin.Context, statusCode int, response map[string]interface{}) {
	var responseType string

	// Проверка ctx
	{
		if ctx == nil {
			err := errors.New("ctx *gin.Context == nil. ")

			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})
			return
		}
	}

	// Чтение ?responseType=
	{
		responseType = ctx.Request.URL.Query().Get("responseType")
	}

	// Чтение response на основе responseType
	{
		switch strings.ToLower(responseType) {
		case "json":
			ctx.JSON(statusCode, response)
		case "json_np":
			ctx.JSONP(statusCode, response)
		case "pure_json":
			ctx.PureJSON(statusCode, response)
		case "secure_json":
			ctx.SecureJSON(statusCode, response)
		case "indented_json":
			ctx.IndentedJSON(statusCode, response)
		case "ascii_json":
			ctx.AsciiJSON(statusCode, response)
		case "yaml":
			ctx.YAML(statusCode, response)
		case "string":
			ctx.String(statusCode, "%+v", response)
		default:
			ctx.IndentedJSON(statusCode, response)
		}
	}

	// Запись ответа
	{
		buff, err := json.Marshal(response)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})
			return
		}
		ctx.Set("response", string(buff))
	}
}

// HttpMiddleware - возвращает middleware обработки системы доступа.
func (system *SystemAccess) HttpMiddleware(locked bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var req *http_request.Request
		var us *user.User
		var tok *token.Token

		// Закрыт ли запрос из кода
		{
			if locked {
				err = errors.New(fmt.Sprintf("Попытка вызова закрытого запроса - %s %s",
					ctx.Request.Method, ctx.Request.URL.Path))
				system.Modules.Logger.WARN(base_logger.Message{
					Sender: system.Utils.Config.SystemAccess.Title,
					Text:   err.Error(),
				})

				system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
					"status": status.Failed,
					"code":   http.StatusBadRequest,
					"error": map[string]interface{}{
						"message": err.Error(),
						"err":     err.Error(),
						"fields":  make(map[string]interface{}),
					},
				})
				ctx.Abort()
				return
			}
		}

		// Получение запроса
		{
			urlPath := ctx.Request.URL.Path

			// Магия со split'ом с конца
			{
				if len(ctx.Params) > 0 {
					urlSplit := strings.Split(strings.Trim(ctx.Request.URL.Path, " "), "/")
					params := ctx.Params

					for i := 0; i < len(urlSplit)/2; i++ {
						j := len(urlSplit) - i - 1
						urlSplit[i], urlSplit[j] = urlSplit[j], urlSplit[i]
					}

					for i := 0; i < len(params)/2; i++ {
						j := len(params) - i - 1
						params[i], params[j] = params[j], params[i]
					}
					for _, param := range params {
						for id, uri := range urlSplit {
							if uri == param.Value {
								urlSplit[id] = fmt.Sprintf(":%s", param.Key)
								break
							}
						}
					}

					for i := 0; i < len(urlSplit)/2; i++ {
						j := len(urlSplit) - i - 1
						urlSplit[i], urlSplit[j] = urlSplit[j], urlSplit[i]
					}

					urlPath = ""
					for _, str := range urlSplit {
						if str != "" {
							urlPath += fmt.Sprintf("/%s", str)
						}
					}
				}
			}

			req, err = system.Databases.Mongo.SystemAccess.HttpRequests.GetByURLAndMethod(ctx.Request.Method, urlPath)
			if err != nil {
				system.Modules.Logger.WARN(base_logger.Message{
					Sender: system.Utils.Config.SystemAccess.Title,
					Text:   err.Error(),
				})

				system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
					"status": status.Failed,
					"code":   http.StatusBadRequest,
					"error": map[string]interface{}{
						"message": "Не найден запрос в базе данных. ",
						"err":     err.Error(),
						"fields":  make(map[string]interface{}),
					},
				})
				ctx.Abort()
				return
			}
		}

		// Закрыт ли запрос
		{
			if req.Locked {
				err = errors.New(fmt.Sprintf("Попытка вызова закрытого запроса - %s %s",
					ctx.Request.Method, ctx.Request.URL.Path))
				system.Modules.Logger.WARN(base_logger.Message{
					Sender: system.Utils.Config.SystemAccess.Title,
					Text:   err.Error(),
				})

				system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
					"status": status.Failed,
					"code":   http.StatusBadRequest,
					"error": map[string]interface{}{
						"message": err.Error(),
						"err":     err.Error(),
						"fields":  make(map[string]interface{}),
					},
				})
				ctx.Abort()
				return
			}
		}

		isDashboard := strings.HasPrefix(req.URL, fmt.Sprintf("/%s", "dashboard"))

		// api
		{
			if !isDashboard {
				// Требуется ли авторизация
				{
					if req.Authorized {
						// Получение токена
						{
							strTok := ctx.Request.Header.Get("Authorization")

							if strTok == "" {
								dataTok, err := ctx.Cookie("token")
								if err != nil {
									system.Modules.Logger.WARN(base_logger.Message{
										Sender: system.Title,
										Text:   err.Error(),
									})

									system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
										"status": status.Failed,
										"code":   http.StatusBadRequest,
										"error": map[string]interface{}{
											"message": err.Error(),
											"err":     err.Error(),
											"fields":  make(map[string]interface{}),
										},
									})
									ctx.Abort()
									return
								}

								tok, err = system.Databases.Mongo.SystemAccess.Tokens.GetByData(dataTok)
								if err != nil {
									system.Modules.Logger.WARN(base_logger.Message{
										Sender: system.Title,
										Text:   err.Error(),
									})

									system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
										"status": status.Failed,
										"code":   http.StatusBadRequest,
										"error": map[string]interface{}{
											"message": "Токен не найден. ",
											"err":     err.Error(),
											"fields":  make(map[string]interface{}),
										},
									})
									ctx.Abort()
									return
								}
							} else {
								dataTok := strings.Replace(strTok, "Bearer ", "", 1)

								tok, err = system.Databases.Mongo.SystemAccess.Tokens.GetByData(dataTok)
								if err != nil {
									system.Modules.Logger.WARN(base_logger.Message{
										Sender: system.Utils.Config.SystemAccess.Title,
										Text:   err.Error(),
									})

									system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
										"status": status.Failed,
										"code":   http.StatusBadRequest,
										"error": map[string]interface{}{
											"message": "Токен не найден. ",
											"err":     err.Error(),
											"fields":  make(map[string]interface{}),
										},
									})
									ctx.Abort()
									return
								}
							}

							if tok.Expire <= time.Now().Unix() {
								err := errors.New("Токен не действителен. ")
								system.Modules.Logger.WARN(base_logger.Message{
									Sender: system.Utils.Config.SystemAccess.Title,
									Text:   err.Error(),
								})

								system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
									"status": status.Failed,
									"code":   http.StatusBadRequest,
									"error": map[string]interface{}{
										"message": err.Error(),
										"err":     err.Error(),
										"fields":  make(map[string]interface{}),
									},
								})
								ctx.Abort()
								return
							}
						}

						// Получение пользователя
						{
							us, err = system.Databases.Mongo.Main.Users.GetByID(tok.Owner)
							if err != nil {
								system.Modules.Logger.WARN(base_logger.Message{
									Sender: system.Utils.Config.SystemAccess.Title,
									Text:   err.Error(),
								})

								system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
									"status": status.Failed,
									"code":   http.StatusBadRequest,
									"error": map[string]interface{}{
										"message": "Не найден пользователь в базе данных. ",
										"err":     err.Error(),
										"fields":  make(map[string]interface{}),
									},
								})
								ctx.Abort()
								return
							}

							ctx.Set("login", us.Login)
						}

						// Проверка на доступ
						{
							found := true

							if us.Login != "admin" {
								found = false
							}

							if !found {
								for _, rlID := range us.Roles {
									rl, err := system.Databases.Mongo.SystemAccess.Roles.GetByID(rlID)
									if err != nil {
										system.Modules.Logger.WARN(base_logger.Message{
											Sender: system.Utils.Config.SystemAccess.Title,
											Text:   err.Error(),
										})

										system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
											"status": status.Failed,
											"code":   http.StatusBadRequest,
											"error": map[string]interface{}{
												"message": "Не найдена роль в базе данных. ",
												"err":     err.Error(),
												"fields":  make(map[string]interface{}),
											},
										})
										ctx.Abort()
										return
									}

									for _, rqID := range rl.HttpRequests {
										if rqID == fmt.Sprintf("%s %s", req.Method, req.URL) {
											found = true
											break
										}
									}

									if !found {
										for _, modID := range rl.Modules {
											mod, err := system.Databases.Mongo.SystemAccess.Modules.GetByID(modID)
											if err != nil {
												system.Modules.Logger.WARN(base_logger.Message{
													Sender: system.Utils.Config.SystemAccess.Title,
													Text:   err.Error(),
												})

												system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
													"status": status.Failed,
													"code":   http.StatusBadRequest,
													"error": map[string]interface{}{
														"message": "Не найден модуль в базе данных. ",
														"err":     err.Error(),
														"fields":  make(map[string]interface{}),
													},
												})
												ctx.Abort()
												return
											}

											for _, rqID := range mod.HttpRequests {
												if rqID == fmt.Sprintf("%s %s", req.Method, req.URL) {
													found = true
													break
												}
											}

											if found {
												break
											}
										}
									}

									if found {
										break
									}
								}
							}

							if !found {
								err := errors.New("Нет доступа. ")
								system.Modules.Logger.WARN(base_logger.Message{
									Sender: system.Utils.Config.SystemAccess.Title,
									Text:   err.Error(),
								})

								system.WriteHttpResponse(ctx, http.StatusBadRequest, gin.H{
									"status": status.Failed,
									"code":   http.StatusBadRequest,
									"error": map[string]interface{}{
										"message": err.Error(),
										"err":     err.Error(),
										"fields":  make(map[string]interface{}),
									},
								})
								ctx.Abort()
								return
							}
						}
					}
				}
			}
		}

		// dashboard
		{
			if isDashboard {
				// Требуется ли авторизация
				{
					if req.Authorized {
						// Получение токена
						{
							strTok := ctx.Request.Header.Get("Authorization")

							if strTok == "" {
								dataTok, err := ctx.Cookie("token")
								if err != nil {
									system.Modules.Logger.WARN(base_logger.Message{
										Sender: system.Title,
										Text:   err.Error(),
									})

									_, err = system.HttpLogout(ctx)
									if err != nil {
										system.Modules.Logger.WARN(base_logger.Message{
											Sender: system.Utils.Config.SystemAccess.Title,
											Text:   err.Error(),
										})
									}

									ctx.Redirect(http.StatusTemporaryRedirect, "/dashboard/auth")
									ctx.Abort()
									return
								}

								tok, err = system.Databases.Mongo.SystemAccess.Tokens.GetByData(dataTok)
								if err != nil {
									system.Modules.Logger.WARN(base_logger.Message{
										Sender: system.Title,
										Text:   err.Error(),
									})

									_, err = system.HttpLogout(ctx)
									if err != nil {
										system.Modules.Logger.WARN(base_logger.Message{
											Sender: system.Utils.Config.SystemAccess.Title,
											Text:   err.Error(),
										})
									}

									ctx.Redirect(http.StatusTemporaryRedirect, "/dashboard/auth")
									ctx.Abort()
									return
								}
							} else {
								dataTok := strings.Replace(strTok, "Bearer ", "", 1)

								tok, err = system.Databases.Mongo.SystemAccess.Tokens.GetByData(dataTok)
								if err != nil {
									system.Modules.Logger.WARN(base_logger.Message{
										Sender: system.Utils.Config.SystemAccess.Title,
										Text:   err.Error(),
									})

									_, err = system.HttpLogout(ctx)
									if err != nil {
										system.Modules.Logger.WARN(base_logger.Message{
											Sender: system.Utils.Config.SystemAccess.Title,
											Text:   err.Error(),
										})
									}

									ctx.Redirect(http.StatusTemporaryRedirect, "/dashboard/auth")
									ctx.Abort()
									return
								}
							}

							if tok.Expire <= time.Now().Unix() {
								err := errors.New("Токен не действителен. ")
								system.Modules.Logger.WARN(base_logger.Message{
									Sender: system.Utils.Config.SystemAccess.Title,
									Text:   err.Error(),
								})

								_, err = system.HttpLogout(ctx)
								if err != nil {
									system.Modules.Logger.WARN(base_logger.Message{
										Sender: system.Utils.Config.SystemAccess.Title,
										Text:   err.Error(),
									})
								}

								ctx.Redirect(http.StatusTemporaryRedirect, "/dashboard/auth")
								ctx.Abort()
								return
							}
						}

						// Получение пользователя
						{
							us, err = system.Databases.Mongo.Main.Users.GetByID(tok.Owner)
							if err != nil {
								system.Modules.Logger.WARN(base_logger.Message{
									Sender: system.Utils.Config.SystemAccess.Title,
									Text:   err.Error(),
								})

								_, err = system.HttpLogout(ctx)
								if err != nil {
									system.Modules.Logger.WARN(base_logger.Message{
										Sender: system.Utils.Config.SystemAccess.Title,
										Text:   err.Error(),
									})
								}

								ctx.Redirect(http.StatusTemporaryRedirect, "/dashboard/auth")
								ctx.Abort()
								return
							}

							ctx.Set("login", us.Login)
						}

						// Проверка на доступ
						{
							found := true

							if us.Login != "admin" {
								found = false
							}

							if !found {
								for _, rlID := range us.Roles {
									rl, err := system.Databases.Mongo.SystemAccess.Roles.GetByID(rlID)
									if err != nil {
										system.Modules.Logger.WARN(base_logger.Message{
											Sender: system.Utils.Config.SystemAccess.Title,
											Text:   err.Error(),
										})

										_, err = system.HttpLogout(ctx)
										if err != nil {
											system.Modules.Logger.WARN(base_logger.Message{
												Sender: system.Utils.Config.SystemAccess.Title,
												Text:   err.Error(),
											})
										}

										ctx.Redirect(http.StatusTemporaryRedirect, "/dashboard/auth")
										ctx.Abort()
										return
									}

									for _, rqID := range rl.HttpRequests {
										if rqID == fmt.Sprintf("%s %s", req.Method, req.URL) {
											found = true
											break
										}
									}

									if !found {
										for _, modID := range rl.Modules {
											mod, err := system.Databases.Mongo.SystemAccess.Modules.GetByID(modID)
											if err != nil {
												system.Modules.Logger.WARN(base_logger.Message{
													Sender: system.Utils.Config.SystemAccess.Title,
													Text:   err.Error(),
												})

												_, err = system.HttpLogout(ctx)
												if err != nil {
													system.Modules.Logger.WARN(base_logger.Message{
														Sender: system.Utils.Config.SystemAccess.Title,
														Text:   err.Error(),
													})
												}

												ctx.Redirect(http.StatusTemporaryRedirect, "/dashboard/auth")
												ctx.Abort()
												return
											}

											for _, rqID := range mod.HttpRequests {
												if rqID == fmt.Sprintf("%s %s", req.Method, req.URL) {
													found = true
													break
												}
											}

											if found {
												break
											}
										}
									}

									if found {
										break
									}
								}
							}

							if !found {
								err = errors.New("Нет доступа. ")
								system.Modules.Logger.WARN(base_logger.Message{
									Sender: system.Utils.Config.SystemAccess.Title,
									Text:   err.Error(),
								})

								_, err = system.HttpLogout(ctx)
								if err != nil {
									system.Modules.Logger.WARN(base_logger.Message{
										Sender: system.Utils.Config.SystemAccess.Title,
										Text:   err.Error(),
									})
								}

								ctx.Redirect(http.StatusTemporaryRedirect, "/dashboard/auth")
								ctx.Abort()
								return
							}
						}
					}
				}
			}
		}

		ctx.Next()
	}
}

// Authorizer - получить authorizer.Authorizer.
func (system *SystemAccess) Authorizer() interfacies.Authorizer {
	return system.Modules.Authorizer
}

// HttpAuth - авторизация пользователя.
func (system *SystemAccess) HttpAuth(ctx *gin.Context, login, password string) (*token.Token, map[string]interface{}, string, error) {
	var tok *token.Token
	var us *user.User

	// Проверка ctx
	{
		if ctx == nil {
			err := errors.New("ctx *gin.Context == nil. ")

			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return nil, make(map[string]interface{}), "Ошибка авторизации пользователя. ", err
		}
	}

	// Получение пользователя по логину
	{
		var err error

		us, err = system.Databases.Mongo.Main.Users.GetByLogin(login)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			fields := map[string]interface{}{
				"login": "Пользователь не найден. ",
			}

			return nil, fields, "Пользователь не найден. ", err
		}
	}

	// Проверка пароля
	{
		if us.Password != password {
			err := errors.New("Неправильный пароль. ")
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			fields := map[string]interface{}{
				"password": "Неправильный пароль. ",
			}

			return nil, fields, "Неправильный пароль. ", err
		}
	}

	// Создание токена
	{
		claims_, data, err := system.Authorizer().SignIn(ctx, us, time.Now().Unix()+int64(system.Utils.Config.SystemAccess.Token.LifeTime))
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return nil, make(map[string]interface{}), "Ошибка создания токена. ", err
		}

		tok = &token.Token{
			ID:      primitive.NewObjectID(),
			Data:    data,
			Owner:   us.ID,
			Created: claims_.StandardClaims.IssuedAt,
			Expire:  claims_.StandardClaims.ExpiresAt,

			Meta: &meta_data.MetaData{
				DateCreated: time.Now(),
				Created:     us.ID,
				DateChanged: time.Now(),
				Changed:     us.ID,
			},
		}
	}

	// Сохранение в бд
	{
		err := system.Databases.Mongo.SystemAccess.Tokens.Add(tok)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return nil, make(map[string]interface{}), "Ошибка сохранения токена. ", err
		}
	}

	// Создание сессии
	{
		sess := &session.Session{
			ID:    primitive.NewObjectID(),
			Token: tok.ID,
			Data:  make(map[string]interface{}),
			Meta: &meta_data.MetaData{
				DateCreated: time.Now(),
				Created:     us.ID,
				DateChanged: time.Now(),
				Changed:     us.ID,
			},
		}

		sess.Data["user_id"] = us.ID.Hex()
		sess.Data["token_data"] = tok.Data
		sess.Data["user_login"] = us.Login

		err := system.Databases.Mongo.System.Sessions.Add(sess)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return nil, make(map[string]interface{}), "Ошибка сохранения сессии. ", err
		}
	}

	// Запись в куки
	{
		err := system.SetTokenInHttpCookie(ctx, tok)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return nil, make(map[string]interface{}), "Ошибка сохранения токена. ", err
		}
	}

	return tok, make(map[string]interface{}), "", nil
}

// HttpLogout - выход пользователя.
func (system *SystemAccess) HttpLogout(ctx *gin.Context) (string, error) {
	var tok *token.Token

	// Проверка ctx
	{
		if ctx == nil {
			err := errors.New("ctx *gin.Context == nil. ")

			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return "Ошибка очищение сессии пользователя. ", err
		}
	}

	// Получение токена
	{
		tokData, err := system.GetTokenFromHttpCookie(ctx)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return "Ошибка получения токена. ", err
		}

		tok, err = system.Databases.Mongo.SystemAccess.Tokens.GetByData(tokData)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return "Ошибка получения токена. ", err
		}
	}

	// Удаление токена
	{
		err := system.Databases.Mongo.SystemAccess.Tokens.RemoveByID(tok.ID)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return "Ошибка удаления токена. ", err
		}
	}

	// Удаление сессии
	{
		sess, err := system.Databases.Mongo.System.Sessions.GetByTokenID(tok.ID)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return "Ошибка получения сессии. ", err
		}

		err = system.Databases.Mongo.System.Sessions.RemoveByID(sess.ID)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Utils.Config.SystemAccess.Title,
				Text:   err.Error(),
			})

			return "Ошибка удаления сессии. ", err
		}
	}

	return "", nil
}

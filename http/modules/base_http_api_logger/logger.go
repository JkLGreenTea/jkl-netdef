package base_http_api_logger

import (
	loggers_cfg "JkLNetDef/engine/config/loggers"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"JkLNetDef/engine/utils/synchronizer"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	Number = 1
)

// Logger - логгер HTTP Api запросов.
type Logger struct {
	Title string

	GlobalLogger    interfacies.Logger
	Config          *loggers_cfg.Api
	ManagerSessions interfacies.ManagerSessions
	Sync            *synchronizer.Synchronizer

	Mx *sync.Mutex
}

// LogFormatterParams - параметры для вывода.
type LogFormatterParams struct {
	Request *http.Request

	// TimeStamp - показывает время после того, как сервер вернет ответ.
	TimeStamp time.Time
	// StatusCode - код ответа HTTP.
	StatusCode int
	// Latency - сколько времени сервер затратил на обработку определенного запроса.
	Latency time.Duration
	// ClientAddr - адрес клиента.
	ClientAddr string
	// Method - метод HTTP, заданный для запроса.
	Method string
	// Path - путь, который запрашивает клиент.
	Path string
	// ErrorMessage - установите, произошла ли ошибка при обработке запроса.
	ErrorMessage string
	// BodySize - размер тела запроса.
	BodySize int
	// Keys - ключи, установленные в контексте запроса.
	Keys map[string]interface{}
	// Login - логин пользователя.
	Login string
	// FormArgs - form-data аргументы.
	FormArgs map[string][]string
	// JsonArgs - json аргументы.
	JsonArgs string
	//ContentType - тип контента запроса
	ContentType string
}

// New - создает логгер запросов
func New(title string, cfg *loggers_cfg.Api, glogger interfacies.Logger, sync_ *synchronizer.Synchronizer, managerSessions interfacies.ManagerSessions) *Logger {
	year, month, day := time.Now().Date()
	path_ := path.Join(cfg.LogFilePath, fmt.Sprintf("%d_%s_%d", year, month.String(), day))
	if ok, _ := exists(path_); ok {
		files, err := ioutil.ReadDir(path_)
		if err != nil {
			glogger.FATAL(base_logger.Message{
				Sender: title,
				Text:   err.Error(),
			})

			return nil
		}
		Number = len(files) + 1
	}

	return &Logger{
		GlobalLogger:    glogger,
		ManagerSessions: managerSessions,
		Config:          cfg,
		Sync:            sync_,

		Mx: new(sync.Mutex),
	}
}

// HttpMiddleware - возвращает middleware обработки лога запросов.
func (logger *Logger) HttpMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Sync.Logger.WaitGroup.Add(1)

		formatParams := new(LogFormatterParams)

		start := time.Now()

		defer func() {
			// Заполнение параметров формата лога.
			{
				formatParams.Request = ctx.Request
				formatParams.Keys = ctx.Keys
				formatParams.TimeStamp = time.Now()
				formatParams.Login = "unknown"
				formatParams.ContentType = ctx.Request.Header.Get("Content-Type")
				formatParams.Method = ctx.Request.Method
				formatParams.StatusCode = ctx.Writer.Status()
				formatParams.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()
				formatParams.BodySize = ctx.Writer.Size()
				formatParams.Path = ctx.Request.URL.Path
				formatParams.Latency = formatParams.TimeStamp.Sub(start)

				if ctx.Request.Header.Get("X-Real-Ip") != "" {
					formatParams.ClientAddr = ctx.Request.Header.Get("X-Real-Ip")
				} else {
					formatParams.ClientAddr = ctx.Request.RemoteAddr
				}
			}

			if ctx.Request.Header.Get("Authorization") != "" {
				userLogin, err := logger.ManagerSessions.GetUserLogin(ctx)
				if err != nil {
					logger.GlobalLogger.WARN(base_logger.Message{
						Sender: logger.Title,
						Text:   err.Error(),
					})
				}

				if userLogin != "" {
					formatParams.Login = userLogin
				}
			}

			switch code := formatParams.StatusCode; {
			case code >= 100 && code < 200 || code >= 200 && code < 300 || code >= 300 && code < 400:
				{
					// Вывод INFO
					{
						logger.Sync.Logger.WaitGroup.Add(2)

						// Вывод в консоль
						{
							go func() {
								defer logger.Sync.Logger.WaitGroup.Done()

								var text string
								// Формирование text
								{
									var pheader, header, addr, workerTime, callerPath, statusCode, httpMethod string

									// Запись данных
									{
										// header
										{
											switch logger.Config.Colors.INFO {
											case "Reset":
												header = color.Ize(color.Reset, "INFO")
											case "Bold":
												header = color.Ize(color.Bold, "INFO")
											case "Red":
												header = color.Ize(color.Red, "INFO")
											case "Green":
												header = color.Ize(color.Green, "INFO")
											case "Yellow":
												header = color.Ize(color.Yellow, "INFO")
											case "Blue":
												header = color.Ize(color.Blue, "INFO")
											case "Purple":
												header = color.Ize(color.Purple, "INFO")
											case "Cyan":
												header = color.Ize(color.Cyan, "INFO")
											case "Gray":
												header = color.Ize(color.Gray, "INFO")
											case "White":
												header = color.Ize(color.White, "INFO")
											}
										}

										// pheader
										{
											switch logger.Config.Colors.ALogColor {
											case "Reset":
												pheader = color.Ize(color.Reset, "A")
											case "Bold":
												pheader = color.Ize(color.Bold, "A")
											case "Red":
												pheader = color.Ize(color.Red, "A")
											case "Green":
												pheader = color.Ize(color.Green, "A")
											case "Yellow":
												pheader = color.Ize(color.Yellow, "A")
											case "Blue":
												pheader = color.Ize(color.Blue, "A")
											case "Purple":
												pheader = color.Ize(color.Purple, "A")
											case "Cyan":
												pheader = color.Ize(color.Cyan, "A")
											case "Gray":
												pheader = color.Ize(color.Gray, "A")
											case "White":
												pheader = color.Ize(color.White, "A")
											}
										}

										// addr
										{
											addr = formatParams.ClientAddr
											if len(formatParams.ClientAddr) < 22 {
												position := true
												for i := 0; i < 22-len(formatParams.ClientAddr); i++ {
													if position {
														addr += " "
													} else {
														addr = " " + addr
													}

													position = !position
												}
											}
										}

										// workerTime
										{
											subTm := formatParams.TimeStamp.Sub(start).String()
											workerTime = subTm

											if len(subTm) < 14 {
												position := true
												for i := 0; i < 14-len(subTm); i++ {
													if position {
														workerTime += " "
													} else {
														workerTime = " " + workerTime
													}

													position = !position
												}
											}
										}

										// callerPath
										{
											if logger.Config.EnableCallerPath.INFO {
												_, file, no, ok := runtime.Caller(1)
												if ok {
													folder, _ := os.Getwd()
													callerPath = fmt.Sprintf("-[%s - %d]", strings.Replace(file, strings.Replace(folder, "\\", "/", -1), "", 1), no)
												}
											}
										}

										// httpMethod
										{
											switch formatParams.Method {
											case http.MethodGet:
												{
													switch logger.Config.Colors.GET {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodPost:
												{
													switch logger.Config.Colors.POST {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodPut:
												{
													switch logger.Config.Colors.PUT {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodDelete:
												{
													switch logger.Config.Colors.DEBUG {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodPatch:
												{
													switch logger.Config.Colors.PATCH {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodOptions:
												{
													switch logger.Config.Colors.OPTIONS {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodHead:
												{
													switch logger.Config.Colors.HEAD {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											}
										}

										// statusCode
										{
											statusCode = strconv.Itoa(code)

											if code >= 100 && code < 200 {
												switch logger.Config.Colors.HTTPCode100 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 200 && code < 300 {
												switch logger.Config.Colors.HTTPCode200 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 300 && code < 400 {
												switch logger.Config.Colors.HTTPCode300 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 400 && code < 500 {
												switch logger.Config.Colors.HTTPCode400 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 500 && code < 600 {
												switch logger.Config.Colors.HTTPCode500 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											}
										}
									}

									{
										text = fmt.Sprintf("[%s][%s]-|%s|-|%s|-[%s]-[%s]%s: %s %s \n",
											pheader,
											header,
											start.Format(logger.Config.TimeFormat),
											addr,
											statusCode,
											workerTime,
											callerPath,
											httpMethod,
											ctx.Request.Host+ctx.Request.URL.Path)
									}
								}

								fmt.Print(text)
							}()
						}

						// Вывод в файл
						{
							go func() {
								defer logger.Sync.Logger.WaitGroup.Done()

								if logger.Config.EnableOutputFile.INFO {
									var text, callerPath string
									var err error
									var responseArgs, requestArgs []byte

									responseArgs_, ok1 := ctx.Get("response")
									requestArgs_, ok2 := ctx.Get("request")

									if ok1 && ok2 {
										responseArgs, err = json.Marshal(responseArgs_)
										requestArgs, err = json.Marshal(requestArgs_)
									}

									// callerPath
									{
										if logger.Config.EnableCallerPath.INFO {
											_, file, no, ok := runtime.Caller(1)
											if ok {
												folder, _ := os.Getwd()
												callerPath = fmt.Sprintf("-[%s - %d]", strings.Replace(file, strings.Replace(folder, "\\", "/", -1), "", 1), no)
											}
										}
									}

									// Формирование text
									{
										text = fmt.Sprintf("[A][INFO]-|%s|-|%s|-[%d]-[%s]%s: %s %s {\"response\": %s, \"request\": %s} \n",
											start.Format(logger.Config.TimeFormat),
											formatParams.ClientAddr,
											formatParams.StatusCode,
											formatParams.Latency.String(),
											callerPath,
											formatParams.Method,
											ctx.Request.Host+ctx.Request.URL.Path,
											string(responseArgs),
											string(requestArgs))
									}

									err = logger.WriteLogFile(text)
									if err != nil {
										return
									}
								}
							}()
						}
					}
				}
			case code >= 400 && code < 500:
				{
					// Вывод WARN
					{
						logger.Sync.Logger.WaitGroup.Add(2)

						// Вывод в консоль
						{
							go func() {
								defer logger.Sync.Logger.WaitGroup.Done()

								var text string
								// Формирование text
								{
									var pheader, header, addr, workerTime, callerPath, statusCode, httpMethod string

									// Запись данных
									{
										// header
										{
											switch logger.Config.Colors.WARN {
											case "Reset":
												header = color.Ize(color.Reset, "WARN")
											case "Bold":
												header = color.Ize(color.Bold, "WARN")
											case "Red":
												header = color.Ize(color.Red, "WARN")
											case "Green":
												header = color.Ize(color.Green, "WARN")
											case "Yellow":
												header = color.Ize(color.Yellow, "WARN")
											case "Blue":
												header = color.Ize(color.Blue, "WARN")
											case "Purple":
												header = color.Ize(color.Purple, "WARN")
											case "Cyan":
												header = color.Ize(color.Cyan, "WARN")
											case "Gray":
												header = color.Ize(color.Gray, "WARN")
											case "White":
												header = color.Ize(color.White, "WARN")
											}
										}

										// pheader
										{
											switch logger.Config.Colors.ALogColor {
											case "Reset":
												pheader = color.Ize(color.Reset, "A")
											case "Bold":
												pheader = color.Ize(color.Bold, "A")
											case "Red":
												pheader = color.Ize(color.Red, "A")
											case "Green":
												pheader = color.Ize(color.Green, "A")
											case "Yellow":
												pheader = color.Ize(color.Yellow, "A")
											case "Blue":
												pheader = color.Ize(color.Blue, "A")
											case "Purple":
												pheader = color.Ize(color.Purple, "A")
											case "Cyan":
												pheader = color.Ize(color.Cyan, "A")
											case "Gray":
												pheader = color.Ize(color.Gray, "A")
											case "White":
												pheader = color.Ize(color.White, "A")
											}
										}

										// addr
										{
											addr = formatParams.ClientAddr
											if len(formatParams.ClientAddr) < 22 {
												position := true
												for i := 0; i < 22-len(formatParams.ClientAddr); i++ {
													if position {
														addr += " "
													} else {
														addr = " " + addr
													}

													position = !position
												}
											}
										}

										// workerTime
										{
											subTm := formatParams.TimeStamp.Sub(start).String()
											workerTime = subTm

											if len(subTm) < 14 {
												position := true
												for i := 0; i < 14-len(subTm); i++ {
													if position {
														workerTime += " "
													} else {
														workerTime = " " + workerTime
													}

													position = !position
												}
											}
										}

										// callerPath
										{
											if logger.Config.EnableCallerPath.WARN {
												_, file, no, ok := runtime.Caller(1)
												if ok {
													folder, _ := os.Getwd()
													callerPath = fmt.Sprintf("-[%s - %d]", strings.Replace(file, strings.Replace(folder, "\\", "/", -1), "", 1), no)
												}
											}
										}

										// httpMethod
										{
											switch formatParams.Method {
											case http.MethodGet:
												{
													switch logger.Config.Colors.GET {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodPost:
												{
													switch logger.Config.Colors.POST {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodPut:
												{
													switch logger.Config.Colors.PUT {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodDelete:
												{
													switch logger.Config.Colors.DEBUG {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodPatch:
												{
													switch logger.Config.Colors.PATCH {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodOptions:
												{
													switch logger.Config.Colors.OPTIONS {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodHead:
												{
													switch logger.Config.Colors.HEAD {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											}
										}

										// statusCode
										{
											statusCode = strconv.Itoa(code)

											if code >= 100 && code < 200 {
												switch logger.Config.Colors.HTTPCode100 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 200 && code < 300 {
												switch logger.Config.Colors.HTTPCode200 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 300 && code < 400 {
												switch logger.Config.Colors.HTTPCode300 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 400 && code < 500 {
												switch logger.Config.Colors.HTTPCode400 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 500 && code < 600 {
												switch logger.Config.Colors.HTTPCode500 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											}
										}
									}

									{
										text = fmt.Sprintf("[%s][%s]-|%s|-|%s|-[%s]-[%s]%s: %s %s \n",
											pheader,
											header,
											start.Format(logger.Config.TimeFormat),
											addr,
											statusCode,
											workerTime,
											callerPath,
											httpMethod,
											ctx.Request.Host+ctx.Request.URL.Path)
									}
								}

								fmt.Print(text)
							}()
						}

						// Вывод в файл
						{
							go func() {
								defer logger.Sync.Logger.WaitGroup.Done()

								if logger.Config.EnableOutputFile.WARN {
									var text, callerPath string
									var err error
									var responseArgs, requestArgs []byte

									responseArgs_, ok1 := ctx.Get("response")
									requestArgs_, ok2 := ctx.Get("request")

									if ok1 && ok2 {
										responseArgs, err = json.Marshal(responseArgs_)
										requestArgs, err = json.Marshal(requestArgs_)
									}

									// callerPath
									{
										if logger.Config.EnableCallerPath.WARN {
											_, file, no, ok := runtime.Caller(1)
											if ok {
												folder, _ := os.Getwd()
												callerPath = fmt.Sprintf("-[%s - %d]", strings.Replace(file, strings.Replace(folder, "\\", "/", -1), "", 1), no)
											}
										}
									}

									// Формирование text
									{
										text = fmt.Sprintf("[A][WARN]-|%s|-|%s|-[%d]-[%s]%s: %s %s {\"response\": %s, \"request\": %s} \n",
											start.Format(logger.Config.TimeFormat),
											formatParams.ClientAddr,
											formatParams.StatusCode,
											formatParams.Latency.String(),
											callerPath,
											formatParams.Method,
											ctx.Request.Host+ctx.Request.URL.Path,
											string(responseArgs),
											string(requestArgs))
									}

									err = logger.WriteLogFile(text)
									if err != nil {
										return
									}
								}
							}()
						}
					}
				}
			case code >= 500 && code < 600:
				{
					// Вывод ERROR
					{
						logger.Sync.Logger.WaitGroup.Add(2)

						// Вывод в консоль
						{
							go func() {
								defer logger.Sync.Logger.WaitGroup.Done()

								var text string
								// Формирование text
								{
									var pheader, header, addr, workerTime, callerPath, statusCode, httpMethod string

									// Запись данных
									{
										// header
										{
											switch logger.Config.Colors.ERROR {
											case "Reset":
												header = color.Ize(color.Reset, "ERROR")
											case "Bold":
												header = color.Ize(color.Bold, "ERROR")
											case "Red":
												header = color.Ize(color.Red, "ERROR")
											case "Green":
												header = color.Ize(color.Green, "ERROR")
											case "Yellow":
												header = color.Ize(color.Yellow, "ERROR")
											case "Blue":
												header = color.Ize(color.Blue, "ERROR")
											case "Purple":
												header = color.Ize(color.Purple, "ERROR")
											case "Cyan":
												header = color.Ize(color.Cyan, "ERROR")
											case "Gray":
												header = color.Ize(color.Gray, "ERROR")
											case "White":
												header = color.Ize(color.White, "ERROR")
											}
										}

										// pheader
										{
											switch logger.Config.Colors.ALogColor {
											case "Reset":
												pheader = color.Ize(color.Reset, "A")
											case "Bold":
												pheader = color.Ize(color.Bold, "A")
											case "Red":
												pheader = color.Ize(color.Red, "A")
											case "Green":
												pheader = color.Ize(color.Green, "A")
											case "Yellow":
												pheader = color.Ize(color.Yellow, "A")
											case "Blue":
												pheader = color.Ize(color.Blue, "A")
											case "Purple":
												pheader = color.Ize(color.Purple, "A")
											case "Cyan":
												pheader = color.Ize(color.Cyan, "A")
											case "Gray":
												pheader = color.Ize(color.Gray, "A")
											case "White":
												pheader = color.Ize(color.White, "A")
											}
										}

										// addr
										{
											addr = formatParams.ClientAddr
											if len(formatParams.ClientAddr) < 22 {
												position := true
												for i := 0; i < 22-len(formatParams.ClientAddr); i++ {
													if position {
														addr += " "
													} else {
														addr = " " + addr
													}

													position = !position
												}
											}
										}

										// workerTime
										{
											subTm := formatParams.TimeStamp.Sub(start).String()
											workerTime = subTm

											if len(subTm) < 14 {
												position := true
												for i := 0; i < 14-len(subTm); i++ {
													if position {
														workerTime += " "
													} else {
														workerTime = " " + workerTime
													}

													position = !position
												}
											}
										}

										// callerPath
										{
											if logger.Config.EnableCallerPath.ERROR {
												_, file, no, ok := runtime.Caller(1)
												if ok {
													folder, _ := os.Getwd()
													callerPath = fmt.Sprintf("-[%s - %d]", strings.Replace(file, strings.Replace(folder, "\\", "/", -1), "", 1), no)
												}
											}
										}

										// httpMethod
										{
											switch formatParams.Method {
											case http.MethodGet:
												{
													switch logger.Config.Colors.GET {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodPost:
												{
													switch logger.Config.Colors.POST {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodPut:
												{
													switch logger.Config.Colors.PUT {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodDelete:
												{
													switch logger.Config.Colors.DEBUG {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodPatch:
												{
													switch logger.Config.Colors.PATCH {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodOptions:
												{
													switch logger.Config.Colors.OPTIONS {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											case http.MethodHead:
												{
													switch logger.Config.Colors.HEAD {
													case "Reset":
														httpMethod = color.Ize(color.Reset, formatParams.Method)
													case "Bold":
														httpMethod = color.Ize(color.Bold, formatParams.Method)
													case "Red":
														httpMethod = color.Ize(color.Red, formatParams.Method)
													case "Green":
														httpMethod = color.Ize(color.Green, formatParams.Method)
													case "Yellow":
														httpMethod = color.Ize(color.Yellow, formatParams.Method)
													case "Blue":
														httpMethod = color.Ize(color.Blue, formatParams.Method)
													case "Purple":
														httpMethod = color.Ize(color.Purple, formatParams.Method)
													case "Cyan":
														httpMethod = color.Ize(color.Cyan, formatParams.Method)
													case "Gray":
														httpMethod = color.Ize(color.Gray, formatParams.Method)
													case "White":
														httpMethod = color.Ize(color.White, formatParams.Method)
													}
												}
											}
										}

										// statusCode
										{
											statusCode = strconv.Itoa(code)

											if code >= 100 && code < 200 {
												switch logger.Config.Colors.HTTPCode100 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 200 && code < 300 {
												switch logger.Config.Colors.HTTPCode200 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 300 && code < 400 {
												switch logger.Config.Colors.HTTPCode300 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 400 && code < 500 {
												switch logger.Config.Colors.HTTPCode400 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											} else if code >= 500 && code < 600 {
												switch logger.Config.Colors.HTTPCode500 {
												case "Reset":
													statusCode = color.Ize(color.Reset, statusCode)
												case "Bold":
													statusCode = color.Ize(color.Bold, statusCode)
												case "Red":
													statusCode = color.Ize(color.Red, statusCode)
												case "Green":
													statusCode = color.Ize(color.Green, statusCode)
												case "Yellow":
													statusCode = color.Ize(color.Yellow, statusCode)
												case "Blue":
													statusCode = color.Ize(color.Blue, statusCode)
												case "Purple":
													statusCode = color.Ize(color.Purple, statusCode)
												case "Cyan":
													statusCode = color.Ize(color.Cyan, statusCode)
												case "Gray":
													statusCode = color.Ize(color.Gray, statusCode)
												case "White":
													statusCode = color.Ize(color.White, statusCode)
												}
											}
										}
									}

									{
										text = fmt.Sprintf("[%s][%s]-|%s|-|%s|-[%s]-[%s]%s: %s %s \n",
											pheader,
											header,
											start.Format(logger.Config.TimeFormat),
											addr,
											statusCode,
											workerTime,
											callerPath,
											httpMethod,
											ctx.Request.Host+ctx.Request.URL.Path)
									}
								}

								fmt.Print(text)
							}()
						}

						// Вывод в файл
						{
							go func() {
								defer logger.Sync.Logger.WaitGroup.Done()

								if logger.Config.EnableOutputFile.ERROR {
									var text, callerPath string
									var err error
									var responseArgs, requestArgs []byte

									responseArgs_, ok1 := ctx.Get("response")
									requestArgs_, ok2 := ctx.Get("request")

									if ok1 && ok2 {
										responseArgs, err = json.Marshal(responseArgs_)
										requestArgs, err = json.Marshal(requestArgs_)
									}

									// callerPath
									{
										if logger.Config.EnableCallerPath.ERROR {
											_, file, no, ok := runtime.Caller(1)
											if ok {
												folder, _ := os.Getwd()
												callerPath = fmt.Sprintf("-[%s - %d]", strings.Replace(file, strings.Replace(folder, "\\", "/", -1), "", 1), no)
											}
										}
									}

									// Формирование text
									{
										text = fmt.Sprintf("[A][ERROR]-|%s|-|%s|-[%d]-[%s]%s: %s %s {\"response\": %s, \"request\": %s} \n",
											start.Format(logger.Config.TimeFormat),
											formatParams.ClientAddr,
											formatParams.StatusCode,
											formatParams.Latency.String(),
											callerPath,
											formatParams.Method,
											ctx.Request.Host+ctx.Request.URL.Path,
											string(responseArgs),
											string(requestArgs))
									}

									err = logger.WriteLogFile(text)
									if err != nil {
										return
									}
								}
							}()
						}
					}
				}
			}
			logger.Sync.Logger.WaitGroup.Done()
		}()

		ctx.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteLogFile - записать лог в лог файл.
func (logger *Logger) WriteLogFile(msg interface{}) error {
	var file_ *os.File
	var err error

	logger.Mx.Lock()
	defer logger.Mx.Unlock()

	// Подготовка
	{
		year, month, day := time.Now().Date()
		path_ := path.Join(logger.Config.LogFilePath, fmt.Sprintf("%d_%s_%d", year, month.String(), day))

		// Проверка существования
		{
			exist, err := exists(path_)
			if err != nil {
				logger.GlobalLogger.ERROR(base_logger.Message{
					Sender: logger.Title,
					Text:   err.Error(),
				})
				return err
			}

			if !exist {
				err := os.Mkdir(path_, 0755)
				if err != nil {
					return err
				}
			}
		}

		// Открытие файла
		{
			file_, err = os.OpenFile(path.Join(path_, fmt.Sprintf("%d_%s_%d_(%d).log", year, month.String(), day, Number)),
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logger.GlobalLogger.ERROR(base_logger.Message{
					Sender: logger.Title,
					Text:   err.Error(),
				})

				return err
			}
		}

		// Закрытие файла
		{
			defer func() {
				err = file_.Close()
				if err != nil {
					logger.GlobalLogger.ERROR(base_logger.Message{
						Sender: logger.Title,
						Text:   err.Error(),
					})
				}
			}()
		}
	}

	// Запись в файл
	{
		switch msg_ := msg.(type) {
		case string:
			if _, err := file_.WriteString(msg_); err != nil {
				logger.GlobalLogger.ERROR(base_logger.Message{
					Sender: logger.Title,
					Text:   err.Error(),
				})

				return err
			}
		}
	}

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

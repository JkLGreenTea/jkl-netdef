package bracnhes

import (
	my_http "JkLNetDef/engine/http/engine"
	"JkLNetDef/engine/http/models/status"
	"JkLNetDef/engine/proxy"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

// InitCaptchaBranch - инициализация ветви капчи.
func InitCaptchaBranch(engine *my_http.Engine, proxyStorage *proxy.Storage) error {
	router := engine.Router.Group("/captcha", "Ветка каптчи. ", "", false)

	// Запросы
	{
		// Страница с каптчей.
		router.GET("/", "", "Страница с каптчей. ", "",
			false, false, false, engine.Modules.SystemAccess.HttpMiddleware(false),
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
					engine.WriteResponseHTML(ctx, http.StatusOK, "captcha.html", gin.H{
						"status": status.Success,
						"code":   http.StatusOK,
						"data":   responseArgs,
					})
				}
			})

		// Получить рандом фото для каптчи
		router.GET("/img", "", " Получить рандом фото для каптчи. ", "",
			false, false, false, engine.Modules.SystemAccess.HttpMiddleware(false),
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

				ctx.FileAttachment(engine.Config.Engine.CaptchaImgDir+fmt.Sprintf("%d.jpg", rand.Intn(399)), "puzzle.jpg")
			})
	}

	return nil
}

package interfacies

import "github.com/gin-gonic/gin"

type HttpServerLogger interface {
	HttpMiddleware() gin.HandlerFunc
	WriteLogFile(msg interface{}) error
}

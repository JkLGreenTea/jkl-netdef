package http

import (
	"JkLNetDef/engine/modules/base_logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	TypeJSON         = "json"
	TypeJSONP        = "json_np"
	TypePureJSON     = "pure_json"
	TypeSecureJSON   = "secure_json"
	TypeIndentedJSON = "indented_json"
	TypeAsciiJSON    = "ascii_json"
	TypeYAML         = "yaml"
	TypeString       = "string"
)

// WriteResponse - запись ответа.
func (engine *Engine) WriteResponse(ctx *gin.Context, statusCode int, response map[string]interface{}) {
	var responseType string

	// Чтение ?responseType=
	{
		responseType = ctx.Request.URL.Query().Get("responseType")
	}

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
		ctx.JSON(statusCode, response)
	}

	buff, err := json.Marshal(response)
	if err != nil {
		engine.Modules.Logger.WARN(base_logger.Message{
			Sender: engine.Title,
			Text:   err.Error(),
		})
	}
	ctx.Set("response", string(buff))
}

// WriteResponseByType - запись ответа определенного типа.
func (engine *Engine) WriteResponseByType(ctx *gin.Context, responseType string, statusCode int, response gin.H) {
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
		ctx.JSON(statusCode, response)
	}

	buff, err := json.Marshal(response)
	if err != nil {
		engine.Modules.Logger.WARN(base_logger.Message{
			Sender: engine.Title,
			Text:   err.Error(),
		})
	}
	ctx.Set("response", string(buff))
}

// WriteResponseHTML - запись HTML ответа.
func (engine *Engine) WriteResponseHTML(ctx *gin.Context, statusCode int, name string, response gin.H) {
	buff, err := json.Marshal(response)
	if err != nil {
		engine.Modules.Logger.WARN(base_logger.Message{
			Sender: engine.Title,
			Text:   err.Error(),
		})
	}
	ctx.Set("response", string(buff))

	ctx.HTML(statusCode, name, response)
}

// WriteResponsePHP - запись php ответа.
func (engine *Engine) WriteResponsePHP(ctx *gin.Context, statusCode int, response []byte) {
	ctx.Data(statusCode, "text/html; charset=utf-8", response)
}

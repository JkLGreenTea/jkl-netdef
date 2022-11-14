package http

import (
	"JkLNetDef/engine/modules/base_logger"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strings"
)

// ReadJSON - прочитать JSON.
func (engine *Engine) ReadJSON(ctx *gin.Context, structure interface{}) error {
	// Чтение данных
	{
		if err := ctx.ShouldBindJSON(structure); err != nil && err.Error() != "EOF" {
			engine.Modules.Logger.WARN(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			return err
		}

		ctx.Set("request", structure)
	}

	return nil
}

// ReadUri - прочитать Uri.
func (engine *Engine) ReadUri(ctx *gin.Context, structure interface{}) error {
	// Чтение данных
	{
		if err := ctx.ShouldBindUri(structure); err != nil && err.Error() != "EOF" {
			engine.Modules.Logger.WARN(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			return err
		}
	}

	return nil
}

// ReadQuery - прочитать Query.
func (engine *Engine) ReadQuery(ctx *gin.Context, structure interface{}) error {
	// Чтение данных
	{
		if err := ctx.ShouldBindQuery(structure); err != nil && err.Error() != "EOF" {
			engine.Modules.Logger.WARN(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			return err
		}

		ctx.Set("request", structure)
	}

	return nil
}

// ReadYAML - прочитать YAML.
func (engine *Engine) ReadYAML(ctx *gin.Context, structure interface{}) error {
	// Чтение данных
	{
		if err := ctx.ShouldBindYAML(structure); err != nil && err.Error() != "EOF" {
			engine.Modules.Logger.WARN(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			return err
		}

		ctx.Set("request", structure)
	}

	return nil
}

// ReadFormMultipart - прочитать FormMultipart.
func (engine *Engine) ReadFormMultipart(ctx *gin.Context, structure interface{}) error {
	// Чтение данных
	{
		if err := ctx.ShouldBindWith(structure, binding.FormMultipart); err != nil && err.Error() != "EOF" {
			engine.Modules.Logger.WARN(base_logger.Message{
				Sender: engine.Title,
				Text:   err.Error(),
			})

			return err
		}

		ctx.Set("request", structure)
	}

	return nil
}

// ReadRequest - прочитать запрос.
func (engine *Engine) ReadRequest(ctx *gin.Context, structure interface{}) error {

	for _, contentType := range strings.Split(ctx.Request.Header.Get("Content-Type"), ";") {
		switch contentType {
		case "application/json", "text/json":
			err := engine.ReadUri(ctx, structure)
			if err != nil {
				return err
			}
			return engine.ReadJSON(ctx, structure)
		case "application/yaml", "text/yaml":
			err := engine.ReadUri(ctx, structure)
			if err != nil {
				return err
			}
			return engine.ReadYAML(ctx, structure)
		case "multipart/form-data":
			err := engine.ReadUri(ctx, structure)
			if err != nil {
				return err
			}
			return engine.ReadFormMultipart(ctx, structure)
		case "":
			err := engine.ReadUri(ctx, structure)
			if err != nil {
				return err
			}
			return nil
		default:
			return errors.New("Недопустимое значение 'Content-Type'. ")
		}
	}

	return errors.New("Недопустимое значение 'Content-Type'. ")
}

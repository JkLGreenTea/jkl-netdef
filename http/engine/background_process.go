package http

import (
	"JkLNetDef/engine/modules/base_logger"
)

// loadDefaultBackgroundProcess - загрузка стандартных фоновых процессов.
func (engine *Engine) loadDefaultBackgroundProcess() error {
	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Начата загрузка стандартных фоновых процессов... ",
	})

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: engine.Title,
		Text:   "Загрузка стандартных фоновых процессов завершена. ",
	})

	return nil
}

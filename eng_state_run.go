package engine

import (
	"JkLNetDef/engine/config"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"JkLNetDef/engine/proxy"
	"errors"
	"fmt"
)

// engineRun - состояние движка - запущен.
type engineStateRun struct {
	*Engine
}

// run - запуск движка.
func (engine *engineStateRun) run() error {
	return errors.New("Движок уже запущен.. ")
}

// stop - остановка движка.
func (engine *engineStateRun) stop() error {
	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: "ENGINE",
		Text:   "Завершение работы движка... ",
	})

	// Остановка прокси
	{
		prxs, err := engine.ProxyStorage.GetAll()
		if err != nil {
			engine.Modules.Logger.ERROR(err.Error())
			return err
		}

		for _, prx := range prxs {
			if err := prx.Stop(); err != nil {
				engine.Modules.Logger.ERROR(err.Error())
				continue
			}
		}
	}

	// Остановка анализа контроллера репутации
	{
		engine.Modules.ControllerReputation.StopAnalise()
	}

	// Остановка API
	{
		if err := engine.Api.Http.Stop(); err != nil {
			engine.Modules.Logger.ERROR(err.Error())
			return err
		}
	}

	// Отключение бд
	{
		if err := engine.Databases.Disconnect(); err != nil {
			engine.Modules.Logger.ERROR(err.Error())
		}
	}

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: "ENGINE",
		Text:   "Движок остановлен. ",
	})

	return nil
}

// newProxy - создание прокси.
func (engine *engineStateRun) newProxy(cfg *config.ProxyConfig) (interfacies.HttpProxy, error) {
	prx, err := proxy.New(cfg, engine.Modules.Blocker, engine.Utils.Synchronizer, engine.Modules.ControllerReputation,
		engine.Modules.Logger, engine.Services)
	if err != nil {
		engine.Modules.Logger.ERROR(err.Error())
		return nil, err
	}

	// Сохранение прокси
	{
		if err = engine.ProxyStorage.Add(prx); err != nil {
			engine.Modules.Logger.ERROR(err.Error())
			return nil, err
		}

		engine.Config.Proxies[prx.GetID().String()] = prx.GetConfig()
	}

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: "ENGINE",
		Text:   fmt.Sprintf("Прокси '%s %s' по протоколу %s создан. ", prx.GetTitle(), prx.GetType(), prx.GetProtocol()),
	})

	return prx, nil
}

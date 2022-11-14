package engine

import (
	"JkLNetDef/engine/config"
	"JkLNetDef/engine/http_api/bracnhes"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"JkLNetDef/engine/proxy"
	"errors"
	"fmt"
	"github.com/xlab/closer"
	"os"
)

type engineStateStop struct {
	*Engine
}

// run - запуск движка.
func (engine *engineStateStop) run() error {
	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: "ENGINE",
		Text:   "Запуск движка... ",
	})

	// Бинд функции остановки
	{
		closer.Bind(func() {
			err := engine.currentState.stop()
			if err != nil {
				engine.Modules.Logger.FATAL(base_logger.Message{
					Sender: "ENGINE",
					Text:   err.Error(),
				})
			}

			os.Exit(0)
		})
	}

	// Запуск контроллера репутации
	{
		engine.Modules.ControllerReputation.Analise()
	}

	// Запуск API
	{
		if err := engine.runApi(); err != err {
			engine.Modules.Logger.ERROR(err.Error())
			return err
		}
	}

	// Запуск слушателей прокси
	{
		proxes, err := engine.ProxyStorage.GetAll()
		if err != nil {
			engine.Modules.Logger.ERROR(err.Error())
			return err
		}

		engine.Utils.Synchronizer.Proxy.WaitGroup.Add(len(proxes))

		for _, un := range proxes {
			go func(un interfacies.HttpProxy) {
				defer engine.Utils.Synchronizer.Proxy.WaitGroup.Done()

				err = un.Listen()
				if err != nil {
					engine.Modules.Logger.ERROR(err.Error())
					return
				}
			}(un)
		}
	}

	engine.currentState = engine.states.Run

	engine.Modules.Logger.INFO(base_logger.Message{
		Sender: "ENGINE",
		Text:   "Движок запущен... ",
	})

	// Ожидание завершения работы
	{
		closer.Hold()
	}

	return nil
}

// runApi - запуск API.
func (engine *engineStateStop) runApi() error {
	// Инициализация ветвей
	{
		if err := bracnhes.InitMainBranch(engine.Api.Http, engine.ProxyStorage); err != nil {
			engine.Modules.Logger.ERROR(err.Error())
			return err
		}

		if err := bracnhes.InitCaptchaBranch(engine.Api.Http, engine.ProxyStorage); err != nil {
			engine.Modules.Logger.ERROR(err.Error())
			return err
		}

		if err := bracnhes.InitHttpProxyBranch(engine.Api.Http, engine.ProxyStorage); err != nil {
			engine.Modules.Logger.ERROR(err.Error())
			return err
		}
	}

	go func() {
		if err := engine.Api.Http.Run(); err != err {
			engine.Modules.Logger.ERROR(err.Error())
		}
	}()

	return nil
}

// stop - остановка движка.
func (engine *engineStateStop) stop() error {
	err := errors.New("Движок уже остановлен. ")
	return err
}

// newProxy - создание прокси.
func (engine *engineStateStop) newProxy(cfg *config.ProxyConfig) (interfacies.HttpProxy, error) {
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

package base_system_access

import (
	"JkLNetDef/engine/http/models/system/system_access/module"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllModules - получить все модули.
func (system *SystemAccess) GetAllModules() ([]*module.Module, string, error) {
	modules, err := system.Databases.Mongo.SystemAccess.Modules.GetAll()
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения модулей. ", err
	}

	return modules, "", nil
}

// GetListModules - получить список всех модулей.
func (system *SystemAccess) GetListModules(search, authorized, locked, relevance string,
	limit, skip int) ([]*module.Module, int64, string, error) {
	var collections int64
	var modules []*module.Module
	var err error

	// Получение
	{
		modules, err = system.Databases.Mongo.SystemAccess.Modules.GetAll()
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, 0, "Ошибка получения модулей. ", err
		}
	}

	// authorized
	{
		if authorized != "" {
			modules = system.filterModulesByAuthorized(modules, authorized)
		}
	}

	// locked
	{
		if locked != "" {
			modules = system.filterModulesByLocked(modules, locked)
		}
	}

	// Релевантность
	{
		if relevance == "" {
			relevance = "new"
		}

		modules = system.filterModulesByRelevance(modules, relevance)
	}

	// Поиск (search)
	{
		if search != "" {
			modules = system.searchModules(modules, search)
		}
	}

	// collections
	{
		collections = int64(len(modules))
	}

	// Обрезать результат (skip/limit)
	{
		if skip > 0 || limit > 0 {
			modules = system.useSkipAndLimitModules(modules, skip, limit)
		}
	}

	return modules, collections, "", nil
}

// GetSelectListModules - получить список всех модулей для селектов.
func (system *SystemAccess) GetSelectListModules(search, authorized, locked,
	relevance string) ([]*module.Module, int64, string, error) {
	var collections int64
	var modules []*module.Module
	var err error

	// Получение
	{
		modules, err = system.Databases.Mongo.SystemAccess.Modules.GetAll()
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, 0, "Ошибка получения модулей. ", err
		}
	}

	// authorized
	{
		if authorized != "" {
			modules = system.filterModulesByAuthorized(modules, authorized)
		}
	}

	// locked
	{
		if locked != "" {
			modules = system.filterModulesByLocked(modules, locked)
		}
	}

	// Релевантность
	{
		if relevance == "" {
			relevance = "new"
		}

		modules = system.filterModulesByRelevance(modules, relevance)
	}

	// Поиск (search)
	{
		if search != "" {
			modules = system.searchModules(modules, search)
		}
	}

	// collections
	{
		collections = int64(len(modules))
	}

	// Откинуть лишние данные
	{
		for _, mod := range modules {
			mod.Meta = nil
			mod.Description = ""
			mod.HttpRequests = nil
		}
	}

	return modules, collections, "", nil
}

// GetModuleByID - получить модуль по ID.
func (system *SystemAccess) GetModuleByID(id primitive.ObjectID) (*module.Module, string, error) {
	mod, err := system.Databases.Mongo.SystemAccess.Modules.GetByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения модуля. ", err
	}

	return mod, "", nil
}

// NewModule - создать новый модуль.
func (system *SystemAccess) NewModule(ctx context.Context, title, description string, locked, authorized bool,
	requests []string) (*module.Module, string, error) {
	mod := &module.Module{
		ID:           primitive.NewObjectID(),
		Title:        title,
		Description:  description,
		Locked:       locked,
		Authorized:   authorized,
		HttpRequests: requests,
	}

	// Meta
	{
		meta, message, err := system.Modules.ManagerMetadata.NewMeta(ctx)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, message, err
		}

		mod.Meta = meta
	}

	// Заполнение пустых полей
	{
		if mod.HttpRequests == nil {
			mod.HttpRequests = make([]string, 0)
		}
	}

	// Сохранение
	{
		err := system.Databases.Mongo.SystemAccess.Modules.Add(mod)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, "Ошибка добавления модуля. ", err
		}
	}

	return mod, "", nil
}

// RemoveModuleByID - удалить модуль по ID.
func (system *SystemAccess) RemoveModuleByID(id primitive.ObjectID) (string, error) {
	err := system.Databases.Mongo.SystemAccess.Modules.RemoveByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return "Ошибка удаления модуля. ", err
	}

	return "", nil
}

// UpdateModuleByID - обновить данные модуля по ID.
func (system *SystemAccess) UpdateModuleByID(id primitive.ObjectID, ctx context.Context, title, description *string,
	locked, authorized *bool, requests []string) (string, error) {
	mod, message, err := system.GetModuleByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return message, err
	}

	system.replacingModuleData(mod, title, description, locked, authorized, requests)

	// Meta
	{
		message, err = system.Modules.ManagerMetadata.UpdateMeta(ctx, mod.Meta)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return message, err
		}
	}

	// Обновление
	{
		err = system.Databases.Mongo.SystemAccess.Modules.UpdateByID(mod)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return "Ошибка добавления модуля. ", err
		}
	}

	return "", nil
}

package base_system_access

import (
	"JkLNetDef/engine/http/models/system/system_access/http_request"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllHttpRequests - получить все http запросы.
func (system *SystemAccess) GetAllHttpRequests() ([]*http_request.Request, string, error) {
	modules, err := system.Databases.Mongo.SystemAccess.HttpRequests.GetAll()
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения запросов. ", err
	}

	return modules, "", nil
}

// GetSelectListHttpRequests - получить список всех http запросов для селектов.
func (system *SystemAccess) GetSelectListHttpRequests(method, authorized, locked, isSystem, isStatic, search,
	relevance string) ([]*http_request.Request, int64, string, error) {
	var collections int64
	var requests []*http_request.Request
	var err error
	var message string

	// Получение
	{
		requests, message, err = system.GetAllHttpRequests()
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, 0, message, err
		}
	}

	// method
	{
		if method != "" {
			requests = system.filterRequestsByMethod(requests, method)
		}
	}

	// authorized
	{
		if authorized != "" {
			requests = system.filterRequestsByAuthorized(requests, authorized)
		}
	}

	// locked
	{
		if locked != "" {
			requests = system.filterRequestsByLocked(requests, locked)
		}
	}

	// isSystem
	{
		if isSystem != "" {
			requests = system.filterRequestsByIsSystem(requests, isSystem)
		}
	}

	// isStatic
	{
		if isStatic != "" {
			requests = system.filterRequestsByIsStatic(requests, isStatic)
		}
	}

	// Поиск (search)
	{
		if search != "" {
			requests = system.searchRequests(requests, search)
		}
	}

	// collections
	{
		collections = int64(len(requests))
	}

	// Релевантность
	{
		if relevance == "" {
			relevance = "new"
		}

		requests = system.filterRequestsByRelevance(requests, relevance)
	}

	// Откинуть лишние данные
	{
		for _, req := range requests {
			req.Meta = nil
			req.Description = ""
		}
	}

	return requests, collections, "", nil
}

// GetListHttpRequests - получить список всех http запросов.
func (system *SystemAccess) GetListHttpRequests(method, authorized, locked, isSystem, isStatic, search,
	relevance string, limit, skip int) ([]*http_request.Request, int64, string, error) {
	var collections int64
	var requests []*http_request.Request
	var err error
	var message string

	// Получение
	{
		requests, message, err = system.GetAllHttpRequests()
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, 0, message, err
		}
	}

	// method
	{
		if method != "" {
			requests = system.filterRequestsByMethod(requests, method)
		}
	}

	// authorized
	{
		if authorized != "" {
			requests = system.filterRequestsByAuthorized(requests, authorized)
		}
	}

	// locked
	{
		if locked != "" {
			requests = system.filterRequestsByLocked(requests, locked)
		}
	}

	// isSystem
	{
		if isSystem != "" {
			requests = system.filterRequestsByIsSystem(requests, isSystem)
		}
	}

	// isStatic
	{
		if isStatic != "" {
			requests = system.filterRequestsByIsStatic(requests, isStatic)
		}
	}

	// Поиск (search)
	{
		if search != "" {
			requests = system.searchRequests(requests, search)
		}
	}

	// collections
	{
		collections = int64(len(requests))
	}

	// Релевантность
	{
		if relevance == "" {
			relevance = "new"
		}

		requests = system.filterRequestsByRelevance(requests, relevance)
	}

	// Обрезать результат (skip/limit)
	{
		if skip > 0 || limit > 0 {
			requests = system.useSkipAndLimitRequests(requests, skip, limit)
		}
	}

	return requests, collections, "", nil
}

// GetHttpRequestByID - получить http запрос по ID.
func (system *SystemAccess) GetHttpRequestByID(id primitive.ObjectID) (*http_request.Request, string, error) {
	req, err := system.Databases.Mongo.SystemAccess.HttpRequests.GetByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения запроса. ", err
	}

	return req, "", nil
}

// NewHttpRequestByID - создать http запрос по ID.
func (system *SystemAccess) NewHttpRequestByID(ctx context.Context, method, url, version string, locked, authorized, isStatic, isSystem bool,
	info, title, description string) (*http_request.Request, string, error) {

	// Проверка Url
	{
		if url[len(url)-1:len(url)] != "/" {
			url += "/"
		}
	}

	req := &http_request.Request{
		ID:          primitive.NewObjectID(),
		Method:      method,
		URL:         url,
		Version:     version,
		Locked:      locked,
		Authorized:  authorized,
		IsStatic:    isStatic,
		IsSystem:    isSystem,
		Info:        info,
		Title:       title,
		Description: description,
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

		req.Meta = meta
	}

	// Сохранение данных запроса
	{
		err := system.Databases.Mongo.SystemAccess.HttpRequests.Add(req)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, "Ошибка добавления запроса. ", err
		}
	}

	return req, "", nil
}

// RemoveHttpRequestByID - удалить http запрос по ID.
func (system *SystemAccess) RemoveHttpRequestByID(id primitive.ObjectID) (string, error) {
	err := system.Databases.Mongo.SystemAccess.HttpRequests.RemoveByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return "Ошибка удаления запроса. ", err
	}

	return "", nil
}

// UpdateHttpRequestByID - обновить данные http запрос по ID.
func (system *SystemAccess) UpdateHttpRequestByID(id primitive.ObjectID, ctx context.Context, method, url, version *string,
	locked, authorized, isStatic, isSystem *bool, info, title, description *string) (string, error) {
	req, message, err := system.GetHttpRequestByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return message, err
	}

	system.replacingHttpRequestData(req, method, url, version, locked, authorized, isStatic, isSystem, info, title, description)

	// Meta
	{
		message, err = system.Modules.ManagerMetadata.UpdateMeta(ctx, req.Meta)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return message, err
		}
	}

	// Обновление в бд
	{
		err = system.Databases.Mongo.SystemAccess.HttpRequests.UpdateByID(req)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return "Ошибка обновления данных запроса. ", err
		}
	}

	return "", nil
}

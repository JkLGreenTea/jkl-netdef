package base_system_access

import (
	"JkLNetDef/engine/http/models/system/system_access/role"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllRoles - получить все роли.
func (system *SystemAccess) GetAllRoles() ([]*role.Role, string, error) {
	roles, err := system.Databases.Mongo.SystemAccess.Roles.GetAll()
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения ролей. ", err
	}

	return roles, "", nil
}

// GetListRoles - получить список всех ролей.
func (system *SystemAccess) GetListRoles(search, relevance string, limit, skip int) ([]*role.Role, int64, string, error) {
	var collections int64
	var roles []*role.Role
	var err error

	// Получение
	{
		roles, err = system.Databases.Mongo.SystemAccess.Roles.GetAll()
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, 0, "Ошибка получения ролей. ", err
		}
	}

	// Поиск (search)
	{
		if search != "" {
			roles = system.searchRoles(roles, search)
		}
	}

	// relevance
	{
		if relevance == "" {
			relevance = "new"
		}

		roles = system.filterRolesByRelevance(roles, relevance)
	}

	// collections
	{
		collections = int64(len(roles))
	}

	// Обрезать результат (skip/limit)
	{
		if skip > 0 || limit > 0 {
			roles = system.useSkipAndLimitRoles(roles, skip, limit)
		}
	}

	return roles, collections, "", nil
}

// GetRoleByID - получить роль по ID.
func (system *SystemAccess) GetRoleByID(id primitive.ObjectID) (*role.Role, string, error) {
	rl, err := system.Databases.Mongo.SystemAccess.Roles.GetByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения роли. ", err
	}

	return rl, "", nil
}

// NewRole - создать новую роль.
func (system *SystemAccess) NewRole(ctx context.Context, title string, requests []string,
	modules []primitive.ObjectID) (*role.Role, string, error) {
	rl := &role.Role{
		ID:           primitive.NewObjectID(),
		Title:        title,
		HttpRequests: requests,
		Modules:      modules,
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

		rl.Meta = meta
	}

	// Заполнение пустых полей
	{
		if rl.HttpRequests == nil {
			rl.HttpRequests = make([]string, 0)
		}

		if rl.Modules == nil {
			rl.Modules = make([]primitive.ObjectID, 0)
		}
	}

	// Сохранение
	{
		err := system.Databases.Mongo.SystemAccess.Roles.Add(rl)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, "Ошибка добавления роли. ", err
		}
	}

	return rl, "", nil
}

// RemoveRoleByID - удалить роль по ID.
func (system *SystemAccess) RemoveRoleByID(id primitive.ObjectID) (string, error) {
	err := system.Databases.Mongo.SystemAccess.Roles.RemoveByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return "Ошибка удаления роли. ", err
	}

	return "", nil
}

// UpdateRoleByID - обновить данные роли по ID.
func (system *SystemAccess) UpdateRoleByID(id primitive.ObjectID, ctx context.Context,
	title *string, requests []string, modules []primitive.ObjectID) (string, error) {
	rl, message, err := system.GetRoleByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return message, err
	}

	system.replacingRoleData(rl, title, requests, modules)

	// Meta
	{
		message, err = system.Modules.ManagerMetadata.UpdateMeta(ctx, rl.Meta)
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
		err = system.Databases.Mongo.SystemAccess.Roles.UpdateByID(rl)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return "Ошибка добавления роли. ", err
		}
	}

	return "", nil
}

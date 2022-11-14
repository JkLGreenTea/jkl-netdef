package base_manager_notifications

import (
	"JkLNetDef/engine/databases"
	"JkLNetDef/engine/http/models/system/notifications"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Manager - менеджер системы уведомлений.
type Manager struct {
	Title     string
	Databases *databases.Databases
	Modules   *Modules
}

// Modules - модули базы данных.
type Modules struct {
	Logger          interfacies.Logger
	ManagerMetadata interfacies.ManagerMetaData
}

// New - создание нового уведомления.
func (manager *Manager) New(ctx context.Context, title, msg string, sender,
	recipient primitive.ObjectID) (*notifications.Notification, string, error) {
	notif := &notifications.Notification{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Message:   msg,
		Date:      time.Now(),
		Recipient: recipient,
		Sender:    sender,
		IsReading: false,
	}

	// Meta
	{
		meta, message, err := manager.Modules.ManagerMetadata.NewMeta(ctx)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, message, err
		}

		notif.Meta = meta
	}

	// Сохранение
	{
		err := manager.Databases.Mongo.System.Notifications.Add(notif)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, "Ошибка добавления уведомления. ", err
		}
	}

	return notif, "", nil
}

// GetByID - получить уведомление по ID из базы данных.
func (manager *Manager) GetByID(id primitive.ObjectID) (*notifications.Notification, string, error) {
	notif, err := manager.Databases.Mongo.System.Notifications.GetByID(id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения уведомления. ", err
	}

	return notif, "", nil
}

// RemoveByID - удаление уведомления по ID.
func (manager *Manager) RemoveByID(id primitive.ObjectID) (string, error) {
	err := manager.Databases.Mongo.System.Notifications.RemoveByID(id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return "Ошибка удаления уведомления. ", err
	}

	return "", nil
}

// GetAll - получение всех уведомлений.
func (manager *Manager) GetAll(noRead bool) ([]*notifications.Notification, string, error) {
	if noRead {
		notifs, err := manager.Databases.Mongo.System.Notifications.GetAllNoRead()
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, "Ошибка получения уведомлений. ", err
		}

		return notifs, "", nil
	}

	notifs, err := manager.Databases.Mongo.System.Notifications.GetAll()
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения уведомлений. ", err
	}

	return notifs, "", nil
}

// GetList - получение списка всех уведомлений.
func (manager *Manager) GetList(noRead bool, search string, skip, limit int) ([]*notifications.Notification, string, error) {
	notifs, message, err := manager.GetAll(noRead)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return nil, message, err
	}

	notifs = manager.search(notifs, search)
	notifs = manager.useSkipAndLimit(notifs, skip, limit)

	return notifs, "", nil
}

// GetListByRecipientID - получение списка всех моих уведомлений.
func (manager *Manager) GetListByRecipientID(id primitive.ObjectID, noRead bool, search string, skip, limit int) ([]*notifications.Notification, string, error) {
	notifs, message, err := manager.GetAllByRecipientID(noRead, id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return nil, message, err
	}

	// search
	{
		if search != "" {
			notifs = manager.search(notifs, search)
		}
	}

	// skip, limit
	{
		if skip != 0 && limit != 0 {
			notifs = manager.useSkipAndLimit(notifs, skip, limit)
		}
	}

	return notifs, "", nil
}

// GetAllByRecipientID - получить все уведомления по ID получателя из базы данных.
func (manager *Manager) GetAllByRecipientID(noRead bool, id primitive.ObjectID) ([]*notifications.Notification, string, error) {
	if noRead {
		notifs, err := manager.Databases.Mongo.System.Notifications.GetByRecipientIDNoRead(id)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, "Ошибка получения уведомлений. ", err
		}

		return notifs, "", nil
	}

	notifs, err := manager.Databases.Mongo.System.Notifications.GetByRecipientID(id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения уведомлений. ", err
	}

	return notifs, "", nil
}

// UpdateByID - обновить данные уведомления в базе данных.
func (manager *Manager) UpdateByID(ctx context.Context, id primitive.ObjectID, title, message *string,
	date *time.Time, recipient, sender primitive.ObjectID, isReading *bool) (string, error) {
	notif, message_, err := manager.GetByID(id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return message_, err
	}

	manager.replacingData(notif, title, message, date, recipient, sender, isReading)

	// Meta
	{
		message_, err = manager.Modules.ManagerMetadata.UpdateMeta(ctx, notif.Meta)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return message_, err
		}
	}

	// Обновление
	{
		err = manager.Databases.Mongo.System.Notifications.UpdateByID(notif)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return "Ошибка обновления данных уведомления. ", err
		}
	}

	return "", nil
}

// UpdateByID - обновить данные уведомления в базе данных.
func (manager *Manager) Read(id primitive.ObjectID, ctx context.Context) (string, error) {
	notif, message_, err := manager.GetByID(id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return message_, err
	}

	notif.IsReading = true

	// Meta
	{
		message_, err = manager.Modules.ManagerMetadata.UpdateMeta(ctx, notif.Meta)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return message_, err
		}
	}

	// Обновление
	{
		err = manager.Databases.Mongo.System.Notifications.UpdateByID(notif)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return "Ошибка обновления данных уведомления. ", err
		}
	}

	return "", nil
}

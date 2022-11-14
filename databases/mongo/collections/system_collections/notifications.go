package system_collections

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	notifications2 "JkLNetDef/engine/http/models/system/notifications"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SysNotifications - запросы системы уведомлений.
type SysNotifications struct {
	Title  string        // Название
	Client *mongo.Client // Клиент бд

	Config *databases_cfg.MongoDB // Конфиг
	Logger interfacies.Logger     // Логгер
}

// Add - добавление уведомления.
func (database *SysNotifications) Add(notif *notifications2.Notification) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionNotifications)

	_, err := collection.InsertOne(context.TODO(), notif)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveByID - удаление уведомления по ID.
func (database *SysNotifications) RemoveByID(id primitive.ObjectID) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionNotifications)

	filter := bson.D{
		{"_id", id},
	}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// GetByID - получить уведомление по ID из базы данных.
func (database *SysNotifications) GetByID(id primitive.ObjectID) (*notifications2.Notification, error) {
	notif := new(notifications2.Notification)
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionNotifications)

	filter := bson.D{
		{"_id", id},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(notif)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return notif, nil
}

// GetByRecipientID - получить уведомления по ID получателя из базы данных.
func (database *SysNotifications) GetByRecipientID(id primitive.ObjectID) ([]*notifications2.Notification, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionNotifications)

	notifs := make([]*notifications2.Notification, 0)

	filter := bson.D{
		{"recipient", id},
	}
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var notif *notifications2.Notification
		if err = cursor.Decode(&notif); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		notifs = append(notifs, notif)
	}

	return notifs, nil
}

// GetByRecipientIDNoRead - получить не прочитанные уведомления по ID получателя из базы данных.
func (database *SysNotifications) GetByRecipientIDNoRead(id primitive.ObjectID) ([]*notifications2.Notification, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionNotifications)

	notifs := make([]*notifications2.Notification, 0)

	filter := bson.D{
		{"recipient", id},
		{"is_reading", false},
	}
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var notif *notifications2.Notification
		if err = cursor.Decode(&notif); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		notifs = append(notifs, notif)
	}

	return notifs, nil
}

// GetAll - получение всех уведомлений.
func (database *SysNotifications) GetAll() ([]*notifications2.Notification, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionNotifications)

	notifs := make([]*notifications2.Notification, 0)

	filter := bson.D{}
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var notif *notifications2.Notification
		if err = cursor.Decode(&notif); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		notifs = append(notifs, notif)
	}

	return notifs, nil
}

// GetAllNoRead - получение всех не прочитанных уведомлений.
func (database *SysNotifications) GetAllNoRead() ([]*notifications2.Notification, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionNotifications)

	notifs := make([]*notifications2.Notification, 0)

	filter := bson.D{
		{"is_reading", false},
	}
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var notif *notifications2.Notification
		if err = cursor.Decode(&notif); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		notifs = append(notifs, notif)
	}

	return notifs, nil
}

// UpdateByID - обновить данные уведомления в базе данных.
func (database *SysNotifications) UpdateByID(notif *notifications2.Notification) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionNotifications)

	filter := bson.D{
		{"_id", notif.ID},
	}

	_, err := collection.ReplaceOne(context.TODO(), filter, notif)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

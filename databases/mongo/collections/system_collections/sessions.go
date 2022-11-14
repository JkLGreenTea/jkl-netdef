package system_collections

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	session2 "JkLNetDef/engine/http/models/system/session"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SysSessions - запросы менеджера сессии.
type SysSessions struct {
	Title  string        // Название
	Client *mongo.Client // Клиент бд

	Config *databases_cfg.MongoDB // Конфиг
	Logger interfacies.Logger     // Логгер
}

// Add - добавление сессии.
func (database *SysSessions) Add(sess *session2.Session) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionSessions)

	_, err := collection.InsertOne(context.TODO(), sess)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveByID - удаление сессии по ID.
func (database *SysSessions) RemoveByID(id primitive.ObjectID) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionSessions)

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

// GetByID - получить сессию по ID из базы данных.
func (database *SysSessions) GetByID(id primitive.ObjectID) (*session2.Session, error) {
	sess := new(session2.Session)
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionSessions)

	filter := bson.D{
		{"_id", id},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(sess)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return sess, nil
}

// GetByTokenID - получить сессию по ID токена из базы данных.
func (database *SysSessions) GetByTokenID(id primitive.ObjectID) (*session2.Session, error) {
	sess := new(session2.Session)
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionSessions)

	filter := bson.D{
		{"token", id},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(sess)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return sess, nil
}

// UpdateByID - обновление сессии.
func (database *SysSessions) UpdateByID(sess *session2.Session) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionSessions)

	filter := bson.D{
		{"_id", sess.ID},
	}

	_, err := collection.ReplaceOne(context.TODO(), filter, sess)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// GetAll - получение всех сессий.
func (database *SysSessions) GetAll() ([]*session2.Session, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionSessions)

	requests := make([]*session2.Session, 0)

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
		var sess *session2.Session
		if err = cursor.Decode(&sess); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		requests = append(requests, sess)
	}

	return requests, nil
}

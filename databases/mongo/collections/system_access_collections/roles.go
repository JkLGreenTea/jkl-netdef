package system_access_collections

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	role2 "JkLNetDef/engine/http/models/system/system_access/role"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SysRoles - запросы ролей.
type SysRoles struct {
	Title  string        // Название
	Client *mongo.Client // Клиент бд

	Config *databases_cfg.MongoDB // Конфиг
	Logger interfacies.Logger     // Логгер
}

// Add - добавление роли.
func (database *SysRoles) Add(rl *role2.Role) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionRoles)

	_, err := collection.InsertOne(context.TODO(), rl)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveByID - удаление роли по ID.
func (database *SysRoles) RemoveByID(id primitive.ObjectID) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionRoles)

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

// GetByID - получить роль по ID из базы данных.
func (database *SysRoles) GetByID(id primitive.ObjectID) (*role2.Role, error) {
	rl := new(role2.Role)
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionRoles)

	filter := bson.D{
		{"_id", id},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(rl)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return rl, nil
}

// GetAll - получение всех ролей.
func (database *SysRoles) GetAll() ([]*role2.Role, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionRoles)

	requests := make([]*role2.Role, 0)

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
		var rl *role2.Role
		if err = cursor.Decode(&rl); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		requests = append(requests, rl)
	}

	return requests, nil
}

// UpdateByID - обновить данные роли в базе данных.
func (database *SysRoles) UpdateByID(rl *role2.Role) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionRoles)

	filter := bson.D{
		{"_id", rl.ID},
	}

	_, err := collection.ReplaceOne(context.TODO(), filter, rl)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// GetCollections - получить кол-во ролей в системе.
func (database *SysRoles) GetCollections() (int64, error) {
	var collections int64

	filter := bson.D{}
	ctx := context.TODO()

	collections, err := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionRoles).CountDocuments(ctx, filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return 0, err
	}

	return collections, nil
}

// GetSlice - получение отрезка ролей из базы данных.
func (database *SysRoles) GetSlice(skip, limit int64) ([]*role2.Role, error) {
	var users []*role2.Role

	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionRoles)

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	if skip > 0 {
		findOptions.SetSkip(skip)
	}
	filter := bson.D{}
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, filter, findOptions)
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
		var rl *role2.Role
		if err = cursor.Decode(&rl); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		users = append(users, rl)
	}

	return users, nil
}

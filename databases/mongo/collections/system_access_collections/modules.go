package system_access_collections

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	module2 "JkLNetDef/engine/http/models/system/system_access/module"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SysModules - запросы модулей.
type SysModules struct {
	Title  string        // Название
	Client *mongo.Client // Клиент бд

	Config *databases_cfg.MongoDB // Конфиг
	Logger interfacies.Logger     // Логгер
}

// Add - добавление запроса.
func (database *SysModules) Add(mod *module2.Module) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionModules)

	_, err := collection.InsertOne(context.TODO(), mod)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveByID - удаление запроса по ID.
func (database *SysModules) RemoveByID(id primitive.ObjectID) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionModules)

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

// GetByID - получить запрос по ID из базы данных.
func (database *SysModules) GetByID(id primitive.ObjectID) (*module2.Module, error) {
	mod := new(module2.Module)
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionModules)

	filter := bson.D{
		{"_id", id},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(mod)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return mod, nil
}

// GetAll - получение всех запросов.
func (database *SysModules) GetAll() ([]*module2.Module, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionModules)

	modules := make([]*module2.Module, 0)

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
		var mod *module2.Module
		if err = cursor.Decode(&mod); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		modules = append(modules, mod)
	}

	return modules, nil
}

// UpdateByID - обновить данные запроса в базе данных.
func (database *SysModules) UpdateByID(mod *module2.Module) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionModules)

	filter := bson.D{
		{"_id", mod.ID},
	}

	_, err := collection.ReplaceOne(context.TODO(), filter, mod)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// GetSlice - получение отрезка запросов из базы данных.
func (database *SysModules) GetSlice(skip, limit int64) ([]*module2.Module, error) {
	var modules []*module2.Module

	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionModules)

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
		var mod *module2.Module
		if err = cursor.Decode(&mod); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		modules = append(modules, mod)
	}

	return modules, nil
}

// GetCollections - получить кол-во запросов в системе.
func (database *SysModules) GetCollections() (int64, error) {
	var collections int64

	filter := bson.D{}
	ctx := context.TODO()

	collections, err := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionModules).CountDocuments(ctx, filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return 0, err
	}

	return collections, nil
}

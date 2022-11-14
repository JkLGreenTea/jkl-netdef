package system_access_collections

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	http_request2 "JkLNetDef/engine/http/models/system/system_access/http_request"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SysHttpRequests - запросы запросов.
type SysHttpRequests struct {
	Title  string        // Название
	Client *mongo.Client // Клиент бд

	Config *databases_cfg.MongoDB // Конфиг
	Logger interfacies.Logger     // Логгер
}

// Add - добавление запроса.
func (database *SysHttpRequests) Add(req *http_request2.Request) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionHttpRequests)

	_, err := collection.InsertOne(context.TODO(), req)
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
func (database *SysHttpRequests) RemoveByID(id primitive.ObjectID) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionHttpRequests)

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
func (database *SysHttpRequests) GetByID(id primitive.ObjectID) (*http_request2.Request, error) {
	req := new(http_request2.Request)
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionHttpRequests)

	filter := bson.D{
		{"_id", id},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(req)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return req, nil
}

// GetByURLAndMethod - получить запрос по методу и url из базы данных.
func (database *SysHttpRequests) GetByURLAndMethod(method, url string) (*http_request2.Request, error) {
	req := new(http_request2.Request)
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionHttpRequests)

	filter := bson.D{
		{"method", method},
		{"url", url},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(req)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return req, nil
}

// GetAll - получение всех запросов.
func (database *SysHttpRequests) GetAll() ([]*http_request2.Request, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionHttpRequests)

	requests := make([]*http_request2.Request, 0)

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
		var req *http_request2.Request
		if err = cursor.Decode(&req); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// UpdateByID - обновить данные запроса в базе данных.
func (database *SysHttpRequests) UpdateByID(req *http_request2.Request) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionHttpRequests)

	filter := bson.D{
		{"_id", req.ID},
	}

	_, err := collection.ReplaceOne(context.TODO(), filter, req)
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
func (database *SysHttpRequests) GetSlice(skip, limit int64) ([]*http_request2.Request, error) {
	var requests []*http_request2.Request

	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionHttpRequests)

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
		var req *http_request2.Request
		if err = cursor.Decode(&req); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// GetCollections - получить кол-во запросов в системе.
func (database *SysHttpRequests) GetCollections() (int64, error) {
	var collections int64

	filter := bson.D{}
	ctx := context.TODO()

	collections, err := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionHttpRequests).CountDocuments(ctx, filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return 0, err
	}

	return collections, nil
}

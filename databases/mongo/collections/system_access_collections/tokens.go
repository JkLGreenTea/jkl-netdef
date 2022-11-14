package system_access_collections

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	token2 "JkLNetDef/engine/http/models/system/system_access/token"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SysTokens - запросы токенов.
type SysTokens struct {
	Title  string        // Название
	Client *mongo.Client // Клиент бд

	Config *databases_cfg.MongoDB // Конфиг
	Logger interfacies.Logger     // Логгер
}

// Add - добавление токена.
func (database *SysTokens) Add(tok *token2.Token) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionTokens)

	_, err := collection.InsertOne(context.TODO(), tok)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveByID - удаление токена по ID.
func (database *SysTokens) RemoveByID(id primitive.ObjectID) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionTokens)

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

// GetAllUserTokens - получение всех токенов пользователя.
func (database *SysTokens) GetAllUserTokens(id primitive.ObjectID) ([]*token2.Token, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionTokens)

	tokens := make([]*token2.Token, 0)

	filter := bson.D{
		{"owner", id},
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
		var tok *token2.Token
		if err = cursor.Decode(&tok); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		tokens = append(tokens, tok)
	}

	return tokens, nil
}

// GetByID - получить токен по ID из базы данных.
func (database *SysTokens) GetByID(id primitive.ObjectID) (*token2.Token, error) {
	tok := new(token2.Token)
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionTokens)

	filter := bson.D{
		{"_id", id},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(tok)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return tok, nil
}

// GetByData - получить токен по data параметру токена из базы данных.
func (database *SysTokens) GetByData(data string) (*token2.Token, error) {
	tok := new(token2.Token)
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionTokens)

	filter := bson.D{
		{"data", data},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(tok)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return tok, nil
}

// GetAll - получение всех токенов.
func (database *SysTokens) GetAll() ([]*token2.Token, error) {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionTokens)

	tokens := make([]*token2.Token, 0)

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
		var tok *token2.Token
		if err = cursor.Decode(&tok); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		tokens = append(tokens, tok)
	}

	return tokens, nil
}

// UpdateByID - обновить данные токена в базе данных.
func (database *SysTokens) UpdateByID(tok *token2.Token) error {
	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionTokens)

	filter := bson.D{
		{"_id", tok.ID},
	}

	_, err := collection.ReplaceOne(context.TODO(), filter, tok)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// GetSlice - получение отрезка токенов из базы данных.
func (database *SysTokens) GetSlice(skip, limit int64) ([]*token2.Token, error) {
	var tokens []*token2.Token

	collection := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionTokens)

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
		var tok *token2.Token
		if err = cursor.Decode(&tok); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		tokens = append(tokens, tok)
	}

	return tokens, nil
}

// GetCollections - получить кол-во токенов в системе.
func (database *SysTokens) GetCollections() (int64, error) {
	var collections int64

	filter := bson.D{}
	ctx := context.TODO()

	collections, err := database.Client.Database(database.Config.System.DatabaseName).
		Collection(database.Config.System.CollectionTokens).CountDocuments(ctx, filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return 0, err
	}

	return collections, nil
}

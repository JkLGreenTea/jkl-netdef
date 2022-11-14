package blocker

import (
	"JkLNetDef/engine/models/blocker/token"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetToken - получение токена.
func (database *Blocker) GetToken(data string) (*token.Token, bool, error) {
	tok := new(token.Token)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionTokens)

	filter := bson.D{
		{"data", data},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(tok)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, false, nil
		}
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, false, err
	}

	return tok, true, nil
}

// AddToken - добавление токена.
func (database *Blocker) AddToken(tok *token.Token) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionTokens)

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

// RemoveToken - удаление токена.
func (database *Blocker) RemoveToken(data string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionTokens)

	filter := bson.D{
		{"data", data},
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

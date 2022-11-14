package blocker

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllClientInClientListOnCaptchaCheck - весь клиентов из списка на прохождение капчи.
func (database *Blocker) GetAllClientInClientListOnCaptchaCheck() ([]*blocker_model.ClientOnCaptchaCheck, error) {
	clients := make([]*blocker_model.ClientOnCaptchaCheck, 0)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionClientListOnCaptchaCheck)

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
		cursor.Close(ctx)
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var client_ *blocker_model.ClientOnCaptchaCheck
		if err = cursor.Decode(&client_); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		clients = append(clients, client_)
	}

	return clients, nil
}

// CheckClientInClientListOnCaptchaCheck - проверка клиента на наличие в списке на прохождение капчи.
func (database *Blocker) CheckClientInClientListOnCaptchaCheck(host, userAgent, requestedUrl string) (bool, error) {
	client_ := new(blocker_model.ClientOnCaptchaCheck)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionClientListOnCaptchaCheck)

	filter := bson.D{
		{"host", host},
		{"user_agent", userAgent},
		{"requested_url", requestedUrl},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(client_)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return false, err
	}

	return true, nil
}

// AddClientInClientListOnCaptchaCheck - добавить клиента в список клиентов на прохождение капчи.
func (database *Blocker) AddClientInClientListOnCaptchaCheck(client_ *blocker_model.ClientOnCaptchaCheck) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionClientListOnCaptchaCheck)

	_, err := collection.InsertOne(context.TODO(), client_)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveClientFromClientListOnCaptchaCheck - клиента из списка клиентов на прохождение капчи.
func (database *Blocker) RemoveClientFromClientListOnCaptchaCheck(host, userAgent string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionClientListOnCaptchaCheck)

	filter := bson.D{
		{"host", host},
		{"user_agent", userAgent},
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

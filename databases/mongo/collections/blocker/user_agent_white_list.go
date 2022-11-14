package blocker

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllUserAgentsFromWhiteList - весь список User-Agent белого списка.
func (database *Blocker) GetAllUserAgentsFromWhiteList() (map[string]*blocker_model.UserAgentInWhiteList, error) {
	agents := make(map[string]*blocker_model.UserAgentInWhiteList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionUserAgentWhiteList)

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
		var agent *blocker_model.UserAgentInWhiteList
		if err = cursor.Decode(&agent); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		agents[agent.ID.Hex()] = agent
	}

	return agents, nil
}

// CheckUserAgentInWhiteList - проверка User-Agent на наличие в белом списке.
func (database *Blocker) CheckUserAgentInWhiteList(agent string) (bool, error) {
	agent_ := new(blocker_model.UserAgentInWhiteList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionUserAgentWhiteList)

	filter := bson.D{
		{"title", agent},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(agent_)
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

// AddUserAgentInWhiteList - добавить User-Agent в белый список.
func (database *Blocker) AddUserAgentInWhiteList(agent *blocker_model.UserAgentInWhiteList) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionUserAgentWhiteList)

	_, err := collection.InsertOne(context.TODO(), agent)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveUserAgentFromWhiteList - удалить User-Agent из белого списка.
func (database *Blocker) RemoveUserAgentFromWhiteList(host string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionUserAgentWhiteList)

	filter := bson.D{
		{"host", host},
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

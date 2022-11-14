package blocker

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllUserAgentsFromBlackList - весь список User-Agent черного списка.
func (database *Blocker) GetAllUserAgentsFromBlackList() (map[string]*blocker_model.UserAgentInBlackList, error) {
	agents := make(map[string]*blocker_model.UserAgentInBlackList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionUserAgentBlackList)

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
		var agent *blocker_model.UserAgentInBlackList
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

// CheckUserAgentInBlackList - проверка User-Agent на наличие в черном списке.
func (database *Blocker) CheckUserAgentInBlackList(agent string) (bool, error) {
	agent_ := new(blocker_model.UserAgentInBlackList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionUserAgentBlackList)

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

// AddUserAgentInBlackList - добавить User-Agent в черный список.
func (database *Blocker) AddUserAgentInBlackList(agent *blocker_model.UserAgentInBlackList) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionUserAgentBlackList)

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

// RemoveUserAgentFromBlackList - удалить User-Agent из черного списка.
func (database *Blocker) RemoveUserAgentFromBlackList(agent string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionUserAgentBlackList)

	filter := bson.D{
		{"title", agent},
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

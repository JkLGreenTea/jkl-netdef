package blocker

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllHostsFromBlackList - весь список IP-адресов черного списка.
func (database *Blocker) GetAllHostsFromBlackList() (map[string]*blocker_model.HostFromBlackList, error) {
	hosts := make(map[string]*blocker_model.HostFromBlackList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostBlackList)

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
		var host_ *blocker_model.HostFromBlackList
		if err = cursor.Decode(&host_); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		hosts[host_.Host] = host_
	}

	return hosts, nil
}

// CheckHostInBlackList - проверка IP на наличие в черном списке.
func (database *Blocker) CheckHostInBlackList(host string) (bool, error) {
	host_ := new(blocker_model.HostFromBlackList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostBlackList)

	filter := bson.D{
		{"host", host},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(host_)
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

// AddHostInBlackList - добавить ip в черный список.
func (database *Blocker) AddHostInBlackList(host *blocker_model.HostFromBlackList) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostBlackList)

	_, err := collection.InsertOne(context.TODO(), host)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveHostFromBlackList - удалить IP из черного списка.
func (database *Blocker) RemoveHostFromBlackList(host string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostBlackList)

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

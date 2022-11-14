package blocker

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllHostsFromHardWhiteList - весь список IP-адресов жесткого белого списка.
func (database *Blocker) GetAllHostsFromHardWhiteList() (map[string]*blocker_model.HostFromWhiteList, error) {
	hosts := make(map[string]*blocker_model.HostFromWhiteList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostHardWhiteList)

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
		var host_ *blocker_model.HostFromWhiteList
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

// CheckHostInHardWhiteList - проверка IP на наличие в жестком белом списке.
func (database *Blocker) CheckHostInHardWhiteList(host string) (bool, error) {
	host_ := new(blocker_model.HostFromWhiteList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostHardWhiteList)

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

// AddHostInHardWhiteList - добавить ip в жесткий белый список.
func (database *Blocker) AddHostInHardWhiteList(host *blocker_model.HostFromWhiteList) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostHardWhiteList)

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

// RemoveHostFromHardWhiteList - удалить IP из жесткого белого списка.
func (database *Blocker) RemoveHostFromHardWhiteList(host string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostHardWhiteList)

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

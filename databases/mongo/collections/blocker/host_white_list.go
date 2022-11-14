package blocker

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllHostsFromWhiteList - весь список IP-адресов белого списка.
func (database *Blocker) GetAllHostsFromWhiteList() (map[string]*blocker_model.HostFromWhiteList, error) {
	hosts := make(map[string]*blocker_model.HostFromWhiteList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostWhiteList)

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

// CheckHostInWhiteList - проверка IP на наличие в белом списке.
func (database *Blocker) CheckHostInWhiteList(host string) (bool, error) {
	host_ := new(blocker_model.HostFromWhiteList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostWhiteList)

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

// AddHostInWhiteList - добавить ip в белый список.
func (database *Blocker) AddHostInWhiteList(host *blocker_model.HostFromWhiteList) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostWhiteList)

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

// RemoveHostFromWhiteList - удалить IP из белого списка.
func (database *Blocker) RemoveHostFromWhiteList(host string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostWhiteList)

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

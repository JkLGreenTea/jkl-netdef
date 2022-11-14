package blocker

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllHostsFromBanList - весь список заблокированных IP-адресов.
func (database *Blocker) GetAllHostsFromBanList() (map[string]*blocker_model.HostFromBanList, error) {
	hosts := make(map[string]*blocker_model.HostFromBanList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostBanList)

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
		var host_ *blocker_model.HostFromBanList
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

// CheckHostInBanList - проверка IP на наличие в списке заблокированных.
func (database *Blocker) CheckHostInBanList(host string) (bool, error) {
	host_ := new(blocker_model.HostFromBanList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostBanList)

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

// AddHostInBanList - добавить ip в список заблокированных.
func (database *Blocker) AddHostInBanList(host *blocker_model.HostFromBanList) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostBanList)

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

// RemoveHostFromBanList - удалить IP из списка заблокированных.
func (database *Blocker) RemoveHostFromBanList(host string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionHostBanList)

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

package blocker

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllLocationsFromBlackList - весь список местоположений черного списка.
func (database *Blocker) GetAllLocationsFromBlackList() (map[string]*blocker_model.LocationFromBlackList, error) {
	locations := make(map[string]*blocker_model.LocationFromBlackList)
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
		var location *blocker_model.LocationFromBlackList
		if err = cursor.Decode(&location); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		locations[location.Location] = location
	}

	return locations, nil
}

// CheckLocationInBlackList - проверка страны на наличие в черном списке.
func (database *Blocker) CheckLocationInBlackList(location string) (bool, error) {
	location_ := new(blocker_model.LocationFromBlackList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionLocationBlackList)

	filter := bson.D{
		{"location", location},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(location_)
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

// AddLocationInBlackList - добавить страну в черный список.
func (database *Blocker) AddLocationInBlackList(location *blocker_model.LocationFromBlackList) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionLocationBlackList)

	_, err := collection.InsertOne(context.TODO(), location)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveLocationFromBlackList - удалить страну из черного списка.
func (database *Blocker) RemoveLocationFromBlackList(location string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionLocationBlackList)

	filter := bson.D{
		{"location", location},
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

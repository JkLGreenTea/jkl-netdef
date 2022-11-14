package blocker

import (
	blocker_model "JkLNetDef/engine/models/blocker"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllLocationsFromWhiteList - весь список местоположений белого списка.
func (database *Blocker) GetAllLocationsFromWhiteList() (map[string]*blocker_model.LocationFromWhiteList, error) {
	locations := make(map[string]*blocker_model.LocationFromWhiteList)
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
		var location *blocker_model.LocationFromWhiteList
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

// CheckLocationInWhiteList - проверка страны на наличие в белом списке.
func (database *Blocker) CheckLocationInWhiteList(location string) (bool, error) {
	location_ := new(blocker_model.LocationFromWhiteList)
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionLocationWhiteList)

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

// AddLocationInWhiteList - добавить страну в белый список.
func (database *Blocker) AddLocationInWhiteList(location *blocker_model.LocationFromWhiteList) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionLocationWhiteList)

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

// RemoveLocationFromWhiteList - удалить страну из белого списка.
func (database *Blocker) RemoveLocationFromWhiteList(location string) error {
	collection := database.Client.Database(database.Config.Blocker.DatabaseName).
		Collection(database.Config.Blocker.CollectionLocationWhiteList)

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

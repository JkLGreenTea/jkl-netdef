package main_collections

import (
	databases_cfg "JkLNetDef/engine/config/databases"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/models/user"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Users - запросы пользователей.
type Users struct {
	Title  string        // Название
	Client *mongo.Client // Клиент бд

	Config *databases_cfg.MongoDB // Конфиг
	Logger interfacies.Logger     // Логгер
}

// Add - добавить пользователя в базу данных.
func (database *Users) Add(us *user.User) error {
	collection := database.Client.Database(database.Config.Main.DatabaseName).
		Collection(database.Config.Main.CollectionUser)

	_, err := collection.InsertOne(context.TODO(), us)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// GetByID - получить пользователя по ID из базы данных.
func (database *Users) GetByID(id primitive.ObjectID) (*user.User, error) {
	us := new(user.User)
	collection := database.Client.Database(database.Config.Main.DatabaseName).
		Collection(database.Config.Main.CollectionUser)

	filter := bson.D{
		{"_id", id},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(us)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return us, nil
}

// GetLoginByID - получить логин пользователя по ID из базы данных.
func (database *Users) GetLoginByID(id primitive.ObjectID) (string, error) {
	us := new(user.User)
	collection := database.Client.Database(database.Config.Main.DatabaseName).
		Collection(database.Config.Main.CollectionUser)

	filter := bson.D{
		{"_id", id},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(us)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return "", err
	}

	return us.Login, nil
}

// GetByLogin - получить пользователя по логину из базы данных.
func (database *Users) GetByLogin(login string) (*user.User, error) {
	us := new(user.User)
	collection := database.Client.Database(database.Config.Main.DatabaseName).
		Collection(database.Config.Main.CollectionUser)

	filter := bson.D{
		{"login", login},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(us)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	return us, nil
}

// GetAll - получение всех пользователей из базы данных.
func (database *Users) GetAll() ([]*user.User, error) {
	var users []*user.User

	collection := database.Client.Database(database.Config.Main.DatabaseName).
		Collection(database.Config.Main.CollectionUser)

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
		err := cursor.Close(ctx)
		if err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var us *user.User
		if err = cursor.Decode(&us); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		users = append(users, us)
	}

	return users, nil
}

// GetCollections - получить кол-во пользователей в системе.
func (database *Users) GetCollections() (int64, error) {
	var collections int64

	filter := bson.D{}
	ctx := context.TODO()

	collections, err := database.Client.Database(database.Config.Main.DatabaseName).
		Collection(database.Config.Main.CollectionUser).CountDocuments(ctx, filter)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return 0, err
	}

	return collections, nil
}

// GetSlice - получение отрезка пользователей из базы данных.
func (database *Users) GetSlice(skip, limit int64) ([]*user.User, error) {
	var users []*user.User

	collection := database.Client.Database(database.Config.Main.DatabaseName).
		Collection(database.Config.Main.CollectionUser)

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	if skip > 0 {
		findOptions.SetSkip(skip)
	}
	filter := bson.D{}
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var us *user.User
		if err = cursor.Decode(&us); err != nil {
			database.Logger.WARN(base_logger.Message{
				Sender: database.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		users = append(users, us)
	}

	return users, nil
}

// UpdateByID - изменение информации пользователя по ID.
func (database *Users) UpdateByID(us *user.User) error {
	collection := database.Client.Database(database.Config.Main.DatabaseName).
		Collection(database.Config.Main.CollectionUser)

	filter := bson.M{
		"_id": us.ID,
	}

	_, err := collection.ReplaceOne(context.TODO(), filter, us)
	if err != nil {
		database.Logger.WARN(base_logger.Message{
			Sender: database.Title,
			Text:   err.Error(),
		})

		return err
	}

	return nil
}

// RemoveByID - удаление пользователя по ID.
func (database *Users) RemoveByID(id primitive.ObjectID) error {
	collection := database.Client.Database(database.Config.Main.DatabaseName).
		Collection(database.Config.Main.CollectionUser)

	filter := bson.D{
		{"_id", id},
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

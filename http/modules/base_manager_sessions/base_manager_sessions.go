package base_manager_sessions

import (
	"JkLNetDef/engine/databases"
	session2 "JkLNetDef/engine/http/models/system/session"
	"JkLNetDef/engine/interfacies"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

// Manager - менеджер сессий.
type Manager struct {
	Title     string
	Databases *databases.Databases

	Modules *Modules
}

// Modules - модули менеджера.
type Modules struct {
	ManagerMetadata      interfacies.ManagerMetaData
	ManagerNotifications interfacies.ManagerNotifications
	Logger               interfacies.Logger
}

// session - получение сессии пользователя.
func (manager *Manager) ginSession(ctx *gin.Context) (*session2.Session, error) {
	strTok := ctx.Request.Header.Get("Authorization")

	if strTok == "" {
		dataTok, err := ctx.Cookie("token")
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		tok, err := manager.Databases.Mongo.SystemAccess.Tokens.GetByData(dataTok)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		sess, err := manager.Databases.Mongo.System.Sessions.GetByTokenID(tok.ID)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		return sess, nil
	} else {
		dataTok := strings.Replace(strTok, "Bearer ", "", 1)

		tok, err := manager.Databases.Mongo.SystemAccess.Tokens.GetByData(dataTok)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		sess, err := manager.Databases.Mongo.System.Sessions.GetByTokenID(tok.ID)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, err
		}

		return sess, nil
	}
}

// save - сохранение сессии пользователя.
func (manager *Manager) save(ctx context.Context, sess *session2.Session) error {
	return manager.Databases.Mongo.System.Sessions.UpdateByID(sess)
}

// Get - получение объекта сессии
func (manager *Manager) Get(ctx context.Context, key string) (interface{}, bool, error) {
	switch ctx_ := ctx.(type) {
	case *gin.Context:
		{
			sess, err := manager.ginSession(ctx_)
			if err != nil {
				manager.Modules.Logger.WARN(base_logger.Message{
					Sender: manager.Title,
					Text:   err.Error(),
				})

				return nil, false, err
			}

			obj, ok := sess.Data[key]

			return obj, ok, nil
		}

	}

	return nil, false, errors.New("Недопустимое значение контекста. ")
}

// Set - запись объекта в сессию.
func (manager *Manager) Set(ctx context.Context, key string, obj interface{}) error {
	switch ctx_ := ctx.(type) {
	case *gin.Context:
		{
			sess, err := manager.ginSession(ctx_)
			if err != nil {
				manager.Modules.Logger.WARN(base_logger.Message{
					Sender: manager.Title,
					Text:   err.Error(),
				})

				return err
			}

			sess.Data[key] = obj

			return manager.save(ctx, sess)
		}
	}

	return errors.New("Недопустимое значение контекста. ")
}

// Delete - удаление объекта из сессии.
func (manager *Manager) Delete(ctx context.Context, key string) error {
	switch ctx_ := ctx.(type) {
	case *gin.Context:
		{
			sess, err := manager.ginSession(ctx_)
			if err != nil {
				manager.Modules.Logger.WARN(base_logger.Message{
					Sender: manager.Title,
					Text:   err.Error(),
				})

				return err
			}

			delete(sess.Data, key)

			return manager.save(ctx, sess)
		}
	}

	return errors.New("Недопустимое значение контекста. ")
}

// GetUserID - получить user_id из сессии.
func (manager *Manager) GetUserID(ctx context.Context) (primitive.ObjectID, error) {
	switch ctx_ := ctx.(type) {
	case *gin.Context:
		{
			userID, ok, err := manager.Get(ctx_, "user_id")
			if err != nil {
				return primitive.ObjectID{}, err
			}

			if !ok {
				return primitive.ObjectID{}, errors.New("Ключ 'user_id' не найден в сессии. ")
			}

			docID, err := primitive.ObjectIDFromHex(userID.(string))
			if err != nil {
				return primitive.ObjectID{}, err
			}

			return docID, nil
		}
	}

	return primitive.ObjectID{}, errors.New("Недопустимое значение контекста. ")
}

// GetUserLogin - получить user_login из сессии.
func (manager *Manager) GetUserLogin(ctx context.Context) (string, error) {
	switch ctx_ := ctx.(type) {
	case *gin.Context:
		{
			userLogin, ok, err := manager.Get(ctx_, "user_login")
			if err != nil {
				return "", err
			}

			if !ok {
				return "", errors.New("Ключ 'user_login' не найден в сессии. ")
			}

			return userLogin.(string), nil
		}
	}

	return "", errors.New("Недопустимое значение контекста. ")
}

// GetSessionByID - получить сессию по Id.
func (manager *Manager) GetSessionByID(id primitive.ObjectID) (*session2.Session, string, error) {
	sess, err := manager.Databases.Mongo.System.Sessions.GetByID(id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения сессии. ", err
	}

	return sess, "", nil
}

// GetAllSessions - получить все сессии.
func (manager *Manager) GetAllSessions() ([]*session2.Session, string, error) {
	sessions, err := manager.Databases.Mongo.System.Sessions.GetAll()
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения сессий. ", err
	}

	return sessions, "", nil
}

// GetListSessions - получить список всех сессий.
func (manager *Manager) GetListSessions(search, noDelete, relevance string,
	limit, skip int) ([]*session2.Session, int64, string, error) {
	var collections int64
	var sessions []*session2.Session
	var err error

	// Получение
	{
		sessions, err = manager.Databases.Mongo.System.Sessions.GetAll()
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return nil, 0, "Ошибка получения сессий. ", err
		}
	}

	// no_delete
	{
		if noDelete != "" {
			sessions = manager.filterSessionsByNoDelete(sessions, noDelete)
		}
	}

	// Релевантность
	{
		if relevance == "" {
			relevance = "new"
		}

		sessions = manager.filterSessionsByRelevance(sessions, relevance)
	}

	// Поиск (search)
	{
		if search != "" {
			sessions = manager.searchSessions(sessions, search)
		}
	}

	// collections
	{
		collections = int64(len(sessions))
	}

	// Обрезать результат (skip/limit)
	{
		if skip > 0 || limit > 0 {
			sessions = manager.useSkipAndLimitSessions(sessions, skip, limit)
		}
	}

	return sessions, collections, "", nil
}

// DeleteSessionByID - удаление сессии по id.
func (manager *Manager) DeleteSessionByID(id primitive.ObjectID) (string, error) {
	sess, message, err := manager.GetSessionByID(id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return message, err
	}

	if sess.NoDelete {
		err := errors.New("Удаление сессии запрещено! ")
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return "Удаление сессии запрещено! ", err
	}

	err = manager.Databases.Mongo.System.Sessions.RemoveByID(id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return "Ошибка удаления сессии. ", err
	}

	return "", nil
}

// UpdateSessionByID - обновление данных сесиии по id.
func (manager *Manager) UpdateSessionByID(id primitive.ObjectID, ctx context.Context, noDelete *bool,
	data map[string]interface{}) (string, error) {
	sess, err := manager.Databases.Mongo.System.Sessions.GetByID(id)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return "Ошибка получения сессии. ", err
	}

	// Замена
	{
		if noDelete != nil {
			sess.NoDelete = *noDelete
		}

		if data != nil {
			sess.Data = data
		}
	}

	// Meta
	{
		message, err := manager.Modules.ManagerMetadata.UpdateMeta(ctx, sess.Meta)
		if err != nil {
			manager.Modules.Logger.WARN(base_logger.Message{
				Sender: manager.Title,
				Text:   err.Error(),
			})

			return message, err
		}
	}

	err = manager.Databases.Mongo.System.Sessions.UpdateByID(sess)
	if err != nil {
		manager.Modules.Logger.WARN(base_logger.Message{
			Sender: manager.Title,
			Text:   err.Error(),
		})

		return "Ошибка удаления сессии. ", err
	}

	return "", nil
}

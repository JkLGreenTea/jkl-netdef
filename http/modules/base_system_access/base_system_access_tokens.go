package base_system_access

import (
	"JkLNetDef/engine/http/models/system/system_access/token"
	"JkLNetDef/engine/models/user"
	"JkLNetDef/engine/modules/base_logger"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// ClearTokenFromHttpCookie - очистить токен из http печенек.
func (system *SystemAccess) ClearTokenFromHttpCookie(ctx *gin.Context) error {
	ctx.SetCookie("token", "", 0,
		"/", system.Utils.Config.Domain, false, false)

	return nil
}

// SetTokenInHttpCookie - установить токен в http печеньки.
func (system *SystemAccess) SetTokenInHttpCookie(ctx *gin.Context, tok *token.Token) error {
	ctx.SetCookie("token", tok.Data, int(system.Utils.Config.SystemAccess.Token.LifeTime),
		"/", system.Utils.Config.Domain, false, false)

	return nil
}

// GetTokenFromHttpCookie - получить токен из http печенек.
func (system *SystemAccess) GetTokenFromHttpCookie(ctx *gin.Context) (string, error) {
	return ctx.Cookie("token")
}

// GetAllTokens - получить все токены.
func (system *SystemAccess) GetAllTokens() ([]*token.Token, string, error) {
	tokens, err := system.Databases.Mongo.SystemAccess.Tokens.GetAll()
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Utils.Config.SystemAccess.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения токена. ", err
	}

	return tokens, "", nil
}

// GetListTokens - получить список всех токенов.
func (system *SystemAccess) GetListTokens(search, relevance, noDelete string, limit, skip int) ([]*token.Token, int64, string, error) {
	var collections int64
	var tokens []*token.Token
	var err error

	// Получение
	{
		tokens, err = system.Databases.Mongo.SystemAccess.Tokens.GetAll()
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, 0, "Ошибка получения токена. ", err
		}
	}

	// NoDelete
	{
		if noDelete != "" {
			tokens = system.filterTokensByNoDelete(tokens, noDelete)
		}
	}

	// Поиск (search)
	{
		if search != "" {
			tokens = system.searchTokens(tokens, search)
		}
	}

	// relevance
	{
		if relevance == "" {
			relevance = "new"
		}

		tokens = system.filterTokensByRelevance(tokens, relevance)
	}

	// collections
	{
		collections = int64(len(tokens))
	}

	// Обрезать результат (skip/limit)
	{
		if skip > 0 || limit > 0 {
			tokens = system.useSkipAndLimitTokens(tokens, skip, limit)
		}
	}

	return tokens, collections, "", nil
}

// GetTokenByID - получить токен по ID.
func (system *SystemAccess) GetTokenByID(id primitive.ObjectID) (*token.Token, string, error) {
	tok, err := system.Databases.Mongo.SystemAccess.Tokens.GetByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения токена. ", err
	}

	return tok, "", nil
}

// GetTokenByData - получить токен по data.
func (system *SystemAccess) GetTokenByData(data string) (*token.Token, string, error) {
	tok, err := system.Databases.Mongo.SystemAccess.Tokens.GetByData(data)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения токена. ", err
	}

	return tok, "", nil
}

// GetAllUserTokens - получить все токены пользователя.
func (system *SystemAccess) GetAllUserTokens(id primitive.ObjectID) ([]*token.Token, string, error) {
	tokens, err := system.Databases.Mongo.SystemAccess.Tokens.GetAllUserTokens(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return nil, "Ошибка получения токена. ", err
	}

	return tokens, "", nil
}

// NewToken - создать новую роль.
func (system *SystemAccess) NewToken(ctx context.Context, owner primitive.ObjectID,
	expire time.Time, noDelete bool) (*token.Token, string, error) {
	var tok *token.Token
	var us *user.User

	// Получение пользователя
	{
		var err error

		us, err = system.Databases.Mongo.Main.Users.GetByID(owner)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, "Ошибка получения пользователя. ", err
		}
	}

	// Проверка expire
	{
		if expire.Unix() == 0 {
			expire = time.Unix(int64(system.Utils.Config.SystemAccess.Token.LifeTime), 0)
		}
	}

	// Создание токена
	{
		claims_, data, err := system.Authorizer().SignIn(ctx, us, expire.Unix())
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, "Ошибка создания токена. ", err
		}

		tok = &token.Token{
			ID:       primitive.NewObjectID(),
			Data:     data,
			Owner:    us.ID,
			Created:  claims_.StandardClaims.IssuedAt,
			Expire:   claims_.StandardClaims.ExpiresAt,
			NoDelete: noDelete,
		}
	}

	// Meta
	{
		meta, message, err := system.Modules.ManagerMetadata.NewMeta(ctx)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, message, err
		}

		tok.Meta = meta
	}

	// Сохранение в бд
	{
		err := system.Databases.Mongo.SystemAccess.Tokens.Add(tok)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return nil, "Ошибка сохранения токена. ", err
		}
	}

	return tok, "", nil
}

// RemoveTokenByID - удалить токен по ID.
func (system *SystemAccess) RemoveTokenByID(id primitive.ObjectID) (string, error) {
	tok, message, err := system.GetTokenByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return message, err
	}

	if tok.NoDelete {
		err := errors.New("Удаление токена запрещено! ")
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return "Удаление токена запрещено! ", err
	}

	err = system.Databases.Mongo.SystemAccess.Tokens.RemoveByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return "Ошибка удаления токена. ", err
	}

	return "", nil
}

// UpdateTokenByID - обновить данные роли по ID.
func (system *SystemAccess) UpdateTokenByID(id primitive.ObjectID, ctx context.Context,
	owner primitive.ObjectID, created, expire *time.Time, noDelete *bool) (string, error) {
	tok, message, err := system.GetTokenByID(id)
	if err != nil {
		system.Modules.Logger.WARN(base_logger.Message{
			Sender: system.Title,
			Text:   err.Error(),
		})

		return message, err
	}

	system.replacingTokenData(tok, owner, created, expire, noDelete)

	// Meta
	{
		message, err = system.Modules.ManagerMetadata.UpdateMeta(ctx, tok.Meta)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return message, err
		}
	}

	// Обновление
	{
		err = system.Databases.Mongo.SystemAccess.Tokens.UpdateByID(tok)
		if err != nil {
			system.Modules.Logger.WARN(base_logger.Message{
				Sender: system.Title,
				Text:   err.Error(),
			})

			return "Ошибка обновления токена. ", err
		}
	}

	return "", nil
}
